# Optimize Review Planning

## Objective
Analyze the functional areas identified in `audit_result/modules.json` and add an optimized review order to each functional area for independent review.

## Input
- File: `audit_result/modules.json`
- Contains: A list of functional areas with their descriptions, dependencies, and relationships

## Task
You will analyze the modules.json file and add an `"order"` field (integer value) to each functional area object. The order should be optimized from least complex (fewer dependencies) to most complex (more dependencies).

## Optimization Principles

1. **Dependency Resolution**: Review dependencies before dependent features
2. **Progressive Complexity**: Start with standalone systems, then build up to complex interconnected features  
3. **Architectural Foundation**: Understand architecture layers before diving into features
4. **Incremental Knowledge**: Each tier builds on knowledge from previous tiers
5. **Efficient Review**: Minimize context switching and repeated analysis

## Expected Ordering Strategy

### Tier Classification
1. **Foundation (Order 1-3)**: Systems with no dependencies (Testing, Documentation, Workflows)
2. **Core Architecture (Order 4-7)**: Minimal dependencies, foundational layers (Theme, Error, Validation, Layout)
3. **Base Components (Order 8-12)**: UI libraries, state management, API layers
4. **Simple Features (Order 13-16)**: Standalone or minimally connected features
5. **Entity Management (Order 17-23)**: Sub-features with moderate complexity
6. **Feature Hubs (Order 24-26)**: Main feature areas with multiple connections
7. **Complex Systems (Order 27-28)**: Highly interconnected core features

## Analysis Steps

1. **Read and Parse**: Load the modules.json file
2. **Dependency Analysis**: 
   - Identify features with no dependencies (start points)
   - Map out dependency chains
   - Count related_areas for each functional area
3. **Complexity Assessment**:
   - Consider area_type (architecture < sub-feature < feature)
   - Evaluate number of related_areas
   - Factor in existing complexity_score if present
4. **Order Assignment**:
   - Assign order 1 to the simplest, most independent area
   - Progress through increasing complexity
   - Ensure dependencies come before dependents
5. **Update JSON**: Add the "order" field to each functional_area object

## Output Requirements

Update the `audit_result/modules.json` file by adding an `"order"` field to each object in the `functional_areas` array. The order should be an integer from 1 to N (where N is the total number of functional areas).

### Example Structure
```json
{
  "name": "Testing Infrastructure",
  "description": "...",
  "area_type": "architecture",
  "order": 1,  // Added field
  "source_references": [...],
  "related_areas": []
}
```

## Benefits of Optimized Order

- **Progressive Understanding**: Reviewers build knowledge incrementally
- **Reduced Cognitive Load**: Start with simpler systems
- **Efficient Time Usage**: Avoid re-analyzing dependencies
- **Better Quality Reviews**: Understanding foundations improves feature reviews
- **Clear Roadmap**: Provides a structured path through the codebase

## Validation Criteria

- Every functional area must have a unique order number
- Dependencies should have lower order numbers than their dependents
- The ordering should follow the tier classification strategy
- Foundation systems (no dependencies) should come first
- Most complex, interconnected systems should come last