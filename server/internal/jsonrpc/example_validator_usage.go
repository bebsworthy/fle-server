package jsonrpc

// This file contains example usage of the validation framework.
// It demonstrates how to use the validator with JSON-RPC types and custom data.

import (
	"fmt"
	"log"
)

// ExampleValidatorUsage demonstrates basic usage of the validation framework.
func ExampleValidatorUsage() {
	// Create a new validator instance
	validator := NewValidator()

	// Example 1: Validating a JSON-RPC Request
	fmt.Println("=== Example 1: Validating JSON-RPC Request ===")
	
	validRequest := &Request{
		JSONRPCVersion: "2.0",
		Method:         "user.login",
		ID:             "req-123",
	}
	
	err := validator.ValidateRequest(validRequest)
	if err != nil {
		log.Printf("Valid request failed validation: %v", err)
	} else {
		fmt.Println("✓ Valid request passed validation")
	}
	
	// Example with invalid request
	invalidRequest := &Request{
		JSONRPCVersion: "1.0", // Invalid version
		Method:         "",    // Empty method
		ID:             "req-123",
	}
	
	err = validator.ValidateRequest(invalidRequest)
	if err != nil {
		fmt.Printf("✗ Invalid request failed validation (as expected): %v\n", err)
		
		// Access detailed error information
		if validationErrs, ok := err.(ValidationErrors); ok {
			for _, validationErr := range validationErrs {
				fmt.Printf("  Field: %s, Tag: %s, Message: %s\n", 
					validationErr.Field, validationErr.Tag, validationErr.Message)
			}
		}
	}

	// Example 2: Validating Session Codes
	fmt.Println("\n=== Example 2: Session Code Validation ===")
	
	validCodes := []string{"happy-panda-42", "blue-river-7", "QUICK-FOX-99"}
	for _, code := range validCodes {
		err := validator.ValidateSessionCode(code)
		if err == nil {
			fmt.Printf("✓ Session code '%s' is valid\n", code)
		} else {
			fmt.Printf("✗ Session code '%s' is invalid: %v\n", code, err)
		}
	}
	
	invalidCodes := []string{"invalid-code", "happy-panda", "happy-panda-100"}
	for _, code := range invalidCodes {
		err := validator.ValidateSessionCode(code)
		if err != nil {
			fmt.Printf("✗ Session code '%s' is invalid (as expected): %v\n", code, err)
		}
	}

	// Example 3: Custom Struct Validation
	fmt.Println("\n=== Example 3: Custom Struct Validation ===")
	
	// Define a custom struct with validation tags
	type LoginRequest struct {
		Username    string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Password    string `json:"password" validate:"required,min=8"`
		SessionCode string `json:"session_code" validate:"required,sessioncode"`
		Email       string `json:"email" validate:"omitempty,email"`
		Age         int    `json:"age" validate:"min=13,max=120"`
	}
	
	// Valid login request
	validLogin := &LoginRequest{
		Username:    "john_doe",
		Password:    "securePassword123",
		SessionCode: "happy-panda-42",
		Email:       "john@example.com",
		Age:         25,
	}
	
	err = validator.Validate(validLogin)
	if err == nil {
		fmt.Println("✓ Valid login request passed validation")
	} else {
		fmt.Printf("✗ Valid login request failed: %v\n", err)
	}
	
	// Invalid login request
	invalidLogin := &LoginRequest{
		Username:    "jo", // Too short
		Password:    "weak", // Too short
		SessionCode: "invalid-code", // Invalid format
		Email:       "not-an-email", // Invalid email
		Age:         12, // Too young
	}
	
	err = validator.Validate(invalidLogin)
	if err != nil {
		fmt.Printf("✗ Invalid login request failed validation (as expected): %v\n", err)
		
		// With fast-fail, only the first error is returned
		if validationErrs, ok := err.(ValidationErrors); ok {
			fmt.Printf("  First validation error - Field: %s, Message: %s\n", 
				validationErrs[0].Field, validationErrs[0].Message)
		}
	}

	// Example 4: Variable Validation
	fmt.Println("\n=== Example 4: Variable Validation ===")
	
	// Validate individual variables
	testEmail := "user@example.com"
	err = validator.ValidateVar(testEmail, "required,email")
	if err == nil {
		fmt.Printf("✓ Email '%s' is valid\n", testEmail)
	}
	
	testAge := 25
	err = validator.ValidateVar(testAge, "min=18,max=65")
	if err == nil {
		fmt.Printf("✓ Age %d is valid\n", testAge)
	}
	
	// Example 5: JSON-RPC Error Validation
	fmt.Println("\n=== Example 5: JSON-RPC Error Validation ===")
	
	validError := &Error{
		Code:    -32601,
		Message: "Method not found",
		Data:    "The requested method 'unknown.method' was not found",
	}
	
	err = validator.ValidateError(validError)
	if err == nil {
		fmt.Println("✓ JSON-RPC error is valid")
	}
	
	invalidError := &Error{
		Code:    -32601,
		Message: "", // Empty message
	}
	
	err = validator.ValidateError(invalidError)
	if err != nil {
		fmt.Printf("✗ Invalid JSON-RPC error failed validation (as expected): %v\n", err)
	}

	// Example 6: Response Validation
	fmt.Println("\n=== Example 6: JSON-RPC Response Validation ===")
	
	successResponse := &Response{
		JSONRPCVersion: "2.0",
		Result:         map[string]interface{}{"status": "success", "user_id": 123},
		ID:             "req-123",
	}
	
	err = validator.ValidateResponse(successResponse)
	if err == nil {
		fmt.Println("✓ Success response is valid")
	}
	
	errorResponse := &Response{
		JSONRPCVersion: "2.0",
		Error:          validError,
		ID:             "req-123",
	}
	
	err = validator.ValidateResponse(errorResponse)
	if err == nil {
		fmt.Println("✓ Error response is valid")
	}
	
	fmt.Println("\n=== Validation Examples Complete ===")
}