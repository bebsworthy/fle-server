// Package jsonrpc provides JSON-RPC 2.0 routing and method dispatch functionality.
package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// HandlerFunc represents a JSON-RPC method handler function.
// It receives a context, parsed params, and returns a result and error.
// The params will be validated according to the registered schema before calling the handler.
type HandlerFunc func(ctx context.Context, params json.RawMessage) (interface{}, error)

// MethodInfo holds metadata about a registered JSON-RPC method.
type MethodInfo struct {
	// Handler is the function that handles requests for this method
	Handler HandlerFunc

	// ParamsSchema defines the validation schema for method parameters.
	// This can be a struct type, interface{}, or nil if no validation is needed.
	ParamsSchema interface{}

	// ResultSchema defines the validation schema for method results.
	// This can be a struct type, interface{}, or nil if no validation is needed.
	ResultSchema interface{}

	// Description provides human-readable documentation for the method
	Description string

	// ValidateParams indicates whether to validate incoming parameters
	ValidateParams bool

	// ValidateResult indicates whether to validate outgoing results
	ValidateResult bool
}

// Router provides JSON-RPC 2.0 method registration and request routing functionality.
// It is thread-safe and supports concurrent request processing with proper synchronization.
type Router struct {
	// methods stores registered method handlers with their metadata
	methods map[string]*MethodInfo

	// validator provides validation functionality for requests and responses
	validator *Validator

	// mutex protects concurrent access to the methods map
	mutex sync.RWMutex
}

// NewRouter creates a new JSON-RPC router with validation support.
func NewRouter() *Router {
	return &Router{
		methods:   make(map[string]*MethodInfo),
		validator: NewValidator(),
	}
}

// RegisterMethod registers a new JSON-RPC method with optional validation schemas.
// The method name should follow JSON-RPC naming conventions.
// Method names beginning with "rpc." are reserved for internal RPC methods.
func (r *Router) RegisterMethod(methodName string, handler HandlerFunc, info *MethodInfo) error {
	if methodName == "" {
		return fmt.Errorf("method name cannot be empty")
	}

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	// Validate method name format
	if err := r.validator.ValidateVar(methodName, "required,min=1"); err != nil {
		return fmt.Errorf("invalid method name: %w", err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if method is already registered
	if _, exists := r.methods[methodName]; exists {
		return fmt.Errorf("method '%s' is already registered", methodName)
	}

	// Create method info with defaults if not provided
	if info == nil {
		info = &MethodInfo{}
	}
	info.Handler = handler

	// Store the method
	r.methods[methodName] = info

	return nil
}

// RegisterMethodWithValidation is a convenience method for registering a method with validation schemas.
func (r *Router) RegisterMethodWithValidation(methodName string, handler HandlerFunc, paramsSchema, resultSchema interface{}, description string) error {
	info := &MethodInfo{
		ParamsSchema:   paramsSchema,
		ResultSchema:   resultSchema,
		Description:    description,
		ValidateParams: paramsSchema != nil,
		ValidateResult: resultSchema != nil,
	}

	return r.RegisterMethod(methodName, handler, info)
}

// RegisterSimpleMethod is a convenience method for registering a method without validation schemas.
func (r *Router) RegisterSimpleMethod(methodName string, handler HandlerFunc, description string) error {
	info := &MethodInfo{
		Description:    description,
		ValidateParams: false,
		ValidateResult: false,
	}

	return r.RegisterMethod(methodName, handler, info)
}

// UnregisterMethod removes a method from the router.
func (r *Router) UnregisterMethod(methodName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.methods[methodName]; !exists {
		return fmt.Errorf("method '%s' is not registered", methodName)
	}

	delete(r.methods, methodName)
	return nil
}

// HasMethod returns true if the specified method is registered.
func (r *Router) HasMethod(methodName string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.methods[methodName]
	return exists
}

// GetMethods returns a list of all registered method names.
func (r *Router) GetMethods() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	methods := make([]string, 0, len(r.methods))
	for methodName := range r.methods {
		methods = append(methods, methodName)
	}

	return methods
}

// GetMethodInfo returns the method information for a registered method.
func (r *Router) GetMethodInfo(methodName string) (*MethodInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	info, exists := r.methods[methodName]
	if !exists {
		return nil, fmt.Errorf("method '%s' is not registered", methodName)
	}

	return info, nil
}

// Route processes a JSON-RPC request and returns a response.
// This method handles request validation, method dispatch, and response formatting.
// It is thread-safe and can be called concurrently.
func (r *Router) Route(ctx context.Context, request *Request) *Response {
	// Validate the request structure
	if err := r.validator.ValidateRequest(request); err != nil {
		return NewErrorResponse(r.createValidationError(err), request.ID)
	}

	// Handle notifications (requests without ID)
	if request.IsNotification() {
		r.routeNotification(ctx, request)
		return nil // No response for notifications
	}

	// Find the method handler
	r.mutex.RLock()
	methodInfo, exists := r.methods[request.Method]
	r.mutex.RUnlock()

	if !exists {
		return NewErrorResponse(ErrMethodNotFound, request.ID)
	}

	// Validate parameters if schema is provided
	if methodInfo.ValidateParams && methodInfo.ParamsSchema != nil {
		if err := r.validateParams(request.Params, methodInfo.ParamsSchema); err != nil {
			return NewErrorResponse(r.createParamsError(err), request.ID)
		}
	}

	// Call the method handler
	result, err := r.callHandler(ctx, methodInfo.Handler, request.Params)
	if err != nil {
		return NewErrorResponse(r.createInternalError(err), request.ID)
	}

	// Validate result if schema is provided
	if methodInfo.ValidateResult && methodInfo.ResultSchema != nil {
		if err := r.validateResult(result, methodInfo.ResultSchema); err != nil {
			return NewErrorResponse(r.createInternalError(fmt.Errorf("result validation failed: %w", err)), request.ID)
		}
	}

	// Create and validate the response
	response := NewResponse(result, request.ID)
	if err := r.validator.ValidateResponse(response); err != nil {
		return NewErrorResponse(r.createInternalError(fmt.Errorf("response validation failed: %w", err)), request.ID)
	}

	return response
}

// routeNotification handles notification requests (requests without ID).
func (r *Router) routeNotification(ctx context.Context, request *Request) {
	// Find the method handler
	r.mutex.RLock()
	methodInfo, exists := r.methods[request.Method]
	r.mutex.RUnlock()

	if !exists {
		// Silently ignore notifications for non-existent methods as per JSON-RPC spec
		return
	}

	// Validate parameters if schema is provided
	if methodInfo.ValidateParams && methodInfo.ParamsSchema != nil {
		if err := r.validateParams(request.Params, methodInfo.ParamsSchema); err != nil {
			// Silently ignore invalid notifications as per JSON-RPC spec
			return
		}
	}

	// Call the method handler (ignore result and errors for notifications)
	_, _ = r.callHandler(ctx, methodInfo.Handler, request.Params)
}

// RouteJSON is a convenience method that accepts JSON bytes and returns JSON response.
// It handles JSON parsing and serialization automatically.
func (r *Router) RouteJSON(ctx context.Context, requestJSON []byte) ([]byte, error) {
	// Parse the request
	var request Request
	if err := json.Unmarshal(requestJSON, &request); err != nil {
		// Return parse error response
		response := NewErrorResponse(ErrParse, nil)
		return json.Marshal(response)
	}

	// Route the request
	response := r.Route(ctx, &request)

	// Handle notifications (no response)
	if response == nil {
		return nil, nil
	}

	// Marshal the response
	responseJSON, err := json.Marshal(response)
	if err != nil {
		// Return internal error if response marshaling fails
		errorResponse := NewErrorResponse(ErrInternal, request.ID)
		responseJSON, _ = json.Marshal(errorResponse)
	}

	return responseJSON, nil
}

// validateParams validates method parameters against the provided schema.
func (r *Router) validateParams(params json.RawMessage, schema interface{}) error {
	if params == nil {
		return nil
	}

	// If schema is a reflect.Type, create an instance
	if schemaType, ok := schema.(reflect.Type); ok {
		// Create a new instance of the schema type
		instance := reflect.New(schemaType).Interface()
		
		// Unmarshal params into the instance
		if err := json.Unmarshal(params, instance); err != nil {
			return fmt.Errorf("failed to parse params: %w", err)
		}

		// Validate the instance
		return r.validator.Validate(instance)
	}

	// If schema is a concrete type, unmarshal and validate directly
	if err := json.Unmarshal(params, schema); err != nil {
		return fmt.Errorf("failed to parse params: %w", err)
	}

	return r.validator.Validate(schema)
}

// validateResult validates method results against the provided schema.
func (r *Router) validateResult(result interface{}, schema interface{}) error {
	if result == nil {
		return nil
	}

	// If schema is a reflect.Type, validate that result matches the type
	if schemaType, ok := schema.(reflect.Type); ok {
		resultType := reflect.TypeOf(result)
		if !resultType.AssignableTo(schemaType) {
			return fmt.Errorf("result type %v is not assignable to schema type %v", resultType, schemaType)
		}
	}

	return r.validator.Validate(result)
}

// callHandler safely calls a method handler with error recovery.
func (r *Router) callHandler(ctx context.Context, handler HandlerFunc, params json.RawMessage) (result interface{}, err error) {
	// Recover from panics in handler code
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handler panic: %v", r)
		}
	}()

	return handler(ctx, params)
}

// createValidationError creates a JSON-RPC error from a validation error.
func (r *Router) createValidationError(err error) *Error {
	if validationErrs, ok := err.(ValidationErrors); ok && len(validationErrs) > 0 {
		return NewErrorWithData(InvalidRequest, "Request validation failed", validationErrs[0].Message)
	}
	return NewErrorWithData(InvalidRequest, "Request validation failed", err.Error())
}

// createParamsError creates a JSON-RPC error for parameter validation failures.
func (r *Router) createParamsError(err error) *Error {
	if validationErrs, ok := err.(ValidationErrors); ok && len(validationErrs) > 0 {
		return NewErrorWithData(InvalidParams, "Parameter validation failed", validationErrs[0].Message)
	}
	return NewErrorWithData(InvalidParams, "Parameter validation failed", err.Error())
}

// createInternalError creates a JSON-RPC internal error.
func (r *Router) createInternalError(err error) *Error {
	return NewErrorWithData(InternalError, "Internal error", err.Error())
}

// Clear removes all registered methods from the router.
// This is useful for testing or dynamic method management.
func (r *Router) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.methods = make(map[string]*MethodInfo)
}

// MethodCount returns the number of registered methods.
func (r *Router) MethodCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.methods)
}