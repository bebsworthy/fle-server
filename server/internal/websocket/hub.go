package websocket

import (
	"log/slog"
	"sync"

	"github.com/fle/server/internal/jsonrpc"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
// It implements the hub pattern for WebSocket connection management as described
// in the design specifications.
type Hub struct {
	// clients holds all currently connected clients
	clients map[*Client]bool

	// sessions maps session codes to their corresponding clients for targeted messaging
	sessions map[string]*Client

	// broadcast channel for broadcasting messages to all connected clients
	broadcast chan []byte

	// register channel for registering new clients
	register chan *Client

	// unregister channel for unregistering clients
	unregister chan *Client

	// mu provides thread-safe access to the clients and sessions maps
	mu sync.RWMutex

	// logger for structured logging
	logger *slog.Logger
}

// Client represents a single WebSocket connection with its associated session.
type Client struct {
	// hub is a reference to the hub managing this client
	hub *Hub

	// conn is the websocket connection
	conn *websocket.Conn

	// send is a buffered channel of outbound messages
	send chan []byte

	// sessionCode is the unique session identifier for this client
	sessionCode string

	// logger for structured logging specific to this client
	logger *slog.Logger

	// jsonrpcRouter handles JSON-RPC method routing for this client
	jsonrpcRouter *jsonrpc.Router
}

// NewHub creates a new Hub instance ready to manage WebSocket connections.
// It initializes all channels and maps required for the hub pattern.
func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		sessions:   make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// NewClient creates a new Client instance with the provided WebSocket connection
// and session code. The client is not automatically registered with the hub.
func NewClient(hub *Hub, conn *websocket.Conn, sessionCode string, logger *slog.Logger, jsonrpcRouter *jsonrpc.Router) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 256), // Buffered channel to prevent blocking
		sessionCode:   sessionCode,
		logger:        logger,
		jsonrpcRouter: jsonrpcRouter,
	}
}

// Run starts the hub's main event loop to handle client registration,
// unregistration, and message broadcasting. This method should be called
// in a separate goroutine as it runs indefinitely.
func (h *Hub) Run() {
	h.logger.Info("WebSocket hub started")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// RegisterClient adds a new client to the hub. This method should be called
// when a new WebSocket connection is established. It registers the client
// both in the general clients map and in the sessions map for targeted messaging.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes a client from the hub. This method should be called
// when a WebSocket connection is closed. It handles cleanup of both the clients
// and sessions maps.
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// SendToSession sends a message to a specific client identified by session code.
// If the session is not found, the message is silently dropped. This method
// is thread-safe and non-blocking.
func (h *Hub) SendToSession(sessionCode string, message []byte) {
	h.mu.RLock()
	client, exists := h.sessions[sessionCode]
	h.mu.RUnlock()

	if !exists {
		h.logger.Warn("attempted to send message to non-existent session",
			"sessionCode", sessionCode)
		return
	}

	select {
	case client.send <- message:
		h.logger.Debug("message sent to session",
			"sessionCode", sessionCode,
			"messageLength", len(message))
	default:
		// Client's send channel is full, close and unregister the client
		h.logger.Warn("client send channel full, unregistering",
			"sessionCode", sessionCode)
		close(client.send)
		h.UnregisterClient(client)
	}
}

// BroadcastMessage sends a message to all connected clients. This method
// is thread-safe and non-blocking.
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// GetClientCount returns the current number of connected clients.
// This method is thread-safe.
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetSessionCodes returns a slice of all active session codes.
// This method is thread-safe and returns a copy to prevent concurrent access issues.
func (h *Hub) GetSessionCodes() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	codes := make([]string, 0, len(h.sessions))
	for code := range h.sessions {
		codes = append(codes, code)
	}
	return codes
}

// HasSession returns true if a client with the given session code is connected.
// This method is thread-safe.
func (h *Hub) HasSession(sessionCode string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.sessions[sessionCode]
	return exists
}

// registerClient is the internal implementation for registering a client.
// It updates both the clients and sessions maps under write lock for thread safety.
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	h.sessions[client.sessionCode] = client
	h.mu.Unlock()

	h.logger.Info("client registered",
		"sessionCode", client.sessionCode,
		"clientCount", len(h.clients))
}

// unregisterClient is the internal implementation for unregistering a client.
// It removes the client from both maps and closes the send channel if it's not already closed.
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		delete(h.sessions, client.sessionCode)
		
		// Close the send channel if it's not already closed
		select {
		case <-client.send:
			// Channel is already closed
		default:
			close(client.send)
		}
	}
	clientCount := len(h.clients)
	h.mu.Unlock()

	h.logger.Info("client unregistered",
		"sessionCode", client.sessionCode,
		"clientCount", clientCount)
}

// broadcastMessage is the internal implementation for broadcasting messages.
// It sends the message to all connected clients. If a client's send channel is full,
// the client is automatically unregistered to prevent blocking other clients.
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	h.logger.Debug("broadcasting message",
		"clientCount", len(clients),
		"messageLength", len(message))

	// Send to all clients without holding the lock
	for _, client := range clients {
		select {
		case client.send <- message:
			// Message sent successfully
		default:
			// Client's send channel is full, close and unregister the client
			h.logger.Warn("client send channel full during broadcast, unregistering",
				"sessionCode", client.sessionCode)
			close(client.send)
			h.UnregisterClient(client)
		}
	}
}