# Implementation Tasks: backend-core

## Available Agents Detected
- **go-developer**: General Go development and server implementation
- **go-project-architect**: Go project architecture and structure design
- **go-test-engineer**: Go testing, benchmarks, and test coverage
- **go-websocket-specialist**: WebSocket protocol and real-time communication expert
- **product-owner-reviewer**: Requirements validation and specification compliance

**Note**: No Go-specific code reviewer detected. Will use go-project-architect for code reviews as it has architecture review capabilities.

## 📋 Review Process Checklist
- ✅ Every development track has CR, PR, and RW tasks
- ✅ Dependencies correctly block parallel work until reviews pass
- ✅ Code reviewers assigned based on available agents
- ✅ Product reviewer is always product-owner-reviewer
- ✅ Rework tasks use the same agent as development
- ✅ Review outputs have specified file paths

## Review and Rework Process (CRITICAL - DO NOT SKIP)

**⚠️ IMPORTANT: This is a LOOP, not a sequence!**

1. After EACH track's development, run Code Review (CR)
2. If CR fails → Rework → Back to step 1
3. If CR passes → Run Product Review (PR)
4. If PR fails → Rework → Back to step 1
5. Only when BOTH reviews pass can dependent tracks start

## Parallel Execution Tracks

### Track A: Foundation and Project Setup (No Dependencies)
> Primary Agent: go-developer

- [x] 1. **Initialize Go module and project structure**
  - Create `server/` directory
  - Initialize Go module: `go mod init github.com/fle/server`
  - Create directory structure: `cmd/server/`, `internal/`, `config/`
  - Add `.gitignore` for Go projects
  - Files to create: `server/go.mod`, `server/.gitignore`
  - _Requirements: 5.1, 5.2, 5.3_
  - _Agent: go-developer_

- [x] 2. **Set up configuration management**
  - Create `internal/config/config.go` with Config struct
  - Implement environment variable loading with defaults
  - Add `.env.example` with all configuration options
  - Support for development/production environments
  - Files to create: `server/internal/config/config.go`, `server/.env.example`
  - _Requirements: 8.1, 8.2, 8.3, 8.4_
  - _Agent: go-developer_

- [x] 3. **Create Makefile and development tools**
  - Create Makefile with build, run, test, clean targets
  - Add test-coverage target
  - Configure golangci-lint for code quality
  - Files to create: `server/Makefile`, `server/.golangci.yml`
  - _Requirements: 5.4, 5.5_
  - _Agent: go-developer_

- [x] CR-A. **Code Review: Foundation Setup**
  - Review Go module initialization and structure
  - Verify configuration management implementation
  - Check Makefile targets and development tools
  - Validate project organization follows Go best practices
  - Ensure all environment variables have sensible defaults
  - Review output saved to: `server/code_review/CR-A.md`
  - _Dependencies: Tasks 1-3_
  - _Agent: go-project-architect_

- [x] PR-A. **Product Review: Track A Foundation**
  - Validate project setup meets requirements 5.1-5.5
  - Verify configuration management meets requirements 8.1-8.4
  - Check that all specified directories and files exist
  - Ensure development environment is properly configured
  - Confirm Makefile supports all required operations
  - Review output saved to: `server/product_review/track-a.md`
  - _Spec References: requirements.md sections 5.x, 8.x; design.md Configuration and Deployment_
  - _Dependencies: CR-A_
  - _Agent: product-owner-reviewer_

- [x] RW-A. **Rework: Address Track A Review Findings**
  - Review findings from `server/code_review/CR-A.md` and/or `server/product_review/track-a.md`
  - Fix any structural or configuration issues
  - Update Makefile if needed
  - Re-run linting and validation
  - Update documentation if needed
  - _Trigger: Only if CR-A or PR-A status is "Requires changes"_
  - _Dependencies: CR-A and/or PR-A (failed)_
  - _Agent: go-developer_

### Track B: HTTP Server and Logging (Dependencies: Track A reviews approved)
> Primary Agent: go-developer

- [ ] 4. **Implement HTTP server with health endpoint**
  - Create `internal/server/server.go` with Server struct
  - Implement `NewServer()` constructor
  - Add `/health` endpoint handler
  - Configure CORS middleware for development
  - Files to create: `server/internal/server/server.go`, `server/internal/server/handlers.go`
  - _Requirements: 3.1, 3.3, 3.4, 3.5_
  - _Dependencies: PR-A (approved)_
  - _Agent: go-developer_

- [x] 5. **Set up structured logging**
  - Create `internal/logger/logger.go` with slog configuration
  - Support JSON format for production, text for development
  - Implement configurable log levels
  - Add request ID generation for tracing
  - Files to create: `server/internal/logger/logger.go`
  - _Requirements: 9.1, 9.2, 9.3, 9.5_
  - _Dependencies: PR-A (approved)_
  - _Agent: go-developer_

- [x] 6. **Create main entry point**
  - Create/Update `cmd/server/main.go`
  - Load configuration from environment
  - Initialize logger
  - Start HTTP server
  - Implement graceful shutdown
  - Files to create/update: `server/cmd/server/main.go`
  - _Requirements: 3.1, 3.5_
  - _Dependencies: Tasks 4-5_
  - _Agent: go-developer_

- [x] CR-B. **Code Review: HTTP Server and Logging**
  - Review HTTP server implementation and middleware
  - Verify logging configuration and patterns
  - Check error handling and graceful shutdown
  - Validate CORS configuration for security
  - Ensure main entry point follows Go conventions
  - Review output saved to: `server/code_review/CR-B.md`
  - _Dependencies: Tasks 4-6_
  - _Agent: go-project-architect_

- [x] PR-B. **Product Review: Track B HTTP Server**
  - Validate HTTP server meets requirements 3.1, 3.3-3.5
  - Verify logging implementation meets requirements 9.1-9.3, 9.5
  - Test health endpoint responds correctly
  - Confirm CORS works for frontend development
  - Check server starts and logs properly
  - Review output saved to: `server/product_review/track-b.md`
  - _Spec References: requirements.md sections 3.x, 9.x; design.md HTTP Server_
  - _Dependencies: CR-B_
  - _Agent: product-owner-reviewer_

- [ ] RW-B. **Rework: Address Track B Review Findings**
  - Review findings from `server/code_review/CR-B.md` and/or `server/product_review/track-b.md`
  - Fix HTTP server issues
  - Improve logging implementation
  - Update error handling
  - Re-test server startup and shutdown
  - _Trigger: Only if CR-B or PR-B status is "Requires changes"_
  - _Dependencies: CR-B and/or PR-B (failed)_
  - _Agent: go-developer_

### Track C: Session Management (Dependencies: Track A reviews approved)
> Primary Agent: go-developer

- [x] 7. **Implement session code generator**
  - Create `internal/session/generator.go`
  - Integrate golang-petname library
  - Add custom number suffix (1-99)
  - Implement case-insensitive validation
  - Files to create: `server/internal/session/generator.go`
  - _Requirements: 2.1, 2.2, 2.6_
  - _Dependencies: PR-A (approved)_
  - _Agent: go-developer_

- [x] 8. **Create session manager**
  - Create `internal/session/manager.go`
  - Implement SessionManager with in-memory storage
  - Add collision detection and retry logic
  - Implement session creation and retrieval
  - Thread-safe with sync.RWMutex
  - Files to create: `server/internal/session/manager.go`, `server/internal/session/types.go`
  - _Requirements: 2.1, 2.3, 2.4, 2.5_
  - _Dependencies: Task 7_
  - _Agent: go-developer_

- [x] CR-C. **Code Review: Session Management**
  - Review session code generation implementation
  - Verify thread-safety of session manager
  - Check collision handling logic
  - Validate case-insensitive comparison
  - Ensure proper use of golang-petname
  - Review output saved to: `server/code_review/CR-C.md`
  - _Dependencies: Tasks 7-8_
  - _Agent: go-project-architect_

- [x] PR-C. **Product Review: Track C Sessions**
  - Validate session generation meets requirements 2.1-2.2, 2.6
  - Verify session manager meets requirements 2.3-2.5
  - Test session code format (adjective-noun-number)
  - Confirm case-insensitive validation works
  - Check collision handling with multiple sessions
  - Review output saved to: `server/product_review/track-c.md`
  - _Spec References: requirements.md sections 2.x; design.md Session Manager_
  - _Dependencies: CR-C_
  - _Agent: product-owner-reviewer_

- [x] RW-C. **Rework: Address Track C Review Findings**
  - Review findings from `server/code_review/CR-C.md` and/or `server/product_review/track-c.md`
  - Fix session generation issues
  - Improve thread-safety if needed
  - Update collision handling
  - Re-test session management
  - _Trigger: Only if CR-C or PR-C status is "Requires changes"_
  - _Dependencies: CR-C and/or PR-C (failed)_
  - _Agent: go-developer_

### Checkpoint Review 1
- [x] CR1. **Comprehensive Review: Foundation and Core Components** ✅ APPROVED
  - Review overall project architecture consistency
  - Validate integration between HTTP server and session management
  - Check that all components follow Go best practices
  - Ensure proper error handling across modules
  - Verify logging is consistent throughout
  - Review output saved to: `server/code_review/checkpoint-1.md`
  - _Dependencies: PR-B (approved), PR-C (approved)_
  - _Agent: go-project-architect_

### Track D: WebSocket Infrastructure (Dependencies: CR1 approved)
> Primary Agent: go-websocket-specialist

- [x] 9. **Implement WebSocket hub**
  - Create `internal/websocket/hub.go`
  - Implement Hub struct with client management
  - Add register/unregister channels
  - Implement broadcast functionality
  - Thread-safe with sync.RWMutex
  - Files to create: `server/internal/websocket/hub.go`
  - _Requirements: 1.6, 4.1, 4.2, 4.4_
  - _Dependencies: CR1_
  - _Agent: go-websocket-specialist_

- [ ] 10. **Create WebSocket client handler**
  - Create `internal/websocket/client.go`
  - Implement Client struct with read/write pumps
  - Add ping/pong heartbeat mechanism
  - Handle graceful disconnection
  - Implement panic recovery
  - Files to create: `server/internal/websocket/client.go`
  - _Requirements: 1.1, 1.5, 4.3, 4.5, 7.2_
  - _Dependencies: Task 9_
  - _Agent: go-websocket-specialist_

- [ ] 11. **Add WebSocket endpoint to server**
  - Update `internal/server/handlers.go`
  - Add `/ws` endpoint with upgrade handler
  - Integrate session management on connection
  - Send welcome message with session code
  - Files to modify: `server/internal/server/handlers.go`, `server/internal/server/server.go`
  - _Requirements: 1.1, 3.2, 4.3_
  - _Dependencies: Tasks 9-10_
  - _Agent: go-websocket-specialist_

- [ ] CR-D. **Code Review: WebSocket Infrastructure**
  - Review hub implementation and concurrency handling
  - Verify client connection lifecycle management
  - Check WebSocket upgrade and error handling
  - Validate heartbeat mechanism
  - Ensure proper panic recovery
  - Review output saved to: `server/code_review/CR-D.md`
  - _Dependencies: Tasks 9-11_
  - _Agent: go-websocket-specialist_

- [ ] PR-D. **Product Review: Track D WebSocket**
  - Validate WebSocket hub meets requirements 1.6, 4.1-4.2, 4.4
  - Verify client handler meets requirements 1.1, 1.5, 4.3, 4.5
  - Test WebSocket connection and upgrade
  - Confirm session integration works
  - Check heartbeat keeps connections alive
  - Review output saved to: `server/product_review/track-d.md`
  - _Spec References: requirements.md sections 1.x, 4.x; design.md WebSocket Hub_
  - _Dependencies: CR-D_
  - _Agent: product-owner-reviewer_

- [ ] RW-D. **Rework: Address Track D Review Findings**
  - Review findings from `server/code_review/CR-D.md` and/or `server/product_review/track-d.md`
  - Fix WebSocket implementation issues
  - Improve connection handling
  - Update heartbeat mechanism if needed
  - Re-test WebSocket connections
  - _Trigger: Only if CR-D or PR-D status is "Requires changes"_
  - _Dependencies: CR-D and/or PR-D (failed)_
  - _Agent: go-websocket-specialist_

### Track E: JSON-RPC and Validation Framework (Dependencies: CR1 approved)
> Primary Agent: go-developer

- [ ] 12. **Implement JSON-RPC message types**
  - Create `internal/jsonrpc/types.go`
  - Define Request, Response, and Error structs
  - Add validation tags for go-playground/validator
  - Define standard error codes
  - Files to create: `server/internal/jsonrpc/types.go`
  - _Requirements: 1.2, 1.4, 6.1_
  - _Dependencies: CR1_
  - _Agent: go-developer_

- [ ] 13. **Create validation framework**
  - Create `internal/jsonrpc/validator.go`
  - Integrate go-playground/validator
  - Add custom validators for session codes
  - Implement fast-fail validation
  - Return detailed error messages
  - Files to create: `server/internal/jsonrpc/validator.go`
  - _Requirements: 6.1-6.9, 7.5_
  - _Dependencies: Task 12_
  - _Agent: go-developer_

- [ ] 14. **Implement JSON-RPC router**
  - Create `internal/jsonrpc/router.go`
  - Implement method registration system
  - Add request routing and validation
  - Handle response validation
  - Implement error responses
  - Files to create: `server/internal/jsonrpc/router.go`
  - _Requirements: 1.2, 1.3, 1.4, 1.7, 6.2, 6.3_
  - _Dependencies: Tasks 12-13_
  - _Agent: go-developer_

- [ ] 15. **Integrate JSON-RPC with WebSocket**
  - Update `internal/websocket/client.go`
  - Add JSON-RPC message processing
  - Connect router to client handler
  - Implement error handling
  - Files to modify: `server/internal/websocket/client.go`
  - _Requirements: 1.2, 1.3, 1.4_
  - _Dependencies: Task 14_
  - _Agent: go-developer_

- [x] CR-E. **Code Review: JSON-RPC and Validation** ✅ APPROVED
  - Review JSON-RPC type definitions and validation tags
  - Verify validation framework implementation
  - Check router method registration system
  - Validate error handling and responses
  - Ensure fast-fail validation works correctly
  - Review output saved to: `server/code_review/CR-E.md`
  - _Dependencies: Tasks 12-15_
  - _Agent: go-project-architect_

- [x] PR-E. **Product Review: Track E JSON-RPC**
  - Validate JSON-RPC types meet requirement 1.2, 1.4, 6.1
  - Verify validation framework meets requirements 6.1-6.9
  - Test router meets requirements 1.3, 1.7, 6.2-6.3
  - Confirm fast-fail validation works
  - Check detailed error messages are returned
  - Review output saved to: `server/product_review/track-e.md`
  - _Spec References: requirements.md sections 1.x, 6.x, 7.x; design.md JSON-RPC Router_
  - _Dependencies: CR-E_
  - _Agent: product-owner-reviewer_

- [ ] RW-E. **Rework: Address Track E Review Findings**
  - Review findings from `server/code_review/CR-E.md` and/or `server/product_review/track-e.md`
  - Fix JSON-RPC implementation issues
  - Improve validation framework
  - Update router if needed
  - Re-test message processing
  - _Trigger: Only if CR-E or PR-E status is "Requires changes"_
  - _Dependencies: CR-E and/or PR-E (failed)_
  - _Agent: go-developer_

### Track F: Testing and Integration (Dependencies: All previous PR reviews approved)
> Primary Agent: go-test-engineer

- [ ] 16. **Write unit tests for session management**
  - Create `internal/session/generator_test.go`
  - Test code generation uniqueness
  - Test case-insensitive validation
  - Test collision handling
  - Create `internal/session/manager_test.go`
  - Test concurrent access
  - Files to create: `server/internal/session/*_test.go`
  - _Requirements: Testing Strategy_
  - _Dependencies: PR-C (approved), PR-D (approved), PR-E (approved)_
  - _Agent: go-test-engineer_

- [ ] 17. **Write WebSocket integration tests**
  - Create `internal/websocket/hub_test.go`
  - Test concurrent connections
  - Test message broadcasting
  - Create `internal/websocket/client_test.go`
  - Test connection lifecycle
  - Files to create: `server/internal/websocket/*_test.go`
  - _Requirements: Testing Strategy_
  - _Dependencies: PR-D (approved)_
  - _Agent: go-test-engineer_

- [ ] 18. **Write JSON-RPC validation tests**
  - Create `internal/jsonrpc/validator_test.go`
  - Test all validation scenarios
  - Test error responses
  - Create `internal/jsonrpc/router_test.go`
  - Test method registration and routing
  - Files to create: `server/internal/jsonrpc/*_test.go`
  - _Requirements: Testing Strategy_
  - _Dependencies: PR-E (approved)_
  - _Agent: go-test-engineer_

- [ ] 19. **Create end-to-end integration test**
  - Create `cmd/server/main_test.go`
  - Test full connection flow
  - Test session creation on connect
  - Test JSON-RPC message processing
  - Test graceful shutdown
  - Files to create: `server/cmd/server/main_test.go`
  - _Requirements: Testing Strategy_
  - _Dependencies: Tasks 16-18_
  - _Agent: go-test-engineer_

- [ ] CR-F. **Code Review: Testing and Integration**
  - Review test coverage and quality
  - Verify all edge cases are tested
  - Check test organization and naming
  - Validate integration test scenarios
  - Ensure tests are maintainable
  - Review output saved to: `server/code_review/CR-F.md`
  - _Dependencies: Tasks 16-19_
  - _Agent: go-project-architect_

- [ ] PR-F. **Product Review: Track F Testing**
  - Validate test coverage meets all requirements
  - Verify unit tests cover critical paths
  - Check integration tests validate full flow
  - Confirm all error scenarios are tested
  - Ensure tests are passing
  - Review output saved to: `server/product_review/track-f.md`
  - _Spec References: design.md Testing Strategy_
  - _Dependencies: CR-F_
  - _Agent: product-owner-reviewer_

- [ ] RW-F. **Rework: Address Track F Review Findings**
  - Review findings from `server/code_review/CR-F.md` and/or `server/product_review/track-f.md`
  - Improve test coverage
  - Fix failing tests
  - Add missing test scenarios
  - Update test documentation
  - _Trigger: Only if CR-F or PR-F status is "Requires changes"_
  - _Dependencies: CR-F and/or PR-F (failed)_
  - _Agent: go-test-engineer_

### Final Review Track

- [ ] CR-FINAL. **Final Comprehensive Code Review**
  - Review entire codebase for consistency
  - Verify all components integrate properly
  - Check security considerations
  - Validate performance implications
  - Ensure production readiness
  - Review output saved to: `server/code_review/final.md`
  - _Dependencies: PR-F (approved)_
  - _Agent: go-project-architect_

- [ ] PR-FINAL. **Final Product Review**
  - Validate all requirements are met
  - Test complete system functionality
  - Verify WebSocket connections work
  - Check session management works
  - Confirm JSON-RPC validation works
  - Review output saved to: `server/product_review/final.md`
  - _Dependencies: CR-FINAL_
  - _Agent: product-owner-reviewer_

## Execution Strategy

### Parallel Groups with Review Gates

1. **Group 1 (Immediate Start)**:
   - Track A: Tasks 1-3 (Foundation setup)

2. **Group 2 (Review Gate 1)**:
   - CR-A → PR-A → (RW-A if needed)

3. **Group 3 (After PR-A Approval - Parallel Execution)**:
   - Track B: Tasks 4-6 (HTTP Server)
   - Track C: Tasks 7-8 (Session Management)

4. **Group 4 (Review Gate 2)**:
   - CR-B → PR-B → (RW-B if needed)
   - CR-C → PR-C → (RW-C if needed)

5. **Group 5 (Checkpoint)**:
   - CR1 (Comprehensive checkpoint review)

6. **Group 6 (After CR1 Approval - Parallel Execution)**:
   - Track D: Tasks 9-11 (WebSocket)
   - Track E: Tasks 12-15 (JSON-RPC)

7. **Group 7 (Review Gate 3)**:
   - CR-D → PR-D → (RW-D if needed)
   - CR-E → PR-E → (RW-E if needed)

8. **Group 8 (After All Implementation Reviews)**:
   - Track F: Tasks 16-19 (Testing)

9. **Group 9 (Review Gate 4)**:
   - CR-F → PR-F → (RW-F if needed)

10. **Group 10 (Final Review)**:
    - CR-FINAL → PR-FINAL

### Agent Utilization Summary
- **Primary Development**: go-developer (11 tasks)
- **WebSocket Specialist**: go-websocket-specialist (3 tasks + 1 review)
- **Testing Specialist**: go-test-engineer (4 tasks)
- **Code Review**: go-project-architect (8 reviews)
- **Product Review**: product-owner-reviewer (7 reviews)
- **Rework**: Original developers (6 conditional tasks)

### Time Estimates
- **Parallel execution time**: ~3-4 days with all agents working
- **Sequential execution time**: ~8-10 days
- **Review overhead**: ~2-3 hours per track (including potential rework)
- **Total tasks**: 35 (19 development, 8 code reviews, 7 product reviews, 1 checkpoint)
- **Conditional tasks**: 6 rework tasks (only if reviews fail)

## Success Metrics
- ✅ All tests passing with >80% coverage
- ✅ WebSocket connections stable under load
- ✅ Session codes are memorable and unique
- ✅ JSON-RPC validation catches all invalid inputs
- ✅ Server handles 1000+ concurrent connections
- ✅ All review gates passed

---

Do the tasks, agent assignments, and review process look good?