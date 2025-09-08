// Package server provides HTTP server functionality for the FLE application.
// It implements a robust HTTP server with configurable endpoints, middleware,
// and graceful shutdown capabilities.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/fle/server/internal/config"
	"github.com/fle/server/internal/jsonrpc"
	"github.com/fle/server/internal/session"
	"github.com/fle/server/internal/websocket"
)

// Server represents the HTTP server instance with its configuration and state.
// It encapsulates all server-related functionality including middleware,
// routing, and lifecycle management.
type Server struct {
	// config holds the server configuration
	config *config.Config

	// httpServer is the underlying HTTP server instance
	httpServer *http.Server

	// router is the HTTP request multiplexer
	router *http.ServeMux

	// logger provides structured logging
	logger *slog.Logger

	// hub manages WebSocket connections
	hub *websocket.Hub

	// sessionManager handles session lifecycle
	sessionManager *session.Manager

	// jsonrpcRouter handles JSON-RPC method routing
	jsonrpcRouter *jsonrpc.Router
}

// NewServer creates and configures a new Server instance.
// It sets up the HTTP server with the provided configuration,
// initializes middleware, and configures routes.
//
// Parameters:
//   - cfg: Configuration for the server
//   - logger: Structured logger for server operations
//
// Returns:
//   - *Server: Configured server instance
//   - error: Error if server creation fails
func NewServer(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Create session manager
	sessionManager := session.NewManager(session.DefaultSessionOptions())

	// Create WebSocket hub
	hub := websocket.NewHub(logger)

	// Create JSON-RPC router
	jsonrpcRouter := jsonrpc.NewRouter()

	// Create the server instance
	server := &Server{
		config:         cfg,
		router:         http.NewServeMux(),
		logger:         logger,
		hub:            hub,
		sessionManager: sessionManager,
		jsonrpcRouter:  jsonrpcRouter,
	}

	// Set up routes
	server.setupRoutes()

	// Set up JSON-RPC methods
	server.setupJSONRPCMethods()

	// Start WebSocket hub
	go server.hub.Run()

	// Create HTTP server with configured parameters
	server.httpServer = &http.Server{
		Addr:         cfg.Address(),
		Handler:      server.setupMiddleware(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("HTTP server created",
		"address", cfg.Address(),
		"environment", cfg.Environment,
		"cors_origin", cfg.CORSOrigin,
	)

	return server, nil
}

// setupRoutes configures all HTTP routes for the server.
// This includes health endpoints and any other application routes.
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("GET /health", s.handleHealth)

	// WebSocket endpoint
	s.router.HandleFunc("GET /ws", s.handleWebSocket)

	s.logger.Debug("Routes configured",
		"routes", []string{"/health", "/ws"},
	)
}

// setupJSONRPCMethods registers all JSON-RPC methods with the router.
func (s *Server) setupJSONRPCMethods() {
	// Register ping method for testing connectivity
	s.jsonrpcRouter.RegisterSimpleMethod("ping", s.handlePing, "Simple ping method for testing JSON-RPC connectivity")
	
	// Register echo method for testing message passing
	s.jsonrpcRouter.RegisterSimpleMethod("echo", s.handleEcho, "Echo method that returns the input parameters")
	
	// Register get session info method
	s.jsonrpcRouter.RegisterSimpleMethod("getSessionInfo", s.handleGetSessionInfo, "Get information about the current WebSocket session")
	
	s.logger.Debug("JSON-RPC methods registered", 
		"methodCount", s.jsonrpcRouter.MethodCount(),
		"methods", s.jsonrpcRouter.GetMethods())
}

// setupMiddleware configures and chains all HTTP middleware.
// This includes CORS, logging, and any other cross-cutting concerns.
func (s *Server) setupMiddleware() http.Handler {
	var handler http.Handler = s.router

	// Apply CORS middleware for development
	if s.config.IsDevelopment() {
		handler = s.corsMiddleware(handler)
	}

	// Apply logging middleware
	handler = s.loggingMiddleware(handler)

	return handler
}

// Start begins listening for HTTP requests on the configured address.
// This method blocks until the server is stopped or encounters an error.
//
// Returns:
//   - error: Error if server fails to start or encounters issues while running
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		"address", s.config.Address(),
		"environment", s.config.Environment,
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the HTTP server.
// It waits for existing connections to close within the provided context timeout.
//
// Parameters:
//   - ctx: Context with timeout for graceful shutdown
//
// Returns:
//   - error: Error if shutdown fails
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")

	// Close session manager
	if s.sessionManager != nil {
		s.sessionManager.Close()
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("HTTP server stopped successfully")
	return nil
}

// Address returns the complete server address.
func (s *Server) Address() string {
	return s.config.Address()
}

// IsRunning returns true if the server is currently running.
// This is determined by checking if the HTTP server is not nil and not in a closed state.
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}

// Handler returns the HTTP handler for the server.
// This is useful for testing with httptest.Server.
func (s *Server) Handler() http.Handler {
	return s.setupMiddleware()
}
