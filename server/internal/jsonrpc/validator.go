// Package jsonrpc provides JSON-RPC 2.0 message types and validation framework.
package jsonrpc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator provides comprehensive validation functionality for JSON-RPC types
// and custom data structures with fast-fail behavior and detailed error reporting.
type Validator struct {
	// validate is the underlying go-playground validator instance
	validate *validator.Validate
}

// ValidationError represents a detailed validation error with field information.
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`

	// Tag is the validation rule that failed (e.g., "required", "min", "max")
	Tag string `json:"tag"`

	// Value is the actual value that failed validation
	Value interface{} `json:"value"`

	// Param is the parameter for the validation rule (e.g., "5" for "min=5")
	Param string `json:"param"`

	// Message is a human-readable error message
	Message string `json:"message"`
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return e.Message
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors.
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "validation errors"
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Message)
	}

	return fmt.Sprintf("validation failed: %s", strings.Join(messages, "; "))
}

// NewValidator creates a new validator instance with custom validators registered.
func NewValidator() *Validator {
	validate := validator.New()

	// Register custom tag name function to use json field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Create validator instance
	v := &Validator{
		validate: validate,
	}

	// Register custom validators
	v.registerCustomValidators()

	return v
}

// registerCustomValidators registers all custom validation functions.
func (v *Validator) registerCustomValidators() {
	// Register session code validator
	v.validate.RegisterValidation("sessioncode", v.validateSessionCode)

	// Register JSON-RPC version validator
	v.validate.RegisterValidation("jsonrpcversion", v.validateJSONRPCVersion)
}

// validateSessionCode validates that a string follows the session code format:
// "adjective-noun-number" where number is 1-99.
// This validator is case-insensitive.
func (v *Validator) validateSessionCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if code == "" {
		return false
	}

	// Convert to lowercase for case-insensitive validation
	normalized := strings.ToLower(strings.TrimSpace(code))

	// Split by dashes
	parts := strings.Split(normalized, "-")

	// Must have exactly 3 parts: adjective-noun-number
	if len(parts) != 3 {
		return false
	}

	// Check that each part is not empty
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return false
		}
	}

	// Check that the last part is a valid number (1-99)
	lastPart := parts[2]
	if len(lastPart) == 0 || len(lastPart) > 2 {
		return false
	}

	// Check if it's a valid number in range 1-99
	var number int
	n, err := fmt.Sscanf(lastPart, "%d", &number)
	if n != 1 || err != nil {
		return false
	}

	return number >= 1 && number <= 99
}

// validateJSONRPCVersion validates that a string is exactly "2.0".
func (v *Validator) validateJSONRPCVersion(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	return version == Version
}

// Validate validates any struct using the configured validation rules.
// It returns detailed validation errors on failure, with fast-fail behavior.
// This means validation stops at the first field that fails validation.
func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	// Convert validator errors to our custom format
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var errors ValidationErrors
		
		for _, validationErr := range validationErrs {
			customErr := ValidationError{
				Field: validationErr.Field(),
				Tag:   validationErr.Tag(),
				Value: validationErr.Value(),
				Param: validationErr.Param(),
				Message: v.buildErrorMessage(validationErr),
			}
			errors = append(errors, customErr)
			
			// Fast-fail: return after first error
			break
		}

		return errors
	}

	// Return original error if it's not a validation error
	return err
}

// ValidateVar validates a single variable using the specified validation rules.
// This is useful for validating individual values outside of struct contexts.
func (v *Validator) ValidateVar(field interface{}, rules string) error {
	err := v.validate.Var(field, rules)
	if err == nil {
		return nil
	}

	// Convert validator errors to our custom format
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var errors ValidationErrors
		
		for _, validationErr := range validationErrs {
			customErr := ValidationError{
				Field: "value", // Generic field name for variable validation
				Tag:   validationErr.Tag(),
				Value: validationErr.Value(),
				Param: validationErr.Param(),
				Message: v.buildErrorMessage(validationErr),
			}
			errors = append(errors, customErr)
			
			// Fast-fail: return after first error
			break
		}

		return errors
	}

	// Return original error if it's not a validation error
	return err
}

// ValidateRequest validates a JSON-RPC 2.0 request message.
func (v *Validator) ValidateRequest(req *Request) error {
	if req == nil {
		return ValidationErrors{
			{
				Field:   "request",
				Tag:     "required",
				Value:   nil,
				Message: "request cannot be nil",
			},
		}
	}
	return v.Validate(req)
}

// ValidateResponse validates a JSON-RPC 2.0 response message.
func (v *Validator) ValidateResponse(resp *Response) error {
	if resp == nil {
		return ValidationErrors{
			{
				Field:   "response",
				Tag:     "required",
				Value:   nil,
				Message: "response cannot be nil",
			},
		}
	}
	return v.Validate(resp)
}

// ValidateError validates a JSON-RPC 2.0 error object.
func (v *Validator) ValidateError(err *Error) error {
	if err == nil {
		return ValidationErrors{
			{
				Field:   "error",
				Tag:     "required",
				Value:   nil,
				Message: "error cannot be nil",
			},
		}
	}
	return v.Validate(err)
}

// ValidateSessionCode validates a session code using the custom session code format.
func (v *Validator) ValidateSessionCode(code string) error {
	return v.ValidateVar(code, "required,sessioncode")
}

// buildErrorMessage creates a human-readable error message from a validation error.
func (v *Validator) buildErrorMessage(validationErr validator.FieldError) string {
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()
	value := validationErr.Value()

	switch tag {
	case "required":
		return fmt.Sprintf("field '%s' is required", field)
	case "eq":
		return fmt.Sprintf("field '%s' must equal '%s', got '%v'", field, param, value)
	case "min":
		if validationErr.Kind() == reflect.String {
			return fmt.Sprintf("field '%s' must be at least %s characters long, got %d", field, param, len(fmt.Sprintf("%v", value)))
		}
		return fmt.Sprintf("field '%s' must be at least %s, got '%v'", field, param, value)
	case "max":
		if validationErr.Kind() == reflect.String {
			return fmt.Sprintf("field '%s' must be at most %s characters long, got %d", field, param, len(fmt.Sprintf("%v", value)))
		}
		return fmt.Sprintf("field '%s' must be at most %s, got '%v'", field, param, value)
	case "len":
		return fmt.Sprintf("field '%s' must be exactly %s in length, got %d", field, param, len(fmt.Sprintf("%v", value)))
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of [%s], got '%v'", field, param, value)
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email address, got '%v'", field, value)
	case "url":
		return fmt.Sprintf("field '%s' must be a valid URL, got '%v'", field, value)
	case "numeric":
		return fmt.Sprintf("field '%s' must be numeric, got '%v'", field, value)
	case "alpha":
		return fmt.Sprintf("field '%s' must contain only alphabetic characters, got '%v'", field, value)
	case "alphanum":
		return fmt.Sprintf("field '%s' must contain only alphanumeric characters, got '%v'", field, value)
	case "sessioncode":
		return fmt.Sprintf("field '%s' must be a valid session code in format 'adjective-noun-number' (e.g., 'happy-panda-42'), got '%v'", field, value)
	case "jsonrpcversion":
		return fmt.Sprintf("field '%s' must be exactly '2.0' for JSON-RPC 2.0 compliance, got '%v'", field, value)
	case "gt":
		return fmt.Sprintf("field '%s' must be greater than %s, got '%v'", field, param, value)
	case "gte":
		return fmt.Sprintf("field '%s' must be greater than or equal to %s, got '%v'", field, param, value)
	case "lt":
		return fmt.Sprintf("field '%s' must be less than %s, got '%v'", field, param, value)
	case "lte":
		return fmt.Sprintf("field '%s' must be less than or equal to %s, got '%v'", field, param, value)
	case "dive":
		return fmt.Sprintf("field '%s' contains invalid nested values", field)
	default:
		return fmt.Sprintf("field '%s' failed validation rule '%s' with value '%v'", field, tag, value)
	}
}

// GetSupportedTags returns a list of all supported validation tags.
func (v *Validator) GetSupportedTags() []string {
	return []string{
		// Standard validation tags
		"required", "omitempty", "eq", "ne", "min", "max", "len",
		"gt", "gte", "lt", "lte", "oneof", "email", "url",
		"numeric", "alpha", "alphanum", "ascii", "printascii",
		"dive", "keys", "endkeys", "uuid", "uuid3", "uuid4", "uuid5",
		
		// Custom validation tags
		"sessioncode", "jsonrpcversion",
	}
}

// Example usage and validation patterns:
//
// Basic struct validation:
//   validator := NewValidator()
//   err := validator.Validate(myStruct)
//   if err != nil {
//       if validationErrs, ok := err.(ValidationErrors); ok {
//           for _, err := range validationErrs {
//               fmt.Printf("Field: %s, Error: %s\n", err.Field, err.Message)
//           }
//       }
//   }
//
// Variable validation:
//   err := validator.ValidateVar("invalid-code", "required,sessioncode")
//   if err != nil {
//       fmt.Println("Session code validation failed:", err)
//   }
//
// JSON-RPC specific validation:
//   request := &Request{...}
//   err := validator.ValidateRequest(request)
//   if err != nil {
//       fmt.Println("Request validation failed:", err)
//   }