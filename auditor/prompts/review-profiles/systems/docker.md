---
type: system
name: docker
keywords: [docker, container, dockerfile, compose, containerization]
priority: 8
---

# Docker Documentation Standards

## Documentation Requirements

### Dockerfile Documentation
```dockerfile
# Base image with specific version for reproducibility
FROM node:18-alpine AS builder

# Metadata
LABEL maintainer="team@example.com"
LABEL description="API service for product management"
LABEL version="1.0.0"

# Install system dependencies
# Required for: node-gyp compilation
RUN apk add --no-cache python3 make g++

# Set working directory
WORKDIR /app

# Copy dependency files first for layer caching
# Changes to package files will invalidate cache
COPY package*.json ./

# Install dependencies
# --ci flag ensures reproducible builds
RUN npm ci --only=production

# Copy application code
COPY . .

# Build application
RUN npm run build

# Production stage - smaller final image
FROM node:18-alpine

# Run as non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

WORKDIR /app

# Copy built application from builder stage
COPY --from=builder --chown=nodejs:nodejs /app/dist ./dist
COPY --from=builder --chown=nodejs:nodejs /app/node_modules ./node_modules

# Switch to non-root user
USER nodejs

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node healthcheck.js

# Document exposed port
EXPOSE 3000

# Start command with proper signal handling
CMD ["node", "--enable-source-maps", "dist/index.js"]
```

### Docker Compose Documentation
```yaml
version: '3.8'

# Service definitions
services:
  # API service - handles HTTP requests
  api:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        # Build-time variables
        NODE_ENV: production
    image: myapp/api:${VERSION:-latest}
    container_name: myapp-api
    restart: unless-stopped
    ports:
      - "3000:3000"  # Host:Container port mapping
    environment:
      # Runtime configuration
      NODE_ENV: production
      DATABASE_URL: postgres://user:pass@db:5432/myapp
      REDIS_URL: redis://cache:6379
    depends_on:
      # Service dependencies - wait for these to start
      db:
        condition: service_healthy
      cache:
        condition: service_started
    networks:
      - backend
    volumes:
      # Mount for uploaded files
      - uploads:/app/uploads
      # Read-only config mount
      - ./config:/app/config:ro
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
    
  # PostgreSQL database
  db:
    image: postgres:15-alpine
    container_name: myapp-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: ${DB_PASSWORD}  # From .env file
    volumes:
      # Persist database data
      - postgres_data:/var/lib/postgresql/data
      # Initialization scripts
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - backend

# Network configuration
networks:
  backend:
    driver: bridge
    name: myapp-backend

# Volume definitions
volumes:
  postgres_data:
    driver: local
  uploads:
    driver: local
```

## Review Checklist

- [ ] Dockerfile uses specific base image versions
- [ ] Multi-stage builds used to minimize image size
- [ ] Each RUN command is documented
- [ ] COPY commands explain what and why
- [ ] Non-root user configured for security
- [ ] Health checks defined
- [ ] Exposed ports documented
- [ ] Build arguments documented
- [ ] docker-compose.yml services are documented
- [ ] Environment variables explained
- [ ] Volume mounts purpose clear
- [ ] Network configuration documented
- [ ] Restart policies specified
- [ ] Resource limits set (if applicable)

## Best Practices

### Layer Optimization
```dockerfile
# Bad - creates many layers
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git

# Good - single layer
RUN apt-get update && \
    apt-get install -y \
      curl \
      git && \
    rm -rf /var/lib/apt/lists/*  # Clean up cache
```

### Security
```dockerfile
# Scan for vulnerabilities
# Document in README: docker scan myimage:latest

# Use minimal base images
FROM alpine:3.18

# Don't run as root
USER 1000:1000

# Use secrets for sensitive data
RUN --mount=type=secret,id=npm_token \
    NPM_TOKEN=$(cat /run/secrets/npm_token) \
    npm install
```

### Build Cache
```dockerfile
# Copy dependency files first
COPY package*.json ./
RUN npm ci

# Then copy source code
# Source changes won't invalidate dependency cache
COPY src/ ./src/
```

## Environment Documentation

### .env.example
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=appuser
DB_PASSWORD=changeme  # Change in production

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# Application Settings
NODE_ENV=development
PORT=3000
LOG_LEVEL=info
```

## Common Issues

### Issue: Large Image Sizes
Images over 1GB due to build artifacts.
**Solution**: Use multi-stage builds, clean package manager cache, use alpine images.

### Issue: Unclear Environment Variables
No documentation of required environment variables.
**Solution**: Provide .env.example and document each variable.

### Issue: Missing Health Checks
Containers without health check definitions.
**Solution**: Add HEALTHCHECK in Dockerfile or healthcheck in docker-compose.yml.