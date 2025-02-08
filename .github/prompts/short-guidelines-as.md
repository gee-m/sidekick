# Short Development Guidelines

> **IMPORTANT**: Use this as a guideline and a template collection. Adapt the templates for your use case, work through the checklist, and document your thought process. Refer to `guidelines.md` for detailed instructions.


## Core Requirements

1. **Location**: Place all new code under appropriate directories in `/internal/`
2. **Architecture**: Follow Clean Architecture principles
3. **Testing**: Co-locate tests with source files
4. **Documentation**: Document public APIs and complex logic
5. **Error Handling**: Use proper error wrapping and types

## Implementation Checklist

### For New Features

1. **Analysis**
   - [ ] Map business rules
   - [ ] Define interfaces
   - [ ] Document events

2. **Code Organization**
   - [ ] Place in correct package
   - [ ] Follow naming conventions
   - [ ] Keep related files together

3. **Quality**
   - [ ] Add unit tests
   - [ ] Implement validation
   - [ ] Handle errors properly
   - [ ] Set up event emission

### For New Appgents

1. **Required Structure**
```
internal/appgents/myappgent/
├── appgent.go       # Main implementation
├── handler.go       # HTTP handlers
├── repository.go    # Data access
└── templates/       # UI templates
```

2. **Core Components**
   - Service interface
   - Event definitions
   - Repository interface
   - HTTP handlers
   - UI templates

## Project Structure
```
.
├── cmd/                     # Application entrypoints
├── internal/               # Private application code
│   ├── appgents/          # Individual appgents
│   ├── core/              # Core business logic
│   ├── platform/          # Infrastructure code
│   └── presentation/      # UI components
├── pkg/                   # Public libraries
├── web/                   # Web assets
│   ├── static/           # Static files
│   └── templates/        # HTML templates
├── tests/                # Integration & E2E tests
└── tools/                # Development tools
```

## Coding Standards

1. **Go Idioms**
   - Use interfaces for abstraction
   - Keep interfaces small
   - Handle errors explicitly
   - Prefer composition over inheritance

2. **Package Design**
   - One package per directory
   - Package name matches directory
   - Internal packages under `internal/`
   - Public packages under `pkg/`

3. **Naming**
   - Use meaningful package names
   - Capitalize exported names
   - Use camelCase for variable names
   - Use PascalCase for types

4. **Comments**
   - Document all exported items
   - Use complete sentences
   - Start with the name of thing described
   - Include examples for complex functions

## Tooling

1. **Required Tools**
   - Go 1.23+
   - Air (live reload)
   - Templ (template compiler)
   - Sqlc (SQL generator)
   - Golangci-lint

2. **Development Commands**
```bash
make dev      # Run development server
make test     # Run tests
make lint     # Run linters
make build    # Build production binary
make migrate  # Run database migrations
```

## Best Practices

1. **Event System**
   - Use for cross-appgent communication
   - Include proper correlation IDs
   - Handle timeouts

2. **Error Handling**
   - Use defined error types
   - Wrap errors with context
   - Proper error recovery

3. **Testing**
   - Unit tests alongside code
   - Integration tests in `tests/`
   - Use interfaces for mocking

4. **Data Privacy**
   - Filter sensitive data
   - Check access permissions
   - Handle data expiration
