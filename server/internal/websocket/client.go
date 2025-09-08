package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/fle/server/internal/jsonrpc"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin during development
		// In production, this should be more restrictive
		return true
	},
}

// ServeWS handles WebSocket requests from the peer and creates a new client
// connection. It upgrades the HTTP connection to WebSocket and registers
// the client with the hub.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, sessionCode string, logger *slog.Logger, router *jsonrpc.Router) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WebSocket upgrade failed", 
			"error", err,
			"sessionCode", sessionCode)
		return
	}

	client := NewClient(hub, conn, sessionCode, logger, router)
	client.hub.RegisterClient(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the WebSocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("panic in readPump",
				"sessionCode", c.sessionCode,
				"panic", r)
		}
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.logger.Debug("pong received", "sessionCode", c.sessionCode)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.conn.SetPingHandler(func(appData string) error {
		c.logger.Debug("ping received", "sessionCode", c.sessionCode)
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.conn.WriteMessage(websocket.PongMessage, []byte(appData)); err != nil {
			c.logger.Warn("failed to send pong", "sessionCode", c.sessionCode, "error", err)
			return err
		}
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Warn("WebSocket connection error",
					"sessionCode", c.sessionCode,
					"error", err)
			} else {
				c.logger.Debug("WebSocket connection closed",
					"sessionCode", c.sessionCode,
					"error", err)
			}
			break
		}

		c.logger.Debug("message received",
			"sessionCode", c.sessionCode,
			"messageLength", len(message))

		// Process the message as JSON-RPC
		c.processJSONRPCMessage(message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("panic in writePump",
				"sessionCode", c.sessionCode,
				"panic", r)
		}
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.logger.Debug("send channel closed, sending close message",
					"sessionCode", c.sessionCode)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error("failed to get next writer",
					"sessionCode", c.sessionCode,
					"error", err)
				return
			}
			
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.logger.Error("failed to close writer",
					"sessionCode", c.sessionCode,
					"error", err)
				return
			}

			c.logger.Debug("message sent",
				"sessionCode", c.sessionCode,
				"messageLength", len(message),
				"additionalMessages", n)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Debug("ping failed, connection likely closed",
					"sessionCode", c.sessionCode,
					"error", err)
				return
			}
			c.logger.Debug("ping sent", "sessionCode", c.sessionCode)
		}
	}
}

// Close gracefully closes the client connection by sending a close message
// and cleaning up resources. This method is safe to call multiple times.
func (c *Client) Close() error {
	c.logger.Debug("closing client connection", "sessionCode", c.sessionCode)
	
	// Send close message to the client
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		c.logger.Warn("failed to send close message",
			"sessionCode", c.sessionCode,
			"error", err)
	}

	// Close the underlying connection
	return c.conn.Close()
}

// Send sends a message to this specific client. This method is thread-safe
// and non-blocking. If the client's send channel is full, the message is dropped.
func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
		c.logger.Debug("message queued for client",
			"sessionCode", c.sessionCode,
			"messageLength", len(message))
	default:
		c.logger.Warn("client send channel full, message dropped",
			"sessionCode", c.sessionCode,
			"messageLength", len(message))
	}
}

// SessionCode returns the session code associated with this client.
func (c *Client) SessionCode() string {
	return c.sessionCode
}

// IsConnected returns true if the WebSocket connection is still active.
// This is a best-effort check and may not be 100% accurate due to the
// asynchronous nature of network connections.
func (c *Client) IsConnected() bool {
	// Try to set a read deadline to test if the connection is still active
	// This is a lightweight way to check connection status
	err := c.conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	return err == nil
}

// processJSONRPCMessage processes incoming WebSocket messages as JSON-RPC requests.
// It parses the message, routes it through the JSON-RPC router, and sends back the response.
func (c *Client) processJSONRPCMessage(message []byte) {
	c.logger.Debug("processing JSON-RPC message",
		"sessionCode", c.sessionCode,
		"message", string(message))

	// Create a context for the request
	ctx := context.Background()
	
	// Check if the router is available
	if c.jsonrpcRouter == nil {
		c.logger.Error("JSON-RPC router not available",
			"sessionCode", c.sessionCode)
		c.sendJSONRPCError(nil, jsonrpc.ErrInternal, "JSON-RPC router not available")
		return
	}

	// Try to route the JSON message through the JSON-RPC router
	responseBytes, err := c.jsonrpcRouter.RouteJSON(ctx, message)
	if err != nil {
		c.logger.Error("failed to route JSON-RPC message",
			"sessionCode", c.sessionCode,
			"error", err,
			"message", string(message))
		c.sendJSONRPCError(nil, jsonrpc.ErrInternal, err.Error())
		return
	}

	// If responseBytes is nil, it was a notification (no response needed)
	if responseBytes == nil {
		c.logger.Debug("JSON-RPC notification processed successfully",
			"sessionCode", c.sessionCode)
		return
	}

	// Send the JSON-RPC response back to the client
	c.logger.Debug("sending JSON-RPC response",
		"sessionCode", c.sessionCode,
		"response", string(responseBytes))

	select {
	case c.send <- responseBytes:
		c.logger.Debug("JSON-RPC response queued for sending",
			"sessionCode", c.sessionCode,
			"responseLength", len(responseBytes))
	default:
		c.logger.Warn("send channel full, dropping JSON-RPC response",
			"sessionCode", c.sessionCode,
			"responseLength", len(responseBytes))
	}
}

// sendJSONRPCError sends a JSON-RPC error response back to the client.
func (c *Client) sendJSONRPCError(id interface{}, rpcError *jsonrpc.Error, details string) {
	// Create error with additional details if provided
	var err *jsonrpc.Error
	if details != "" {
		err = jsonrpc.NewErrorWithData(rpcError.Code, rpcError.Message, details)
	} else {
		err = rpcError
	}
	
	// Create error response
	response := jsonrpc.NewErrorResponse(err, id)
	
	// Marshal to JSON
	responseBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		c.logger.Error("failed to marshal JSON-RPC error response",
			"sessionCode", c.sessionCode,
			"error", marshalErr)
		return
	}

	// Send error response
	select {
	case c.send <- responseBytes:
		c.logger.Debug("JSON-RPC error response sent",
			"sessionCode", c.sessionCode,
			"errorCode", err.Code,
			"errorMessage", err.Message)
	default:
		c.logger.Warn("send channel full, dropping JSON-RPC error response",
			"sessionCode", c.sessionCode,
			"errorCode", err.Code)
	}
}