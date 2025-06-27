# Variables
BINARY_NAME=dbbackup
VERSION ?= dev
BUILD_DIR=./dist
GO_VERSION=1.24

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-w -s -X main.version=$(VERSION)"

.PHONY: all build clean test test-verbose coverage deps fmt vet lint check fix help install dev-tools

# Default target
all: check build

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development tools
dev-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u honnef.co/go/tools/cmd/staticcheck
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Dependencies
deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Check formatting
fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@files=$$($(GOFMT) -s -l . | grep -v vendor/ | grep -v .git/); \
	if [ -n "$$files" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$files"; \
		echo "Run 'make fmt' to fix formatting"; \
		exit 1; \
	fi
	@echo "All files are properly formatted"

# Vet code
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

# Test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Test with verbose output
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race ./...

# Test coverage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint with golangci-lint
lint: ## Run golangci-lint linter
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run 'make dev-tools' to install it."; \
		exit 1; \
	fi

# Alternative linting with staticcheck
lint-staticcheck: ## Run staticcheck linter
	@echo "Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed. Run 'make dev-tools' to install it."; \
		exit 1; \
	fi

# Check imports
imports: ## Fix imports with goimports
	@echo "Fixing imports..."
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Run 'make dev-tools' to install it."; \
		exit 1; \
	fi

# All CI checks (same as GitHub Actions)
check: fmt-check vet test ## Run all CI checks locally
	@echo "All checks passed! ✅"

# Lint with minimal config (ignore common defer Close issues)
lint-minimal: ## Run basic linting (ignore defer Close errors)
	@echo "Running minimal linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --enable=unused,staticcheck,gosimple,ineffassign,govet; \
	else \
		echo "golangci-lint not installed. Run 'make dev-tools' to install it."; \
		exit 1; \
	fi

# Fix common issues
fix: fmt imports ## Fix formatting and imports
	@echo "Code fixed! ✅"

# Build for current platform
build: deps ## Build binary for current platform
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go

# Build for all platforms
build-all: deps ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	./scripts/build.sh $(VERSION) $(BUILD_DIR)

# Install to $GOPATH/bin
install: ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) main.go

# Development workflow
dev: clean deps fix check build ## Complete development workflow
	@echo "Development build complete! ✅"

# CI simulation (exactly what GitHub Actions runs)
ci: deps fmt-check vet test ## Simulate CI pipeline locally
	@echo "CI simulation complete! ✅"

# Release preparation
release-prep: clean deps fix check test build-all ## Prepare for release
	@echo "Release preparation complete! ✅"
	@echo "Built binaries are in $(BUILD_DIR)/"

# Docker build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Docker build for multiple architectures
docker-buildx: ## Build multi-arch Docker image
	@echo "Building multi-arch Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t $(BINARY_NAME):$(VERSION) .

# Security scan
security: ## Run Go security scanner
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Dependency check
deps-check: ## Check for dependency updates
	@echo "Checking for dependency updates..."
	$(GOCMD) list -u -m all

# Tidy dependencies
deps-tidy: ## Tidy and vendor dependencies
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Show project info
info: ## Show project information
	@echo "Project: $(BINARY_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Go version: $(shell $(GOCMD) version)"
	@echo "Build dir: $(BUILD_DIR)"
	@echo "GOPATH: $(GOPATH)"
	@echo "GOOS: $(shell $(GOCMD) env GOOS)"
	@echo "GOARCH: $(shell $(GOCMD) env GOARCH)"

# Quick checks (faster than full CI)
quick-check: fmt-check vet ## Quick checks for development
	@echo "Quick checks passed! ✅"