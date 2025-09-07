# Architecture & Design Pattern Analysis

## Architecture Evaluation Framework

### 1. Layered Architecture Assessment

**Common Layers** (identify and validate):
- **Presentation Layer**: UI components, views, controllers
- **Application Layer**: Use cases, application services, DTOs
- **Domain Layer**: Business logic, entities, domain services
- **Infrastructure Layer**: Database, external services, frameworks

**Layer Violations to Detect**:
- Presentation accessing database directly
- Domain depending on infrastructure
- Circular dependencies between layers
- Business logic in wrong layer
- Data access in presentation layer

### 2. Component Architecture Patterns

**Microservices Indicators**:
- Service boundaries alignment with business capabilities
- Inter-service communication patterns
- Data consistency strategies
- Service autonomy level

**Monolithic Patterns**:
- Module boundaries and cohesion
- Package/namespace organization
- Shared state management
- Cross-cutting concerns handling

**Frontend Patterns**:
- Component composition (Atomic Design)
- State management architecture
- Data flow patterns (unidirectional, bidirectional)
- Separation of container/presentational components

### 3. Dependency Management

**Good Practices**:
- Dependency Injection usage
- Interface-based programming
- Loose coupling patterns
- Clear dependency direction
- Dependency inversion implementation

**Issues to Flag**:
- Circular dependencies
- Hidden dependencies
- Tight coupling
- Concrete class dependencies
- Global state usage
- Service locator anti-pattern

### 4. Data Architecture

**Data Access Patterns**:
- Repository Pattern implementation
- Unit of Work pattern
- Active Record vs Data Mapper
- CQRS (Command Query Responsibility Segregation)
- Event Sourcing indicators

**Database Design Issues**:
- N+1 query problems
- Missing indexes indicators
- Over-normalization/under-normalization
- Missing foreign key constraints
- Inappropriate NoSQL usage

### 5. Integration Patterns

**API Design**:
- RESTful principles adherence
- GraphQL schema design
- API versioning strategy
- Contract-first development
- Error response consistency

**Messaging Patterns**:
- Message Queue usage
- Pub/Sub implementation
- Event-driven architecture
- Saga pattern implementation
- Eventual consistency handling

### 6. Cross-Cutting Concerns

**Properly Handled**:
- Logging strategy
- Error handling consistency
- Security implementation
- Caching strategy
- Configuration management
- Monitoring/metrics collection

**Common Issues**:
- Scattered logging logic
- Inconsistent error handling
- Security logic duplication
- Cache invalidation problems
- Hard-coded configurations

### 7. Modularity & Cohesion

**High Cohesion Indicators**:
- Related functionality grouped
- Clear module boundaries
- Single purpose modules
- Minimal inter-module dependencies

**Low Cohesion Warning Signs**:
- Scattered related logic
- Feature envy between modules
- Unclear module purposes
- Excessive module dependencies

### 8. Scalability Patterns

**Horizontal Scaling Ready**:
- Stateless services
- Session externalization
- Database connection pooling
- Cache-aside pattern
- Load balancing considerations

**Scalability Issues**:
- In-memory state
- File system dependencies
- Single points of failure
- Synchronous long-running operations
- Missing circuit breakers

## Framework-Specific Patterns

### React/Vue/Angular
- Component lifecycle management
- State management (Redux, Vuex, NgRx)
- Routing architecture
- Form handling patterns
- Side effects management

### Node.js/Express
- Middleware organization
- Route structuring
- Error middleware
- Async handling patterns
- Cluster usage

### Spring/Django/Rails
- MVC implementation
- Service layer patterns
- ORM usage patterns
- Dependency injection
- AOP (Aspect-Oriented Programming)

### Microservices
- Service discovery
- Circuit breakers
- API Gateway pattern
- Distributed tracing
- Service mesh usage

## Architecture Debt Indicators

### High Priority
- No clear architecture
- Mixed concerns in components
- Circular dependencies
- Missing abstraction layers
- Hardcoded integrations

### Medium Priority
- Inconsistent patterns
- Partial pattern implementation
- Minor layer violations
- Missing some cross-cutting concerns
- Suboptimal data access

### Low Priority
- Style inconsistencies
- Minor coupling issues
- Optimization opportunities
- Pattern modernization options
- Documentation gaps

## Refactoring Recommendations

### Structural Refactoring
- Extract layers
- Introduce interfaces
- Apply dependency injection
- Separate concerns
- Create facades

### Pattern Introduction
- Implement Repository pattern
- Add Service layer
- Introduce DTOs
- Apply Strategy pattern
- Use Factory pattern

### Modularization
- Extract modules
- Define clear boundaries
- Reduce coupling
- Increase cohesion
- Create shared libraries