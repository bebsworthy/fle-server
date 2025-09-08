package jsonrpc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// TestRequest_JSON tests JSON marshaling and unmarshaling of Request.
func TestRequest_JSON(t *testing.T) {
	tests := []struct {
		name        string
		request     *Request
		expectedJSON string
	}{
		{
			name: "request with params and ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "subtract",
				Params:         json.RawMessage(`{"minuend": 42, "subtrahend": 23}`),
				ID:             1,
			},
			expectedJSON: `{"jsonrpc":"2.0","method":"subtract","params":{"minuend": 42, "subtrahend": 23},"id":1}`,
		},
		{
			name: "notification without ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "update",
				Params:         json.RawMessage(`[1, 2, 3, 4, 5]`),
			},
			expectedJSON: `{"jsonrpc":"2.0","method":"update","params":[1, 2, 3, 4, 5]}`,
		},
		{
			name: "request without params",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "get_data",
				ID:             "get_data_id",
			},
			expectedJSON: `{"jsonrpc":"2.0","method":"get_data","id":"get_data_id"}`,
		},
		{
			name: "request with string ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "test",
				ID:             "string-id-123",
			},
			expectedJSON: `{"jsonrpc":"2.0","method":"test","id":"string-id-123"}`,
		},
		{
			name: "request with null ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "test",
				ID:             nil,
			},
			expectedJSON: `{"jsonrpc":"2.0","method":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonBytes, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Normalize whitespace for comparison
			normalizedExpected := strings.ReplaceAll(tt.expectedJSON, " ", "")
			normalizedActual := strings.ReplaceAll(string(jsonBytes), " ", "")

			if normalizedActual != normalizedExpected {
				t.Errorf("JSON marshaling mismatch.\nExpected: %s\nGot: %s", tt.expectedJSON, string(jsonBytes))
			}

			// Test unmarshaling
			var parsedRequest Request
			err = json.Unmarshal(jsonBytes, &parsedRequest)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Compare fields
			if parsedRequest.JSONRPCVersion != tt.request.JSONRPCVersion {
				t.Errorf("JSONRPCVersion mismatch: expected %s, got %s", tt.request.JSONRPCVersion, parsedRequest.JSONRPCVersion)
			}

			if parsedRequest.Method != tt.request.Method {
				t.Errorf("Method mismatch: expected %s, got %s", tt.request.Method, parsedRequest.Method)
			}

			// JSON unmarshaling converts numbers to float64, so compare carefully
			if !compareIDs(parsedRequest.ID, tt.request.ID) {
				t.Errorf("ID mismatch: expected %v, got %v", tt.request.ID, parsedRequest.ID)
			}

			// Compare params (raw message comparison)
			if !equalRawMessage(parsedRequest.Params, tt.request.Params) {
				t.Errorf("Params mismatch: expected %s, got %s", string(tt.request.Params), string(parsedRequest.Params))
			}
		})
	}
}

// TestRequest_IsNotification tests the IsNotification method.
func TestRequest_IsNotification(t *testing.T) {
	tests := []struct {
		name       string
		request    *Request
		isNotification bool
	}{
		{
			name: "request with numeric ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "test",
				ID:             1,
			},
			isNotification: false,
		},
		{
			name: "request with string ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "test",
				ID:             "test-id",
			},
			isNotification: false,
		},
		{
			name: "notification without ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "notify",
			},
			isNotification: true,
		},
		{
			name: "notification with nil ID",
			request: &Request{
				JSONRPCVersion: "2.0",
				Method:         "notify",
				ID:             nil,
			},
			isNotification: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.request.IsNotification() != tt.isNotification {
				t.Errorf("IsNotification() = %v, expected %v", tt.request.IsNotification(), tt.isNotification)
			}
		})
	}
}

// TestResponse_JSON tests JSON marshaling and unmarshaling of Response.
func TestResponse_JSON(t *testing.T) {
	tests := []struct {
		name        string
		response    *Response
		expectedJSON string
	}{
		{
			name: "success response",
			response: &Response{
				JSONRPCVersion: "2.0",
				Result:         "success",
				ID:             1,
			},
			expectedJSON: `{"jsonrpc":"2.0","result":"success","id":1}`,
		},
		{
			name: "error response",
			response: &Response{
				JSONRPCVersion: "2.0",
				Error: &Error{
					Code:    -32601,
					Message: "Method not found",
				},
				ID: 1,
			},
			expectedJSON: `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
		},
		{
			name: "response with complex result",
			response: &Response{
				JSONRPCVersion: "2.0",
				Result: map[string]interface{}{
					"name":  "John",
					"age":   30,
					"items": []int{1, 2, 3},
				},
				ID: "complex-result",
			},
			expectedJSON: `{"jsonrpc":"2.0","result":{"name":"John","age":30,"items":[1,2,3]},"id":"complex-result"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonBytes, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			// Test unmarshaling
			var parsedResponse Response
			err = json.Unmarshal(jsonBytes, &parsedResponse)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Compare basic fields
			if parsedResponse.JSONRPCVersion != tt.response.JSONRPCVersion {
				t.Errorf("JSONRPCVersion mismatch: expected %s, got %s", tt.response.JSONRPCVersion, parsedResponse.JSONRPCVersion)
			}

			// JSON unmarshaling converts numbers to float64, so compare carefully
			if !compareIDs(parsedResponse.ID, tt.response.ID) {
				t.Errorf("ID mismatch: expected %v, got %v", tt.response.ID, parsedResponse.ID)
			}

			// Check error consistency
			if tt.response.Error != nil && parsedResponse.Error == nil {
				t.Error("Error was lost during marshaling/unmarshaling")
			} else if tt.response.Error == nil && parsedResponse.Error != nil {
				t.Error("Unexpected error appeared during marshaling/unmarshaling")
			}
		})
	}
}

// TestResponse_IsError tests the IsError method.
func TestResponse_IsError(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		isError  bool
	}{
		{
			name: "success response",
			response: &Response{
				JSONRPCVersion: "2.0",
				Result:         "success",
				ID:             1,
			},
			isError: false,
		},
		{
			name: "error response",
			response: &Response{
				JSONRPCVersion: "2.0",
				Error: &Error{
					Code:    -32601,
					Message: "Method not found",
				},
				ID: 1,
			},
			isError: true,
		},
		{
			name: "response with nil error",
			response: &Response{
				JSONRPCVersion: "2.0",
				Result:         nil,
				Error:          nil,
				ID:             1,
			},
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.response.IsError() != tt.isError {
				t.Errorf("IsError() = %v, expected %v", tt.response.IsError(), tt.isError)
			}

			// IsSuccess should be the opposite of IsError
			if tt.response.IsSuccess() != !tt.isError {
				t.Errorf("IsSuccess() = %v, expected %v", tt.response.IsSuccess(), !tt.isError)
			}
		})
	}
}

// TestError_Error tests the Error method implementation.
func TestError_Error(t *testing.T) {
	tests := []struct {
		name        string
		error       *Error
		expectedMsg string
	}{
		{
			name: "error without data",
			error: &Error{
				Code:    -32601,
				Message: "Method not found",
			},
			expectedMsg: "JSON-RPC error -32601: Method not found",
		},
		{
			name: "error with string data",
			error: &Error{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Parameter 'name' is required",
			},
			expectedMsg: "JSON-RPC error -32602: Invalid params (data: Parameter 'name' is required)",
		},
		{
			name: "error with complex data",
			error: &Error{
				Code:    -32603,
				Message: "Internal error",
				Data:    map[string]string{"detail": "Database connection failed"},
			},
			expectedMsg: "JSON-RPC error -32603: Internal error (data: map[detail:Database connection failed])",
		},
		{
			name: "error with nil data",
			error: &Error{
				Code:    -32700,
				Message: "Parse error",
				Data:    nil,
			},
			expectedMsg: "JSON-RPC error -32700: Parse error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.error.Error() != tt.expectedMsg {
				t.Errorf("Error() = %q, expected %q", tt.error.Error(), tt.expectedMsg)
			}
		})
	}
}

// TestStandardErrors tests the predefined standard errors.
func TestStandardErrors(t *testing.T) {
	tests := []struct {
		name          string
		error         *Error
		expectedCode  int
		expectedMsg   string
	}{
		{"Parse Error", ErrParse, ParseError, "Parse error"},
		{"Invalid Request", ErrInvalidRequest, InvalidRequest, "Invalid Request"},
		{"Method Not Found", ErrMethodNotFound, MethodNotFound, "Method not found"},
		{"Invalid Params", ErrInvalidParams, InvalidParams, "Invalid params"},
		{"Internal Error", ErrInternal, InternalError, "Internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.error.Code != tt.expectedCode {
				t.Errorf("Code = %d, expected %d", tt.error.Code, tt.expectedCode)
			}
			if tt.error.Message != tt.expectedMsg {
				t.Errorf("Message = %q, expected %q", tt.error.Message, tt.expectedMsg)
			}
		})
	}
}

// TestNewError tests error creation functions.
func TestNewError(t *testing.T) {
	code := -32001
	message := "Custom error"
	
	err := NewError(code, message)
	if err.Code != code {
		t.Errorf("Code = %d, expected %d", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("Message = %q, expected %q", err.Message, message)
	}
	if err.Data != nil {
		t.Errorf("Data = %v, expected nil", err.Data)
	}
}

// TestNewErrorWithData tests error creation with data.
func TestNewErrorWithData(t *testing.T) {
	code := -32001
	message := "Custom error"
	data := map[string]string{"field": "value"}
	
	err := NewErrorWithData(code, message, data)
	if err.Code != code {
		t.Errorf("Code = %d, expected %d", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("Message = %q, expected %q", err.Message, message)
	}
	// Use reflect.DeepEqual for map comparison
	if !reflect.DeepEqual(err.Data, data) {
		t.Errorf("Data = %v, expected %v", err.Data, data)
	}
}

// TestNewRequest tests request creation.
func TestNewRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		params        interface{}
		id            interface{}
		expectError   bool
	}{
		{
			name:   "simple request",
			method: "test.method",
			params: nil,
			id:     1,
			expectError: false,
		},
		{
			name:   "request with params",
			method: "subtract",
			params: map[string]int{"minuend": 42, "subtrahend": 23},
			id:     "subtract-call",
			expectError: false,
		},
		{
			name:   "request with array params",
			method: "sum",
			params: []int{1, 2, 3, 4, 5},
			id:     2,
			expectError: false,
		},
		{
			name:   "request with invalid params (unmarshalable)",
			method: "test",
			params: make(chan int), // channels cannot be marshaled
			id:     3,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(tt.method, tt.params, tt.id)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if req.JSONRPCVersion != Version {
				t.Errorf("JSONRPCVersion = %s, expected %s", req.JSONRPCVersion, Version)
			}
			if req.Method != tt.method {
				t.Errorf("Method = %s, expected %s", req.Method, tt.method)
			}
			if req.ID != tt.id {
				t.Errorf("ID = %v, expected %v", req.ID, tt.id)
			}

			// If params were provided, verify they can be unmarshaled
			if tt.params != nil {
				var params interface{}
				err := json.Unmarshal(req.Params, &params)
				if err != nil {
					t.Errorf("Failed to unmarshal params: %v", err)
				}
			}
		})
	}
}

// TestNewNotification tests notification creation.
func TestNewNotification(t *testing.T) {
	method := "notify.update"
	params := map[string]string{"status": "updated"}

	req, err := NewNotification(method, params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.JSONRPCVersion != Version {
		t.Errorf("JSONRPCVersion = %s, expected %s", req.JSONRPCVersion, Version)
	}
	if req.Method != method {
		t.Errorf("Method = %s, expected %s", req.Method, method)
	}
	if req.ID != nil {
		t.Errorf("ID = %v, expected nil", req.ID)
	}
	if !req.IsNotification() {
		t.Error("Request should be a notification")
	}
}

// TestNewResponse tests response creation.
func TestNewResponse(t *testing.T) {
	result := "success"
	id := "test-id"

	resp := NewResponse(result, id)

	if resp.JSONRPCVersion != Version {
		t.Errorf("JSONRPCVersion = %s, expected %s", resp.JSONRPCVersion, Version)
	}
	if resp.Result != result {
		t.Errorf("Result = %v, expected %v", resp.Result, result)
	}
	if resp.ID != id {
		t.Errorf("ID = %v, expected %v", resp.ID, id)
	}
	if resp.Error != nil {
		t.Errorf("Error = %v, expected nil", resp.Error)
	}
	if !resp.IsSuccess() {
		t.Error("Response should be a success response")
	}
}

// TestNewErrorResponse tests error response creation.
func TestNewErrorResponse(t *testing.T) {
	err := &Error{
		Code:    -32601,
		Message: "Method not found",
	}
	id := "error-id"

	resp := NewErrorResponse(err, id)

	if resp.JSONRPCVersion != Version {
		t.Errorf("JSONRPCVersion = %s, expected %s", resp.JSONRPCVersion, Version)
	}
	if resp.Result != nil {
		t.Errorf("Result = %v, expected nil", resp.Result)
	}
	if resp.ID != id {
		t.Errorf("ID = %v, expected %v", resp.ID, id)
	}
	if resp.Error != err {
		t.Errorf("Error = %v, expected %v", resp.Error, err)
	}
	if !resp.IsError() {
		t.Error("Response should be an error response")
	}
}

// TestErrorCodeValidation tests error code validation functions.
func TestErrorCodeValidation(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		isValid      bool
		isReserved   bool
		isServerError bool
	}{
		{"Parse Error", ParseError, true, true, false},
		{"Invalid Request", InvalidRequest, true, true, false},
		{"Method Not Found", MethodNotFound, true, true, false},
		{"Invalid Params", InvalidParams, true, true, false},
		{"Internal Error", InternalError, true, true, false},
		{"Server Error Start", ServerErrorStart, true, true, true},
		{"Server Error End", ServerErrorEnd, true, true, true},
		{"Server Error Middle", -32050, true, true, true},
		{"Application Error", 1001, true, false, false},
		{"Negative Application Error", -1001, true, false, false},
		{"Zero", 0, true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsValidErrorCode(tt.code) != tt.isValid {
				t.Errorf("IsValidErrorCode(%d) = %v, expected %v", tt.code, IsValidErrorCode(tt.code), tt.isValid)
			}
			if IsReservedErrorCode(tt.code) != tt.isReserved {
				t.Errorf("IsReservedErrorCode(%d) = %v, expected %v", tt.code, IsReservedErrorCode(tt.code), tt.isReserved)
			}
			if IsServerErrorCode(tt.code) != tt.isServerError {
				t.Errorf("IsServerErrorCode(%d) = %v, expected %v", tt.code, IsServerErrorCode(tt.code), tt.isServerError)
			}
		})
	}
}

// TestVersion tests the version constant.
func TestVersion(t *testing.T) {
	if Version != "2.0" {
		t.Errorf("Version = %s, expected '2.0'", Version)
	}
}

// Helper function to compare JSON raw messages.
func equalRawMessage(a, b json.RawMessage) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	// Normalize by unmarshaling and remarshaling
	var aVal, bVal interface{}
	
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	
	if err := json.Unmarshal(a, &aVal); err != nil {
		return string(a) == string(b)
	}
	if err := json.Unmarshal(b, &bVal); err != nil {
		return string(a) == string(b)
	}
	
	aBytes, _ := json.Marshal(aVal)
	bBytes, _ := json.Marshal(bVal)
	
	return string(aBytes) == string(bBytes)
}

// Helper function to compare IDs (handles JSON number conversion to float64)
func compareIDs(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	// Handle numeric ID comparison (JSON converts numbers to float64)
	aFloat, aIsFloat := a.(float64)
	bFloat, bIsFloat := b.(float64)
	aInt, aIsInt := a.(int)
	bInt, bIsInt := b.(int)
	
	if aIsFloat && bIsInt {
		return aFloat == float64(bInt)
	}
	if aIsInt && bIsFloat {
		return float64(aInt) == bFloat
	}
	if aIsInt && bIsInt {
		return aInt == bInt
	}
	if aIsFloat && bIsFloat {
		return aFloat == bFloat
	}
	
	// For non-numeric types, use direct comparison
	return a == b
}