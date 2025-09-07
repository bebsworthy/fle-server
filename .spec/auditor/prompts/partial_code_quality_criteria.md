# Code Quality Assessment Criteria

## Evaluation Framework

### 1. Complexity Analysis
- **Cyclomatic Complexity**: Count decision points (if, for, while, case, catch, &&, ||, ?:)
  - Low: 1-5 (simple, easy to test)
  - Medium: 6-10 (moderate, requires careful testing)
  - High: 11-20 (complex, difficult to test)
  - Very High: 20+ (highly complex, refactoring recommended)

- **Cognitive Complexity**: Mental effort to understand code
  - Nesting depth penalties
  - Break in linear flow penalties
  - Language construct penalties

### 2. SOLID Principles Compliance
- **Single Responsibility**: Each module/class has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Subtypes must be substitutable
- **Interface Segregation**: No forced implementation of unused methods
- **Dependency Inversion**: Depend on abstractions, not concretions

### 3. Code Duplication
- **Type 1**: Exact code clones
- **Type 2**: Renamed/parameterized clones
- **Type 3**: Modified clones with changed statements
- **Type 4**: Semantic clones (different code, same logic)

### 4. Design Patterns & Anti-patterns

**Good Patterns to Identify**:
- Factory, Strategy, Observer, Decorator
- Repository, Service Layer, DTO
- Dependency Injection
- Command Query Separation

**Anti-patterns to Flag**:
- God Object/Class (doing too much)
- Spaghetti Code (tangled control flow)
- Copy-Paste Programming
- Magic Numbers/Strings
- Long Parameter Lists (>3-4 parameters)
- Feature Envy (method uses another class's data excessively)
- Inappropriate Intimacy (classes know too much about each other)

### 5. Error Handling Quality
- **Coverage**: All error cases handled
- **Specificity**: Specific exceptions vs generic catch-all
- **Recovery**: Graceful degradation strategies
- **Logging**: Appropriate error logging
- **User Feedback**: Clear error messages

### 6. Code Maintainability Indicators
- **Naming Quality**: Clear, descriptive, consistent naming
- **Function Length**: <20 lines preferred, <50 acceptable
- **File Length**: <300 lines preferred, <500 acceptable
- **Comment Quality**: Why, not what
- **Documentation**: Public APIs documented
- **Test Coverage**: Unit tests present and meaningful

### 7. Performance Considerations
- **Algorithm Efficiency**: O(n) vs O(n²) vs O(n³)
- **Database Queries**: N+1 problems, missing indexes
- **Memory Management**: Leaks, excessive allocations
- **Caching Strategy**: Appropriate use of caching
- **Async Patterns**: Proper async/await usage

### 8. Security Indicators
- **Input Validation**: All inputs validated
- **SQL Injection**: Parameterized queries used
- **XSS Prevention**: Output encoding applied
- **Authentication**: Proper auth checks
- **Authorization**: Role-based access control
- **Secrets Management**: No hardcoded secrets
- **Dependencies**: Known vulnerabilities check

## Severity Classification

### Critical (Must Fix)
- Security vulnerabilities
- Data loss risks
- System crashes
- Infinite loops
- Memory leaks
- Exposed secrets

### High (Should Fix Soon)
- Performance bottlenecks
- Missing error handling
- Very high complexity (>20)
- Significant duplication (>50 lines)
- SOLID violations with high impact
- Missing critical tests

### Medium (Plan to Fix)
- Moderate complexity (11-20)
- Moderate duplication (20-50 lines)
- Inconsistent patterns
- Missing documentation
- Long methods (>50 lines)
- Magic numbers/strings

### Low (Nice to Fix)
- Minor naming issues
- Minor duplication (<20 lines)
- Missing optional documentation
- Style inconsistencies
- Opportunities for optimization
- Minor refactoring candidates

## Quick Win Identification

Look for improvements that are:
- **Low Effort + High Impact**: Priority 1
- **Low Effort + Medium Impact**: Priority 2
- **Medium Effort + High Impact**: Priority 3

Examples:
- Extract magic numbers to constants
- Add missing error handling
- Extract duplicate code to utilities
- Add missing null checks
- Rename unclear variables
- Split large functions