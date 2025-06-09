# Contributing to Harbor CLI

We love your input! We want to make contributing to Harbor CLI as easy and transparent as possible.

## Development Setup

### Prerequisites

- Go 1.24.3 or higher
- Make
- Docker (for testing)
- golangci-lint (for linting)

### Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/hrbcli.git
   cd hrbcli
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/pascal71/hrbcli.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Install development tools:
   ```bash
   make tools
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write clear, concise commit messages
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests

```bash
# Run unit tests
make test

# Run integration tests (requires Harbor instance)
make integration-test

# Run linter
make lint
```

### 4. Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

## Code Style

### Go Code

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Keep functions small and focused
- Add comments for exported functions

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to create project: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Command Structure

```go
var projectCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new project",
    Long:  `Create a new project in Harbor with the specified name.`,
    Example: `  # Create a public project
  hrbcli project create myproject --public

  # Create a private project with storage quota
  hrbcli project create myproject --storage-limit 10G`,
    Args: cobra.ExactArgs(1),
    RunE: createProject,
}
```

## Testing

### Unit Tests

```go
func TestCreateProject(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {
            name:    "valid project name",
            args:    []string{"myproject"},
            wantErr: false,
        },
        {
            name:    "invalid project name",
            args:    []string{"my project"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Use test fixtures
- Clean up resources after tests
- Use unique names to avoid conflicts

## Documentation

### Code Comments

```go
// CreateProject creates a new project in Harbor.
// It validates the project name and checks for duplicates before creation.
func CreateProject(name string, public bool) error {
    // Implementation
}
```

### Command Documentation

- Update help text for new commands
- Add examples for complex commands
- Update README.md for significant features

## Pull Request Process

1. Update documentation
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Create pull request with clear description
6. Link related issues

### PR Title Format

```
type: brief description

Examples:
feat: add webhook management commands
fix: handle pagination in project list
docs: update installation instructions
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Build process or auxiliary tool changes

## Reporting Issues

### Bug Reports

Include:
- Harbor CLI version
- Harbor server version
- Command that caused the issue
- Expected behavior
- Actual behavior
- Steps to reproduce

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Additional context

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Respect different viewpoints and experiences

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
