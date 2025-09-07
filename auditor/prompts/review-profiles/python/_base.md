---
type: language
name: python
keywords: [python, py, python3]
priority: 2
---

# Python Documentation Standards

## Documentation Requirements

### Module Documentation
- **Module Docstring**: First statement explaining module purpose
- **Imports**: Grouped and ordered (standard library, third-party, local)
- **Module-level Variables**: Document constants and configuration
- **Public API**: Explicitly define `__all__` for public exports

### Function Documentation
All public functions must have docstrings containing:
- **Summary**: One-line description
- **Parameters**: Name, type, and description
- **Returns**: Type and description
- **Raises**: Exceptions that may be raised
- **Examples**: Usage examples for complex functions

### Class Documentation
Classes should document:
- **Purpose**: What the class represents
- **Attributes**: Public and protected attributes
- **Methods**: Public interface documentation
- **Usage**: Example instantiation and common patterns

### Type Hints
- Use type hints for all function signatures
- Import from `typing` for complex types
- Document generic types and type variables
- Use `Optional[]` for nullable parameters

## Review Checklist

- [ ] Module has descriptive docstring
- [ ] All public functions have docstrings
- [ ] Docstrings follow chosen convention (Google/NumPy/Sphinx)
- [ ] Type hints present for parameters and returns
- [ ] Complex types are documented
- [ ] Exceptions are documented in docstrings
- [ ] Examples provided for complex functionality
- [ ] Class attributes documented
- [ ] Private vs public interface is clear
- [ ] Dependencies and requirements documented

## Docstring Conventions

### Google Style
```python
def function(param1: str, param2: int) -> bool:
    """Summary line.
    
    Longer description if needed.
    
    Args:
        param1: Description of param1.
        param2: Description of param2.
    
    Returns:
        Description of return value.
    
    Raises:
        ValueError: If param2 is negative.
    """
```

### NumPy Style
```python
def function(param1: str, param2: int) -> bool:
    """
    Summary line.
    
    Longer description if needed.
    
    Parameters
    ----------
    param1 : str
        Description of param1.
    param2 : int
        Description of param2.
    
    Returns
    -------
    bool
        Description of return value.
    """
```

## Best Practices

### Documentation Tools
- Use `sphinx` for generating documentation
- Configure `autodoc` for automatic API docs
- Use `napoleon` for Google/NumPy style support
- Consider `mkdocs` for simpler projects

### Type Annotations
```python
from typing import List, Optional, Union, Dict, Any

def process_data(
    items: List[str],
    config: Optional[Dict[str, Any]] = None
) -> Union[str, None]:
    """Process data with optional configuration."""
```

### Async Documentation
- Document async/await patterns
- Explain coroutine behavior
- Document event loop requirements

## Common Issues

### Issue: Missing Type Hints
Functions without type annotations.
**Solution**: Add type hints to all public functions, use `mypy` for validation.

### Issue: Inconsistent Docstring Style
Mixed Google and NumPy styles.
**Solution**: Choose one convention and enforce with linting tools.

### Issue: Undocumented Side Effects
Functions that modify state without documentation.
**Solution**: Clearly document any side effects or state modifications.