# Standard Documentation Structure

## Required Sections for Feature/Component Documentation

### 1. Header Section (Frontmatter)
```markdown
---
title: [Area Name]
category: [feature|sub-feature|architecture]
tags: [relevant, tags]
last_updated: [YYYY-MM-DD]
revision_history: # Optional, for updates
  - [YYYY-MM-DD]: [change summary]
---
```

### 2. Overview Section
- **Brief description** (2-3 sentences)
- **Purpose and business value**
- **Key responsibilities**
- **Target audience** (who uses this component)

### 3. Architecture & Design
- **Component structure**
- **Design patterns used**
- **Data flow** (if applicable)
- **State management approach**
- **Use mermaid diagrams** where helpful:
  ```mermaid
  graph TD
    A[Component] --> B[Dependency]
  ```

### 4. Dependencies
- **External libraries used**
  - Library name: version - purpose
- **Internal dependencies** (other functional areas)
- **Required configurations**
- **Environment requirements**

### 5. API Documentation (if applicable)
- **Public interfaces**
- **Method signatures**
- **Parameters and return types**
- **Usage examples**
- **Error responses**

Example format:
```typescript
methodName(param1: Type, param2?: Type): ReturnType
```
- **Description**: What this method does
- **Parameters**:
  - `param1`: Description
  - `param2` (optional): Description
- **Returns**: Description of return value
- **Example**: Code snippet

### 6. Implementation Details
- **Key algorithms or logic**
- **Important business rules**
- **Performance considerations**
- **Error handling approach**
- **Security considerations**

### 7. Configuration
- **Environment variables**
  ```bash
  VARIABLE_NAME=value # Description
  ```
- **Configuration files**
  - File path and format
  - Key settings explained
- **Feature flags**
  - Flag name: purpose and values

### 8. Usage Examples
- **Common use cases**
- **Code snippets** (from actual source)
- **Integration examples**
- **Best practices**
- **Anti-patterns to avoid**

### 9. Testing
- **Testing approach**
- **How to run tests**
  ```bash
  npm test [specific-test-pattern]
  ```
- **Key test scenarios**
- **Test coverage expectations**
- **Mock/stub requirements**

### 10. Troubleshooting (optional but recommended)
- **Common issues and solutions**
- **Debug tips**
- **Logging locations**
- **Performance profiling**

### 11. Related Documentation
- **Links to related functional areas**
- **External documentation references**
- **API specifications**
- **Design documents**

## Documentation Standards

### File Size Guidelines
- Target: 3-10KB for optimal LLM consumption
- Split large topics into multiple focused documents
- Use links for detailed sub-topics

### Code Examples
- **MUST** be from actual source code
- Include file path references: `// from src/component/file.ts`
- Keep examples concise but complete
- Show both basic and advanced usage

### Formatting Guidelines
- Use proper markdown syntax
- Include syntax highlighting for code blocks
- Use tables for structured data
- Add visual breaks between major sections
- Use bold for **important terms**
- Use inline code for `variables` and `methods`

### Update Guidelines
When updating existing documentation:
- Preserve sections marked with `<!-- manually-added -->`
- Update timestamps in frontmatter
- Add entries to revision_history
- Mark deprecated features appropriately
- Include migration guides when needed

### Quality Checklist
- [ ] All sections present and complete
- [ ] Code examples tested and working
- [ ] Links validated
- [ ] Technical accuracy verified
- [ ] Business context included
- [ ] Target audience considered
- [ ] File size optimized (3-10KB)