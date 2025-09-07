---
type: framework
name: django
keywords: [python, django, web, orm, mvc]
extends:
  - patterns/rest-api.md
priority: 10
---

# Django Documentation Standards

## Documentation Requirements

### Model Documentation
```python
class Product(models.Model):
    """
    Represents a product in the catalog.
    
    Relationships:
        - Category: Many-to-one via category_id
        - Reviews: One-to-many via Review.product
    """
    name = models.CharField(
        max_length=200,
        help_text="Product display name"
    )
    price = models.DecimalField(
        max_digits=10,
        decimal_places=2,
        help_text="Price in USD"
    )
    created_at = models.DateTimeField(
        auto_now_add=True,
        help_text="Timestamp when product was created"
    )
    
    class Meta:
        verbose_name = "Product"
        verbose_name_plural = "Products"
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['name']),
        ]
    
    def __str__(self):
        return self.name
```

### View Documentation
```python
class ProductViewSet(viewsets.ModelViewSet):
    """
    API endpoint for product management.
    
    Provides CRUD operations for products with filtering and pagination.
    
    Permissions:
        - List/Retrieve: Any user
        - Create/Update/Delete: Staff only
    """
    queryset = Product.objects.all()
    serializer_class = ProductSerializer
    permission_classes = [IsAuthenticatedOrReadOnly]
    filterset_fields = ['category', 'price_range']
    
    def perform_create(self, serializer):
        """
        Set the owner to the current user on creation.
        """
        serializer.save(owner=self.request.user)
```

### Serializer Documentation
```python
class ProductSerializer(serializers.ModelSerializer):
    """
    Serializer for Product model.
    
    Includes calculated fields for display price and availability.
    """
    display_price = serializers.SerializerMethodField(
        help_text="Formatted price with currency symbol"
    )
    
    class Meta:
        model = Product
        fields = ['id', 'name', 'price', 'display_price', 'created_at']
        read_only_fields = ['created_at']
    
    def get_display_price(self, obj):
        """Format price for display."""
        return f"${obj.price}"
```

## Review Checklist

- [ ] Models have docstrings explaining purpose and relationships
- [ ] Model fields use `help_text` parameter
- [ ] Meta classes document ordering and indexes
- [ ] Views/ViewSets document permissions and behavior
- [ ] Serializers document field transformations
- [ ] URL patterns are documented
- [ ] Custom managers and querysets documented
- [ ] Signals documented with sender and receiver
- [ ] Middleware purpose and order documented
- [ ] Settings are documented with comments

## Django-Specific Patterns

### URL Configuration
```python
urlpatterns = [
    # Product endpoints
    path('products/', ProductListView.as_view(), name='product-list'),
    path('products/<int:pk>/', ProductDetailView.as_view(), name='product-detail'),
    
    # Category endpoints
    path('categories/', include('categories.urls')),
]
```

### Custom Managers
```python
class PublishedManager(models.Manager):
    """
    Manager that returns only published items.
    
    Usage:
        Product.published.all()  # Only published products
    """
    def get_queryset(self):
        return super().get_queryset().filter(status='published')
```

### Admin Configuration
```python
@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    """
    Admin interface for Product management.
    
    Features:
        - Inline editing of variants
        - Bulk actions for publishing/unpublishing
        - Advanced filtering by date and category
    """
    list_display = ['name', 'price', 'category', 'created_at']
    list_filter = ['category', 'created_at']
    search_fields = ['name', 'description']
```

## Common Issues

### Issue: Missing Model Relationships Documentation
Models without clear relationship documentation.
**Solution**: Document all ForeignKey, ManyToMany relationships in class docstring.

### Issue: Undocumented Permissions
Views without clear permission requirements.
**Solution**: Document permission classes and requirements in view docstrings.

### Issue: No Migration Documentation
Migrations without explanation of changes.
**Solution**: Add comments in migration files explaining the purpose.