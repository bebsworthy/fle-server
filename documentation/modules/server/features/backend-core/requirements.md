# Requirements: backend-core

## Introduction

The backend-core feature provides the foundational Go server infrastructure for the FLE platform. This includes basic HTTP server setup, WebSocket connection handling with JSON-RPC protocol support, and session management. This is strictly the core infrastructure layer - content management, activity generation, and LLM integration will be handled in separate features.

### Research Context

Based on research findings, the implementation will use Gorilla WebSocket with a hub pattern for connection management and human-friendly session codes that are easy to share and remember.

## Functional Requirements

### Requirement 1: WebSocket Server Infrastructure

**User Story:** As a frontend application, I want to connect to the backend via WebSocket, so that I can maintain a persistent bidirectional communication channel.

#### Acceptance Criteria
1. WHEN a client connects to the WebSocket endpoint THEN the server SHALL accept the connection and establish a persistent channel
2. WHEN a client sends a JSON-RPC 2.0 request THEN the server SHALL parse and validate both the request format and payload
3. WHEN a valid JSON-RPC request is received THEN the server SHALL route it to the appropriate handler and return a JSON-RPC response
4. WHEN an invalid JSON-RPC request is received THEN the server SHALL return a JSON-RPC error response with appropriate error code
5. WHEN a client disconnects THEN the server SHALL clean up resources and close the connection gracefully
6. IF multiple clients connect simultaneously THEN the server SHALL manage each connection independently using a hub pattern
7. The server SHALL provide a framework for registering methods with their validation schemas (actual methods defined in future features)

### Requirement 2: Session Management

**User Story:** As a new visitor, I want to receive a memorable session code automatically, so that I can easily share or remember my learning session.

#### Acceptance Criteria
1. WHEN a new client connects without a session code THEN the server SHALL generate a human-friendly session code (e.g., "happy-panda-42" or "blue-river-7")
2. WHEN generating a session code THEN the server SHALL use a format of adjective-noun-number for readability
3. WHEN a client connects with an existing session code THEN the server SHALL validate and restore their session
4. IF a client provides an invalid session code THEN the server SHALL generate a new session
5. WHEN a session is created THEN the server SHALL return the session code to the client for URL and localStorage storage
6. Session codes SHALL be case-insensitive for user convenience

### Requirement 3: HTTP Server Setup

**User Story:** As a developer, I want a properly configured HTTP server, so that I can serve static files and handle WebSocket upgrades.

#### Acceptance Criteria
1. WHEN the server starts THEN it SHALL listen on a configurable port (default 8080)
2. WHEN a request arrives at /ws THEN the server SHALL upgrade to WebSocket connection
3. WHEN a request arrives at /health THEN the server SHALL return a health check response
4. WHEN serving HTTP THEN the server SHALL support CORS for frontend development
5. WHEN the server starts THEN it SHALL log the listening address and port

### Requirement 4: Connection Management

**User Story:** As a system, I want to manage multiple WebSocket connections efficiently, so that many users can connect simultaneously.

#### Acceptance Criteria
1. WHEN multiple clients connect THEN the server SHALL manage each connection in a separate goroutine
2. WHEN using the hub pattern THEN the server SHALL broadcast messages efficiently to relevant connections
3. WHEN a connection is established THEN the server SHALL track it with its session code
4. WHEN a client disconnects THEN the server SHALL clean up resources and remove from active connections
5. The server SHALL implement ping/pong heartbeat to detect stale connections

### Requirement 5: Basic Project Structure

**User Story:** As a developer, I want a well-organized Go project structure, so that the codebase is maintainable and extensible.

#### Acceptance Criteria
1. WHEN the project is created THEN it SHALL follow standard Go project layout
2. WHEN organizing code THEN the server SHALL use internal packages for non-exported functionality
3. WHEN defining types THEN the server SHALL create a models package for shared data structures
4. The project SHALL include a Makefile for common operations (build, run, test)
5. The project SHALL use Go modules for dependency management

### Requirement 6: Message and Payload Validation Framework

**User Story:** As a system operator, I want a validation framework for all WebSocket messages and their payloads, so that invalid data is rejected immediately and the system remains stable.

#### Acceptance Criteria
1. WHEN a message is received THEN the server SHALL validate it conforms to JSON-RPC 2.0 specification before processing
2. The server SHALL provide a framework where methods can register their expected schemas for params and results
3. WHEN a registered method is called THEN the server SHALL validate the params payload against the method's schema
4. The validation framework SHALL support:
   - Required fields and optional fields
   - Data types (string, number, boolean, object, array)
   - String format patterns
   - Numeric ranges where applicable
   - Array length constraints
   - Nested object validation
5. WHEN validation fails THEN the server SHALL immediately return a JSON-RPC error (-32602 Invalid params) with details about what failed
6. WHEN sending responses THEN the server SHALL validate the result payload conforms to the method's response schema
7. IF a payload contains unexpected fields THEN the server SHALL reject it (no additional properties allowed)
8. WHEN invalid JSON is received THEN the server SHALL return parse error (-32700) without attempting to process
9. The server SHALL use a validation library (e.g., go-playground/validator) to ensure consistent validation rules

### Requirement 7: Error Handling

**User Story:** As a developer, I want proper error handling, so that issues can be debugged easily.

#### Acceptance Criteria
1. WHEN any error occurs THEN the server SHALL log it with appropriate context
2. WHEN a panic occurs in a connection handler THEN the server SHALL recover and continue serving other connections
3. WHEN JSON-RPC errors occur THEN the server SHALL return properly formatted error responses
4. Error messages SHALL be descriptive but not expose sensitive information
5. WHEN validation fails THEN errors SHALL indicate which field or constraint was violated

## Non-Functional Requirements

### Requirement 8: Configuration

**User Story:** As a developer, I want environment-based configuration, so that I can deploy to different environments.

#### Acceptance Criteria
1. WHEN the server starts THEN it SHALL read configuration from environment variables
2. Configuration SHALL include port, host, CORS settings, and log level
3. WHEN configuration is missing THEN the server SHALL use sensible defaults
4. The server SHALL support a .env file for local development

### Requirement 9: Logging

**User Story:** As a developer, I want structured logging, so that I can monitor and debug the application.

#### Acceptance Criteria
1. WHEN the server runs THEN it SHALL use structured logging (JSON format in production)
2. WHEN in development THEN the server SHALL use human-readable log format
3. Log levels SHALL be configurable (debug, info, warn, error)
4. WHEN connections are established/closed THEN the server SHALL log with session codes
5. The server SHALL include request IDs for tracing related operations

## Summary

The backend-core feature focuses solely on establishing the foundational Go server infrastructure:
- Basic HTTP server with WebSocket support
- Human-friendly session codes (e.g., "happy-panda-42")
- JSON-RPC message handling with strict validation
- Fast-fail message validation on all inputs and outputs
- Connection management with hub pattern
- Proper project structure and configuration
- Structured logging and error handling

This provides the core infrastructure upon which future features (content management, activity generation, LLM integration, etc.) will be built.

---

Do the requirements look good? If so, we can move on to the design.