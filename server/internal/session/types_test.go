package session

import (
	"errors"
	"testing"
)

func TestSessionError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SessionError
		expected string
	}{
		{
			name: "error without wrapped error",
			err: &SessionError{
				Code:    "TEST_ERROR",
				Message: "test error message",
				Err:     nil,
			},
			expected: "test error message",
		},
		{
			name: "error with wrapped error",
			err: &SessionError{
				Code:    "TEST_ERROR",
				Message: "test error message",
				Err:     errors.New("wrapped error"),
			},
			expected: "test error message: wrapped error",
		},
		{
			name: "predefined session not found error",
			err:  ErrSessionNotFound,
			expected: "session not found",
		},
		{
			name: "predefined session expired error",
			err:  ErrSessionExpired,
			expected: "session has expired",
		},
		{
			name: "predefined invalid session code error",
			err:  ErrInvalidSessionCode,
			expected: "invalid session code format",
		},
		{
			name: "predefined code generation failed error",
			err:  ErrCodeGenerationFailed,
			expected: "failed to generate unique session code after maximum retries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("SessionError.Error() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestSessionError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		err      *SessionError
		expected error
	}{
		{
			name: "error without wrapped error",
			err: &SessionError{
				Code:    "TEST_ERROR",
				Message: "test error message",
				Err:     nil,
			},
			expected: nil,
		},
		{
			name: "error with wrapped error",
			err: &SessionError{
				Code:    "TEST_ERROR",
				Message: "test error message",
				Err:     errors.New("wrapped error"),
			},
			expected: errors.New("wrapped error"),
		},
		{
			name: "predefined errors have no wrapped error",
			err:  ErrSessionNotFound,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Unwrap()
			if tt.expected == nil && result != nil {
				t.Errorf("SessionError.Unwrap() = %v, expected nil", result)
			} else if tt.expected != nil && result == nil {
				t.Errorf("SessionError.Unwrap() = nil, expected %v", tt.expected)
			} else if tt.expected != nil && result != nil && result.Error() != tt.expected.Error() {
				t.Errorf("SessionError.Unwrap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultSessionOptions(t *testing.T) {
	options := DefaultSessionOptions()
	
	if options == nil {
		t.Fatal("DefaultSessionOptions() returned nil")
	}
	
	if options.MaxRetries != 10 {
		t.Errorf("DefaultSessionOptions().MaxRetries = %d, expected 10", options.MaxRetries)
	}
	
	if options.SessionTimeout.Hours() != 24 {
		t.Errorf("DefaultSessionOptions().SessionTimeout = %v, expected 24h", options.SessionTimeout)
	}
	
	if options.InitialData == nil {
		t.Error("DefaultSessionOptions().InitialData should not be nil")
	}
	
	if len(options.InitialData) != 0 {
		t.Errorf("DefaultSessionOptions().InitialData should be empty, got %d items", len(options.InitialData))
	}
}