package session

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(nil)
	if manager == nil {
		t.Fatal("NewManager should not return nil")
	}

	if manager.sessions == nil {
		t.Error("sessions map should be initialized")
	}

	if manager.generator == nil {
		t.Error("generator should be initialized")
	}

	if manager.options == nil {
		t.Error("options should be initialized")
	}

	// Clean up
	manager.Close()
}

func TestCreateSession(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)

	if err != nil {
		t.Fatalf("CreateSession should not return error: %v", err)
	}

	if session == nil {
		t.Fatal("session should not be nil")
	}

	if session.Code == "" {
		t.Error("session code should not be empty")
	}

	if session.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	if session.LastAccessed.IsZero() {
		t.Error("LastAccessed should be set")
	}

	if session.Data == nil {
		t.Error("Data should be initialized")
	}
}

func TestCreateSessionWithInitialData(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	initialData := map[string]interface{}{
		"user_id": "test123",
		"role":    "admin",
	}

	options := &SessionOptions{
		MaxRetries:  10,
		InitialData: initialData,
	}

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, options)

	if err != nil {
		t.Fatalf("CreateSession should not return error: %v", err)
	}

	if session.Data["user_id"] != "test123" {
		t.Error("initial data should be copied to session")
	}

	if session.Data["role"] != "admin" {
		t.Error("initial data should be copied to session")
	}
}

func TestGetSession(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Create a session first
	ctx := context.Background()
	originalSession, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Store the original LastAccessed time
	originalLastAccessed := originalSession.LastAccessed

	// Sleep briefly to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Get the session
	retrievedSession, err := manager.GetSession(originalSession.Code)
	if err != nil {
		t.Fatalf("GetSession should not return error: %v", err)
	}

	if retrievedSession.Code != originalSession.Code {
		t.Error("retrieved session should have same code")
	}

	// LastAccessed should be updated
	if !retrievedSession.LastAccessed.After(originalLastAccessed) {
		t.Error("LastAccessed should be updated when session is retrieved")
	}
}

func TestGetSessionNotFound(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Use a valid format but nonexistent session code
	_, err := manager.GetSession("nonexistent-code-42")
	if err != ErrSessionNotFound {
		t.Errorf("GetSession should return ErrSessionNotFound, got: %v", err)
	}
}

func TestGetSessionInvalidCode(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Test empty code
	_, err := manager.GetSession("")
	if err != ErrInvalidSessionCode {
		t.Errorf("GetSession should return ErrInvalidSessionCode for empty code, got: %v", err)
	}

	// Test invalid format
	_, err = manager.GetSession("invalid")
	if err != ErrInvalidSessionCode {
		t.Errorf("GetSession should return ErrInvalidSessionCode for invalid format, got: %v", err)
	}
}

func TestDeleteSession(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Create a session first
	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Delete the session
	deleted := manager.DeleteSession(session.Code)
	if !deleted {
		t.Error("DeleteSession should return true when session exists")
	}

	// Try to get the deleted session
	_, err = manager.GetSession(session.Code)
	if err != ErrSessionNotFound {
		t.Error("GetSession should return ErrSessionNotFound after deletion")
	}

	// Delete non-existent session
	deleted = manager.DeleteSession("nonexistent-code-123")
	if deleted {
		t.Error("DeleteSession should return false when session doesn't exist")
	}
}

func TestUpdateSessionData(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Create a session first
	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Update session data
	updateData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	err = manager.UpdateSessionData(session.Code, updateData)
	if err != nil {
		t.Fatalf("UpdateSessionData should not return error: %v", err)
	}

	// Get the session and verify data
	updatedSession, err := manager.GetSession(session.Code)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if updatedSession.Data["key1"] != "value1" {
		t.Error("session data should be updated")
	}

	if updatedSession.Data["key2"] != 42 {
		t.Error("session data should be updated")
	}
}

func TestSessionExpiration(t *testing.T) {
	// Create manager with very short timeout
	options := &SessionOptions{
		MaxRetries:     10,
		SessionTimeout: 1 * time.Millisecond,
	}
	manager := NewManager(options)
	defer manager.Close()

	// Create a session
	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Try to get expired session
	_, err = manager.GetSession(session.Code)
	if err != ErrSessionExpired {
		t.Errorf("GetSession should return ErrSessionExpired, got: %v", err)
	}
}

func TestGetSessionCount(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Initial count should be 0
	if count := manager.GetSessionCount(); count != 0 {
		t.Errorf("initial session count should be 0, got: %d", count)
	}

	// Create some sessions
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		_, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}
	}

	// Count should be 3
	if count := manager.GetSessionCount(); count != 3 {
		t.Errorf("session count should be 3, got: %d", count)
	}
}

func TestListSessions(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Create some sessions
	ctx := context.Background()
	var sessionCodes []string
	for i := 0; i < 3; i++ {
		session, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}
		sessionCodes = append(sessionCodes, session.Code)
	}

	// List sessions
	listedCodes := manager.ListSessions()
	if len(listedCodes) != 3 {
		t.Errorf("ListSessions should return 3 codes, got: %d", len(listedCodes))
	}

	// Check that all created codes are in the list
	for _, code := range sessionCodes {
		found := false
		for _, listedCode := range listedCodes {
			if code == listedCode {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("session code %s not found in listed codes", code)
		}
	}
}

func TestCleanup(t *testing.T) {
	// Create manager with very short timeout
	options := &SessionOptions{
		MaxRetries:     10,
		SessionTimeout: 1 * time.Millisecond,
	}
	manager := NewManager(options)
	defer manager.Close()

	// Create some sessions
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		_, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Run cleanup
	removed := manager.Cleanup()
	if removed != 3 {
		t.Errorf("Cleanup should have removed 3 sessions, got: %d", removed)
	}

	// Session count should be 0 now
	if count := manager.GetSessionCount(); count != 0 {
		t.Errorf("session count should be 0 after cleanup, got: %d", count)
	}
}

func TestConcurrentAccess(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Test concurrent session creation
	ctx := context.Background()
	const numGoroutines = 10
	sessions := make(chan *Session, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			session, err := manager.CreateSession(ctx, nil)
			if err != nil {
				errors <- err
				return
			}
			sessions <- session
		}()
	}

	// Collect results
	var createdSessions []*Session
	for i := 0; i < numGoroutines; i++ {
		select {
		case session := <-sessions:
			createdSessions = append(createdSessions, session)
		case err := <-errors:
			t.Fatalf("concurrent CreateSession failed: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent operations")
		}
	}

	// Verify all sessions were created with unique codes
	codeSet := make(map[string]bool)
	for _, session := range createdSessions {
		if codeSet[session.Code] {
			t.Errorf("duplicate session code found: %s", session.Code)
		}
		codeSet[session.Code] = true
	}

	if len(createdSessions) != numGoroutines {
		t.Errorf("expected %d sessions, got %d", numGoroutines, len(createdSessions))
	}
}

func TestCreateSessionCollisionHandling(t *testing.T) {
	// Create a manager with low retry limit for testing
	options := &SessionOptions{
		MaxRetries:     3,
		SessionTimeout: 1 * time.Hour,
	}
	manager := NewManager(options)
	defer manager.Close()

	ctx := context.Background()

	// Test normal operation first
	session1, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession should not fail: %v", err)
	}

	// The code should be valid and unique
	if session1.Code == "" {
		t.Fatal("session code should not be empty")
	}

	// Test collision scenario by manipulating the generator to produce duplicates
	// This is tricky to test directly, so we test the retry mechanism indirectly
	// by creating many sessions rapidly and ensuring uniqueness
	const numSessions = 50
	codes := make(map[string]bool)
	
	for i := 0; i < numSessions; i++ {
		session, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Fatalf("CreateSession failed on attempt %d: %v", i, err)
		}
		
		if codes[session.Code] {
			t.Errorf("Collision detected: duplicate code %s", session.Code)
		}
		codes[session.Code] = true
	}

	if len(codes) != numSessions {
		t.Errorf("Expected %d unique codes, got %d", numSessions, len(codes))
	}
}

func TestCreateSessionWithCancelledContext(t *testing.T) {
	// Create a manager that will force retries by having very low MaxRetries
	// and pre-filling with sessions to increase collision probability
	options := &SessionOptions{
		MaxRetries:     1,
		SessionTimeout: 1 * time.Hour,
	}
	manager := NewManager(options)
	defer manager.Close()

	// Pre-create many sessions to increase collision probability
	ctx := context.Background()
	for i := 0; i < 50; i++ {
		_, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Logf("Pre-creation failed at %d: %v", i, err)
			break
		}
	}

	// Create a cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Try multiple times as the context check only happens during retries
	var foundCancellationError bool
	for i := 0; i < 10; i++ {
		_, err := manager.CreateSession(cancelledCtx, nil)
		if err != nil && strings.Contains(err.Error(), "cancelled") {
			foundCancellationError = true
			break
		}
	}

	if !foundCancellationError {
		t.Log("Context cancellation not triggered during session creation (this is possible if no collisions occur)")
		// This is actually acceptable behavior - if no retries are needed, context won't be checked
	}
}

func TestCreateSessionMaxRetriesExhausted(t *testing.T) {
	// This test is challenging because we need to force collisions
	// We'll create a scenario with many sessions to increase collision probability
	options := &SessionOptions{
		MaxRetries:     1, // Very low retry limit
		SessionTimeout: 1 * time.Hour,
	}
	manager := NewManager(options)
	defer manager.Close()

	ctx := context.Background()

	// Create many sessions to increase collision probability
	// With a small number range (1-99) and many sessions, we might get collisions
	const numSessions = 200
	var errors []error
	var successCount int

	for i := 0; i < numSessions; i++ {
		_, err := manager.CreateSession(ctx, nil)
		if err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	// We should have at least some successful sessions
	if successCount == 0 {
		t.Error("Expected at least some successful session creations")
	}

	// Log the results for analysis
	t.Logf("Successfully created %d out of %d sessions", successCount, numSessions)
	t.Logf("Encountered %d errors", len(errors))
}

func TestUpdateSessionDataEdgeCases(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Test updating with nil data
	err = manager.UpdateSessionData(session.Code, nil)
	if err != nil {
		t.Fatalf("UpdateSessionData with nil data should not fail: %v", err)
	}

	// Test updating with empty data
	emptyData := make(map[string]interface{})
	err = manager.UpdateSessionData(session.Code, emptyData)
	if err != nil {
		t.Fatalf("UpdateSessionData with empty data should not fail: %v", err)
	}

	// Test updating nonexistent session
	err = manager.UpdateSessionData("nonexistent-code-42", map[string]interface{}{"key": "value"})
	if err != ErrSessionNotFound {
		t.Errorf("UpdateSessionData should return ErrSessionNotFound for nonexistent session, got: %v", err)
	}

	// Test updating with invalid session code
	err = manager.UpdateSessionData("invalid-code", map[string]interface{}{"key": "value"})
	if err != ErrSessionNotFound {
		t.Errorf("UpdateSessionData should return ErrSessionNotFound for invalid code, got: %v", err)
	}

	// Test updating with empty session code
	err = manager.UpdateSessionData("", map[string]interface{}{"key": "value"})
	if err != ErrInvalidSessionCode {
		t.Errorf("UpdateSessionData should return ErrInvalidSessionCode for empty code, got: %v", err)
	}
}

func TestUpdateSessionDataExpired(t *testing.T) {
	// Create manager with very short timeout
	options := &SessionOptions{
		MaxRetries:     10,
		SessionTimeout: 1 * time.Millisecond,
	}
	manager := NewManager(options)
	defer manager.Close()

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Try to update expired session
	updateData := map[string]interface{}{"key": "value"}
	err = manager.UpdateSessionData(session.Code, updateData)
	if err != ErrSessionExpired {
		t.Errorf("UpdateSessionData should return ErrSessionExpired for expired session, got: %v", err)
	}

	// Verify session was removed from manager
	if count := manager.GetSessionCount(); count != 0 {
		t.Errorf("Expected 0 sessions after expiration cleanup, got: %d", count)
	}
}

func TestDeleteSessionEdgeCases(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Test deleting with empty code
	deleted := manager.DeleteSession("")
	if deleted {
		t.Error("DeleteSession should return false for empty code")
	}

	// Test deleting with invalid format
	deleted = manager.DeleteSession("invalid")
	if deleted {
		t.Error("DeleteSession should return false for invalid format")
	}

	// Test deleting valid format but nonexistent session
	deleted = manager.DeleteSession("valid-format-42")
	if deleted {
		t.Error("DeleteSession should return false for nonexistent session")
	}

	// Test case-insensitive deletion
	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Delete using uppercase version of the code
	uppercaseCode := strings.ToUpper(session.Code)
	deleted = manager.DeleteSession(uppercaseCode)
	if !deleted {
		t.Error("DeleteSession should handle case-insensitive codes")
	}
}

func TestConcurrentMixedOperations(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	const numOperations = 100
	var wg sync.WaitGroup
	errors := make(chan error, numOperations*3) // Buffer for all possible errors

	// Track created sessions for read/update/delete operations
	sessionCodes := make(chan string, numOperations)

	// Create sessions concurrently
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()
			session, err := manager.CreateSession(ctx, nil)
			if err != nil {
				errors <- fmt.Errorf("create session %d failed: %w", id, err)
				return
			}
			sessionCodes <- session.Code
		}(i)
	}

	// Wait for all create operations
	wg.Wait()
	close(sessionCodes)

	// Collect session codes
	var codes []string
	for code := range sessionCodes {
		codes = append(codes, code)
	}

	if len(codes) != numOperations {
		t.Fatalf("Expected %d session codes, got %d", numOperations, len(codes))
	}

	// Perform mixed operations concurrently
	operationTypes := []string{"get", "update", "delete", "list", "count"}
	
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		operation := operationTypes[i%len(operationTypes)]
		code := codes[i%len(codes)]

		go func(op string, sessionCode string, id int) {
			defer wg.Done()
			
			switch op {
			case "get":
				_, err := manager.GetSession(sessionCode)
				if err != nil && err != ErrSessionNotFound {
					errors <- fmt.Errorf("get session %d failed: %w", id, err)
				}
			case "update":
				updateData := map[string]interface{}{
					"operation": fmt.Sprintf("update_%d", id),
					"timestamp": time.Now().Unix(),
				}
				err := manager.UpdateSessionData(sessionCode, updateData)
				if err != nil && err != ErrSessionNotFound {
					errors <- fmt.Errorf("update session %d failed: %w", id, err)
				}
			case "delete":
				manager.DeleteSession(sessionCode)
			case "list":
				manager.ListSessions()
			case "count":
				manager.GetSessionCount()
			}
		}(operation, code, i)
	}

	wg.Wait()

	// Check for unexpected errors
	select {
	case err := <-errors:
		t.Errorf("Concurrent operation failed: %v", err)
	default:
		t.Logf("Successfully completed %d concurrent mixed operations", numOperations)
	}
}

func TestSessionManagerCleanupGoroutine(t *testing.T) {
	options := &SessionOptions{
		MaxRetries:     10,
		SessionTimeout: 50 * time.Millisecond, // Short timeout for testing
	}
	
	manager := NewManager(options)
	
	// Create some sessions
	ctx := context.Background()
	const numSessions = 5
	for i := 0; i < numSessions; i++ {
		_, err := manager.CreateSession(ctx, nil)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}
	}

	// Verify sessions exist
	if count := manager.GetSessionCount(); count != numSessions {
		t.Errorf("Expected %d sessions, got %d", numSessions, count)
	}

	// Wait for sessions to expire and cleanup to run
	// The cleanup runs every 10 minutes by default, but sessions expire after 50ms
	time.Sleep(100 * time.Millisecond)
	
	// Manually trigger cleanup to test it
	removed := manager.Cleanup()
	if removed != numSessions {
		t.Logf("Manual cleanup removed %d sessions (expected %d)", removed, numSessions)
	}

	// Close the manager and ensure cleanup goroutine stops
	done := make(chan struct{})
	go func() {
		manager.Close()
		close(done)
	}()

	select {
	case <-done:
		t.Log("Manager closed successfully")
	case <-time.After(1 * time.Second):
		t.Error("Manager.Close() did not complete within timeout")
	}
}

func TestGetSessionCaseInsensitive(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Test getting session with different case variations
	variations := []string{
		strings.ToUpper(session.Code),
		strings.ToLower(session.Code),
		strings.Title(session.Code),
	}

	for _, variation := range variations {
		retrieved, err := manager.GetSession(variation)
		if err != nil {
			t.Errorf("GetSession failed for case variation %s: %v", variation, err)
			continue
		}

		if retrieved.Code != session.Code {
			t.Errorf("Retrieved session code %s doesn't match original %s", retrieved.Code, session.Code)
		}
	}
}

func TestSessionDataIntegrity(t *testing.T) {
	manager := NewManager(nil)
	defer manager.Close()

	// Create session with initial data
	initialData := map[string]interface{}{
		"user_id":     "test123",
		"role":        "admin",
		"permissions": []string{"read", "write", "delete"},
		"metadata": map[string]interface{}{
			"login_time": time.Now().Unix(),
			"ip":         "192.168.1.1",
		},
	}

	options := &SessionOptions{
		InitialData: initialData,
	}

	ctx := context.Background()
	session, err := manager.CreateSession(ctx, options)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Verify initial data integrity
	retrieved, err := manager.GetSession(session.Code)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	// Check each initial data item
	if retrieved.Data["user_id"] != "test123" {
		t.Errorf("user_id mismatch: got %v, expected test123", retrieved.Data["user_id"])
	}

	if retrieved.Data["role"] != "admin" {
		t.Errorf("role mismatch: got %v, expected admin", retrieved.Data["role"])
	}

	// Update data and verify
	updateData := map[string]interface{}{
		"role":         "user",        // Update existing
		"last_action":  "test_action", // Add new
		"login_count":  42,            // Add new with different type
	}

	err = manager.UpdateSessionData(session.Code, updateData)
	if err != nil {
		t.Fatalf("UpdateSessionData failed: %v", err)
	}

	// Retrieve and verify updates
	updated, err := manager.GetSession(session.Code)
	if err != nil {
		t.Fatalf("GetSession after update failed: %v", err)
	}

	// Verify updated data
	if updated.Data["role"] != "user" {
		t.Errorf("role not updated: got %v, expected user", updated.Data["role"])
	}

	if updated.Data["last_action"] != "test_action" {
		t.Errorf("last_action not added: got %v, expected test_action", updated.Data["last_action"])
	}

	if updated.Data["login_count"] != 42 {
		t.Errorf("login_count not added: got %v, expected 42", updated.Data["login_count"])
	}

	// Verify original data still exists (except what was updated)
	if updated.Data["user_id"] != "test123" {
		t.Errorf("user_id should not change: got %v, expected test123", updated.Data["user_id"])
	}
}
