// Package session provides session management and code generation functionality.
package session

import (
	"time"
)

// Session represents an active session with its metadata.
type Session struct {
	// Code is the unique human-friendly session identifier (e.g., "happy-panda-42")
	Code string `json:"code"`

	// CreatedAt is the timestamp when the session was created
	CreatedAt time.Time `json:"created_at"`

	// LastAccessed is the timestamp of the last access to this session
	LastAccessed time.Time `json:"last_accessed"`

	// Data is a generic map for storing session-specific data
	Data map[string]interface{} `json:"data,omitempty"`
}

// SessionError represents errors related to session operations.
type SessionError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface for SessionError.
func (e *SessionError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error.
func (e *SessionError) Unwrap() error {
	return e.Err
}

// Common session error types
var (
	// ErrSessionNotFound is returned when a session with the given code is not found
	ErrSessionNotFound = &SessionError{
		Code:    "SESSION_NOT_FOUND",
		Message: "session not found",
	}

	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = &SessionError{
		Code:    "SESSION_EXPIRED",
		Message: "session has expired",
	}

	// ErrInvalidSessionCode is returned when the session code format is invalid
	ErrInvalidSessionCode = &SessionError{
		Code:    "INVALID_SESSION_CODE",
		Message: "invalid session code format",
	}

	// ErrCodeGenerationFailed is returned when session code generation fails after retries
	ErrCodeGenerationFailed = &SessionError{
		Code:    "CODE_GENERATION_FAILED",
		Message: "failed to generate unique session code after maximum retries",
	}
)

// SessionOptions contains configuration options for session creation.
type SessionOptions struct {
	// MaxRetries is the maximum number of retries for generating a unique session code
	MaxRetries int

	// SessionTimeout is the duration after which a session expires
	SessionTimeout time.Duration

	// InitialData is the initial data to store with the session
	InitialData map[string]interface{}
}

// DefaultSessionOptions returns the default session configuration.
func DefaultSessionOptions() *SessionOptions {
	return &SessionOptions{
		MaxRetries:     10,
		SessionTimeout: 24 * time.Hour, // 24 hours default timeout
		InitialData:    make(map[string]interface{}),
	}
}
