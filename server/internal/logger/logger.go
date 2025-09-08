// Package logger provides structured logging functionality for the FLE server.
// It implements configurable logging with JSON format for production and
// human-readable format for development using Go's structured logging (slog).
package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/fle/server/internal/config"
)

// Logger wraps slog.Logger with additional functionality for structured logging.
// It provides methods for generating request IDs and creating child loggers with context.
type Logger struct {
	*slog.Logger
	config *config.Config
}

// Options configures logger behavior.
// It allows customization of output destination and format.
type Options struct {
	// Output is the destination for log messages. If nil, os.Stderr is used.
	Output io.Writer

	// AddSource includes source code position in log records.
	AddSource bool

	// ReplaceAttr allows customization of log attributes before output.
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

var (
	// defaultLogger is the global logger instance used by package-level functions.
	// This global is necessary to provide convenient package-level logging functions
	// while maintaining a single configured logger instance throughout the application.
	//nolint:gochecknoglobals
	defaultLogger *Logger

	// initOnce ensures the default logger is initialized only once.
	// This global is necessary for thread-safe singleton initialization.
	//nolint:gochecknoglobals
	initOnce sync.Once
)

// New creates a new Logger instance based on the provided configuration.
// The logger format (JSON or text) is determined by the environment setting.
// Log level is configured based on the config.LogLevel setting.
func New(cfg *config.Config, opts ...Options) (*Logger, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Use default options if none provided
	var options Options
	if len(opts) > 0 {
		options = opts[0]
	}

	// Set default output if not specified
	output := options.Output
	if output == nil {
		output = os.Stderr
	}

	// Create handler options with configured level
	handlerOpts := &slog.HandlerOptions{
		Level:       cfg.LogLevelSlog(),
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}

	var handler slog.Handler

	// Choose handler based on environment
	if cfg.IsProduction() {
		// JSON format for production (structured logging for log aggregation systems)
		handler = slog.NewJSONHandler(output, handlerOpts)
	} else {
		// Text format for development (human-readable console output)
		handler = slog.NewTextHandler(output, handlerOpts)
	}

	slogLogger := slog.New(handler)

	logger := &Logger{
		Logger: slogLogger,
		config: cfg,
	}

	return logger, nil
}

// Init initializes the global logger with the provided configuration.
// This should be called once at application startup.
// Subsequent calls are ignored (safe to call multiple times).
func Init(cfg *config.Config, opts ...Options) error {
	var err error
	initOnce.Do(func() {
		defaultLogger, err = New(cfg, opts...)
	})
	return err
}

// Default returns the global logger instance.
// If the logger hasn't been initialized with Init(), it panics.
// This is intentional to catch configuration errors early.
func Default() *Logger {
	if defaultLogger == nil {
		panic("logger not initialized: call logger.Init() first")
	}
	return defaultLogger
}

const (
	// RequestIDBytes is the number of random bytes used for request ID generation.
	RequestIDBytes = 8
)

// GenerateRequestID creates a unique request ID for tracing related operations.
// The request ID is a cryptographically secure random 16-byte hex string.
// This can be used to correlate log entries for a single request across components.
func GenerateRequestID() string {
	// Generate random data (16 hex characters)
	bytes := make([]byte, RequestIDBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fall back to timestamp-based ID if crypto/rand fails
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// WithRequestID returns a new logger that includes the request ID in all log entries.
// This creates a child logger that automatically adds the request ID as context.
func (l *Logger) WithRequestID(requestID string) *Logger {
	if requestID == "" {
		requestID = GenerateRequestID()
	}

	childLogger := l.With(slog.String("request_id", requestID))
	return &Logger{
		Logger: childLogger,
		config: l.config,
	}
}

// WithSessionCode returns a new logger that includes the session code in all log entries.
// This is useful for tracking operations related to a specific session.
func (l *Logger) WithSessionCode(sessionCode string) *Logger {
	childLogger := l.With(slog.String("session_code", sessionCode))
	return &Logger{
		Logger: childLogger,
		config: l.config,
	}
}

// WithComponent returns a new logger that includes the component name in all log entries.
// This helps identify which part of the system generated the log entry.
func (l *Logger) WithComponent(component string) *Logger {
	childLogger := l.With(slog.String("component", component))
	return &Logger{
		Logger: childLogger,
		config: l.config,
	}
}

// WithContext returns a new logger that includes arbitrary key-value pairs in all log entries.
// This is useful for adding custom context to log entries.
func (l *Logger) WithContext(attrs ...slog.Attr) *Logger {
	childLogger := l.With(slog.Group("context", attrsToAny(attrs)...))
	return &Logger{
		Logger: childLogger,
		config: l.config,
	}
}

// WithFields returns a new logger with additional structured fields.
// This is a convenience method for adding multiple key-value pairs.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return l.WithContext(attrs...)
}

// LogConnection logs WebSocket connection events with standardized format.
// This provides consistent logging for connection lifecycle events.
func (l *Logger) LogConnection(ctx context.Context, event, sessionCode, remoteAddr string) {
	l.InfoContext(ctx, "websocket connection event",
		slog.String("event", event),
		slog.String("session_code", sessionCode),
		slog.String("remote_addr", remoteAddr),
		slog.String("component", "websocket"),
	)
}

// LogRequest logs HTTP request events with standardized format.
// This provides consistent logging for HTTP request processing.
func (l *Logger) LogRequest(
	ctx context.Context, method, path, remoteAddr string, statusCode int, duration time.Duration,
) {
	l.InfoContext(ctx, "http request",
		slog.String("method", method),
		slog.String("path", path),
		slog.String("remote_addr", remoteAddr),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		slog.String("component", "http"),
	)
}

// LogError logs an error with additional context.
// This provides structured error logging with stack traces when appropriate.
func (l *Logger) LogError(ctx context.Context, err error, msg string, attrs ...slog.Attr) {
	allAttrs := append([]slog.Attr{
		slog.String("error", err.Error()),
	}, attrs...)

	l.LogAttrs(ctx, slog.LevelError, msg, allAttrs...)
}

// Package-level convenience functions that use the default logger

// Debug logs a debug message using the default logger.
func Debug(msg string, attrs ...slog.Attr) {
	Default().LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
}

// Info logs an info message using the default logger.
func Info(msg string, attrs ...slog.Attr) {
	Default().LogAttrs(context.TODO(), slog.LevelInfo, msg, attrs...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, attrs ...slog.Attr) {
	Default().LogAttrs(context.TODO(), slog.LevelWarn, msg, attrs...)
}

// Error logs an error message using the default logger.
func Error(msg string, attrs ...slog.Attr) {
	Default().LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
}

// ErrorWithErr logs an error with an error value using the default logger.
func ErrorWithErr(err error, msg string, attrs ...slog.Attr) {
	Default().LogError(context.TODO(), err, msg, attrs...)
}

// WithRequestID returns a logger with request ID using the default logger.
func WithRequestID(requestID string) *Logger {
	return Default().WithRequestID(requestID)
}

// WithSessionCode returns a logger with session code using the default logger.
func WithSessionCode(sessionCode string) *Logger {
	return Default().WithSessionCode(sessionCode)
}

// WithComponent returns a logger with component name using the default logger.
func WithComponent(component string) *Logger {
	return Default().WithComponent(component)
}

// attrsToAny converts []slog.Attr to []any for use with slog.Group.
func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, attr := range attrs {
		result[i] = attr
	}
	return result
}

// IsDebugEnabled returns true if debug level logging is enabled.
func (l *Logger) IsDebugEnabled() bool {
	return l.Enabled(context.TODO(), slog.LevelDebug)
}

// IsInfoEnabled returns true if info level logging is enabled.
func (l *Logger) IsInfoEnabled() bool {
	return l.Enabled(context.TODO(), slog.LevelInfo)
}

// IsWarnEnabled returns true if warn level logging is enabled.
func (l *Logger) IsWarnEnabled() bool {
	return l.Enabled(context.TODO(), slog.LevelWarn)
}

// IsErrorEnabled returns true if error level logging is enabled.
func (l *Logger) IsErrorEnabled() bool {
	return l.Enabled(context.TODO(), slog.LevelError)
}
