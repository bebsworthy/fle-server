# Documentation Quality Guidelines

## Rating Scale

Documentation quality must be rated on a scale of 1 to 5:

- **1 (BAD)**: Critical documentation missing, severe organizational issues, or mostly outdated/incorrect
- **2 (POOR)**: Major gaps in documentation, poor organization, significant accuracy issues
- **3 (ADEQUATE)**: Basic documentation present but lacking depth, some organizational issues
- **4 (GOOD)**: Comprehensive documentation with minor gaps, well-organized, accurate
- **5 (EXCELLENT)**: Complete, well-structured, optimized documentation following all best practices

## Evaluation Criteria

### 1. Content Completeness (30% weight)
- **Architecture Overview**: Clear system architecture documentation
- **Component Documentation**: Each functional area has proper documentation
- **API Documentation**: All public APIs are documented with examples
- **Setup Instructions**: Clear installation and configuration guides
- **Dependencies**: All dependencies and their purposes documented
- **Business Logic**: Key business rules and workflows explained

### 2. Organization & Structure (25% weight)
- **Logical Hierarchy**: Documentation follows intuitive structure
- **Navigation**: Easy to find specific information
- **Cross-referencing**: Proper linking between related documents
- **Consistency**: Uniform format and style across documents
- **Index/TOC**: Comprehensive table of contents or index

### 3. LLM Optimization (20% weight)
- **File Size**: Individual files between 2-10KB for optimal context
- **Chunking**: Complex topics split into digestible sections
- **Self-contained**: Each file can be understood independently
- **Metadata**: Clear frontmatter with tags and categories
- **Summaries**: TL;DR sections for quick comprehension

### 4. Code Documentation Best Practices (15% weight)
- **Code Comments**: Inline comments for complex logic
- **Type Documentation**: All types and interfaces documented
- **Examples**: Code examples for common use cases
- **Error Handling**: Documentation of error states and handling
- **Performance Notes**: Documentation of performance considerations

### 5. Maintenance & Updates (10% weight)
- **Version Information**: Documentation versioned with code
- **Change Log**: Updates tracked and documented
- **Deprecation Notices**: Clear marking of deprecated features
- **Review Dates**: Last review/update dates on documents
- **Contribution Guidelines**: How to update documentation

## Required Documentation Structure

```
project/
├── README.md                 # Project overview and quick start
├── docs/
│   ├── README.md            # Index to all the document with summary description
│   ├── architecture/        # System architecture documentation
│   │   ├── overview.md
│   │   ├── data-flow.md
│   │   └── components.md
│   ├── api/                 # API documentation
│   │   ├── rest-api.md
│   │   └── graphql.md
│   ├── guides/              # How-to guides
│   │   ├── setup.md
│   │   ├── deployment.md
│   │   └── troubleshooting.md
│   ├── reference/           # Technical reference
│   │   ├── configuration.md
│   │   ├── cli.md
│   │   └── environment.md
│   └── development/         # Development documentation
│       ├── contributing.md
│       ├── testing.md
│       └── code-style.md
```

## Documentation Anti-patterns to Avoid

1. **Monolithic Files**: Single files over 1000 lines
2. **Outdated Examples**: Code examples that no longer work
3. **Missing Context**: Documentation assuming too much prior knowledge
4. **No Visual Aids**: Complex concepts without diagrams
5. **Circular References**: Documents that reference each other without clear entry points
6. **Jargon Overload**: Excessive technical terms without explanations
7. **No Search Optimization**: Missing keywords and tags for searchability

## Specific Requirements for Quality Ratings

### For Rating 5 (EXCELLENT)
- All evaluation criteria score 90% or higher
- Documentation is fully automated/generated where possible
- Includes interactive examples or playgrounds
- Multiple learning paths for different audiences
- Comprehensive troubleshooting section

### For Rating 4 (GOOD)
- All evaluation criteria score 70% or higher
- No critical documentation gaps
- Clear and consistent structure
- Most best practices followed

### For Rating 3 (ADEQUATE)
- All evaluation criteria score 50% or higher
- Basic documentation exists for all major components
- Some organizational structure present
- Key functionality documented

### For Rating 2 (POOR)
- Some documentation exists but with major gaps
- Poor organization makes information hard to find
- Many best practices ignored

### For Rating 1 (BAD)
- Critical documentation missing
- What exists is outdated or incorrect
- No clear organization
- Blocks new developers from understanding the system