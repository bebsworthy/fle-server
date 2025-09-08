# Research: backend-core

## Overview
Research findings for implementing the core backend infrastructure for FLE, focusing on basic Go server setup with WebSocket support, JSON-RPC protocol handling, and session management.

## Architecture Context
**Note**: No global architecture.md exists yet. This is a greenfield project with specifications defined in README.md and supporting documentation.

## Feature Scope Analysis

### Core Responsibilities
The backend-core must provide:
1. Basic HTTP server with configurable port
2. WebSocket endpoint with connection upgrade
3. JSON-RPC 2.0 protocol handling
4. Human-friendly session code generation
5. Connection management using hub pattern
6. Structured logging and error handling

## Technical Stack Research

### Gorilla WebSocket Implementation
- **Pattern**: Use `github.com/gorilla/websocket` for WebSocket handling
- **Upgrader**: Configure with CheckOrigin for CORS support during development
- **Connection Management**: Hub pattern for broadcasting and connection tracking
- **Consideration**: Gorilla is in maintenance mode, but stable and widely used

### JSON-RPC Options
- **Option 1**: Implement minimal custom JSON-RPC 2.0 handler (recommended for simplicity)
- **Option 2**: Use `github.com/sourcegraph/jsonrpc2` if more features needed later
- **Message Structure**: Standard request/response with id, method, params, result/error

### Message Validation Strategy
- **Validation Library Options**:
  - `github.com/go-playground/validator/v10` - Struct tag-based validation
  - Custom validation with strict type checking
  - JSON Schema validation with `github.com/xeipuuv/gojsonschema`
- **Approach**: Validate at message boundary before any processing
- **Fast Fail**: Return error immediately on first validation failure
- **Schema Definition**: Define strict schemas for each message type

### Session Code Generation

#### Available Libraries

1. **golang-petname** (`github.com/dustinkirkland/golang-petname`)
   - Most mature and widely used option
   - Generates combinations of adverb-adjective-animal
   - 2 words = adjective-animal (e.g., "wiggly-yellowtail")
   - Customizable separator
   - Over 2.8 trillion unique combinations with 3 words
   - **Note**: Doesn't include numbers by default

2. **namegenerator** (`github.com/0x6flab/namegenerator`)
   - Supports adding suffixes including numbers
   - Based on real names from US Census
   - Can generate formats like "John-Smith-123"

3. **Custom Implementation**
   - Combine golang-petname with random number suffix
   - Example: `fmt.Sprintf("%s-%d", petname.Generate(2, "-"), rand.Intn(100))`
   - Results in: "happy-panda-42", "swift-eagle-7"

#### Recommended Approach
```go
// Using golang-petname with custom number suffix
import "github.com/dustinkirkland/golang-petname"

func generateSessionCode() string {
    // Generate adjective-animal
    base := petname.Generate(2, "-")
    // Add random number 1-99
    number := rand.Intn(99) + 1
    return fmt.Sprintf("%s-%d", base, number)
}
```

- **Format**: adjective-animal-number (e.g., "happy-panda-42")
- **Uniqueness**: Check for collisions in active sessions map, regenerate if needed
- **Case Handling**: Store and compare as lowercase for user convenience
- **Collision Space**: With ~1300 adjectives × ~870 animals × 99 numbers ≈ 111 million combinations

## Implementation Patterns

### Directory Structure Recommendation
```
server/
├── cmd/
│   └── server/          # Main entry point
├── internal/
│   ├── websocket/      # WebSocket hub and connection management
│   ├── jsonrpc/        # JSON-RPC message handling
│   ├── session/        # Session code generation and management
│   └── models/         # Shared data structures
├── config/             # Configuration files
├── .env.example        # Example environment variables
├── Makefile           # Build and run commands
└── go.mod
```

### Hub Pattern for WebSocket
```
Hub responsibilities:
- Register new connections
- Unregister disconnecting clients
- Broadcast messages to connections
- Track connections by session code
- Handle connection-specific messaging
```

### Key Dependencies to Consider
- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/gorilla/mux` - HTTP routing (or just use stdlib)
- `github.com/joho/godotenv` - .env file support
- `github.com/rs/zerolog` or `log/slog` - Structured logging
- `github.com/go-playground/validator/v10` - Struct validation (recommended for fast fail)
- `github.com/dustinkirkland/golang-petname` - Human-friendly session code generation

## Implementation Details

### Session Code Generation Strategy (Updated)
Using golang-petname library with custom number suffix:
- **Library**: `github.com/dustinkirkland/golang-petname` for base generation
- **Format**: adjective-animal-number
- **Number Range**: 1-99 for simplicity
- **Examples**: 
  - brave-lion-7
  - gentle-moon-42
  - swift-eagle-13
- **Collision Handling**: Track active sessions in memory map
- **Total Combinations**: ~111 million (sufficient for our use case)

### WebSocket Connection Lifecycle
1. Client connects to `/ws` endpoint
2. Server upgrades HTTP to WebSocket
3. Server generates or validates session code
4. Connection registered in hub with session
5. Bidirectional JSON-RPC communication
6. On disconnect, cleanup and unregister

### JSON-RPC Message Flow with Validation
1. Receive raw message from WebSocket
2. Parse JSON (fail fast on invalid JSON)
3. Validate JSON-RPC 2.0 structure (id, jsonrpc, method, params)
4. Validate method exists and params match expected schema
5. Route to handler only if validation passes
6. Process request (placeholder for now)
7. Validate response structure before sending
8. Send validated response through WebSocket

### Validation Points
- **Input Validation**:
  - Valid JSON format
  - JSON-RPC 2.0 compliance (jsonrpc: "2.0" required)
  - Method name exists in registry
  - **Payload/Params validation**:
    - Type checking (string, number, bool, object, array)
    - Required vs optional fields
    - String format validation (e.g., session codes match pattern)
    - Numeric constraints (min, max, ranges)
    - Array constraints (min/max length, item types)
    - Nested object validation
    - No additional properties (strict mode)
  
- **Output Validation**:
  - Response matches JSON-RPC 2.0 spec
  - Result or error, never both
  - Error codes follow specification
  - **Result payload validation**:
    - Conforms to method's response schema
    - All required fields present
    - Correct data types
    - No extra fields

### Validation Implementation Strategy
```go
// Example: Framework for registering methods with validation
type MethodHandler struct {
    Name         string
    ParamsType   reflect.Type  // Type with validation tags
    ResultType   reflect.Type  // Type with validation tags
    Handler      HandlerFunc
}

// Example validation tags (methods defined in future features)
type SomeMethodParams struct {
    Field1 string `json:"field1" validate:"required,min=3"`
    Field2 int    `json:"field2" validate:"required,min=0,max=100"`
}
```

- Provide a registration system for methods
- Each method defines its params and result types with validation tags
- Use struct tags for declarative validation
- Custom validators for domain-specific rules
- Validate at unmarshaling time for fast fail
- Return detailed validation errors with field paths

## Development Priorities

### Phase 1: Basic Infrastructure (Current)
1. HTTP server setup with health endpoint
2. WebSocket endpoint with upgrade
3. Session code generation
4. Basic hub implementation
5. JSON-RPC message parsing
6. Structured logging

### Future Phases (Out of Scope)
- Content loading and caching
- File-based storage implementation
- Activity generation
- LLM integration
- Progress tracking

## Configuration Approach

### Environment Variables
```
PORT=8080
HOST=0.0.0.0
CORS_ORIGIN=http://localhost:3000
LOG_LEVEL=debug
ENV=development
```

### Development vs Production
- Development: Human-readable logs, CORS enabled, debug level
- Production: JSON logs, CORS restricted, info level

## Testing Strategy

### Unit Tests
- Session code generation and uniqueness
- JSON-RPC message parsing and validation
- Hub connection management

### Integration Tests
- WebSocket connection flow
- Session persistence across reconnect
- Concurrent connection handling

## Security Considerations

### For This Phase
- CORS configuration for development
- No sensitive data in session codes
- Input validation for JSON-RPC messages
- Panic recovery in connection handlers

### Future Considerations (Not in Scope)
- Rate limiting
- Authentication upgrades
- HTTPS/WSS in production

## Error Handling Patterns

### Connection Errors
- Log with session context
- Clean disconnect handling
- Graceful degradation

### JSON-RPC Errors
- Standard error codes (-32700 to -32603)
- Custom application errors (starting at -32000)
- Descriptive messages without exposing internals

## Next Steps
With this focused research complete, we can proceed with implementing just the core infrastructure components that will serve as the foundation for all future features.