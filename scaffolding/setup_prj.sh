#!/bin/bash

# Harbor CLI Project Setup Script
# This script creates the complete project structure for hrbcli

set -e

PROJECT_NAME="hrbcli"
MODULE_NAME="github.com/pascal71/hrbcli"

echo "🚀 Setting up Harbor CLI project structure..."

# Create root directory if not exists
if [ ! -d "$PROJECT_NAME" ]; then
	mkdir "$PROJECT_NAME"
fi

cd "$PROJECT_NAME"

# Create directory structure
echo "📁 Creating directory structure..."

# Main directories
mkdir -p cmd/hrbcli
mkdir -p pkg/{api,harbor,config,output,utils}
mkdir -p internal/version
mkdir -p scripts
mkdir -p docs/examples
mkdir -p test/{integration,e2e,fixtures}
mkdir -p .github/workflows

# Command directories
mkdir -p cmd

echo "📝 Creating README.md..."
cat >README.md <<'EOF'
# Harbor CLI (hrbcli)

A powerful command-line interface for [Harbor](https://goharbor.io/) container registry, written in Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/pascal71/hrbcli)](https://goreportcard.com/report/github.com/pascal71/hrbcli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Features

- 🚀 **Full Harbor API Coverage** - Manage projects, repositories, users, replications, and more
- 🔧 **Multiple Output Formats** - Table, JSON, and YAML output support
- 🔐 **Secure Authentication** - Support for basic auth and credential storage
- 📦 **Multi-Architecture** - Binaries for Linux, macOS, and Windows (amd64/arm64)
- 🎨 **Interactive Mode** - Prompts for missing required information
- 📚 **Comprehensive Documentation** - Built-in help for all commands
- 🔄 **Shell Completions** - Bash, Zsh, and Fish shell completions

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew tap pascal71/hrbcli
brew install hrbcli
```

### Using Go

```bash
go install github.com/pascal71/hrbcli/cmd/hrbcli@latest
```

### Download Binary

Download the latest release from the [releases page](https://github.com/pascal71/hrbcli/releases).

### Build from Source

```bash
git clone https://github.com/pascal71/hrbcli.git
cd hrbcli
make build
```

## Quick Start

### Configure Harbor Connection

```bash
# Interactive configuration
hrbcli config init

# Or set directly
hrbcli config set harbor_url https://harbor.example.com
hrbcli config set username admin
```

### Basic Commands

```bash
# List all projects
hrbcli project list

# Create a new project
hrbcli project create myproject --public

# List repositories in a project
hrbcli repo list myproject

# List tags for a repository
hrbcli repo tags myproject/myapp

# Delete a repository
hrbcli repo delete myproject/myapp:v1.0.0

# Get system information
hrbcli system info
```

## Configuration

Harbor CLI can be configured through:

1. **Configuration file** (`~/.hrbcli.yaml`)
2. **Environment variables** (prefixed with `HARBOR_`)
3. **Command-line flags**

### Configuration File Example

```yaml
harbor_url: https://harbor.example.com
username: admin
output_format: table
insecure: false
```

### Environment Variables

```bash
export HARBOR_URL=https://harbor.example.com
export HARBOR_USERNAME=admin
export HARBOR_PASSWORD=secretpassword
```

## Documentation

- [Command Reference](docs/COMMANDS.md) - Detailed documentation for all commands
- [Examples](docs/EXAMPLES.md) - Real-world usage examples
- [Design Document](docs/DESIGN.md) - Architecture and design decisions
- [Contributing](docs/CONTRIBUTING.md) - How to contribute to the project

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Harbor](https://goharbor.io/) - The cloud native registry
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
EOF

echo "📝 Creating DESIGN.md..."
cat >docs/DESIGN.md <<'EOF'
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
EOF

echo "📝 Creating CONTRIBUTING.md..."
cat >docs/CONTRIBUTING.md <<'EOF'
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
EOF

echo "📝 Creating COMMANDS.md..."
cat >docs/COMMANDS.md <<'EOF'
# Harbor CLI Command Reference

## Global Flags

These flags are available for all commands:

```
--config string        Config file (default $HOME/.hrbcli.yaml)
--debug               Enable debug output
--harbor-url string   Harbor server URL
--insecure            Skip TLS certificate verification
--no-color            Disable colored output
-o, --output string   Output format (table|json|yaml) (default "table")
--password string     Harbor password
--username string     Harbor username
```

## Commands

### Project Management

#### `hrbcli project list`

List all projects accessible to the user.

```bash
# List all projects
hrbcli project list

# List projects with details
hrbcli project list --detail

# Output as JSON
hrbcli project list -o json

# Filter by name
hrbcli project list --name-filter "prod*"
```

#### `hrbcli project create`

Create a new project.

```bash
# Create a public project
hrbcli project create myproject --public

# Create with storage quota (in bytes, or use K, M, G, T)
hrbcli project create myproject --storage-limit 10G

# Create with member limit
hrbcli project create myproject --member-limit 50
```

#### `hrbcli project delete`

Delete a project.

```bash
# Delete a project
hrbcli project delete myproject

# Force delete without confirmation
hrbcli project delete myproject --force
```

#### `hrbcli project update`

Update project settings.

```bash
# Make project public
hrbcli project update myproject --public=true

# Update storage quota
hrbcli project update myproject --storage-limit 20G

# Enable content trust
hrbcli project update myproject --enable-content-trust
```

### Repository Management

#### `hrbcli repo list`

List repositories in a project.

```bash
# List all repositories in a project
hrbcli repo list myproject

# List with details (size, tags count)
hrbcli repo list myproject --detail

# Filter by name
hrbcli repo list myproject --filter "app*"
```

#### `hrbcli repo delete`

Delete a repository.

```bash
# Delete entire repository
hrbcli repo delete myproject/myapp

# Delete specific tag
hrbcli repo delete myproject/myapp:v1.0.0

# Force delete without confirmation
hrbcli repo delete myproject/myapp --force
```

#### `hrbcli repo tags`

List tags for a repository.

```bash
# List all tags
hrbcli repo tags myproject/myapp

# List with details (size, scan status)
hrbcli repo tags myproject/myapp --detail

# Filter by tag name
hrbcli repo tags myproject/myapp --filter "v1.*"
```

### Artifact Management

#### `hrbcli artifact list`

List artifacts in a repository.

```bash
# List all artifacts
hrbcli artifact list myproject/myapp

# List with vulnerabilities
hrbcli artifact list myproject/myapp --with-scan-overview

# List with labels
hrbcli artifact list myproject/myapp --with-label
```

#### `hrbcli artifact scan`

Trigger vulnerability scan.

```bash
# Scan specific artifact
hrbcli artifact scan myproject/myapp@sha256:abc123

# Scan all artifacts in repository
hrbcli artifact scan myproject/myapp --all
```

#### `hrbcli artifact copy`

Copy artifacts between projects.

```bash
# Copy specific artifact
hrbcli artifact copy myproject/myapp:v1.0 targetproject/myapp:v1.0

# Copy with all tags
hrbcli artifact copy myproject/myapp targetproject/myapp --all-tags
```

### User Management

#### `hrbcli user list`

List Harbor users.

```bash
# List all users
hrbcli user list

# Search by username
hrbcli user list --search "john"

# List with details
hrbcli user list --detail
```

#### `hrbcli user create`

Create a new user.

```bash
# Create user
hrbcli user create john --email john@example.com --realname "John Doe"

# Create admin user
hrbcli user create admin-user --email admin@example.com --admin

# With specific password
hrbcli user create john --email john@example.com --password "secretpass"
```

#### `hrbcli user delete`

Delete a user.

```bash
# Delete user
hrbcli user delete john

# Force delete
hrbcli user delete john --force
```

### System Administration

#### `hrbcli system info`

Get system information.

```bash
# Get general info
hrbcli system info

# Include storage info
hrbcli system info --with-storage

# Output as YAML
hrbcli system info -o yaml
```

#### `hrbcli system health`

Check system health.

```bash
# Check overall health
hrbcli system health

# Check specific component
hrbcli system health --component core
hrbcli system health --component jobservice
```

#### `hrbcli system gc`

Manage garbage collection.

```bash
# Schedule garbage collection
hrbcli system gc schedule

# Get GC history
hrbcli system gc history

# Get GC job details
hrbcli system gc status <job-id>
```

### Replication

#### `hrbcli replication list`

List replication policies.

```bash
# List all policies
hrbcli replication list

# Filter by name
hrbcli replication list --name-filter "prod*"
```

#### `hrbcli replication create`

Create replication policy.

```bash
# Create push-based replication
hrbcli replication create prod-sync \
  --source myproject \
  --destination https://harbor2.example.com \
  --destination-namespace myproject

# Create pull-based replication
hrbcli replication create prod-pull \
  --source https://harbor2.example.com/myproject \
  --destination myproject \
  --direction pull
```

#### `hrbcli replication execute`

Execute replication manually.

```bash
# Execute replication
hrbcli replication execute prod-sync

# Dry run
hrbcli replication execute prod-sync --dry-run
```

### Configuration

#### `hrbcli config init`

Initialize configuration interactively.

```bash
hrbcli config init
```

#### `hrbcli config set`

Set configuration values.

```bash
# Set Harbor URL
hrbcli config set harbor_url https://harbor.example.com

# Set default output format
hrbcli config set output_format json

# Set default project
hrbcli config set default_project library
```

#### `hrbcli config get`

Get configuration values.

```bash
# Get specific value
hrbcli config get harbor_url

# Get all values
hrbcli config get
```

#### `hrbcli config list`

List all configuration.

```bash
hrbcli config list
```

## Examples

### Complete Workflow Examples

#### Setting up a new project

```bash
# Create project
hrbcli project create production --public --storage-limit 100G

# Add user to project
hrbcli project member add production john --role developer

# Create replication from dev to production
hrbcli replication create dev-to-prod \
  --source development \
  --destination production \
  --filter "name=release/*"
```

#### Repository management workflow

```bash
# List repositories
hrbcli repo list myproject

# Check tags
hrbcli repo tags myproject/webapp

# Scan for vulnerabilities
hrbcli artifact scan myproject/webapp:latest

# Copy to production
hrbcli artifact copy myproject/webapp:v1.2.3 production/webapp:v1.2.3

# Clean up old tags
hrbcli repo delete myproject/webapp:old-version
```

#### Security scanning workflow

```bash
# Scan all repositories in project
for repo in $(hrbcli repo list myproject -o json | jq -r '.[].name'); do
  hrbcli artifact scan "$repo" --all
done

# Check scan results
hrbcli artifact list myproject/webapp --with-scan-overview

# Export vulnerability report
hrbcli artifact vulnerabilities myproject/webapp:latest -o json > vulns.json
```
EOF

echo "📝 Creating EXAMPLES.md..."
cat >docs/examples/EXAMPLES.md <<'EOF'
# Harbor CLI Examples

## Authentication and Configuration

### First-time Setup

```bash
# Interactive setup
$ hrbcli config init
Harbor URL: https://harbor.example.com
Username: admin
Password: ********
Verify configuration? (y/n): y
✓ Configuration saved to ~/.hrbcli.yaml
✓ Successfully connected to Harbor

# Manual configuration
$ hrbcli config set harbor_url https://harbor.example.com
$ hrbcli config set username admin
$ export HARBOR_PASSWORD=secretpassword
```

### Multiple Harbor Instances

```bash
# Production Harbor
$ hrbcli --config ~/.hrbcli-prod.yaml project list

# Development Harbor
$ hrbcli --config ~/.hrbcli-dev.yaml project list

# Using environment variables
$ HARBOR_URL=https://dev.harbor.example.com hrbcli project list
```

## Project Management

### Create Projects with Different Settings

```bash
# Simple project creation
$ hrbcli project create my-app

# Public project with quotas
$ hrbcli project create shared-images \
    --public \
    --storage-limit 50G \
    --member-limit 100

# Project with security settings
$ hrbcli project create secure-app \
    --enable-content-trust \
    --prevent-vulnerable \
    --severity-threshold high \
    --auto-scan
```

### Manage Project Members

```bash
# Add user as developer
$ hrbcli project member add my-app john --role developer

# Add user as admin
$ hrbcli project member add my-app alice --role admin

# List project members
$ hrbcli project member list my-app

# Remove member
$ hrbcli project member remove my-app john
```

## Repository Operations

### Working with Repositories

```bash
# List all repositories with sizes
$ hrbcli repo list my-app --detail

NAME                SIZE      TAGS  PULLS  LAST MODIFIED
my-app/frontend     1.2 GB    15    1234   2 hours ago
my-app/backend      856 MB    23    5678   1 hour ago
my-app/database     2.1 GB    8     910    3 days ago

# Search repositories
$ hrbcli repo list my-app --filter "*frontend*"

# Get repository info
$ hrbcli repo info my-app/frontend
```

### Tag Management

```bash
# List tags with vulnerability status
$ hrbcli repo tags my-app/frontend --detail

TAG      SIZE    VULNERABILITIES  SCAN STATUS  CREATED
latest   245 MB  High: 2, Med: 5  Finished     2 hours ago
v2.1.0   245 MB  High: 0, Med: 3  Finished     1 day ago
v2.0.0   238 MB  High: 5, Med: 8  Finished     1 week ago

# Delete old tags
$ hrbcli repo delete my-app/frontend:v1.0.0
$ hrbcli repo delete my-app/frontend:v1.1.0

# Bulk delete tags
$ hrbcli repo tags my-app/frontend --filter "v1.*" -o json | \
    jq -r '.[] | .name' | \
    xargs -I {} hrbcli repo delete my-app/frontend:{}
```

## Security Scanning

### Scan Artifacts

```bash
# Scan single artifact
$ hrbcli artifact scan my-app/frontend:latest

# Scan all artifacts in repository
$ hrbcli artifact scan my-app/frontend --all

# Wait for scan to complete
$ hrbcli artifact scan my-app/frontend:latest --wait
```

### View Vulnerabilities

```bash
# List vulnerabilities
$ hrbcli artifact vulnerabilities my-app/frontend:latest

SEVERITY  CVE             PACKAGE         VERSION  FIXED VERSION
Critical  CVE-2021-12345  openssl         1.0.1    1.0.2
High      CVE-2021-12346  libcurl         7.1.0    7.2.0
High      CVE-2021-12347  nginx           1.18.0   1.19.0

# Export vulnerability report
$ hrbcli artifact vulnerabilities my-app/frontend:latest -o json > vulns.json

# Check if image is safe to deploy
$ hrbcli artifact vulnerabilities my-app/frontend:latest --severity high
$ echo $?  # Exit code 0 if no high/critical vulnerabilities
```

## Replication

### Set Up Replication

```bash
# Create push replication to remote Harbor
$ hrbcli replication create prod-sync \
    --source my-app \
    --destination https://prod-harbor.example.com \
    --destination-namespace production \
    --trigger "manual,schedule:0 2 * * *" \
    --filter "name:*/release-*"

# Create pull replication from Docker Hub
$ hrbcli replication create dockerhub-mirror \
    --source https://hub.docker.com \
    --source-filter "nginx,alpine,ubuntu" \
    --destination mirror \
    --direction pull \
    --trigger "schedule:0 */6 * * *"
```

### Monitor Replication

```bash
# List executions
$ hrbcli replication executions prod-sync

ID    STATUS      START TIME           END TIME
123   Succeeded   2024-01-10 02:00:00  2024-01-10 02:15:00
122   Failed      2024-01-09 02:00:00  2024-01-09 02:05:00

# Get execution details
$ hrbcli replication execution 122

# Get execution logs
$ hrbcli replication logs 122
```

## Automation Scripts

### Promote Images Through Environments

```bash
#!/bin/bash
# promote.sh - Promote image from dev to staging to prod

IMAGE=$1
VERSION=$2

# Scan in dev
echo "Scanning image in development..."
hrbcli artifact scan dev/${IMAGE}:${VERSION} --wait

# Check vulnerabilities
if hrbcli artifact vulnerabilities dev/${IMAGE}:${VERSION} --severity high; then
    echo "✓ No high/critical vulnerabilities found"
else
    echo "✗ High/critical vulnerabilities found, aborting"
    exit 1
fi

# Copy to staging
echo "Promoting to staging..."
hrbcli artifact copy dev/${IMAGE}:${VERSION} staging/${IMAGE}:${VERSION}

# Tag as staging-latest
hrbcli artifact tag staging/${IMAGE}:${VERSION} staging-latest

# After testing, promote to production
read -p "Promote to production? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    hrbcli artifact copy staging/${IMAGE}:${VERSION} prod/${IMAGE}:${VERSION}
    hrbcli artifact tag prod/${IMAGE}:${VERSION} latest
    echo "✓ Promoted to production"
fi
```

### Cleanup Old Images

```bash
#!/bin/bash
# cleanup.sh - Remove old images keeping last N tags

PROJECT=$1
REPO=$2
KEEP=5

# Get all tags sorted by creation date
TAGS=$(hrbcli repo tags ${PROJECT}/${REPO} -o json | \
    jq -r 'sort_by(.created) | reverse | .[].name')

# Keep the latest N tags, delete the rest
echo "$TAGS" | tail -n +$((KEEP+1)) | while read tag; do
    echo "Deleting ${PROJECT}/${REPO}:${tag}"
    hrbcli repo delete ${PROJECT}/${REPO}:${tag} --force
done
```

### Generate Security Report

```bash
#!/bin/bash
# security-report.sh - Generate security report for all images

OUTPUT="security-report-$(date +%Y%m%d).csv"
echo "Project,Repository,Tag,Critical,High,Medium,Low" > $OUTPUT

hrbcli project list -o json | jq -r '.[].name' | while read project; do
    hrbcli repo list $project -o json | jq -r '.[].name' | while read repo; do
        hrbcli repo tags $repo -o json | jq -r '.[].name' | while read tag; do
            VULNS=$(hrbcli artifact vulnerabilities ${repo}:${tag} -o json | \
                jq -r '[.[] | .severity] | group_by(.) | map({(.[0]): length}) | add')
            
            CRITICAL=$(echo $VULNS | jq -r '.Critical // 0')
            HIGH=$(echo $VULNS | jq -r '.High // 0')
            MEDIUM=$(echo $VULNS | jq -r '.Medium // 0')
            LOW=$(echo $VULNS | jq -r '.Low // 0')
            
            echo "$project,$repo,$tag,$CRITICAL,$HIGH,$MEDIUM,$LOW" >> $OUTPUT
        done
    done
done

echo "Report saved to $OUTPUT"
```

## Advanced Usage

### Using with CI/CD

```yaml
# .gitlab-ci.yml example
variables:
  HARBOR_URL: https://harbor.example.com
  
stages:
  - build
  - scan
  - deploy

build:
  stage: build
  script:
    - docker build -t ${HARBOR_URL}/my-app/backend:${CI_COMMIT_SHA} .
    - docker push ${HARBOR_URL}/my-app/backend:${CI_COMMIT_SHA}

scan:
  stage: scan
  script:
    - hrbcli artifact scan my-app/backend:${CI_COMMIT_SHA} --wait
    - hrbcli artifact vulnerabilities my-app/backend:${CI_COMMIT_SHA} --severity high
  allow_failure: false

deploy:
  stage: deploy
  script:
    - hrbcli artifact tag my-app/backend:${CI_COMMIT_SHA} latest
    - kubectl set image deployment/backend backend=${HARBOR_URL}/my-app/backend:latest
  only:
    - main
```

### JSON Processing with jq

```bash
# Get total storage used by project
$ hrbcli project get my-app -o json | jq '.current_usage.storage'

# List projects over quota
$ hrbcli project list -o json | \
    jq '.[] | select(.current_usage.storage > .quota.storage) | .name'

# Find images without recent scans
$ hrbcli repo list my-app -o json | \
    jq '.[] | select(.scan_overview.scan_status != "finished") | .name'

# Export artifact list with specific fields
$ hrbcli artifact list my-app/backend -o json | \
    jq '.[] | {digest: .digest, tags: .tags, size: .size, created: .created}'
```
EOF

echo "📝 Creating Makefile..."
cat >Makefile <<'EOF'
# Harbor CLI Makefile

# Variables
BINARY_NAME := hrbcli
MODULE_NAME := github.com/pascal71/hrbcli
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X ${MODULE_NAME}/internal/version.Version=${VERSION} -X ${MODULE_NAME}/internal/version.BuildTime=${BUILD_TIME}"

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Build targets
TARGETS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=${MODULE_NAME}/internal/version.Version=$(VERSION) -X=${MODULE_NAME}/internal/version.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: clean build

# Build for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME) ./cmd/hrbcli

# Build for all platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	@for target in $(TARGETS); do \
		GOOS=$$(echo $$target | cut -d/ -f1) \
		GOARCH=$$(echo $$target | cut -d/ -f2) \
		OUTPUT=$(GOBIN)/$(BINARY_NAME)-$$(echo $$target | tr / -); \
		if [ $$GOOS = "windows" ]; then OUTPUT="$$OUTPUT.exe"; fi; \
		echo "Building $$target..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $$OUTPUT ./cmd/hrbcli; \
	done

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Run integration tests
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	go test -v -tags=integration ./test/integration/...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Install tools
.PHONY: tools
tools:
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@rm -f coverage.out

# Install locally
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(GOBIN)/$(BINARY_NAME) /usr/local/bin/

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)

# Generate completions
.PHONY: completions
completions: build
	@echo "Generating shell completions..."
	@mkdir -p completions
	@$(GOBIN)/$(BINARY_NAME) completion bash > completions/$(BINARY_NAME).bash
	@$(GOBIN)/$(BINARY_NAME) completion zsh > completions/_$(BINARY_NAME)
	@$(GOBIN)/$(BINARY_NAME) completion fish > completions/$(BINARY_NAME).fish

# Run the application
.PHONY: run
run: build
	@$(GOBIN)/$(BINARY_NAME)

# Create release
.PHONY: release
release:
	@echo "Creating release..."
	goreleaser release --clean

# Show help
.PHONY: help
help:
	@echo "Harbor CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-all      Build for all platforms"
	@echo "  make test           Run tests"
	@echo "  make integration    Run integration tests"
	@echo "  make lint           Run linter"
	@echo "  make tools          Install development tools"
	@echo "  make clean          Clean build artifacts"
	@echo "  make install        Install locally"
	@echo "  make completions    Generate shell completions"
	@echo "  make release        Create release with goreleaser"
	@echo "  make help           Show this help message"
EOF

echo "📝 Creating go.mod..."
cat >go.mod <<'EOF'
module github.com/pascal71/hrbcli

go 1.24.3

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/fatih/color v1.16.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/briandowns/spinner v1.23.0
	github.com/manifoldco/promptui v0.9.0
	gopkg.in/yaml.v3 v3.0.1
)
EOF

echo "📝 Creating .gitignore..."
cat >.gitignore <<'EOF'
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
hrbcli

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Go workspace file
go.work

# Dependency directories
vendor/

# IDE specific files
.idea/
.vscode/
*.swp
*.swo
*~

# OS specific files
.DS_Store
Thumbs.db

# Build artifacts
dist/
completions/

# Configuration files with secrets
.hrbcli.yaml
*.env

# Temporary files
*.tmp
*.bak
*.log
EOF

echo "📝 Creating LICENSE..."
cat >LICENSE <<'EOF'
Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright 2024 Pascal71

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
EOF

echo "📝 Creating .goreleaser.yml..."
cat >.goreleaser.yml <<'EOF'
project_name: hrbcli

before:
  hooks:
    - go mod tidy

builds:
  - id: hrbcli
    main: ./cmd/hrbcli
    binary: hrbcli
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X github.com/pascal71/hrbcli/internal/version.Version={{.Version}}
      - -X github.com/pascal71/hrbcli/internal/version.BuildTime={{.Date}}

archives:
  - id: hrbcli
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
    files:
      - README.md
      - LICENSE
      - docs/*

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
EOF

echo "📝 Creating GitHub Actions workflows..."

cat >.github/workflows/build.yml <<'EOF'
name: Build

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.24.3'
    
    - name: Build
      run: make build
    
    - name: Test
      run: make test
    
    - name: Lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
EOF

cat >.github/workflows/release.yml <<'EOF'
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.24.3'
    
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
EOF

echo "📝 Creating initial Go files..."

# Create cmd/hrbcli/main.go
cat >cmd/hrbcli/main.go <<'EOF'
package main

import (
	"os"

	"github.com/pascal71/hrbcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
EOF

# Create cmd/root.go placeholder
cat >cmd/root.go <<'EOF'
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "hrbcli",
	Short: "Harbor CLI - Manage Harbor from the command line",
	Long:  `Harbor CLI (hrbcli) is a command-line interface for Harbor registry.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	// Add commands here
}

func initConfig() {
	viper.SetEnvPrefix("HARBOR")
	viper.AutomaticEnv()
}
EOF

# Create internal/version/version.go
cat >internal/version/version.go <<'EOF'
package version

var (
	Version   = "dev"
	BuildTime = "unknown"
)
EOF

echo "✅ Project structure created successfully!"
echo ""
echo "📋 Next steps:"
echo "1. cd $PROJECT_NAME"
echo "2. go mod tidy"
echo "3. make build"
echo ""
echo "📚 Documentation:"
echo "- README.md: Project overview and quick start"
echo "- docs/DESIGN.md: Architecture and design decisions"
echo "- docs/COMMANDS.md: Complete command reference"
echo "- docs/CONTRIBUTING.md: Contribution guidelines"
echo "- docs/examples/EXAMPLES.md: Real-world usage examples"
echo ""
echo "🚀 Start implementing commands in the cmd/ directory!"
EOF

# Make the script executable
chmod +x setup-project.sh
