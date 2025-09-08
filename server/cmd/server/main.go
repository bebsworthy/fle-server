// Package main provides the entry point for the FLE server.
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fle/server/internal/config"
	"github.com/fle/server/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Set up structured logging
	logger := setupLogger(cfg)

	logger.Info("FLE Server starting",
		"address", cfg.Address(),
		"environment", cfg.Environment,
		"log_level", cfg.LogLevel,
		"cors_origin", cfg.CORSOrigin,
	)

	// Create and configure the server
	srv, err := server.NewServer(cfg, logger)
	if err != nil {
		logger.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		sig := <-sigChan
		logger.Info("Received shutdown signal", "signal", sig)

		// Cancel the context to trigger graceful shutdown
		cancel()
	}()

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Wait for either an error or shutdown signal
	select {
	case err := <-errChan:
		if err != nil {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("Shutdown signal received, stopping server...")

		// Create a timeout context for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Stop(shutdownCtx); err != nil {
			logger.Error("Failed to stop server gracefully", "error", err)
			os.Exit(1)
		}

		logger.Info("Server stopped successfully")
	}
}

// setupLogger creates and configures a structured logger based on the configuration.
func setupLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevelSlog(),
	}

	var handler slog.Handler
	if cfg.IsDevelopment() {
		// Use text handler for better readability in development
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// Use JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
