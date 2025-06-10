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
- 🚚 **Distribution Management** - Manage preheat providers and policies
- 🛠 **System Configuration** - View and update Harbor system settings
- 🔌 **Registry Endpoint Management** - Configure external registries for replication or proxy cache
- 📊 **Job Service Monitoring** - Inspect worker pools and queue lengths

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

# Get repository details
hrbcli repo get myproject/myapp

# List tags for a repository
hrbcli repo tags myproject/myapp

# Delete a repository
hrbcli repo delete myproject/myapp:v1.0.0

# Get system information
hrbcli system info

# Show Harbor statistics
hrbcli system statistics


# Show job service dashboard
hrbcli jobservice dashboard


```

### Scanner Commands

```bash
hrbcli scanner scan <project>
hrbcli scanner running <project>
hrbcli scanner reports <project> --summary
hrbcli scanner reports <project> --sort repo
```

### Distribution Commands

```bash
hrbcli distribution providers <project>
hrbcli distribution policies <project>
```

See [docs/COMMANDS.md](docs/COMMANDS.md) for more details.

## Configuration

Harbor CLI can be configured through:

1. **Configuration file** (`~/.hrbcli.yaml`)
2. **Environment variables** (`HARBOR_URL`, `HARBOR_USERNAME`, `HARBOR_PASSWORD`)
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

`HARBOR_PASSWORD` can hold either your Harbor account password or a robot account token. Set it as an environment variable to avoid storing credentials in your configuration file.

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
