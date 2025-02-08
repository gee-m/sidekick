# Development Guidelines

> **IMPORTANT**: This document serves as both a guideline and a template collection. When adding new functionality, you should:
> 1. Review the relevant templates in this document
> 2. Copy and adapt the templates for your specific use case
> 3. Work through the checklist for your change
> 4. Document your thought process in comments or a separate design doc
>
> For LLMs: Please explicitly walk through the checklist items and explain your reasoning for each decision when proposing changes.

When developing Aimia, strict adherence to domain-driven design principles is non-negotiable. Every feature must begin with thorough domain modeling through interfaces, with all dependencies explicitly injected and abstracted. The codebase must maintain a clear separation of concerns through a rigorously organized directory structure, where domain logic, infrastructure concerns, and presentation layers are distinctly separated. All components must be defined first through their interfaces before implementation, enabling easy testing and future modifications. HTMX interactions must be handled through clean, server-side templates with clear component boundaries, and all state changes must be propagated through a well-defined event system. The application must use Go's strong type system to enforce business rules at compile time, and all error handling must be explicit and domain-aware. Every package must have a clear, single responsibility, and cross-cutting concerns must be handled through middleware and interceptors. Data access must be abstracted through repository interfaces, with all database operations isolated from business logic. All HTTP handlers must be thin adapters that delegate to domain services, and WebSocket communication must be handled through a dedicated event bus. This architecture is not optional - it's a requirement for maintaining scalability and preventing technical debt in the Aimia platform.

- Don't do import like `github.com/yourusername/appgents`, use `appgents` instead.

## Directory Structure

```
appgents/
├── cmd/                    # Application entrypoints
│   └── server/            # Main server
│       ├── main.go
│       ├── main_test.go
│       └── routes.go
├── internal/              # Internal packages
│   ├── core/             # Core services
│   │   ├── auth/         # Authentication
│   │   ├── bus/          # Event system
│   │   ├── registry/     # Appgent registry
│   │   └── theme/        # Theme management
│   ├── appgents/         # Individual appgents
│   │   ├── playtomic/
│   │   ├── schedule/
│   │   ├── dailylog/
│   │   └── oura/
│   ├── platform/         # Infrastructure
│   │   ├── database/
│   │   │   ├── migrations/  # Database migrations
│   │   │   │   ├── files/  # SQL migration files
│   │   │   │   └── migrate.go
│   │   │   └── database.go
│   │   ├── http/
│   │   └── notification/
│   └── presentation/     # UI components
├── web/                  # Web assets
│   ├── templates/        # HTML templates
│   └── static/          # Static assets
└── tests/               # Integration and E2E tests only
    ├── integration/
    └── e2e/
```

## Implementation Checklist

### When Adding a New Feature

#### 1. Domain Analysis
- [ ] Map out business rules and constraints
- [ ] Identify domain entities and relationships
- [ ] Define interfaces for required services
- [ ] Document domain events that will be emitted
- [ ] Consider error cases and validation rules

#### 2. Component Location
- [ ] Determine if feature belongs in existing appgent or needs new one
- [ ] Place code in appropriate package under `internal/`
- [ ] Follow package naming conventions
- [ ] Keep related files together

#### 3. Database Changes
- [ ] Create new migration files if needed
- [ ] Follow table naming conventions
- [ ] Include proper constraints and indexes
- [ ] Add rollback migrations
- [ ] Update schema documentation

#### 4. Interface Design
- [ ] Define clear service interfaces
- [ ] Document method contracts
- [ ] Consider error types
- [ ] Plan for future extensibility

#### 5. Implementation
- [ ] Add unit tests alongside code
- [ ] Implement validation logic
- [ ] Add proper error handling
- [ ] Set up event emission
- [ ] Document public APIs

#### 6. UI Integration
- [ ] Create/update HTMX templates
- [ ] Add WebSocket handlers if needed
- [ ] Implement proper swap strategies
- [ ] Consider theme support

### When Adding a New Appgent

#### 1. Structure Setup
```
internal/appgents/myappgent/
├── appgent.go         # Main implementation
├── appgent_test.go    # Unit tests
├── handler.go         # HTTP handlers
├── repository.go      # Data access
├── events.go          # Event definitions
└── templates/         # UI templates
    └── view.templ
```

#### 2. Required Components
```go
// internal/appgents/playtomic/interface.go
type PlaytomicService interface {
    // Core business functionality
    CreateWatch(ctx context.Context, req WatchRequest) (string, error)
    RemoveWatch(ctx context.Context, watchID string) error
    ListWatches(ctx context.Context, userID string) ([]Watch, error)

    // Real-time monitoring
    StartMonitoring(watchID string) error
    StopMonitoring(watchID string) error

    // Notifications
    NotifyAvailability(ctx context.Context, court Court) error
}

// internal/appgents/playtomic/appgent.go
type PlaytomicAppgent struct {
    service  PlaytomicService  // Uses interface instead of concrete dependencies
    bus      event.Bus
}

func (p *PlaytomicAppgent) Handle(ctx context.Context, action string, payload []byte) (templ.Component, error) {
    switch action {
    case "watch":
        var req WatchRequest
        if err := json.Unmarshal(payload, &req); err != nil {
            return nil, fmt.Errorf("invalid request payload: %w", err)
        }

        watchID, err := p.service.CreateWatch(ctx, req)
        if err != nil {
            return nil, fmt.Errorf("creating watch: %w", err)
        }

        if err := p.service.StartMonitoring(watchID); err != nil {
            return nil, fmt.Errorf("starting monitoring: %w", err)
        }

        watches, err := p.service.ListWatches(ctx, getUserID(ctx))
        if err != nil {
            return nil, fmt.Errorf("listing watches: %w", err)
        }

        return PlaytomicView(PlaytomicViewProps{
            Watches: watches,
        }), nil
    }
    return nil, fmt.Errorf("unknown action: %s", action)
}

// Service implementation
type playtomicService struct {
    db       *sql.DB
    client   *playtomic.Client
    notifier notification.Service
}

// Easy to test with mocks
func TestPlaytomicAppgent_Handle(t *testing.T) {
    mockService := &mockPlaytomicService{}
    mockBus := &mockEventBus{}

    appgent := &PlaytomicAppgent{
        service: mockService,
        bus:     mockBus,
    }

    // Test cases... in _test.go
}

// handler.go
type Handler struct {
    appgent *MyAppgent
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// repository.go
type Repository interface {
    Store(ctx context.Context, data *MyData) error
    Find(ctx context.Context, id string) (*MyData, error)
}
```

## Error Handling

### 1. Error Types
```go
// internal/core/errors/errors.go
type ErrorType string

const (
    ErrValidation   ErrorType = "VALIDATION"
    ErrNotFound     ErrorType = "NOT_FOUND"
    ErrUnauthorized ErrorType = "UNAUTHORIZED"
    ErrConflict     ErrorType = "CONFLICT"
    ErrInternal     ErrorType = "INTERNAL"
)

type Error struct {
    Type      ErrorType
    Message   string
    Code      string
    Original  error
    Context   map[string]interface{}
}
```

### 2. Error Usage Pattern
```go
// Example of proper error wrapping
func (r *Repository) FindByID(ctx context.Context, id string) (*Domain, error) {
    if err := validateID(id); err != nil {
        return nil, errors.Wrap(err, ErrValidation, "invalid id format")
    }

    domain, err := r.db.QueryRow(ctx, "SELECT * FROM domains WHERE id = $1", id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.Wrap(err, ErrNotFound, "domain not found")
        }
        return nil, errors.Wrap(err, ErrInternal, "database error")
    }

    return domain, nil
}
```

## Event System

### 1. Event Definition
```go
// internal/core/bus/event.go
type Event struct {
    Topic     string
    Payload   interface{}
    Timestamp time.Time
    UserID    string
}
```

### 2. Event Usage
```go
// Publishing events
func (s *Service) ProcessAction(ctx context.Context) error {
    // ... business logic ...

    event := Event{
        Topic:   "action.completed",
        Payload: result,
        UserID:  getUserID(ctx),
    }

    if err := s.bus.Publish(ctx, event); err != nil {
        return fmt.Errorf("publishing event: %w", err)
    }

    return nil
}

// Subscribing to events
func (a *Appgent) Subscribe(bus EventBus) {
    bus.Subscribe("action.completed", a.handleAction)
}
```

## Database Migrations

### 1. Migration Structure
```sql
-- migrations/YYYYMMDDHHMMSS_name.up.sql
BEGIN;

CREATE TABLE example (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Core Fields
    name TEXT NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    -- Foreign Keys
    user_id UUID NOT NULL REFERENCES users(id),
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Constraints
    CONSTRAINT example_name_unique UNIQUE (name)
);

-- Indexes
CREATE INDEX example_user_id_idx ON example(user_id);

-- Triggers
CREATE TRIGGER update_example_updated_at
    BEFORE UPDATE ON example
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

COMMIT;
```

### 2. Migration Guidelines
- One change per migration
- Always include rollback migration
- Use transactions
- Add proper constraints
- Consider indexes
- Document changes

## UI Integration

### Template Organization
```
web/templates/               # Base templates and layouts
├── layout/
│   ├── base.templ          # Base layout template
│   └── theme.templ         # Theme provider
└── components/             # Shared components
    ├── navbar.templ
    └── forms.templ

internal/appgents/          # Appgent-specific templates
└── playtomic/
    └── templates/
        └── view.templ      # Appgent-specific views
```

### 1. HTMX Templates
```html
// internal/appgents/playtomic/templates/view.templ
templ PlaytomicView(props ViewProps) {
    <form
        hx-post="/api/action"
        hx-swap="outerHTML"
        class="space-y-4"
    >
        <div class="form-group">
            <label class="block text-sm font-medium">
                {props.Label}
            </label>
            <input
                type="text"
                name={props.Name}
                class="mt-1 block w-full rounded-md"
                required
            />
        </div>
        <button type="submit" class="btn-primary">
            Submit
        </button>
    </form>
}
```

### 2. WebSocket Integration
```go
type WSManager struct {
    mu       sync.RWMutex
    clients  map[string]map[*Client]bool
}

func (m *WSManager) SendToUser(userID string, event interface{}) {
    m.mu.RLock()
    clients := m.clients[userID]
    m.mu.RUnlock()

    data, err := json.Marshal(event)
    if err != nil {
        log.Printf("error marshaling event: %v", err)
        return
    }

    for client := range clients {
        client.send <- data
    }
}
```

## Testing Strategy

### 1. Unit Tests
Co-located with source files:
```go
// internal/core/auth/service_test.go
func TestSignUp(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        password string
        wantErr bool
    }{
        // Test cases
    }
    // Implementation
}
```

### 2. Integration Tests
```go
// tests/integration/db/repository_test.go
func TestRepository_Integration(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(t)
    repo := database.NewRepository(db)
    // Test implementation
}
```

### 3. E2E Tests
```go
// tests/e2e/scenarios/flow_test.go
func TestCompleteFlow(t *testing.T) {
    app := helpers.NewTestApplication(t)
    client := helpers.NewTestClient(t)
    // Test implementation
}
```

## Logging Strategy

```go
// internal/platform/logging/logger.go
type LogEntry struct {
    Timestamp   time.Time
    Level       string
    Service     string
    Operation   string
    TraceID     string
    RequestID   string
    UserID      string
    Message     string
    Data        interface{}
    Error       *ErrorInfo
}

// Usage
logger.Info("processing_action", LogEntry{
    Operation: "action_name",
    Data: map[string]interface{}{
        "input": input,
        "result": result,
    },
})
```

## Intelligence Layer Integration

### Knowledge Base Integration
```go
// internal/core/knowledge/contributor.go
type KnowledgeContributor interface {
    // ContributeKnowledge adds appgent-specific data to the knowledge base
    ContributeKnowledge(ctx context.Context, data *Knowledge) error

    // GetRelevantKnowledge retrieves knowledge relevant to this appgent
    GetRelevantKnowledge(ctx context.Context, query *KnowledgeQuery) (*Knowledge, error)
}

// Example implementation
type PhotoMemoryContributor struct {
    kb        knowledge.Base
    inference inference.Engine
}

func (p *PhotoMemoryContributor) ContributeKnowledge(ctx context.Context, data *Knowledge) error {
    // 1. Transform photo analysis into knowledge format
    knowledge := transformToKnowledge(data)

    // 2. Add metadata
    knowledge.Source = "photo_memory"
    knowledge.Confidence = calculateConfidence(data)
    knowledge.Timestamp = time.Now()

    // 3. Store in knowledge base
    return p.kb.Store(ctx, knowledge)
}
```

### AI Orchestration
```go
// internal/core/ai/orchestrator.go
type Orchestrator interface {
    // ProcessInsight handles new insights from appgents
    ProcessInsight(ctx context.Context, insight *Insight) error

    // GetRecommendations gets AI recommendations for an appgent
    GetRecommendations(ctx context.Context, appgentID string) ([]Recommendation, error)

    // LearnPreference records user preferences for learning
    LearnPreference(ctx context.Context, preference *UserPreference) error
}

// Usage in appgent
func (a *Appgent) handleUserAction(ctx context.Context, action *Action) error {
    // Record user preference
    preference := &UserPreference{
        UserID: action.UserID,
        Action: action.Type,
        Context: action.Context,
    }

    if err := a.orchestrator.LearnPreference(ctx, preference); err != nil {
        log.Printf("failed to learn preference: %v", err)
        // Continue processing - learning errors shouldn't block user actions
    }

    // Process action
    return a.processAction(ctx, action)
}
```



### Cross-Appgent Communication

> Cross-appgent communication follows a request/response pattern using the event bus. Here's how it works:

1. Each appgent subscribes to its own request topics during initialization
   e.g., LocationMemory subscribes to "location_memory.context.request"

2. When an appgent needs data from another appgent:
   - It generates a unique RequestID
   - Creates a response channel for this specific request
   - Publishes a request event to the target appgent's topic
   - Waits for response on the channel with a timeout

3. The target appgent:
   - Receives the request event
   - Processes it
   - Publishes a response event to its response topic
   - The response includes the original RequestID for correlation

4. The requesting appgent:
   - Receives the response event (it subscribed during init)
   - Matches the RequestID to find the waiting channel
   - Sends the response to the channel
   - The original caller receives the response

Example Topics:
- Request:  "{target_appgent}.{action}.request"
  e.g., "location_memory.context.request"
- Response: "{target_appgent}.{action}.response"
  e.g., "location_memory.context.response"

This pattern allows for:
- Async communication with sync-like usage
- Proper request/response correlation
- Timeout handling
- System-wide monitoring
- Easy testing
*/

```go
// internal/core/bus/events.go
type RequestEvent struct {
    RequestID string                 // UUID for correlating request/response
    Source    string                 // Source appgent ID
    Target    string                 // Target appgent ID
    Type      string                 // Request type
    Data      interface{}            // Request payload
    Metadata  map[string]interface{} // Additional context
}

type ResponseEvent struct {
    RequestID string                 // Matching UUID from request
    Source    string                 // Source appgent ID
    Target    string                 // Target appgent ID
    Type      string                 // Response type
    Data      interface{}            // Response payload
    Error     error                  // Error if request failed
}

// Example: LocationMemory appgent handling requests
type LocationMemory struct {
    bus event.Bus
}

func (l *LocationMemory) Init() error {
    // Subscribe to incoming context requests
    return l.bus.Subscribe("location_memory.context.request", l.handleContextRequest)
}

func (l *LocationMemory) handleContextRequest(ctx context.Context, event interface{}) error {
    req, ok := event.(*RequestEvent)
    if !ok {
        return fmt.Errorf("invalid request event type: %T", event)
    }

    // Process the request
    locationContext, err := l.getLocationContext(ctx, req.Data)

    // Build response
    response := &ResponseEvent{
        RequestID: req.RequestID,  // Correlate with request
        Source:    "location_memory",
        Target:    req.Source,
        Type:      "context_response",
        Data:      locationContext,
        Error:     err,
    }

    // Publish response
    return l.bus.Publish(ctx, "location_memory.context.response", response)
}

// Example: PhotoMemory appgent making requests
type PhotoMemory struct {
    bus        event.Bus
    responses  map[string]chan *ResponseEvent
    mu         sync.RWMutex
}

func (p *PhotoMemory) Init() error {
    // Subscribe to location context responses
    return p.bus.Subscribe("location_memory.context.response", p.handleLocationResponse)
}

func (p *PhotoMemory) analyzeLocation(ctx context.Context, photo *Photo) error {
    requestID := uuid.New().String()

    // Create response channel
    respChan := make(chan *ResponseEvent, 1)
    p.mu.Lock()
    p.responses[requestID] = respChan
    p.mu.Unlock()

    // Clean up when done
    defer func() {
        p.mu.Lock()
        delete(p.responses, requestID)
        p.mu.Unlock()
        close(respChan)
    }()

    // Publish request event
    request := &RequestEvent{
        RequestID: requestID,
        Source:    "photo_memory",
        Target:    "location_memory",
        Type:      "context_request",
        Data:      photo.Location,
        Metadata: map[string]interface{}{
            "photo_id": photo.ID,
            "timestamp": photo.Timestamp,
        },
    }

    if err := p.bus.Publish(ctx, "location_memory.context.request", request); err != nil {
        return fmt.Errorf("publishing location request: %w", err)
    }

    // Wait for response with timeout
    select {
    case resp := <-respChan:
        if resp.Error != nil {
            return fmt.Errorf("location request failed: %w", resp.Error)
        }
        return p.enhanceAnalysis(photo, resp.Data)
    case <-ctx.Done():
        return fmt.Errorf("location request timeout: %w", ctx.Err())
    }
}

func (p *PhotoMemory) handleLocationResponse(ctx context.Context, event interface{}) error {
    resp, ok := event.(*ResponseEvent)
    if !ok {
        return fmt.Errorf("invalid response event type: %T", event)
    }

    p.mu.RLock()
    respChan, exists := p.responses[resp.RequestID]
    p.mu.RUnlock()

    if !exists {
        // Request might have timed out, log and ignore
        log.Printf("received response for unknown request: %s", resp.RequestID)
        return nil
    }

    // Send response to waiting request
    select {
    case respChan <- resp:
        return nil
    default:
        return fmt.Errorf("response channel full for request: %s", resp.RequestID)
    }
}
```

This pattern:
- Uses the central event bus for all inter-appgent communication
- Maintains request/response correlation through UUIDs
- Handles timeouts and cleanup properly
- Provides context through metadata
- Allows for async communication while providing sync-like usage
- Enables system-wide event monitoring and logging

### External Service Integration
```go
// internal/platform/external/client.go
type ServiceClient interface {
    // Execute handles API calls with retry and rate limiting
    Execute(ctx context.Context, req *Request) (*Response, error)

    // RefreshAuth refreshes service authentication
    RefreshAuth(ctx context.Context) error

    // CheckQuota verifies request against quota
    CheckQuota(ctx context.Context, req *Request) error
}

// Example implementation
type PlaytomicClient struct {
    httpClient *http.Client
    rateLimit  *rate.Limiter
    quotaMgr   *QuotaManager
}

func (c *PlaytomicClient) Execute(ctx context.Context, req *Request) (*Response, error) {
    // 1. Check quota
    if err := c.quotaMgr.CheckQuota(ctx, req); err != nil {
        return nil, fmt.Errorf("quota exceeded: %w", err)
    }

    // 2. Apply rate limiting
    if err := c.rateLimit.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit wait: %w", err)
    }

    // 3. Execute with retry
    var resp *Response
    err := retry.Do(func() error {
        var err error
        resp, err = c.executeRequest(ctx, req)
        return err
    }, retry.Attempts(3))

    return resp, err
}
```

### Real-time Processing
```go
// internal/core/processing/scheduler.go
type ProcessingScheduler interface {
    // ScheduleTask schedules a processing task
    ScheduleTask(ctx context.Context, task *Task) error

    // ExecuteImmediately runs a task immediately
    ExecuteImmediately(ctx context.Context, task *Task) error
}

// Usage example
func (p *PhotoMemory) handleNewPhoto(ctx context.Context, photo *Photo) error {
    // Quick analysis for immediate feedback
    basicAnalysis := p.quickAnalyze(photo)

    // Schedule deep analysis
    task := &Task{
        Type: "deep_photo_analysis",
        Data: photo,
        Priority: PriorityBackground,
    }

    if err := p.scheduler.ScheduleTask(ctx, task); err != nil {
        log.Printf("failed to schedule analysis: %v", err)
        // Continue with basic analysis
    }

    return p.updateUI(ctx, basicAnalysis)
}
```

### Personal Data Management
```go
// internal/core/privacy/manager.go
type PrivacyManager interface {
    // FilterSensitiveData removes sensitive data before storage
    FilterSensitiveData(ctx context.Context, data interface{}) (interface{}, error)

    // CheckDataAccess verifies access permissions
    CheckDataAccess(ctx context.Context, userID string, data interface{}) error

    // CleanupExpiredData removes expired personal data
    CleanupExpiredData(ctx context.Context) error
}

// Usage in appgent
func (a *Appgent) storeUserData(ctx context.Context, data interface{}) error {
    // 1. Filter sensitive data
    filtered, err := a.privacy.FilterSensitiveData(ctx, data)
    if err != nil {
        return fmt.Errorf("filtering sensitive data: %w", err)
    }

    // 2. Store filtered data
    return a.repo.Store(ctx, filtered)
}
```

## UI Organization

```
web/
├── templates/           # All template files
│   ├── layout/         # Base layouts and shared structures
│   │   └── base.templ  # Base HTML template
│   ├── components/     # Reusable UI components
│   │   ├── forms/
│   │   └── nav/
│   └── auth/           # Auth-specific templates
│       └── login.templ
└── static/             # Static assets
    ├── css/           # Custom CSS files
    ├── js/            # Custom JavaScript files
    └── images/        # Image assets
```

### Template Guidelines

1. Use Tailwind CSS for styling
   - Prefer utility classes over custom CSS
   - Use consistent color schemes and spacing
   - Follow responsive design patterns

2. Template Organization
   - Keep templates close to their handlers
   - Use composition over inheritance
   - Break down complex templates into components

3. HTMX Usage
   - Use for dynamic updates
   - Follow progressive enhancement
   - Keep endpoints focused and simple

4. JavaScript
   - Minimize custom JS
   - Use HTMX where possible
   - Keep complex interactions isolated

Remember:
1. Always check actual implementation files
2. Follow existing patterns
3. Document your changes
4. Consider all components (domain, database, events, UI)
5. Write tests alongside code
6. Use proper error handling
7. Consider performance implications
8. Respect privacy and data isolation
9. Consider ML model caching
10. Handle cross-appgent dependencies
