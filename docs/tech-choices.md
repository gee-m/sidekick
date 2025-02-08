# Technical Choices

## Core Stack

### Backend
- **Framework**: Chi (github.com/go-chi/chi/v5)
  - Lightweight HTTP router with great middleware support
  - Clean WebSocket integration
  - Easy request context handling

### Frontend
- **Templating**: templ (github.com/a-h/templ)
  - Type-safe templating
  - Native HTMX integration
  - Component-based approach
  - Hot reload support

### Database
- **PostgreSQL** with sqlx
  - Session storage
  - Event persistence
  - Domain data
  - Schema migrations with golang-migrate

## Authentication

Session-based authentication stored in Postgres:
- Sessions table with:
  - session_id (UUID)
  - user_id (UUID)
  - created_at (timestamp)
  - expires_at (timestamp)
  - metadata (JSONB)
- Cookie-based session tracking
- Middleware for auth checks
- Automatic cleanup of expired sessions

## Event System

In-memory event bus with interfaces for future scaling:
- Topic-based pub/sub
- Sync and async event handling
- Support for request/response patterns
- Interface design allows future upgrade to NATS/Redis

## WebSocket Management

Clean WebSocket implementation:
- Connection pooling per user
- Heartbeat monitoring
- Automatic reconnection
- Message queuing
- Room-based broadcasting

## Development Environment

### Local Setup
- Docker Compose:
  - API service
  - PostgreSQL
  - Debug tools
- Air for hot reload
- Make commands for common tasks

### Testing
- Standard `testing` package
- Integration tests with testcontainers
- E2E tests with fully composed environment

## Observability

- **Logging**: Zap (uber-go/zap)
  - Structured logging
  - Development/Production modes
  - Log rotation

- **Metrics & Tracing**: OpenTelemetry
  - Basic metrics setup
  - Tracing support when needed
  - Local development insights

## Configuration

Environment-based configuration:
- Typed config structs
- Environment variable parsing
- Validation on startup
- Development defaults

## Deployment

Render.com deployment:
- Docker-based deploys
- Auto-scaling support
- Built-in SSL
- Database management
- Zero-downtime updates

## Project Structure

Monorepo with clear boundaries:
- cmd/ - Entry points
- internal/ - Private packages
- pkg/ - Public packages
- web/ - Frontend assets
- docs/ - Documentation
- scripts/ - Build/deploy scripts

## Future Considerations

- Migration to distributed event system
- Caching layer (Redis)
- Service mesh for multi-region
- CDN integration
- Rate limiting