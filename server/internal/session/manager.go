// Package session provides session management and code generation functionality.
package session

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Manager provides thread-safe session management with in-memory storage.
type Manager struct {
	// sessions stores active sessions with their codes as keys
	sessions map[string]*Session

	// generator handles session code generation and validation
	generator *Generator

	// mutex provides thread-safe access to the sessions map
	mutex sync.RWMutex

	// options contains session configuration
	options *SessionOptions

	// cleanupInterval is how often expired sessions are cleaned up
	cleanupInterval time.Duration

	// stopCleanup is used to signal the cleanup goroutine to stop
	stopCleanup chan struct{}

	// cleanupDone signals when the cleanup goroutine has stopped
	cleanupDone chan struct{}
}

// NewManager creates a new session manager with the given options.
// If options is nil, default options will be used.
func NewManager(options *SessionOptions) *Manager {
	if options == nil {
		options = DefaultSessionOptions()
	}

	manager := &Manager{
		sessions:        make(map[string]*Session),
		generator:       NewGenerator(),
		options:         options,
		cleanupInterval: 10 * time.Minute, // Clean up every 10 minutes
		stopCleanup:     make(chan struct{}),
		cleanupDone:     make(chan struct{}),
	}

	// Start background cleanup goroutine
	go manager.cleanupExpiredSessions()

	return manager
}

// CreateSession creates a new session with a unique code.
// It will retry code generation up to MaxRetries times if collisions occur.
// Returns the created session or an error if unique code generation fails.
func (m *Manager) CreateSession(ctx context.Context, options *SessionOptions) (*Session, error) {
	if options == nil {
		options = m.options
	}

	// Generate unique session code with collision detection
	var code string
	var collision bool

	for attempt := 0; attempt <= options.MaxRetries; attempt++ {
		code = m.generator.GenerateCode()

		// Validate the generated code format
		if !m.generator.IsValidFormat(code) {
			continue // Try again with a new code
		}

		// Normalize the code for consistent storage
		normalizedCode := m.generator.NormalizeCode(code)

		// Check for collision
		m.mutex.RLock()
		_, collision = m.sessions[normalizedCode]
		m.mutex.RUnlock()

		if !collision {
			// No collision, we can use this code
			code = normalizedCode
			break
		}

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("session creation cancelled: %w", ctx.Err())
		default:
			// Continue with next attempt
		}
	}

	// If we still have a collision after all retries, return an error
	if collision {
		return nil, ErrCodeGenerationFailed
	}

	// Create the session
	now := time.Now()
	session := &Session{
		Code:         code,
		CreatedAt:    now,
		LastAccessed: now,
		Data:         make(map[string]interface{}),
	}

	// Copy initial data if provided
	if options.InitialData != nil {
		for k, v := range options.InitialData {
			session.Data[k] = v
		}
	}

	// Store the session
	m.mutex.Lock()
	m.sessions[code] = session
	m.mutex.Unlock()

	return session, nil
}

// GetSession retrieves a session by its code.
// Returns ErrSessionNotFound if the session doesn't exist.
// Returns ErrSessionExpired if the session has expired.
// Updates the LastAccessed timestamp if the session is found and valid.
func (m *Manager) GetSession(code string) (*Session, error) {
	if code == "" {
		return nil, ErrInvalidSessionCode
	}

	// Validate code format
	if !m.generator.IsValidFormat(code) {
		return nil, ErrInvalidSessionCode
	}

	// Normalize the code
	normalizedCode := m.generator.NormalizeCode(code)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[normalizedCode]
	if !exists {
		return nil, ErrSessionNotFound
	}

	// Check if session has expired
	if m.isExpired(session) {
		// Remove expired session
		delete(m.sessions, normalizedCode)
		return nil, ErrSessionExpired
	}

	// Update last accessed time
	session.LastAccessed = time.Now()

	return session, nil
}

// DeleteSession removes a session by its code.
// Returns true if the session was found and deleted, false otherwise.
func (m *Manager) DeleteSession(code string) bool {
	if code == "" {
		return false
	}

	// Validate code format for consistency with GetSession
	if !m.generator.IsValidFormat(code) {
		return false
	}

	normalizedCode := m.generator.NormalizeCode(code)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.sessions[normalizedCode]
	if exists {
		delete(m.sessions, normalizedCode)
	}

	return exists
}

// UpdateSessionData updates the data for a session.
// Returns ErrSessionNotFound if the session doesn't exist.
// Returns ErrSessionExpired if the session has expired.
func (m *Manager) UpdateSessionData(code string, data map[string]interface{}) error {
	if code == "" {
		return ErrInvalidSessionCode
	}

	normalizedCode := m.generator.NormalizeCode(code)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[normalizedCode]
	if !exists {
		return ErrSessionNotFound
	}

	// Check if session has expired
	if m.isExpired(session) {
		delete(m.sessions, normalizedCode)
		return ErrSessionExpired
	}

	// Update session data
	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}

	for k, v := range data {
		session.Data[k] = v
	}

	// Update last accessed time
	session.LastAccessed = time.Now()

	return nil
}

// GetSessionCount returns the current number of active sessions.
func (m *Manager) GetSessionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.sessions)
}

// ListSessions returns a slice of all active session codes.
// This is useful for debugging and monitoring purposes.
func (m *Manager) ListSessions() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	codes := make([]string, 0, len(m.sessions))
	for code := range m.sessions {
		codes = append(codes, code)
	}

	return codes
}

// Cleanup removes all expired sessions.
// Returns the number of sessions that were removed.
func (m *Manager) Cleanup() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	removed := 0
	for code, session := range m.sessions {
		if m.isExpired(session) {
			delete(m.sessions, code)
			removed++
		}
	}

	return removed
}

// Close stops the background cleanup goroutine and cleans up resources.
// This should be called when the session manager is no longer needed.
func (m *Manager) Close() {
	close(m.stopCleanup)
	<-m.cleanupDone
}

// isExpired checks if a session has expired based on the session timeout.
// This method assumes the caller holds the appropriate lock.
func (m *Manager) isExpired(session *Session) bool {
	return time.Since(session.LastAccessed) > m.options.SessionTimeout
}

// cleanupExpiredSessions runs in a background goroutine to periodically
// remove expired sessions from memory.
func (m *Manager) cleanupExpiredSessions() {
	defer close(m.cleanupDone)

	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.Cleanup()
		case <-m.stopCleanup:
			return
		}
	}
}
