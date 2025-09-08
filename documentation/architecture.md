# Architecture Documentation

## System Overview

FLE (Flexible Learning Environment) is currently a **single-module architecture** consisting of a Go-based backend server that provides WebSocket-based real-time communication with JSON-RPC 2.0 protocol support. The system is designed to support a language learning platform with extensible architecture for future frontend modules.

### Architecture Style
- **Current State**: Monolithic backend server
- **Communication Pattern**: WebSocket with JSON-RPC 2.0
- **Deployment Model**: Single deployable unit
- **Data Storage**: In-memory (session management)

## Modules and Services

### 1. Server Module (`/server`)

#### Technology Stack
- **Language**: Go 1.24.5
- **Primary Framework**: Gorilla WebSocket (v1.5.3)
- **Validation**: go-playground/validator (v10.27.0)
- **Session Codes**: golang-petname (for human-friendly identifiers)
- **Testing**: testify (v1.11.1)
- **Build Tool**: Make
- **Linting**: golangci-lint (extensive configuration)

#### Module Structure
```
server/
├── cmd/server/           # Application entry point
│   └── main.go          # Main server initialization
├── internal/            # Private packages
│   ├── config/         # Configuration management
│   ├── jsonrpc/        # JSON-RPC implementation
│   ├── logger/         # Structured logging (slog)
│   ├── server/         # HTTP server and handlers
│   ├── session/        # Session management
│   └── websocket/      # WebSocket hub and client
├── bin/                # Build output (gitignored)
├── code_review/        # Code review documentation
├── product_review/     # Product review documentation
├── Makefile           # Build automation (40+ targets)
├── go.mod             # Go module definition
└── .env.example       # Environment configuration template
```

#### Key Components

##### Configuration (`internal/config`)
- Environment-based configuration with sensible defaults
- Support for development/production modes
- Validation of all configuration values
- Thread-safe immutable configuration after loading

##### WebSocket Hub (`internal/websocket`)
- **Hub Pattern**: Central message routing for efficient broadcasting
- **Client Management**: Thread-safe client registration/unregistration
- **Connection Lifecycle**: Separate goroutines for read/write pumps
- **Heartbeat**: Ping/pong mechanism (54-second intervals)
- **Capacity**: Designed for 1000+ concurrent connections

##### JSON-RPC Router (`internal/jsonrpc`)
- Full JSON-RPC 2.0 specification compliance
- Method registration system with validation schemas
- Fast-fail validation with detailed error messages
- Thread-safe concurrent request handling
- Custom validators for session codes

##### Session Management (`internal/session`)
- Human-friendly session codes (e.g., "happy-panda-42")
- In-memory storage with thread-safe access (sync.RWMutex)
- Collision detection and retry logic
- Automatic cleanup of expired sessions
- Case-insensitive session code validation

##### HTTP Server (`internal/server`)
- Graceful shutdown with context cancellation
- CORS middleware for frontend development
- Health endpoint for monitoring
- WebSocket upgrade handling
- Structured request logging

##### Logging (`internal/logger`)
- Structured logging with slog
- JSON format in production, text in development
- Request ID generation for tracing
- Context-aware logging methods
- Global logger singleton pattern

#### API Endpoints

##### HTTP Endpoints
- `GET /health` - Health check endpoint returning JSON status
- `GET /ws` - WebSocket upgrade endpoint

##### JSON-RPC Methods (via WebSocket)
- `ping` - Server health check with timestamp
- `echo` - Echo back request parameters
- `getSessionInfo` - Get current session information

#### Testing Strategy
- **Unit Tests**: 94.8% coverage for session management
- **Integration Tests**: 77.5% coverage for WebSocket components
- **End-to-End Tests**: Complete connection flow testing
- **Race Detection**: All tests pass with `-race` flag
- **Benchmarks**: Performance testing for critical paths

## Data Architecture

### Current Implementation
- **Session Data**: In-memory map with thread-safe access
- **Connection State**: Hub-managed client registry
- **Message Queue**: Buffered channels (256 capacity) per client

### Data Flow
1. Client connects via WebSocket to `/ws` endpoint
2. Server creates/restores session with unique code
3. JSON-RPC requests flow through validation framework
4. Router dispatches to registered method handlers
5. Responses sent back through WebSocket connection

## Development and Deployment

### Build System
- **Make Targets**: 40+ targets for complete development workflow
  - `make build` - Build server binary
  - `make test` - Run all tests
  - `make test-coverage` - Generate coverage reports
  - `make lint` - Run comprehensive linting
  - `make run` - Run server locally
  - `make docker-build` - Build Docker image

### Configuration Management
- Environment variable based configuration
- `.env.example` template with documentation
- Sensible defaults for all settings
- Separate development/production configurations

### Development Workflow
1. Local development with hot reload support
2. Comprehensive linting with golangci-lint
3. Automated testing with coverage reporting
4. Code review process documented
5. Product review for requirements validation

## Security Architecture

### Current Implementation
- **CORS**: Configurable origin restrictions
- **Validation**: Comprehensive input validation on all JSON-RPC requests
- **Error Handling**: No sensitive information in error messages
- **Panic Recovery**: Graceful handling prevents cascade failures
- **Rate Limiting**: Connection limits (configurable)

### Authentication & Authorization
- **Session-based**: Automatic session code generation
- **Stateless**: No authentication currently implemented
- **Future**: Planned email registration and 2FA support

## Performance and Scalability

### Current Capabilities
- **Concurrent Connections**: 1000+ WebSocket connections
- **Message Processing**: Sub-millisecond JSON-RPC handling
- **Memory Usage**: Efficient hub pattern with minimal overhead
- **CPU Usage**: Separate goroutines prevent blocking

### Bottlenecks
- In-memory session storage (not persistent)
- Single server instance (no horizontal scaling)
- No caching layer
- No database (all state in memory)

## Monitoring and Observability

### Logging
- Structured JSON logs in production
- Human-readable text logs in development
- Request/response logging with timing
- Error tracking with stack traces
- Connection lifecycle events

### Health Checks
- `/health` endpoint for monitoring
- Version information in health response
- Environment information included

### Metrics
- Connection count tracking
- Session count monitoring
- Request/response timing (in logs)

## Technical Debt and TODOs

### Identified Issues
1. **Version Hardcoded**: Version "1.0.0" hardcoded in health endpoint (should be build-time injection)
2. **Context TODOs**: Multiple `context.TODO()` usages in logger package
3. **Message Size Limit**: 512-byte limit may be too restrictive for some JSON-RPC payloads
4. **No Persistence**: All data lost on server restart
5. **No Database**: No persistent storage implementation

### Future Improvements
- Database integration (PostgreSQL planned)
- Horizontal scaling support
- Message queue integration
- Caching layer
- Authentication system
- Frontend modules (React, Android planned)

## Architecture Decisions

### Design Patterns
1. **Hub Pattern**: Chosen for WebSocket management for efficient broadcasting
2. **Repository Pattern**: Not yet implemented (future database layer)
3. **Singleton Logger**: Global logger instance for consistent logging
4. **Factory Pattern**: Used for creating sessions and clients

### Technology Choices
1. **Go**: Chosen for performance and concurrent connection handling
2. **Gorilla WebSocket**: Mature, well-tested WebSocket library
3. **JSON-RPC 2.0**: Structured protocol for API communication
4. **slog**: Structured logging for better observability
5. **In-Memory Storage**: Simple start, plan to add persistence

## Module Dependencies

### External Dependencies
- `github.com/gorilla/websocket` - WebSocket protocol
- `github.com/go-playground/validator/v10` - Validation framework
- `github.com/dustinkirkland/golang-petname` - Human-friendly names
- `github.com/stretchr/testify` - Testing assertions

### Internal Dependencies
```
cmd/server
    ↓
internal/server ←→ internal/websocket
    ↓               ↓
internal/config    internal/jsonrpc
    ↓               ↓
internal/logger    internal/session
```

## Deployment Architecture

### Current State
- Single binary deployment
- Environment variable configuration
- No container orchestration
- Manual deployment process

### Production Readiness
- ✅ Graceful shutdown
- ✅ Health checks
- ✅ Structured logging
- ✅ Error recovery
- ✅ Configuration management
- ❌ Horizontal scaling
- ❌ Load balancing
- ❌ Service discovery
- ❌ Database persistence
- ❌ Backup/recovery

## Future Architecture

### Planned Modules
1. **Frontend-React**: React-based web application
2. **Frontend-Android**: Native Android application
3. **Database Layer**: PostgreSQL with GORM
4. **API Gateway**: For routing and load balancing
5. **Authentication Service**: User management and auth

### Migration Path
1. Add PostgreSQL for session persistence
2. Implement user authentication
3. Add React frontend module
4. Implement content management
5. Add activity generation with LLM integration
6. Scale horizontally with load balancer

---

*Last Updated: 2025-09-08*
*Analysis Scope: Global architecture analysis of FLE platform*
*Current Implementation: Backend-core server module only*