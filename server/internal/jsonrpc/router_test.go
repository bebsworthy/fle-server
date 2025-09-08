package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewRouter tests router creation.
func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("NewRouter() returned nil")
	}
	
	if router.MethodCount() != 0 {
		t.Errorf("Expected 0 methods, got %d", router.MethodCount())
	}
}

// TestRegisterMethod tests method registration.
func TestRegisterMethod(t *testing.T) {
	router := NewRouter()
	
	// Test registering a simple method
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	err := router.RegisterSimpleMethod("test.method", handler, "Test method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	if !router.HasMethod("test.method") {
		t.Error("Method should be registered")
	}
	
	if router.MethodCount() != 1 {
		t.Errorf("Expected 1 method, got %d", router.MethodCount())
	}
}

// TestRegisterMethodErrors tests method registration error cases.
func TestRegisterMethodErrors(t *testing.T) {
	router := NewRouter()
	
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	// Test empty method name
	err := router.RegisterSimpleMethod("", handler, "Test")
	if err == nil {
		t.Error("Should fail with empty method name")
	}
	
	// Test nil handler
	err = router.RegisterSimpleMethod("test.method", nil, "Test")
	if err == nil {
		t.Error("Should fail with nil handler")
	}
	
	// Test duplicate registration
	_ = router.RegisterSimpleMethod("test.method", handler, "Test")
	err = router.RegisterSimpleMethod("test.method", handler, "Test")
	if err == nil {
		t.Error("Should fail with duplicate method name")
	}
}

// TestRouteSimpleMethod tests routing a simple method call.
func TestRouteSimpleMethod(t *testing.T) {
	router := NewRouter()
	
	// Register a test method
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return map[string]string{"message": "hello"}, nil
	}
	
	err := router.RegisterSimpleMethod("test.hello", handler, "Test method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Create a test request
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "test.hello",
		ID:             "test-123",
	}
	
	// Route the request
	response := router.Route(context.Background(), request)
	
	// Verify the response
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if response.IsError() {
		t.Fatalf("Expected success response, got error: %v", response.Error)
	}
	
	if response.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got %v", response.ID)
	}
	
	result, ok := response.Result.(map[string]string)
	if !ok {
		t.Errorf("Expected map[string]string result, got %T", response.Result)
	} else if result["message"] != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", result["message"])
	}
}

// TestRouteMethodNotFound tests routing with unknown method.
func TestRouteMethodNotFound(t *testing.T) {
	router := NewRouter()
	
	// Create a request for non-existent method
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "unknown.method",
		ID:             "test-123",
	}
	
	// Route the request
	response := router.Route(context.Background(), request)
	
	// Verify error response
	if response == nil {
		t.Fatal("Expected error response, got nil")
	}
	
	if !response.IsError() {
		t.Fatal("Expected error response, got success")
	}
	
	if response.Error.Code != MethodNotFound {
		t.Errorf("Expected error code %d, got %d", MethodNotFound, response.Error.Code)
	}
}

// TestRouteInvalidRequest tests routing with invalid request.
func TestRouteInvalidRequest(t *testing.T) {
	router := NewRouter()
	
	// Create an invalid request (missing method)
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "", // Invalid empty method
		ID:             "test-123",
	}
	
	// Route the request
	response := router.Route(context.Background(), request)
	
	// Verify error response
	if response == nil {
		t.Fatal("Expected error response, got nil")
	}
	
	if !response.IsError() {
		t.Fatal("Expected error response, got success")
	}
	
	if response.Error.Code != InvalidRequest {
		t.Errorf("Expected error code %d, got %d", InvalidRequest, response.Error.Code)
	}
}

// TestRouteNotification tests routing notifications (requests without ID).
func TestRouteNotification(t *testing.T) {
	router := NewRouter()
	
	// Register a test method
	called := false
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		called = true
		return "success", nil
	}
	
	err := router.RegisterSimpleMethod("test.notification", handler, "Test notification")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Create a notification (no ID)
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "test.notification",
		// No ID = notification
	}
	
	// Route the notification
	response := router.Route(context.Background(), request)
	
	// Verify no response for notification
	if response != nil {
		t.Error("Expected nil response for notification, got response")
	}
	
	// Verify handler was called
	if !called {
		t.Error("Handler should have been called for notification")
	}
}

// TestRouteJSON tests the JSON convenience method.
func TestRouteJSON(t *testing.T) {
	router := NewRouter()
	
	// Register a test method
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return map[string]string{"result": "ok"}, nil
	}
	
	err := router.RegisterSimpleMethod("test.json", handler, "Test JSON method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Create JSON request
	requestJSON := []byte(`{
		"jsonrpc": "2.0",
		"method": "test.json",
		"id": "json-test"
	}`)
	
	// Route the JSON request
	responseJSON, err := router.RouteJSON(context.Background(), requestJSON)
	if err != nil {
		t.Fatalf("RouteJSON failed: %v", err)
	}
	
	// Parse and verify the response
	var response Response
	err = json.Unmarshal(responseJSON, &response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if response.IsError() {
		t.Fatalf("Expected success response, got error: %v", response.Error)
	}
	
	if response.ID != "json-test" {
		t.Errorf("Expected ID 'json-test', got %v", response.ID)
	}
}

// TestRouteJSONParseError tests JSON parsing error handling.
func TestRouteJSONParseError(t *testing.T) {
	router := NewRouter()
	
	// Invalid JSON
	requestJSON := []byte(`{"invalid": json}`)
	
	// Route the invalid JSON
	responseJSON, err := router.RouteJSON(context.Background(), requestJSON)
	if err != nil {
		t.Fatalf("RouteJSON should handle parse errors, got: %v", err)
	}
	
	// Parse and verify the error response
	var response Response
	err = json.Unmarshal(responseJSON, &response)
	if err != nil {
		t.Fatalf("Failed to parse error response JSON: %v", err)
	}
	
	if !response.IsError() {
		t.Fatal("Expected error response for invalid JSON")
	}
	
	if response.Error.Code != ParseError {
		t.Errorf("Expected parse error code %d, got %d", ParseError, response.Error.Code)
	}
}

// TestUnregisterMethod tests method unregistration.
func TestUnregisterMethod(t *testing.T) {
	router := NewRouter()
	
	// Register a method
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	err := router.RegisterSimpleMethod("test.method", handler, "Test method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	if !router.HasMethod("test.method") {
		t.Error("Method should be registered")
	}
	
	// Unregister the method
	err = router.UnregisterMethod("test.method")
	if err != nil {
		t.Fatalf("Failed to unregister method: %v", err)
	}
	
	if router.HasMethod("test.method") {
		t.Error("Method should be unregistered")
	}
	
	if router.MethodCount() != 0 {
		t.Errorf("Expected 0 methods after unregistration, got %d", router.MethodCount())
	}
}

// TestGetMethods tests getting all registered methods.
func TestGetMethods(t *testing.T) {
	router := NewRouter()
	
	// Register multiple methods
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	methods := []string{"test.method1", "test.method2", "test.method3"}
	for _, method := range methods {
		err := router.RegisterSimpleMethod(method, handler, "Test method")
		if err != nil {
			t.Fatalf("Failed to register method %s: %v", method, err)
		}
	}
	
	// Get all methods
	registeredMethods := router.GetMethods()
	
	if len(registeredMethods) != len(methods) {
		t.Errorf("Expected %d methods, got %d", len(methods), len(registeredMethods))
	}
	
	// Verify all methods are present (order may vary)
	methodMap := make(map[string]bool)
	for _, method := range registeredMethods {
		methodMap[method] = true
	}
	
	for _, expectedMethod := range methods {
		if !methodMap[expectedMethod] {
			t.Errorf("Method %s not found in registered methods", expectedMethod)
		}
	}
}

// TestRouteWithValidation tests routing with parameter validation.
func TestRouteWithValidation(t *testing.T) {
	router := NewRouter()
	
	// Define a parameter schema
	type TestParams struct {
		Name string `json:"name" validate:"required,min=2"`
		Age  int    `json:"age" validate:"required,gt=0"`
	}
	
	// Register a method with validation
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		var p TestParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"greeting": "Hello, " + p.Name,
			"age":      p.Age,
		}, nil
	}
	
	err := router.RegisterMethodWithValidation(
		"test.validate",
		handler,
		reflect.TypeOf(TestParams{}), // params schema as type
		nil,                          // no result schema
		"Test method with validation",
	)
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Test with valid params
	t.Run("valid params", func(t *testing.T) {
		params := map[string]interface{}{
			"name": "Alice",
			"age":  25,
		}
		paramsJSON, _ := json.Marshal(params)
		
		request := &Request{
			JSONRPCVersion: "2.0",
			Method:         "test.validate",
			Params:         paramsJSON,
			ID:             "valid-test",
		}
		
		response := router.Route(context.Background(), request)
		if response.IsError() {
			t.Fatalf("Expected success, got error: %v", response.Error)
		}
	})
	
	// Test with invalid params (missing required field)
	t.Run("invalid params - missing name", func(t *testing.T) {
		params := map[string]interface{}{
			"age": 25,
			// name is missing
		}
		paramsJSON, _ := json.Marshal(params)
		
		request := &Request{
			JSONRPCVersion: "2.0",
			Method:         "test.validate",
			Params:         paramsJSON,
			ID:             "invalid-test",
		}
		
		response := router.Route(context.Background(), request)
		if !response.IsError() {
			t.Fatalf("Expected validation error, got success response with result: %v", response.Result)
		}
		
		if response.Error.Code != InvalidParams {
			t.Errorf("Expected InvalidParams error (%d), got %d: %s", InvalidParams, response.Error.Code, response.Error.Message)
		}
	})
	
	// Test with invalid params (validation rule violation)
	t.Run("invalid params - age too low", func(t *testing.T) {
		params := map[string]interface{}{
			"name": "Bob",
			"age":  0, // violates gt=0 rule
		}
		paramsJSON, _ := json.Marshal(params)
		
		request := &Request{
			JSONRPCVersion: "2.0",
			Method:         "test.validate",
			Params:         paramsJSON,
			ID:             "invalid-age-test",
		}
		
		response := router.Route(context.Background(), request)
		if !response.IsError() {
			t.Fatal("Expected validation error")
		}
		
		if response.Error.Code != InvalidParams {
			t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
		}
	})
}

// TestRoutePanicRecovery tests panic recovery in handlers.
func TestRoutePanicRecovery(t *testing.T) {
	router := NewRouter()
	
	// Register a method that panics
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		panic("something went wrong!")
	}
	
	err := router.RegisterSimpleMethod("test.panic", handler, "Panicking method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "test.panic",
		ID:             "panic-test",
	}
	
	response := router.Route(context.Background(), request)
	
	// Should recover from panic and return internal error
	if !response.IsError() {
		t.Fatal("Expected error response from panicking handler")
	}
	
	if response.Error.Code != InternalError {
		t.Errorf("Expected InternalError, got %d", response.Error.Code)
	}
	
	// Error message should indicate a panic
	if !strings.Contains(response.Error.Data.(string), "panic") {
		t.Error("Expected error data to mention panic")
	}
}

// TestConcurrentRouting tests concurrent request processing.
func TestConcurrentRouting(t *testing.T) {
	router := NewRouter()
	
	// Register a method that takes some time
	requestCount := int32(0)
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		atomic.AddInt32(&requestCount, 1)
		// Simulate some processing time
		select {
		case <-time.After(10 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return "processed", nil
	}
	
	err := router.RegisterSimpleMethod("test.concurrent", handler, "Concurrent test method")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Run multiple concurrent requests
	const numRequests = 10
	var wg sync.WaitGroup
	responses := make([]*Response, numRequests)
	
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()
			
			request := &Request{
				JSONRPCVersion: "2.0",
				Method:         "test.concurrent",
				ID:             fmt.Sprintf("concurrent-%d", index),
			}
			
			responses[index] = router.Route(context.Background(), request)
		}(i)
	}
	
	wg.Wait()
	
	// Verify all requests were processed successfully
	for i, response := range responses {
		if response.IsError() {
			t.Errorf("Request %d failed: %v", i, response.Error)
		}
		if response.Result != "processed" {
			t.Errorf("Request %d: expected 'processed', got %v", i, response.Result)
		}
	}
	
	// Verify all handlers were called
	finalCount := atomic.LoadInt32(&requestCount)
	if int(finalCount) != numRequests {
		t.Errorf("Expected %d handler calls, got %d", numRequests, finalCount)
	}
}

// TestMethodInfo tests getting method information.
func TestMethodInfo(t *testing.T) {
	router := NewRouter()
	
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	description := "Test method for info"
	err := router.RegisterSimpleMethod("test.info", handler, description)
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Test getting existing method info
	info, err := router.GetMethodInfo("test.info")
	if err != nil {
		t.Fatalf("Failed to get method info: %v", err)
	}
	
	if info.Description != description {
		t.Errorf("Expected description '%s', got '%s'", description, info.Description)
	}
	
	// Test getting non-existent method info
	_, err = router.GetMethodInfo("non.existent")
	if err == nil {
		t.Error("Expected error for non-existent method")
	}
}

// TestClearRouter tests clearing all methods.
func TestClearRouter(t *testing.T) {
	router := NewRouter()
	
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	// Register some methods
	methods := []string{"test.method1", "test.method2", "test.method3"}
	for _, method := range methods {
		err := router.RegisterSimpleMethod(method, handler, "Test method")
		if err != nil {
			t.Fatalf("Failed to register method %s: %v", method, err)
		}
	}
	
	if router.MethodCount() != len(methods) {
		t.Errorf("Expected %d methods before clear, got %d", len(methods), router.MethodCount())
	}
	
	// Clear all methods
	router.Clear()
	
	if router.MethodCount() != 0 {
		t.Errorf("Expected 0 methods after clear, got %d", router.MethodCount())
	}
	
	// Verify methods are no longer registered
	for _, method := range methods {
		if router.HasMethod(method) {
			t.Errorf("Method %s should not exist after clear", method)
		}
	}
}

// TestUnregisterNonExistentMethod tests unregistering a method that doesn't exist.
func TestUnregisterNonExistentMethod(t *testing.T) {
	router := NewRouter()
	
	err := router.UnregisterMethod("non.existent")
	if err == nil {
		t.Error("Expected error when unregistering non-existent method")
	}
}

// TestRouteNotificationForNonExistentMethod tests notifications for non-existent methods.
func TestRouteNotificationForNonExistentMethod(t *testing.T) {
	router := NewRouter()
	
	// Create a notification for non-existent method
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "non.existent",
		// No ID = notification
	}
	
	// Route the notification - should silently ignore
	response := router.Route(context.Background(), request)
	
	// Should return nil (no response for notifications)
	if response != nil {
		t.Error("Expected nil response for notification to non-existent method")
	}
}

// TestRouteNotificationWithInvalidParams tests notifications with invalid parameters.
func TestRouteNotificationWithInvalidParams(t *testing.T) {
	router := NewRouter()
	
	// Define a parameter schema
	type TestParams struct {
		Name string `json:"name" validate:"required"`
	}
	
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		return "success", nil
	}
	
	err := router.RegisterMethodWithValidation(
		"test.notification",
		handler,
		reflect.TypeOf(TestParams{}), // params schema as type
		nil,                          // no result schema
		"Test notification with validation",
	)
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Create a notification with invalid params
	params := map[string]interface{}{
		// name is missing - should fail validation
	}
	paramsJSON, _ := json.Marshal(params)
	
	request := &Request{
		JSONRPCVersion: "2.0",
		Method:         "test.notification",
		Params:         paramsJSON,
		// No ID = notification
	}
	
	// Route the notification - should silently ignore validation errors
	response := router.Route(context.Background(), request)
	
	// Should return nil (no response for notifications, even with validation errors)
	if response != nil {
		t.Error("Expected nil response for notification with invalid params")
	}
}

// TestRouteJSONNotification tests JSON routing for notifications.
func TestRouteJSONNotification(t *testing.T) {
	router := NewRouter()
	
	called := false
	handler := func(ctx context.Context, params json.RawMessage) (interface{}, error) {
		called = true
		return "success", nil
	}
	
	err := router.RegisterSimpleMethod("test.notify", handler, "Test notification")
	if err != nil {
		t.Fatalf("Failed to register method: %v", err)
	}
	
	// Create JSON notification
	requestJSON := []byte(`{
		"jsonrpc": "2.0",
		"method": "test.notify"
	}`)
	
	// Route the JSON notification
	responseJSON, err := router.RouteJSON(context.Background(), requestJSON)
	if err != nil {
		t.Fatalf("RouteJSON failed: %v", err)
	}
	
	// Should return nil for notifications
	if responseJSON != nil {
		t.Error("Expected nil response JSON for notification")
	}
	
	// Verify handler was called
	if !called {
		t.Error("Handler should have been called for notification")
	}
}