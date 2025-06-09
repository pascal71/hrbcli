# Harbor CLI Design Document

## Overview

Harbor CLI (hrbcli) is a command-line interface for Harbor container registry that provides comprehensive access to Harbor's functionality through its REST API.

## Architecture

### Core Components

```
┌─────────────────┐
│   CLI Layer     │  Commands, flags, help text
├─────────────────┤
│  Business Logic │  Harbor operations, validation
├─────────────────┤
│   API Client    │  HTTP client, auth, retries
├─────────────────┤
│  Configuration  │  Config file, env vars, flags
├─────────────────┤
│     Output      │  Formatters (table, JSON, YAML)
└─────────────────┘
```

### Package Structure

- **cmd/**: Command definitions using Cobra
- **pkg/api/**: Low-level API client
- **pkg/harbor/**: High-level Harbor operations
- **pkg/config/**: Configuration management
- **pkg/output/**: Output formatting
- **pkg/utils/**: Shared utilities

## Design Principles

### 1. Separation of Concerns

- Commands (cmd) handle CLI interaction
- Business logic (pkg/harbor) handles Harbor operations
- API client (pkg/api) handles HTTP communication

### 2. Consistent Command Structure

All commands follow the pattern:
```
hrbcli <resource> <action> [arguments] [flags]
```

Examples:
- `hrbcli project create myproject`
- `hrbcli repo list myproject`
- `hrbcli user delete john`

### 3. Configuration Hierarchy

Configuration precedence (highest to lowest):
1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

### 4. Error Handling

- Consistent error messages
- Appropriate exit codes
- Debug mode for troubleshooting
- Human-readable error messages

### 5. Output Flexibility

Support multiple output formats:
- **Table**: Human-readable default
- **JSON**: Machine-readable, parseable
- **YAML**: Human and machine readable

## API Client Design

### Authentication

```go
type Client struct {
    BaseURL    string
    Username   string
    Password   string
    HTTPClient *http.Client
}
```

### Retry Logic

- Exponential backoff for transient errors
- Configurable retry count
- Respect rate limits

### Error Handling

```go
type APIError struct {
    Code    int
    Message string
    Details map[string]interface{}
}
```

## Command Structure

### Resource-Based Organization

```
hrbcli
├── project
│   ├── list
│   ├── create
│   ├── update
│   └── delete
├── repo
│   ├── list
│   ├── delete
│   └── tags
├── user
│   ├── list
│   ├── create
│   ├── update
│   └── delete
└── system
    ├── info
    ├── health
    └── gc
```

### Common Flags

Global flags available to all commands:
- `--harbor-url`: Harbor server URL
- `--username`: Username for authentication
- `--password`: Password for authentication
- `--output/-o`: Output format
- `--insecure`: Skip TLS verification
- `--debug`: Enable debug output

## Security Considerations

### Credential Storage

- Never store passwords in plain text
- Support reading password from environment
- Optional integration with system keychain
- Secure configuration file permissions

### TLS/SSL

- Verify certificates by default
- `--insecure` flag for testing only
- Support custom CA certificates

## Testing Strategy

### Unit Tests

- Test individual functions
- Mock external dependencies
- Achieve >80% code coverage

### Integration Tests

- Test against real Harbor instance
- Use Docker for test environment
- Test all API endpoints

### End-to-End Tests

- Test complete workflows
- Verify command output
- Test error scenarios

## Performance Considerations

### Pagination

- Handle large result sets
- Implement pagination for list operations
- Progress indicators for long operations

### Caching

- Cache project/repo information
- Invalidate cache on modifications
- Optional offline mode

## Extensibility

### Plugin System (Future)

- Allow custom commands
- Plugin discovery mechanism
- Plugin API versioning

### Custom Output Formats

- Template-based output
- Custom formatters
- Export capabilities

## Release Process

### Versioning

- Semantic versioning (MAJOR.MINOR.PATCH)
- Git tags for releases
- Changelog maintenance

### Distribution

- Multi-architecture binaries
- Container images
- Package managers (Homebrew, APT, YUM)

### Compatibility

- Support Harbor API v2.0+
- Graceful degradation for older versions
- Version detection and warnings
