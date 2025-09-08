// Package server provides HTTP handlers and middleware for the FLE application.
package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fle/server/internal/websocket"
)

// HealthResponse represents the structure of the health check response.
// It provides information about the server's operational status.
type HealthResponse struct {
	// Status indicates the overall health status
	Status string `json:"status"`

	// Timestamp is the time when the health check was performed
	Timestamp time.Time `json:"timestamp"`

	// Version can be used to identify the server version
	Version string `json:"version,omitempty"`

	// Environment indicates the current deployment environment
	Environment string `json:"environment"`
}

// handleHealth handles GET requests to the /health endpoint.
// It returns a JSON response indicating the server's health status.
// This endpoint is used for health checks by load balancers and monitoring systems.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:      "healthy",
		Timestamp:   time.Now().UTC(),
		Version:     "1.0.0", // TODO: This should come from build information
		Environment: s.config.Environment,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("Failed to encode health response",
			"error", err,
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.logger.Debug("Health check completed",
		"remote_addr", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"),
	)
}

// WelcomeMessage represents the welcome message sent to newly connected WebSocket clients.
type WelcomeMessage struct {
	Type        string `json:"type"`
	SessionCode string `json:"session_code"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
}

// handleWebSocket handles WebSocket upgrade requests.
// It creates or restores a session, upgrades the connection, and sends a welcome message.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Check if this is a WebSocket upgrade request
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Expected WebSocket upgrade", http.StatusBadRequest)
		return
	}

	// Try to get session code from query parameters or create a new session
	sessionCode := r.URL.Query().Get("session")

	if sessionCode != "" {
		// Try to restore existing session
		if existingSession, err := s.sessionManager.GetSession(sessionCode); err == nil {
			sessionCode = existingSession.Code
			s.logger.Debug("Restored existing session",
				"sessionCode", sessionCode,
				"remote_addr", r.RemoteAddr)
		} else {
			s.logger.Debug("Session not found or expired, creating new session",
				"requested_session", sessionCode,
				"error", err,
				"remote_addr", r.RemoteAddr)
			// Create new session if existing one is not found or expired
			sessionCode = ""
		}
	}

	if sessionCode == "" {
		// Create a new session
		newSession, err := s.sessionManager.CreateSession(context.Background(), nil)
		if err != nil {
			s.logger.Error("Failed to create session",
				"error", err,
				"remote_addr", r.RemoteAddr)
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		sessionCode = newSession.Code
		s.logger.Debug("Created new session",
			"sessionCode", sessionCode,
			"remote_addr", r.RemoteAddr)
	}

	// Upgrade HTTP connection to WebSocket
	websocket.ServeWS(s.hub, w, r, sessionCode, s.logger, s.jsonrpcRouter)

	// Send welcome message after connection is established
	// Note: We need to wait a moment for the connection to be fully established
	go func() {
		time.Sleep(100 * time.Millisecond) // Brief delay to ensure connection is ready

		welcomeMsg := WelcomeMessage{
			Type:        "welcome",
			SessionCode: sessionCode,
			Message:     "WebSocket connection established successfully",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}

		msgBytes, err := json.Marshal(welcomeMsg)
		if err != nil {
			s.logger.Error("Failed to marshal welcome message",
				"sessionCode", sessionCode,
				"error", err)
			return
		}

		// Send welcome message to the specific session
		s.hub.SendToSession(sessionCode, msgBytes)
		s.logger.Debug("Welcome message sent",
			"sessionCode", sessionCode)
	}()

	s.logger.Info("WebSocket connection established",
		"sessionCode", sessionCode,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"))
}

// corsMiddleware adds CORS headers to responses for development environments.
// This allows the frontend development server to communicate with the backend.
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", s.config.CORSOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests and responses with structured logging.
// It captures request details and response status for monitoring and debugging.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process the request
		next.ServeHTTP(wrapper, r)

		// Log the request
		duration := time.Since(start)
		s.logger.Info("HTTP request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapper.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.Header.Get("User-Agent"),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
// It also implements http.Hijacker to support WebSocket upgrades.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker for WebSocket upgrade support.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("responseWriter does not support hijacking")
}

// JSON-RPC Method Handlers

// handlePing handles the "ping" JSON-RPC method.
func (s *Server) handlePing(ctx context.Context, params json.RawMessage) (interface{}, error) {
	s.logger.Debug("JSON-RPC ping method called")
	
	return map[string]interface{}{
		"pong":      true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"server":    "fle-server",
	}, nil
}

// handleEcho handles the "echo" JSON-RPC method.
func (s *Server) handleEcho(ctx context.Context, params json.RawMessage) (interface{}, error) {
	s.logger.Debug("JSON-RPC echo method called", "params", string(params))
	
	// If no params provided, return empty object
	if len(params) == 0 {
		return map[string]interface{}{}, nil
	}
	
	// Parse params as generic JSON
	var result interface{}
	if err := json.Unmarshal(params, &result); err != nil {
		return nil, fmt.Errorf("failed to parse echo params: %w", err)
	}
	
	return result, nil
}

// handleGetSessionInfo handles the "getSessionInfo" JSON-RPC method.
func (s *Server) handleGetSessionInfo(ctx context.Context, params json.RawMessage) (interface{}, error) {
	s.logger.Debug("JSON-RPC getSessionInfo method called")
	
	// For now, return basic info about connected sessions
	// In a real implementation, this would extract session info from context
	return map[string]interface{}{
		"totalSessions": s.hub.GetClientCount(),
		"activeSessions": s.hub.GetSessionCodes(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, nil
}
