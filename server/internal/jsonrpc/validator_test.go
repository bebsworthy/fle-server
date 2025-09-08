package jsonrpc

import (
	"strings"
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Fatal("NewValidator() returned nil")
	}
	if validator.validate == nil {
		t.Fatal("validator.validate is nil")
	}
}

func TestValidateRequest_Valid(t *testing.T) {
	validator := NewValidator()
	
	req := &Request{
		JSONRPCVersion: "2.0",
		Method:         "test.method",
		ID:             "test-id",
	}
	
	err := validator.ValidateRequest(req)
	if err != nil {
		t.Fatalf("ValidateRequest failed for valid request: %v", err)
	}
}

func TestValidateRequest_InvalidVersion(t *testing.T) {
	validator := NewValidator()
	
	req := &Request{
		JSONRPCVersion: "1.0", // Invalid version
		Method:         "test.method",
		ID:             "test-id",
	}
	
	err := validator.ValidateRequest(req)
	if err == nil {
		t.Fatal("ValidateRequest should have failed for invalid version")
	}
	
	// Check that we get detailed error information
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) == 0 {
			t.Fatal("Expected validation errors")
		}
		
		firstErr := validationErrs[0]
		if firstErr.Field != "jsonrpc" {
			t.Errorf("Expected field 'jsonrpc', got '%s'", firstErr.Field)
		}
		if firstErr.Tag != "eq" {
			t.Errorf("Expected tag 'eq', got '%s'", firstErr.Tag)
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

func TestValidateRequest_MissingMethod(t *testing.T) {
	validator := NewValidator()
	
	req := &Request{
		JSONRPCVersion: "2.0",
		Method:         "", // Empty method
		ID:             "test-id",
	}
	
	err := validator.ValidateRequest(req)
	if err == nil {
		t.Fatal("ValidateRequest should have failed for empty method")
	}
	
	// Check that we get detailed error information
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) == 0 {
			t.Fatal("Expected validation errors")
		}
		
		firstErr := validationErrs[0]
		if firstErr.Field != "method" {
			t.Errorf("Expected field 'method', got '%s'", firstErr.Field)
		}
		// For empty string, "required" validation fails first before "min"
		if firstErr.Tag != "required" {
			t.Errorf("Expected tag 'required', got '%s'", firstErr.Tag)
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

func TestValidateSessionCode_Valid(t *testing.T) {
	validator := NewValidator()
	
	validCodes := []string{
		"happy-panda-42",
		"blue-river-7",
		"quick-fox-99",
		"LOUD-BEAR-1", // Test case insensitivity
		"  space-wolf-33  ", // Test trimming
	}
	
	for _, code := range validCodes {
		err := validator.ValidateSessionCode(code)
		if err != nil {
			t.Errorf("ValidateSessionCode failed for valid code '%s': %v", code, err)
		}
	}
}

func TestValidateSessionCode_Invalid(t *testing.T) {
	validator := NewValidator()
	
	invalidCodes := []string{
		"", // Empty
		"happy-panda", // Missing number
		"happy-panda-0", // Number too low
		"happy-panda-100", // Number too high
		"happy-panda-abc", // Non-numeric suffix
		"happy-panda-42-extra", // Too many parts
		"single", // Single word
		"happy--42", // Empty middle part
		"-panda-42", // Empty first part
		"happy-", // Incomplete
	}
	
	for _, code := range invalidCodes {
		err := validator.ValidateSessionCode(code)
		if err == nil {
			t.Errorf("ValidateSessionCode should have failed for invalid code '%s'", code)
		}
	}
}

func TestValidateVar_SessionCode(t *testing.T) {
	validator := NewValidator()
	
	// Test valid session code
	err := validator.ValidateVar("happy-panda-42", "sessioncode")
	if err != nil {
		t.Errorf("ValidateVar failed for valid session code: %v", err)
	}
	
	// Test invalid session code
	err = validator.ValidateVar("invalid-code", "sessioncode")
	if err == nil {
		t.Error("ValidateVar should have failed for invalid session code")
	}
	
	// Test with required rule
	err = validator.ValidateVar("", "required,sessioncode")
	if err == nil {
		t.Error("ValidateVar should have failed for empty required session code")
	}
}

func TestValidateError_Valid(t *testing.T) {
	validator := NewValidator()
	
	err := &Error{
		Code:    -32601,
		Message: "Method not found",
	}
	
	validationErr := validator.ValidateError(err)
	if validationErr != nil {
		t.Fatalf("ValidateError failed for valid error: %v", validationErr)
	}
}

func TestValidateError_Invalid(t *testing.T) {
	validator := NewValidator()
	
	err := &Error{
		Code:    -32601,
		Message: "", // Empty message should fail
	}
	
	validationErr := validator.ValidateError(err)
	if validationErr == nil {
		t.Fatal("ValidateError should have failed for empty message")
	}
}

func TestValidateResponse_Valid(t *testing.T) {
	validator := NewValidator()
	
	resp := &Response{
		JSONRPCVersion: "2.0",
		Result:         "success",
		ID:             "test-id",
	}
	
	err := validator.ValidateResponse(resp)
	if err != nil {
		t.Fatalf("ValidateResponse failed for valid response: %v", err)
	}
}

func TestFastFailValidation(t *testing.T) {
	validator := NewValidator()
	
	// Create a request with multiple validation errors
	req := &Request{
		JSONRPCVersion: "1.0", // Invalid version
		Method:         "",    // Empty method (also invalid)
		ID:             "test-id",
	}
	
	err := validator.ValidateRequest(req)
	if err == nil {
		t.Fatal("ValidateRequest should have failed")
	}
	
	// With fast-fail, we should only get ONE error (the first field that failed)
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) != 1 {
			t.Errorf("Expected exactly 1 validation error (fast-fail), got %d", len(validationErrs))
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

func TestGetSupportedTags(t *testing.T) {
	validator := NewValidator()
	
	tags := validator.GetSupportedTags()
	if len(tags) == 0 {
		t.Error("GetSupportedTags() returned empty slice")
	}
	
	// Check that our custom tags are included
	customTags := []string{"sessioncode", "jsonrpcversion"}
	for _, customTag := range customTags {
		found := false
		for _, tag := range tags {
			if tag == customTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Custom tag '%s' not found in supported tags", customTag)
		}
	}
}

func TestValidationErrorMessages(t *testing.T) {
	validator := NewValidator()
	
	testCases := []struct {
		name          string
		req           *Request
		expectedField string
		expectedTag   string
	}{
		{
			name: "required field",
			req: &Request{
				JSONRPCVersion: "", // Required field missing
				Method:         "test",
				ID:             "test-id",
			},
			expectedField: "jsonrpc",
			expectedTag:   "required",
		},
		{
			name: "minimum length",
			req: &Request{
				JSONRPCVersion: "2.0",
				Method:         "", // Empty string fails "required" first
				ID:             "test-id",
			},
			expectedField: "method",
			expectedTag:   "required", // "required" validation fails before "min"
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateRequest(tc.req)
			if err == nil {
				t.Fatal("Expected validation error")
			}
			
			if validationErrs, ok := err.(ValidationErrors); ok {
				if len(validationErrs) == 0 {
					t.Fatal("Expected validation errors")
				}
				
				firstErr := validationErrs[0]
				if firstErr.Field != tc.expectedField {
					t.Errorf("Expected field '%s', got '%s'", tc.expectedField, firstErr.Field)
				}
				if firstErr.Tag != tc.expectedTag {
					t.Errorf("Expected tag '%s', got '%s'", tc.expectedTag, firstErr.Tag)
				}
				if firstErr.Message == "" {
					t.Error("Expected non-empty error message")
				}
			} else {
				t.Fatalf("Expected ValidationErrors, got %T", err)
			}
		})
	}
}

// TestValidateRequest_NilRequest tests validation with nil request.
func TestValidateRequest_NilRequest(t *testing.T) {
	validator := NewValidator()
	
	err := validator.ValidateRequest(nil)
	if err == nil {
		t.Fatal("ValidateRequest should have failed for nil request")
	}
	
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) == 0 {
			t.Fatal("Expected validation errors")
		}
		
		firstErr := validationErrs[0]
		if firstErr.Field != "request" {
			t.Errorf("Expected field 'request', got '%s'", firstErr.Field)
		}
		if firstErr.Tag != "required" {
			t.Errorf("Expected tag 'required', got '%s'", firstErr.Tag)
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

// TestValidateResponse_NilResponse tests validation with nil response.
func TestValidateResponse_NilResponse(t *testing.T) {
	validator := NewValidator()
	
	err := validator.ValidateResponse(nil)
	if err == nil {
		t.Fatal("ValidateResponse should have failed for nil response")
	}
	
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) == 0 {
			t.Fatal("Expected validation errors")
		}
		
		firstErr := validationErrs[0]
		if firstErr.Field != "response" {
			t.Errorf("Expected field 'response', got '%s'", firstErr.Field)
		}
		if firstErr.Tag != "required" {
			t.Errorf("Expected tag 'required', got '%s'", firstErr.Tag)
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

// TestValidateError_NilError tests validation with nil error.
func TestValidateError_NilError(t *testing.T) {
	validator := NewValidator()
	
	err := validator.ValidateError(nil)
	if err == nil {
		t.Fatal("ValidateError should have failed for nil error")
	}
	
	if validationErrs, ok := err.(ValidationErrors); ok {
		if len(validationErrs) == 0 {
			t.Fatal("Expected validation errors")
		}
		
		firstErr := validationErrs[0]
		if firstErr.Field != "error" {
			t.Errorf("Expected field 'error', got '%s'", firstErr.Field)
		}
		if firstErr.Tag != "required" {
			t.Errorf("Expected tag 'required', got '%s'", firstErr.Tag)
		}
	} else {
		t.Fatalf("Expected ValidationErrors, got %T", err)
	}
}

// TestValidateVar_AllTags tests variable validation with various tags.
func TestValidateVar_AllTags(t *testing.T) {
	validator := NewValidator()
	
	tests := []struct {
		name    string
		value   interface{}
		rules   string
		expectError bool
	}{
		{"required valid", "test", "required", false},
		{"required invalid", "", "required", true},
		{"min valid", "hello", "min=3", false},
		{"min invalid", "hi", "min=3", true},
		{"max valid", "hi", "max=5", false},
		{"max invalid", "toolong", "max=5", true},
		{"numeric valid", "123", "numeric", false},
		{"numeric invalid", "abc", "numeric", true},
		{"alpha valid", "abc", "alpha", false},
		{"alpha invalid", "abc123", "alpha", true},
		{"email valid", "test@example.com", "email", false},
		{"email invalid", "notanemail", "email", true},
		{"url valid", "https://example.com", "url", false},
		{"url invalid", "notaurl", "url", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateVar(tt.value, tt.rules)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error for value %v with rules %s", tt.value, tt.rules)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error for value %v with rules %s: %v", tt.value, tt.rules, err)
				}
			}
		})
	}
}

// TestCustomValidators tests all custom validation functions.
func TestCustomValidators(t *testing.T) {
	validator := NewValidator()
	
	// Test jsonrpcversion validator
	t.Run("jsonrpcversion", func(t *testing.T) {
		validValues := []string{"2.0"}
		invalidValues := []string{"", "1.0", "2.1", "3.0", "2", "2.0.0"}
		
		for _, value := range validValues {
			err := validator.ValidateVar(value, "jsonrpcversion")
			if err != nil {
				t.Errorf("ValidateVar failed for valid jsonrpcversion '%s': %v", value, err)
			}
		}
		
		for _, value := range invalidValues {
			err := validator.ValidateVar(value, "jsonrpcversion")
			if err == nil {
				t.Errorf("ValidateVar should have failed for invalid jsonrpcversion '%s'", value)
			}
		}
	})
	
	// Test sessioncode validator - comprehensive edge cases
	t.Run("sessioncode", func(t *testing.T) {
		validCodes := []string{
			"happy-panda-1",
			"blue-river-99",
			"LOUD-BEAR-50",
			"  space-wolf-33  ",
			"quick-fox-7",
		}
		
		invalidCodes := []string{
			"",                    // empty
			"happy-panda",         // missing number
			"happy-panda-0",       // number too low
			"happy-panda-100",     // number too high
			"happy-panda-abc",     // non-numeric suffix
			"happy-panda-42-extra", // too many parts
			"single",              // single word
			"happy--42",           // empty middle part
			"-panda-42",           // empty first part
			"happy-",              // incomplete
			"happy-panda-",        // incomplete number
			"happy-panda--42",     // empty middle with extra dash
			"   ",                 // only whitespace
			"-",                   // single dash
			"--",                  // double dash
			"a-b-c",               // non-numeric last part
		}
		
		for _, code := range validCodes {
			err := validator.ValidateVar(code, "sessioncode")
			if err != nil {
				t.Errorf("ValidateVar failed for valid sessioncode '%s': %v", code, err)
			}
		}
		
		for _, code := range invalidCodes {
			err := validator.ValidateVar(code, "sessioncode")
			if err == nil {
				t.Errorf("ValidateVar should have failed for invalid sessioncode '%s'", code)
			}
		}
	})
}

// TestValidationErrors_Error tests the ValidationErrors.Error() method.
func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   ValidationErrors
		expected string
	}{
		{
			name:     "empty errors",
			errors:   ValidationErrors{},
			expected: "validation errors",
		},
		{
			name: "single error",
			errors: ValidationErrors{
				{Message: "Field 'name' is required"},
			},
			expected: "validation failed: Field 'name' is required",
		},
		{
			name: "multiple errors",
			errors: ValidationErrors{
				{Message: "Field 'name' is required"},
				{Message: "Field 'age' must be greater than 0"},
			},
			expected: "validation failed: Field 'name' is required; Field 'age' must be greater than 0",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errors.Error() != tt.expected {
				t.Errorf("ValidationErrors.Error() = %q, expected %q", tt.errors.Error(), tt.expected)
			}
		})
	}
}

// TestValidationError_Error tests the ValidationError.Error() method.
func TestValidationError_Error(t *testing.T) {
	err := ValidationError{Message: "Field is required"}
	if err.Error() != "Field is required" {
		t.Errorf("ValidationError.Error() = %q, expected %q", err.Error(), "Field is required")
	}
}

// TestBuildErrorMessage tests all error message building scenarios.
func TestBuildErrorMessage(t *testing.T) {
	validator := NewValidator()
	
	// Test different validation scenarios to trigger different error messages
	tests := []struct {
		name          string
		value         interface{}
		rules         string
		expectContains string
	}{
		{"eq validation", "1.0", "eq=2.0", "must equal '2.0'"},
		{"min string validation", "ab", "min=3", "must be at least 3 characters"},
		{"max string validation", "toolong", "max=5", "must be at most 5 characters"},
		{"numeric validation", "abc", "numeric", "must be numeric"},
		{"alpha validation", "abc123", "alpha", "must contain only alphabetic characters"},
		{"email validation", "notanemail", "email", "must be a valid email address"},
		{"url validation", "notaurl", "url", "must be a valid URL"},
		{"sessioncode validation", "invalid", "sessioncode", "must be a valid session code"},
		{"jsonrpcversion validation", "1.0", "jsonrpcversion", "must be exactly '2.0'"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateVar(tt.value, tt.rules)
			if err == nil {
				t.Fatalf("Expected validation error for %v with rules %s", tt.value, tt.rules)
			}
			
			if validationErrs, ok := err.(ValidationErrors); ok && len(validationErrs) > 0 {
				message := validationErrs[0].Message
				if !strings.Contains(message, tt.expectContains) {
					t.Errorf("Error message %q does not contain expected text %q", message, tt.expectContains)
				}
			} else {
				t.Fatalf("Expected ValidationErrors, got %T", err)
			}
		})
	}
}