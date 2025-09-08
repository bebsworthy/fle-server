// Package jsonrpc provides JSON-RPC 2.0 message types and error codes.
//
// This package implements JSON-RPC 2.0 specification as defined at:
// https://www.jsonrpc.org/specification
//
// It provides strongly-typed message structures for request, response,
// and error handling with comprehensive validation support.
package jsonrpc

import (
	"encoding/json"
	"fmt"
)

const (
	// Version defines the JSON-RPC 2.0 version string
	Version = "2.0"
)

// Request represents a JSON-RPC 2.0 request message.
//
// According to the specification, a request object has the following members:
// - jsonrpc: A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
// - method: A String containing the name of the method to be invoked.
// - params: A Structured value that holds the parameter values to be used during the invocation of the method. This member MAY be omitted.
// - id: An identifier established by the Client. MUST contain a String, Number, or NULL value if included.
type Request struct {
	// JSONRPCVersion specifies the version of the JSON-RPC protocol.
	// Must be exactly "2.0" for JSON-RPC 2.0 compliance.
	JSONRPCVersion string `json:"jsonrpc" validate:"required,eq=2.0"`

	// Method contains the name of the method to be invoked.
	// Method names that begin with "rpc." are reserved for rpc-internal methods.
	Method string `json:"method" validate:"required,min=1"`

	// Params holds the parameter values to be used during method invocation.
	// This can be an object, array, or null. It may be omitted entirely.
	Params json.RawMessage `json:"params,omitempty"`

	// ID is an identifier established by the client.
	// It can be a string, number, or null. If omitted, the request is a notification.
	ID interface{} `json:"id,omitempty"`
}

// IsNotification returns true if this request is a notification
// (has no ID and expects no response).
func (r *Request) IsNotification() bool {
	return r.ID == nil
}

// Response represents a JSON-RPC 2.0 response message.
//
// According to the specification, when a rpc call is made, the Server MUST reply
// with a Response, except for in the case of Notifications.
// The Response is expressed as a single JSON Object, with the following members:
// - jsonrpc: A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
// - result: Required on success. This member MUST NOT exist if there was an error invoking the method.
// - error: Required on error. This member MUST NOT exist if there was no error triggered during invocation.
// - id: This will be the same as the value of the id member in the Request Object.
type Response struct {
	// JSONRPCVersion specifies the version of the JSON-RPC protocol.
	// Must be exactly "2.0" for JSON-RPC 2.0 compliance.
	JSONRPCVersion string `json:"jsonrpc" validate:"required,eq=2.0"`

	// Result contains the result of the method invocation.
	// This member is REQUIRED on success and MUST NOT exist if there was an error.
	Result interface{} `json:"result,omitempty"`

	// Error contains the error object if an error occurred during invocation.
	// This member is REQUIRED on error and MUST NOT exist if there was no error.
	Error *Error `json:"error,omitempty"`

	// ID is the same as the value of the id member in the Request Object.
	// If there was an error in detecting the id in the Request object, it MUST be Null.
	ID interface{} `json:"id"`
}

// IsError returns true if this response contains an error.
func (r *Response) IsError() bool {
	return r.Error != nil
}

// IsSuccess returns true if this response contains a successful result.
func (r *Response) IsSuccess() bool {
	return r.Error == nil
}

// Error represents a JSON-RPC 2.0 error object.
//
// According to the specification, when a rpc call encounters an error,
// the Response Object MUST contain the error member with a value that is an Object
// with the following members:
// - code: A Number that indicates the error type that occurred.
// - message: A String providing a short description of the error.
// - data: A Primitive or Structured value that contains additional information about the error.
type Error struct {
	// Code indicates the error type that occurred.
	// This MUST be an integer. Error codes from and including -32768 to -32000 are reserved for pre-defined errors.
	Code int `json:"code" validate:"required"`

	// Message provides a short description of the error.
	// The message SHOULD be a concise phrase that describes the error.
	Message string `json:"message" validate:"required,min=1"`

	// Data contains additional information about the error.
	// This may be omitted. The value can be a Primitive or Structured value.
	Data interface{} `json:"data,omitempty"`
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("JSON-RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// Standard JSON-RPC 2.0 error codes as defined in the specification.
const (
	// ParseError indicates invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	ParseError = -32700

	// InvalidRequest indicates the JSON sent is not a valid Request object.
	InvalidRequest = -32600

	// MethodNotFound indicates the method does not exist / is not available.
	MethodNotFound = -32601

	// InvalidParams indicates invalid method parameter(s).
	InvalidParams = -32602

	// InternalError indicates an internal JSON-RPC error.
	InternalError = -32603

	// ServerErrorStart is the start of the range for implementation-defined server-errors.
	// Error codes from -32099 to -32000 are reserved for implementation-defined server-errors.
	ServerErrorStart = -32099

	// ServerErrorEnd is the end of the range for implementation-defined server-errors.
	// Error codes from -32099 to -32000 are reserved for implementation-defined server-errors.
	ServerErrorEnd = -32000
)

// Standard error messages for predefined error codes.
var (
	// ErrParse represents a parse error (-32700).
	ErrParse = &Error{
		Code:    ParseError,
		Message: "Parse error",
	}

	// ErrInvalidRequest represents an invalid request error (-32600).
	ErrInvalidRequest = &Error{
		Code:    InvalidRequest,
		Message: "Invalid Request",
	}

	// ErrMethodNotFound represents a method not found error (-32601).
	ErrMethodNotFound = &Error{
		Code:    MethodNotFound,
		Message: "Method not found",
	}

	// ErrInvalidParams represents an invalid params error (-32602).
	ErrInvalidParams = &Error{
		Code:    InvalidParams,
		Message: "Invalid params",
	}

	// ErrInternal represents an internal error (-32603).
	ErrInternal = &Error{
		Code:    InternalError,
		Message: "Internal error",
	}
)

// NewError creates a new JSON-RPC error with the given code and message.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithData creates a new JSON-RPC error with the given code, message, and data.
func NewErrorWithData(code int, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewRequest creates a new JSON-RPC 2.0 request.
func NewRequest(method string, params interface{}, id interface{}) (*Request, error) {
	var paramsBytes json.RawMessage
	if params != nil {
		bytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}
		paramsBytes = bytes
	}

	return &Request{
		JSONRPCVersion: Version,
		Method:         method,
		Params:         paramsBytes,
		ID:             id,
	}, nil
}

// NewNotification creates a new JSON-RPC 2.0 notification (request without ID).
func NewNotification(method string, params interface{}) (*Request, error) {
	return NewRequest(method, params, nil)
}

// NewResponse creates a new JSON-RPC 2.0 success response.
func NewResponse(result interface{}, id interface{}) *Response {
	return &Response{
		JSONRPCVersion: Version,
		Result:         result,
		ID:             id,
	}
}

// NewErrorResponse creates a new JSON-RPC 2.0 error response.
func NewErrorResponse(err *Error, id interface{}) *Response {
	return &Response{
		JSONRPCVersion: Version,
		Error:          err,
		ID:             id,
	}
}

// IsValidErrorCode returns true if the given error code is valid according to JSON-RPC 2.0 specification.
// Valid error codes are:
// - Pre-defined errors: -32768 to -32000 (reserved)
// - Implementation-defined server errors: -32099 to -32000 (reserved for server)
// - Application-defined errors: any other integer
func IsValidErrorCode(code int) bool {
	// All integer values are valid error codes in JSON-RPC 2.0
	// The specification only reserves certain ranges but doesn't invalidate others
	return true
}

// IsReservedErrorCode returns true if the given error code is reserved by the JSON-RPC 2.0 specification.
func IsReservedErrorCode(code int) bool {
	return code >= -32768 && code <= -32000
}

// IsServerErrorCode returns true if the given error code is in the server error range.
func IsServerErrorCode(code int) bool {
	return code >= ServerErrorStart && code <= ServerErrorEnd
}
