# Documentation Importance Scoring Criteria

## Importance Score Scale

Each functional area should be assigned a documentation importance score from 1 to 5:

- **1 (NOT IMPORTANT)**: Documentation would provide minimal value
- **2 (NICE TO HAVE)**: Documentation would be helpful but not essential
- **3 (IMPORTANT)**: Documentation should be created when resources allow
- **4 (VERY IMPORTANT)**: Documentation is highly recommended
- **5 (MANDATORY)**: Documentation is critical and must be created

## Evaluation Criteria for Importance Scoring

### 1. Complexity Factor (Weight: 25%)
Consider the complexity_score from modules.json:
- **Score 5**: Complexity 5 (highly complex, extensive dependencies)
- **Score 4**: Complexity 4 (many dependencies, complex logic)
- **Score 3**: Complexity 3 (moderate complexity)
- **Score 2**: Complexity 2 (straightforward logic)
- **Score 1**: Complexity 1 (simple, standalone)

### 2. Business Criticality (Weight: 20%)
Assess the business impact:
- **Score 5**: Core business logic, revenue-impacting, mission-critical
- **Score 4**: Important business features, customer-facing
- **Score 3**: Supporting business functions
- **Score 2**: Internal tools or utilities
- **Score 1**: Nice-to-have features, experimental

### 3. Change Frequency (Weight: 15%)
How often is this area likely to change:
- **Score 5**: Frequently modified (weekly/monthly)
- **Score 4**: Regularly updated (quarterly)
- **Score 3**: Occasionally changed (bi-annually)
- **Score 2**: Rarely modified (annually)
- **Score 1**: Stable, unlikely to change

### 4. Onboarding Impact (Weight: 15%)
How important for new team members:
- **Score 5**: Essential for understanding the system
- **Score 4**: Very helpful for quick onboarding
- **Score 3**: Useful for intermediate understanding
- **Score 2**: Helpful for advanced users
- **Score 1**: Specialized knowledge, rarely needed

### 5. Integration Points (Weight: 10%)
Number of other areas that depend on this:
- **Score 5**: Many areas depend on this (5+ related_areas)
- **Score 4**: Several dependencies (3-4 related_areas)
- **Score 3**: Some dependencies (2 related_areas)
- **Score 2**: Few dependencies (1 related_area)
- **Score 1**: No dependencies (standalone)

### 6. Security & Performance Impact (Weight: 10%)
Risk and performance considerations:
- **Score 5**: Security-critical or performance bottleneck
- **Score 4**: Security-relevant or performance-sensitive
- **Score 3**: Some security or performance implications
- **Score 2**: Minor security or performance considerations
- **Score 1**: No security or performance concerns

### 7. Innovation & Uniqueness (Weight: 5%)
Custom or unique implementation:
- **Score 5**: Highly custom, innovative solution
- **Score 4**: Significant customization
- **Score 3**: Some unique aspects
- **Score 2**: Mostly standard implementation
- **Score 1**: Standard, well-known patterns

## Documentation Priority Matrix

Based on the importance score, prioritize documentation efforts:

### Priority 1: MANDATORY (Score 5)
- Must be documented immediately
- Block releases without documentation
- Examples: Authentication, Payment processing, Core APIs

### Priority 2: VERY IMPORTANT (Score 4)
- Document in current sprint/cycle
- Should not accumulate technical debt
- Examples: Main user features, Data models, Key workflows

### Priority 3: IMPORTANT (Score 3)
- Plan for documentation in next cycle
- Document before major releases
- Examples: Admin features, Reporting tools, Integration layers

### Priority 4: NICE TO HAVE (Score 2)
- Document when time permits
- Can be community-contributed
- Examples: Utilities, Helper functions, Test utilities

### Priority 5: NOT IMPORTANT (Score 1)
- Self-documenting code is sufficient
- Document only if specifically requested
- Examples: Simple CRUD, Standard configurations, Boilerplate

## Special Considerations

### Automatic Score Adjustments

Increase score by 1 if:
- Area has NO documentation references in modules.json
- Area type is "architectural-layer" (infrastructure is critical)
- Confidence level is "low" (unclear areas need documentation)

Decrease score by 1 if:
- Area already has comprehensive documentation
- Area uses well-documented third-party libraries
- Area follows standard, well-known patterns

### Minimum Documentation Requirements

Regardless of score, ALL areas should have at least:
- A brief comment or README explaining purpose
- Function/method signatures documented
- Critical business rules noted

Areas with score 3+ should have:
- Comprehensive README
- API documentation
- Usage examples
- Architecture decisions recorded