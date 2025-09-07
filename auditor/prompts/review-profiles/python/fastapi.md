---
type: framework
name: fastapi
keywords: [python, fastapi, api, web, async, pydantic]
extends:
  - patterns/rest-api.md
priority: 10
---

# FastAPI Documentation Standards

## Documentation Requirements

### Endpoint Documentation
FastAPI provides automatic documentation, but enhance it with:
```python
@app.get(
    "/items/{item_id}",
    response_model=ItemResponse,
    summary="Get an item by ID",
    description="Retrieve a single item by its unique identifier with full details",
    response_description="The requested item",
    responses={
        404: {"description": "Item not found"},
        403: {"description": "Not authorized to view this item"}
    }
)
async def get_item(
    item_id: int = Path(..., description="The ID of the item to retrieve", ge=1),
    include_details: bool = Query(False, description="Include extended details")
) -> ItemResponse:
    """
    Internal documentation for developers.
    This won't appear in OpenAPI docs.
    """
```

### Pydantic Models
Document all models with field descriptions:
```python
from pydantic import BaseModel, Field

class Item(BaseModel):
    """Represents a product item in the system."""
    
    id: int = Field(..., description="Unique identifier", example=123)
    name: str = Field(..., description="Item display name", min_length=1, max_length=100)
    price: float = Field(..., description="Price in USD", gt=0, example=29.99)
    tags: List[str] = Field(default=[], description="Categorization tags")
    
    class Config:
        schema_extra = {
            "example": {
                "id": 1,
                "name": "Premium Widget",
                "price": 99.99,
                "tags": ["electronics", "premium"]
            }
        }
```

### Dependency Injection
Document dependencies clearly:
```python
async def get_current_user(
    token: str = Depends(oauth2_scheme)
) -> User:
    """
    Validate OAuth2 token and return current user.
    
    Used as a dependency in protected endpoints.
    """
```

### Background Tasks
```python
@app.post("/send-notification/")
async def send_notification(
    email: str,
    background_tasks: BackgroundTasks
):
    """
    Queue a notification email to be sent asynchronously.
    
    The email will be sent in the background after response is returned.
    """
    background_tasks.add_task(send_email, email=email)
```

## Review Checklist

- [ ] All endpoints have `summary` and `description`
- [ ] Path parameters use `Path()` with descriptions
- [ ] Query parameters use `Query()` with descriptions
- [ ] Request bodies use Pydantic models with Field descriptions
- [ ] Response models are explicit with `response_model`
- [ ] Error responses documented in `responses` dict
- [ ] Pydantic models have class docstrings
- [ ] All fields use `Field()` with descriptions
- [ ] Examples provided via `schema_extra` or `example`
- [ ] Dependencies are documented
- [ ] Background tasks explain async behavior
- [ ] WebSocket endpoints document protocol
- [ ] OpenAPI docs accessible at `/docs`

## Best Practices

### Response Models
```python
class ItemResponse(BaseModel):
    """Response model for item endpoints."""
    data: Item
    meta: ResponseMeta
    
@app.get("/items/{id}", response_model=ItemResponse)
```

### Status Codes
```python
from fastapi import status

@app.post(
    "/items/",
    status_code=status.HTTP_201_CREATED,
    response_model=Item
)
```

### Validation Documentation
```python
class CreateItemRequest(BaseModel):
    name: str = Field(
        ...,
        min_length=3,
        max_length=50,
        description="Product name (3-50 characters)",
        regex="^[a-zA-Z0-9 ]+$"
    )
```

### API Versioning
```python
# Document versioning strategy
v1_router = APIRouter(prefix="/api/v1", tags=["v1"])
v2_router = APIRouter(prefix="/api/v2", tags=["v2"])
```

## OpenAPI Customization

```python
app = FastAPI(
    title="My API",
    description="API for managing items",
    version="1.0.0",
    docs_url="/documentation",
    redoc_url="/redoc",
    openapi_tags=[
        {
            "name": "items",
            "description": "Operations with items"
        }
    ]
)
```

## Common Issues

### Issue: Missing Field Descriptions
Pydantic models without Field descriptions.
**Solution**: Always use `Field()` with description parameter.

### Issue: Unclear Validation Rules
Complex validation without documentation.
**Solution**: Document validation in Field description and use clear parameter names.

### Issue: No Response Examples
Endpoints without example responses.
**Solution**: Add `schema_extra` with examples to Pydantic models.