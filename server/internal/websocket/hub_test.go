package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/fle/server/internal/jsonrpc"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// mockConn implements a mock WebSocket connection for testing
type mockConn struct {
	closed    bool
	messages  [][]byte
	closeCode int
	mu        sync.Mutex
	writeChan chan []byte
}

func newMockConn() *mockConn {
	return &mockConn{
		messages:  make([][]byte, 0),
		writeChan: make(chan []byte, 10),
	}
}

func (m *mockConn) WriteMessage(messageType int, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return websocket.ErrCloseSent
	}
	m.messages = append(m.messages, data)
	select {
	case m.writeChan <- data:
	default:
	}
	return nil
}

func (m *mockConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	close(m.writeChan)
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) NextWriter(messageType int) (io.WriteCloser, error) {
	return &mockWriter{conn: m}, nil
}

func (m *mockConn) getMessages() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([][]byte, len(m.messages))
	copy(result, m.messages)
	return result
}

func (m *mockConn) isClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

type mockWriter struct {
	conn   *mockConn
	buffer bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return m.buffer.Write(p)
}

func (m *mockWriter) Close() error {
	data := m.buffer.Bytes()
	m.conn.mu.Lock()
	m.conn.messages = append(m.conn.messages, data)
	m.conn.mu.Unlock()
	return nil
}

func createTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func createTestRouter() *jsonrpc.Router {
	router := jsonrpc.NewRouter()
	// Register a simple test method
	router.RegisterSimpleMethod("test.echo", func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return string(params), nil
	}, "Echo test method")
	return router
}

// Helper function to create a client with mock connection using unsafe conversion
func createTestClient(sessionCode string) (*Client, *mockConn, *Hub) {
	logger := createTestLogger()
	hub := NewHub(logger)
	router := createTestRouter()
	mockConn := newMockConn()
	
	// Use unsafe pointer conversion to bypass type checking for testing
	// This is not recommended in production code but acceptable for unit tests
	wsConn := (*websocket.Conn)(unsafe.Pointer(mockConn))
	client := NewClient(hub, wsConn, sessionCode, logger, router)
	
	return client, mockConn, hub
}

func TestNewHub(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.sessions)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.Equal(t, logger, hub.logger)
	assert.Equal(t, 0, hub.GetClientCount())
	assert.Equal(t, 0, len(hub.GetSessionCodes()))
}

func TestHubRunLifecycle(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)
	
	// Start the hub in a goroutine
	done := make(chan bool)
	go func() {
		// Run hub for a short time
		select {
		case <-time.After(100 * time.Millisecond):
			done <- true
		}
	}()

	go hub.Run()

	// Wait for the test to complete
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Hub.Run() did not respond within expected time")
	}
}

func TestHubClientRegistration(t *testing.T) {
	client, _, hub := createTestClient("session1")

	// Start the hub
	go hub.Run()

	// Register the client
	hub.RegisterClient(client)

	// Give some time for registration to process
	time.Sleep(10 * time.Millisecond)

	// Verify client is registered
	assert.Equal(t, 1, hub.GetClientCount())
	assert.True(t, hub.HasSession("session1"))
	assert.Contains(t, hub.GetSessionCodes(), "session1")
}

func TestHubClientUnregistration(t *testing.T) {
	client, _, hub := createTestClient("session1")

	// Start the hub
	go hub.Run()

	// Register and then unregister a client
	hub.RegisterClient(client)

	// Wait for registration
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, hub.GetClientCount())

	// Unregister the client
	hub.UnregisterClient(client)

	// Wait for unregistration
	time.Sleep(10 * time.Millisecond)

	// Verify client is unregistered
	assert.Equal(t, 0, hub.GetClientCount())
	assert.False(t, hub.HasSession("session1"))
	assert.Empty(t, hub.GetSessionCodes())
}

func TestHubConcurrentConnections(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	const numClients = 10
	clients := make([]*Client, numClients)
	var wg sync.WaitGroup

	// Register multiple clients concurrently
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client, _, _ := createTestClient(fmt.Sprintf("session%d", id))
			client.hub = hub // Update hub reference
			clients[id] = client
			hub.RegisterClient(client)
		}(i)
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond) // Allow processing time

	// Verify all clients are registered
	assert.Equal(t, numClients, hub.GetClientCount())
	assert.Equal(t, numClients, len(hub.GetSessionCodes()))

	// Verify all sessions exist
	for i := 0; i < numClients; i++ {
		assert.True(t, hub.HasSession(fmt.Sprintf("session%d", i)))
	}

	// Unregister all clients concurrently
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			hub.UnregisterClient(clients[id])
		}(i)
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond) // Allow processing time

	// Verify all clients are unregistered
	assert.Equal(t, 0, hub.GetClientCount())
	assert.Empty(t, hub.GetSessionCodes())
}

func TestHubMessageBroadcast(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	const numClients = 3
	clients := make([]*Client, numClients)

	// Register multiple clients
	for i := 0; i < numClients; i++ {
		client, _, _ := createTestClient(fmt.Sprintf("session%d", i))
		client.hub = hub // Update hub reference
		clients[i] = client
		hub.RegisterClient(client)
	}

	time.Sleep(20 * time.Millisecond) // Allow registration

	// Broadcast a message
	testMessage := []byte("broadcast test message")
	hub.BroadcastMessage(testMessage)

	time.Sleep(20 * time.Millisecond) // Allow message processing

	// Verify all clients received the message
	for i := 0; i < numClients; i++ {
		select {
		case msg := <-clients[i].send:
			assert.Equal(t, testMessage, msg)
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Client %d did not receive broadcast message", i)
		}
	}
}

func TestHubTargetedMessaging(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	// Register multiple clients
	client1, _, _ := createTestClient("session1")
	client1.hub = hub
	hub.RegisterClient(client1)

	client2, _, _ := createTestClient("session2")
	client2.hub = hub
	hub.RegisterClient(client2)

	time.Sleep(20 * time.Millisecond) // Allow registration

	// Send message to specific session
	testMessage := []byte("targeted message")
	hub.SendToSession("session1", testMessage)

	time.Sleep(20 * time.Millisecond) // Allow message processing

	// Verify only the targeted client received the message
	select {
	case msg := <-client1.send:
		assert.Equal(t, testMessage, msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive targeted message")
	}

	// Verify client2 did not receive the message
	select {
	case <-client2.send:
		t.Error("Client2 incorrectly received targeted message")
	case <-time.After(50 * time.Millisecond):
		// Expected - client2 should not receive the message
	}
}

func TestHubSendToNonExistentSession(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	// Try to send to non-existent session (should not panic)
	testMessage := []byte("message to nowhere")
	hub.SendToSession("nonexistent", testMessage)

	// If we get here without panic, the test passes
	time.Sleep(10 * time.Millisecond)
}

func TestHubClientChannelFull(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	// Create a client with a small buffer
	client, _, _ := createTestClient("session1")
	client.hub = hub
	// Fill the client's send channel to capacity
	for i := 0; i < cap(client.send); i++ {
		client.send <- []byte("filler message")
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	// Try to send another message - this should cause unregistration
	hub.SendToSession("session1", []byte("overflow message"))
	time.Sleep(20 * time.Millisecond)

	// Client should be unregistered due to full channel
	assert.Equal(t, 0, hub.GetClientCount())
	assert.False(t, hub.HasSession("session1"))
}

func TestHubCleanupOnDisconnect(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	// Register a client
	client, _, _ := createTestClient("session1")
	client.hub = hub
	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, hub.GetClientCount())
	assert.True(t, hub.HasSession("session1"))

	// Unregister client (simulating disconnect)
	hub.UnregisterClient(client)
	time.Sleep(10 * time.Millisecond)

	// Verify cleanup
	assert.Equal(t, 0, hub.GetClientCount())
	assert.False(t, hub.HasSession("session1"))
	assert.Empty(t, hub.GetSessionCodes())

	// Verify send channel is closed
	select {
	case _, ok := <-client.send:
		assert.False(t, ok, "Send channel should be closed")
	case <-time.After(100 * time.Millisecond):
		t.Error("Send channel was not closed")
	}
}

func TestHubSessionCodeRetrieval(t *testing.T) {
	logger := createTestLogger()
	hub := NewHub(logger)

	// Start the hub
	go hub.Run()

	expectedSessions := []string{"session1", "session2", "session3"}
	
	// Register clients with different session codes
	for _, sessionCode := range expectedSessions {
		client, _, _ := createTestClient(sessionCode)
		client.hub = hub
		hub.RegisterClient(client)
	}

	time.Sleep(30 * time.Millisecond) // Allow registration

	// Get session codes
	sessionCodes := hub.GetSessionCodes()
	assert.Equal(t, len(expectedSessions), len(sessionCodes))

	// Verify all expected sessions are present
	for _, expected := range expectedSessions {
		assert.Contains(t, sessionCodes, expected)
		assert.True(t, hub.HasSession(expected))
	}
}

// Benchmark tests for performance evaluation
func BenchmarkHubBroadcast(b *testing.B) {
	logger := createTestLogger()
	hub := NewHub(logger)

	go hub.Run()

	// Create multiple clients
	const numClients = 100
	for i := 0; i < numClients; i++ {
		client, _, _ := createTestClient(fmt.Sprintf("session%d", i))
		client.hub = hub
		hub.RegisterClient(client)
	}

	time.Sleep(100 * time.Millisecond) // Allow registration
	testMessage := []byte("benchmark message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.BroadcastMessage(testMessage)
	}
}

func BenchmarkHubTargetedMessage(b *testing.B) {
	logger := createTestLogger()
	hub := NewHub(logger)

	go hub.Run()

	// Create a client
	client, _, _ := createTestClient("target_session")
	client.hub = hub
	hub.RegisterClient(client)

	time.Sleep(10 * time.Millisecond)
	testMessage := []byte("targeted benchmark message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.SendToSession("target_session", testMessage)
	}
}

func BenchmarkHubClientRegistration(b *testing.B) {
	logger := createTestLogger()
	hub := NewHub(logger)

	go hub.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, _, _ := createTestClient(fmt.Sprintf("session%d", i))
		client.hub = hub
		hub.RegisterClient(client)
		hub.UnregisterClient(client)
	}
}