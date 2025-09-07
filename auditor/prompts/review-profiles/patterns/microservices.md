---
type: pattern
name: microservices
keywords: [microservice, service, distributed, messaging, event-driven]
priority: 5
matches:
  area_type: [service, microservice, worker]
  name_contains: [service, worker, processor, handler]
---

# Microservices Documentation Pattern

## Documentation Requirements

### Service Overview
- **Service Purpose**: Single responsibility and boundaries
- **Dependencies**: Other services this service depends on
- **Consumers**: Services that depend on this service
- **SLA**: Availability, latency, throughput guarantees
- **Deployment**: How and where the service is deployed

### Communication Patterns
- **Synchronous APIs**: REST, gRPC endpoints
- **Asynchronous Messaging**: Events published and consumed
- **Data Contracts**: Schemas for all messages
- **Error Handling**: Retry policies, dead letter queues
- **Circuit Breakers**: Failure handling patterns

### Data Management
- **Data Store**: Database or storage used
- **Data Ownership**: What data this service owns
- **Data Consistency**: Eventual consistency patterns
- **Caching Strategy**: What is cached and TTLs

## Review Checklist

- [ ] Service purpose follows single responsibility principle
- [ ] API contracts are versioned
- [ ] Event schemas are documented
- [ ] Error handling and retry policies documented
- [ ] Health check endpoints provided
- [ ] Metrics and monitoring documented
- [ ] Deployment configuration documented
- [ ] Security boundaries defined
- [ ] Data ownership is clear
- [ ] Integration patterns documented

## Best Practices

### Service Design
- Keep services loosely coupled
- Design for failure (timeouts, retries, circuit breakers)
- Implement health checks and readiness probes
- Use correlation IDs for distributed tracing

### Documentation Standards
- Document service boundaries clearly
- Maintain API and event catalogs
- Include sequence diagrams for complex flows
- Document deployment topology

### Communication
- Prefer asynchronous communication
- Use API gateways for external access
- Implement service discovery
- Document message ordering guarantees

## Common Issues

### Issue: Unclear Service Boundaries
Services have overlapping responsibilities.
**Solution**: Document clear bounded contexts and responsibilities.

### Issue: Missing Error Handling
No documentation on failure scenarios.
**Solution**: Document failure modes, retry policies, and compensation patterns.

### Issue: Undocumented Dependencies
Hidden dependencies between services.
**Solution**: Maintain dependency graph and document all interactions.