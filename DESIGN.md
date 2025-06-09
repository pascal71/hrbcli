# Harbor CLI Project Structure

```
hrbcli/
├── cmd/
│   ├── hrbcli/
│   │   └── main.go              # Main entry point
│   ├── root.go                  # Root command setup
│   ├── project.go               # Project-related commands
│   ├── repository.go            # Repository-related commands
│   ├── user.go                  # User management commands
│   ├── system.go                # System administration commands
│   ├── config.go                # Configuration management commands
│   ├── artifact.go              # Artifact management commands
│   ├── replication.go           # Replication policy commands
│   ├── scanner.go               # Vulnerability scanner commands
│   ├── webhook.go               # Webhook management commands
│   └── version.go               # Version command
├── pkg/
│   ├── api/
│   │   ├── client.go            # Harbor API client
│   │   ├── auth.go              # Authentication handling
│   │   ├── errors.go            # API error handling
│   │   └── types.go             # API request/response types
│   ├── harbor/
│   │   ├── project.go           # Project API operations
│   │   ├── repository.go        # Repository API operations
│   │   ├── user.go              # User API operations
│   │   ├── system.go            # System API operations
│   │   ├── artifact.go          # Artifact API operations
│   │   ├── replication.go       # Replication API operations
│   │   ├── scanner.go           # Scanner API operations
│   │   └── webhook.go           # Webhook API operations
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── validate.go          # Configuration validation
│   ├── output/
│   │   ├── formatter.go         # Output formatting (table, json, yaml)
│   │   ├── printer.go           # Output printing utilities
│   │   └── colors.go            # Terminal color support
│   └── utils/
│       ├── prompt.go            # Interactive prompts
│       ├── spinner.go           # Progress indicators
│       └── helpers.go           # General utility functions
├── internal/
│   └── version/
│       └── version.go           # Version information
├── scripts/
│   ├── build.sh                 # Multi-arch build script
│   ├── install.sh               # Installation script
│   └── completions.sh           # Shell completion generation
├── docs/
│   ├── README.md                # Main documentation
│   ├── COMMANDS.md              # Command reference
│   ├── EXAMPLES.md              # Usage examples
│   └── CONTRIBUTING.md          # Contribution guidelines
├── test/
│   ├── integration/             # Integration tests
│   ├── e2e/                     # End-to-end tests
│   └── fixtures/                # Test data
├── .github/
│   └── workflows/
│       ├── build.yml            # CI/CD pipeline
│       ├── release.yml          # Release automation
│       └── test.yml             # Test automation
├── Dockerfile                   # Container build
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── .goreleaser.yml              # GoReleaser configuration
├── .gitignore
└── LICENSE
```

## Directory Explanations

### `/cmd`
Contains all command definitions using Cobra. Each file represents a major command group:
- **hrbcli/main.go**: Application entry point, minimal logic
- **root.go**: Root command setup, global flags, initialization
- **project.go**: Project management commands (create, list, delete, update)
- **repository.go**: Repository operations (list, delete, scan, retag)
- **user.go**: User management (create, delete, update, set-admin)
- **system.go**: System administration (health, info, gc, quotas)
- **artifact.go**: Artifact operations (list, delete, copy, scan)
- **replication.go**: Replication policies and executions
- **scanner.go**: Vulnerability scanning operations
- **webhook.go**: Webhook policy management
- **config.go**: CLI configuration management
- **version.go**: Version information display

### `/pkg`
Reusable packages that could theoretically be imported by other projects:
- **api/**: Low-level API client implementation
  - `client.go`: HTTP client with retry logic, timeout handling
  - `auth.go`: Basic auth, token management
  - `errors.go`: API error types and handling
  - `types.go`: Shared data structures
- **harbor/**: High-level Harbor API operations organized by resource
- **config/**: Configuration file handling and validation
- **output/**: Output formatting for different formats (table, JSON, YAML)
- **utils/**: Common utilities like prompts, progress bars

### `/internal`
Private application code that shouldn't be imported by other projects:
- **version/**: Build-time version information injection

### `/scripts`
Build and deployment scripts:
- **build.sh**: Multi-architecture build script using Go cross-compilation
- **install.sh**: Installation script for different platforms
- **completions.sh**: Generate shell completions for bash, zsh, fish

### `/docs`
Comprehensive documentation:
- **README.md**: Getting started, installation, basic usage
- **COMMANDS.md**: Detailed command reference with examples
- **EXAMPLES.md**: Real-world usage scenarios
- **CONTRIBUTING.md**: Development setup, coding standards

### `/test`
Testing infrastructure:
- **integration/**: Tests against real Harbor instance
- **e2e/**: End-to-end workflow tests
- **fixtures/**: Test data and mock responses

## Key Design Principles

1. **Separation of Concerns**: Commands (cmd) are separate from business logic (pkg)
2. **Testability**: Interfaces for all major components enable easy mocking
3. **Extensibility**: New commands can be added without modifying existing code
4. **Configuration**: Support for config file, environment variables, and flags
5. **Output Flexibility**: Support multiple output formats (table, JSON, YAML)
6. **Error Handling**: Consistent error messages and exit codes
7. **Interactive Mode**: Prompts for missing required information
8. **Offline Help**: Comprehensive help text for all commands

## Configuration Hierarchy

1. Command-line flags (highest priority)
2. Environment variables (HARBOR_*)
3. Configuration file (~/.hrbcli.yaml)
4. Default values (lowest priority)

## Example Configuration File

```yaml
# ~/.hrbcli.yaml
harbor_url: https://harbor.example.com
username: admin
# Password can be stored in HARBOR_PASSWORD env var for security
insecure: false
output_format: table
default_project: library
```

## Build Targets

The Makefile should support:
- `make build`: Build for current platform
- `make build-all`: Build for all platforms
- `make test`: Run unit tests
- `make integration-test`: Run integration tests
- `make install`: Install locally
- `make clean`: Clean build artifacts
- `make release`: Create release artifacts
