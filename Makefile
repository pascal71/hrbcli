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
	@echo "Running gofmt..."
	@gofmt -w $(shell git ls-files '*.go')
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
