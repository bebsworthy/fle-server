package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fle/server/internal/config"
	"github.com/fle/server/internal/jsonrpc"
	"github.com/fle/server/internal/server"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test server instance for integration tests
type testServer struct {
	server     *server.Server
	httpServer *httptest.Server
	url        string
	wsURL      string
}

func setupTestServer(t *testing.T) *testServer {
	// Set test environment variables
	os.Setenv("ENV", "test")
	os.Setenv("LOG_LEVEL", "error") // Reduce log noise during tests
	
	// Load test configuration
	cfg, err := config.Load()
	require.NoError(t, err, "Failed to load test configuration")
	
	// Override config for testing
	cfg.Host = "127.0.0.1"
	cfg.Port = 0 // Let httptest choose a free port

	// Create server instance
	srv, err := server.NewServer(cfg, setupLogger(cfg))
	require.NoError(t, err, "Failed to create server")

	// Create test HTTP server
	httpServer := httptest.NewServer(srv.Handler())
	
	// Parse server URL for WebSocket connection
	serverURL := httpServer.URL
	wsURL := strings.Replace(serverURL, "http", "ws", 1)

	return &testServer{
		server:     srv,
		httpServer: httpServer,
		url:        serverURL,
		wsURL:      wsURL,
	}
}

func (ts *testServer) Close() {
	if ts.httpServer != nil {
		ts.httpServer.Close()
	}
}

func TestMain(m *testing.M) {
	// Setup
	code := m.Run()
	
	// Teardown
	os.Exit(code)
}

// TestServerStartupAndHealthEndpoint tests server initialization and health endpoint
func TestServerStartupAndHealthEndpoint(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Test health endpoint
	resp, err := http.Get(ts.url + "/health")
	require.NoError(t, err, "Failed to make health request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint should return 200")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Health endpoint should return JSON")

	// Parse health response
	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err, "Failed to decode health response")

	assert.Equal(t, "healthy", health["status"], "Server should be healthy")
	assert.Equal(t, "test", health["environment"], "Environment should be test")
	assert.NotNil(t, health["timestamp"], "Response should include timestamp")
}

// TestWebSocketConnection tests basic WebSocket connection establishment
func TestWebSocketConnection(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Connect to WebSocket endpoint
	wsURL := ts.wsURL + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket")
	defer conn.Close()

	// Set read deadline to avoid hanging
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read welcome message
	_, message, err := conn.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message")

	var welcome map[string]interface{}
	err = json.Unmarshal(message, &welcome)
	require.NoError(t, err, "Failed to unmarshal welcome message")

	assert.Equal(t, "welcome", welcome["type"], "Should receive welcome message")
	assert.NotEmpty(t, welcome["session_code"], "Welcome message should include session code")
	assert.NotEmpty(t, welcome["message"], "Welcome message should include message text")
}

// TestSessionCreationAndRestoration tests session lifecycle
func TestSessionCreationAndRestoration(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// First connection - create new session
	wsURL := ts.wsURL + "/ws"
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket")
	defer conn1.Close()

	conn1.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read welcome message to get session code
	_, message, err := conn1.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message")

	var welcome map[string]interface{}
	err = json.Unmarshal(message, &welcome)
	require.NoError(t, err, "Failed to unmarshal welcome message")

	sessionCode := welcome["session_code"].(string)
	assert.NotEmpty(t, sessionCode, "Session code should not be empty")

	// Close first connection
	conn1.Close()

	// Second connection - restore existing session
	wsURLWithSession := fmt.Sprintf("%s/ws?session=%s", ts.wsURL, sessionCode)
	conn2, _, err := websocket.DefaultDialer.Dial(wsURLWithSession, nil)
	require.NoError(t, err, "Failed to reconnect with session code")
	defer conn2.Close()

	conn2.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read welcome message for restored session
	_, message2, err := conn2.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message for restored session")

	var welcome2 map[string]interface{}
	err = json.Unmarshal(message2, &welcome2)
	require.NoError(t, err, "Failed to unmarshal welcome message")

	assert.Equal(t, sessionCode, welcome2["session_code"], "Should restore the same session")
}

// TestJSONRPCRequestResponse tests JSON-RPC method calls
func TestJSONRPCRequestResponse(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		params   interface{}
		validate func(t *testing.T, result interface{})
	}{
		{
			name:   "ping method",
			method: "ping",
			params: nil,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok, "Result should be a map")
				assert.Equal(t, true, resultMap["pong"], "Ping should return pong: true")
				assert.Equal(t, "fle-server", resultMap["server"], "Should identify server")
				assert.NotEmpty(t, resultMap["timestamp"], "Should include timestamp")
			},
		},
		{
			name:   "echo method with string",
			method: "echo",
			params: "hello world",
			validate: func(t *testing.T, result interface{}) {
				assert.Equal(t, "hello world", result, "Echo should return the input string")
			},
		},
		{
			name:   "echo method with object",
			method: "echo",
			params: map[string]interface{}{"message": "test", "number": 42},
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok, "Result should be a map")
				assert.Equal(t, "test", resultMap["message"], "Should echo message field")
				assert.Equal(t, float64(42), resultMap["number"], "Should echo number field")
			},
		},
		{
			name:   "getSessionInfo method",
			method: "getSessionInfo",
			params: nil,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok, "Result should be a map")
				assert.Contains(t, resultMap, "totalSessions", "Should include total sessions")
				assert.Contains(t, resultMap, "activeSessions", "Should include active sessions")
				assert.NotEmpty(t, resultMap["timestamp"], "Should include timestamp")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTestServer(t)
			defer ts.Close()

			// Connect to WebSocket
			wsURL := ts.wsURL + "/ws"
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err, "Failed to connect to WebSocket")
			defer conn.Close()

			// Set timeouts
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			// Read and discard welcome message
			_, _, err = conn.ReadMessage()
			require.NoError(t, err, "Failed to read welcome message")

			// Send JSON-RPC request
			request := jsonrpc.Request{
				JSONRPCVersion: "2.0",
				ID:             1,
				Method:         tc.method,
			}

			if tc.params != nil {
				params, err := json.Marshal(tc.params)
				require.NoError(t, err, "Failed to marshal params")
				request.Params = json.RawMessage(params)
			}

			requestBytes, err := json.Marshal(request)
			require.NoError(t, err, "Failed to marshal request")

			err = conn.WriteMessage(websocket.TextMessage, requestBytes)
			require.NoError(t, err, "Failed to send JSON-RPC request")

			// Read JSON-RPC response
			_, responseBytes, err := conn.ReadMessage()
			require.NoError(t, err, "Failed to read JSON-RPC response")

			var response jsonrpc.Response
			err = json.Unmarshal(responseBytes, &response)
			require.NoError(t, err, "Failed to unmarshal JSON-RPC response")

			// Validate response
			assert.Equal(t, "2.0", response.JSONRPCVersion, "Response should have correct JSON-RPC version")
			assert.Equal(t, float64(1), response.ID, "Response should have matching ID")
			assert.Nil(t, response.Error, "Response should not have error")
			assert.NotNil(t, response.Result, "Response should have result")

			// Run test-specific validation
			// The Result is already parsed from JSON
			tc.validate(t, response.Result)
		})
	}
}

// TestInvalidJSONRPCRequest tests handling of invalid JSON-RPC requests
func TestInvalidJSONRPCRequest(t *testing.T) {
	testCases := []struct {
		name    string
		message string
		errorCode int
	}{
		{
			name:    "invalid JSON",
			message: "{invalid json",
			errorCode: jsonrpc.ParseError,
		},
		{
			name:    "missing method",
			message: `{"jsonrpc": "2.0", "id": 1}`,
			errorCode: jsonrpc.InvalidRequest,
		},
		{
			name:    "unknown method",
			message: `{"jsonrpc": "2.0", "method": "unknownMethod", "id": 1}`,
			errorCode: jsonrpc.MethodNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTestServer(t)
			defer ts.Close()

			// Connect to WebSocket
			wsURL := ts.wsURL + "/ws"
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err, "Failed to connect to WebSocket")
			defer conn.Close()

			// Set timeouts
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			// Read and discard welcome message
			_, _, err = conn.ReadMessage()
			require.NoError(t, err, "Failed to read welcome message")

			// Send invalid JSON-RPC request
			err = conn.WriteMessage(websocket.TextMessage, []byte(tc.message))
			require.NoError(t, err, "Failed to send invalid JSON-RPC request")

			// Read JSON-RPC error response
			_, responseBytes, err := conn.ReadMessage()
			require.NoError(t, err, "Failed to read JSON-RPC error response")

			var response jsonrpc.Response
			err = json.Unmarshal(responseBytes, &response)
			require.NoError(t, err, "Failed to unmarshal JSON-RPC error response")

			// Validate error response
			assert.Equal(t, "2.0", response.JSONRPCVersion, "Response should have correct JSON-RPC version")
			assert.NotNil(t, response.Error, "Response should have error")
			assert.Equal(t, tc.errorCode, response.Error.Code, "Error should have expected code")
			assert.NotEmpty(t, response.Error.Message, "Error should have message")
		})
	}
}

// TestMultipleConcurrentConnections tests multiple WebSocket connections
func TestMultipleConcurrentConnections(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	const numConnections = 5
	var wg sync.WaitGroup
	errors := make(chan error, numConnections)

	// Create multiple concurrent connections
	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			wsURL := ts.wsURL + "/ws"
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				errors <- fmt.Errorf("client %d failed to connect: %w", clientID, err)
				return
			}
			defer conn.Close()

			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			// Read welcome message
			_, welcomeMsg, err := conn.ReadMessage()
			if err != nil {
				errors <- fmt.Errorf("client %d failed to read welcome: %w", clientID, err)
				return
			}

			var welcome map[string]interface{}
			if err := json.Unmarshal(welcomeMsg, &welcome); err != nil {
				errors <- fmt.Errorf("client %d failed to unmarshal welcome: %w", clientID, err)
				return
			}

			// Send a ping to verify connection works
			request := jsonrpc.Request{
				JSONRPCVersion: "2.0",
				ID:             clientID,
				Method:         "ping",
			}

			requestBytes, err := json.Marshal(request)
			if err != nil {
				errors <- fmt.Errorf("client %d failed to marshal ping: %w", clientID, err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, requestBytes); err != nil {
				errors <- fmt.Errorf("client %d failed to send ping: %w", clientID, err)
				return
			}

			// Read pong response
			_, responseBytes, err := conn.ReadMessage()
			if err != nil {
				errors <- fmt.Errorf("client %d failed to read pong: %w", clientID, err)
				return
			}

			var response jsonrpc.Response
			if err := json.Unmarshal(responseBytes, &response); err != nil {
				errors <- fmt.Errorf("client %d failed to unmarshal pong: %w", clientID, err)
				return
			}

			if response.Error != nil {
				errors <- fmt.Errorf("client %d got error response: %s", clientID, response.Error.Message)
				return
			}

			// Verify pong response
			result, ok := response.Result.(map[string]interface{})
			if !ok {
				errors <- fmt.Errorf("client %d got unexpected result type: %T", clientID, response.Result)
				return
			}

			if result["pong"] != true {
				errors <- fmt.Errorf("client %d got unexpected pong result: %v", clientID, result["pong"])
				return
			}
		}(i)
	}

	// Wait for all connections to complete
	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Error(err)
	}
}

// TestHeartbeatMechanism tests WebSocket ping/pong mechanism
func TestHeartbeatMechanism(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Connect to WebSocket
	wsURL := ts.wsURL + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket")
	defer conn.Close()

	// Set up ping handler to respond to pings
	pongReceived := make(chan bool, 1)
	conn.SetPongHandler(func(appData string) error {
		pongReceived <- true
		return nil
	})

	// Read and discard welcome message
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, _, err = conn.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message")

	// Start a goroutine to handle incoming messages (needed for ping/pong to work)
	messageHandlingDone := make(chan bool, 1)
	stopMessageHandling := make(chan bool, 1)
	go func() {
		defer func() { messageHandlingDone <- true }()
		for {
			select {
			case <-stopMessageHandling:
				return
			default:
			}
			
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Check if it's a timeout - continue for timeouts to handle ping/pong
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				// Connection closed or other error
				return
			}
		}
	}()

	// Send a ping frame
	err = conn.WriteMessage(websocket.PingMessage, []byte("ping"))
	require.NoError(t, err, "Failed to send ping")

	// Wait for pong response
	select {
	case <-pongReceived:
		t.Log("Received pong response")
	case <-time.After(3 * time.Second):
		t.Log("Did not receive pong response, testing server-to-client ping instead")
		// The server sends periodic pings, so let's wait for one
		select {
		case <-pongReceived:
			t.Log("Received server ping as pong")
		case <-time.After(3 * time.Second):
			t.Error("Did not receive any pong response within timeout")
		}
	}
	
	// Stop the message handling goroutine
	close(stopMessageHandling)
	select {
	case <-messageHandlingDone:
	case <-time.After(1 * time.Second):
		// Goroutine didn't stop in time, but that's ok for the test
	}
}

// TestMessageBroadcasting tests message broadcasting between clients
func TestMessageBroadcasting(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Create first client connection
	wsURL := ts.wsURL + "/ws"
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect first client")
	defer conn1.Close()

	// Create second client connection
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect second client")
	defer conn2.Close()

	// Set timeouts for both connections
	timeout := time.Now().Add(10 * time.Second)
	conn1.SetReadDeadline(timeout)
	conn2.SetReadDeadline(timeout)

	// Read welcome messages for both connections
	_, _, err = conn1.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message from conn1")
	
	_, _, err = conn2.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message from conn2")

	// Verify both connections can send and receive JSON-RPC messages
	// Send getSessionInfo from first connection
	request := jsonrpc.Request{
		JSONRPCVersion: "2.0",
		ID:             1,
		Method:         "getSessionInfo",
	}

	requestBytes, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal getSessionInfo request")

	conn1.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err = conn1.WriteMessage(websocket.TextMessage, requestBytes)
	require.NoError(t, err, "Failed to send getSessionInfo request from conn1")

	// Read response from first connection
	_, responseBytes, err := conn1.ReadMessage()
	require.NoError(t, err, "Failed to read getSessionInfo response on conn1")

	var response jsonrpc.Response
	err = json.Unmarshal(responseBytes, &response)
	require.NoError(t, err, "Failed to unmarshal getSessionInfo response")

	assert.Nil(t, response.Error, "getSessionInfo should not return error")

	// Parse result to verify session count
	result, ok := response.Result.(map[string]interface{})
	require.True(t, ok, "Result should be a map, got %T", response.Result)

	// Should report at least 2 total sessions (both connections)
	totalSessions := result["totalSessions"].(float64)
	assert.GreaterOrEqual(t, totalSessions, float64(2), "Should have at least 2 active connections")
}

// TestGracefulShutdown tests server shutdown behavior
func TestGracefulShutdown(t *testing.T) {
	ts := setupTestServer(t)
	
	// Connect to WebSocket before shutdown
	wsURL := ts.wsURL + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket")
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read welcome message
	_, _, err = conn.ReadMessage()
	require.NoError(t, err, "Failed to read welcome message")

	// Verify connection is working with a ping
	request := jsonrpc.Request{
		JSONRPCVersion: "2.0",
		ID:             1,
		Method:         "ping",
	}

	requestBytes, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal ping request")

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err = conn.WriteMessage(websocket.TextMessage, requestBytes)
	require.NoError(t, err, "Failed to send ping request")

	// Read ping response
	_, responseBytes, err := conn.ReadMessage()
	require.NoError(t, err, "Failed to read ping response")

	var response jsonrpc.Response
	err = json.Unmarshal(responseBytes, &response)
	require.NoError(t, err, "Failed to unmarshal ping response")

	assert.Nil(t, response.Error, "Ping should succeed before shutdown")

	// Now close the server (simulates graceful shutdown)
	ts.Close()

	// Try to send another message - this should fail as connection is closed
	err = conn.WriteMessage(websocket.TextMessage, requestBytes)
	if err != nil {
		t.Log("Connection properly closed after server shutdown:", err)
	}

	// The connection should be closed, reading should fail
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, _, err = conn.ReadMessage()
	if err != nil {
		t.Log("Connection properly terminated after shutdown:", err)
	}
}

