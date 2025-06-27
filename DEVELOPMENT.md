# Development Guide

## Prerequisites

- Go 1.21 or later
- Make
- Docker (for container builds)

## Quick Start

```bash
# Clone and setup
git clone <repo>
cd dbbackup-cli

# Install development tools (optional)
make dev-tools

# Run all CI checks locally
make ci

# Build for current platform
make build

# Development workflow (fix, check, build)
make dev
```

## Makefile Commands

### Essential Commands

```bash
make help           # Show all available commands
make ci             # Run full CI simulation (same as GitHub Actions)
make dev            # Complete development workflow
make build          # Build binary for current platform
```

### Code Quality

```bash
make fmt            # Format Go code
make fmt-check      # Check if code is formatted
make vet            # Run go vet
make check          # Run all CI checks (fmt-check + vet + test)
make fix            # Fix formatting and imports
```

### Building

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Install to $GOPATH/bin
```

### Docker

```bash
make docker-build      # Build Docker image
make docker-buildx     # Build multi-arch Docker image
```

### Dependencies

```bash
make deps           # Download dependencies
make deps-tidy      # Tidy dependencies
make deps-check     # Check for updates
```

## Pre-commit Workflow

Before committing changes, always run:

```bash
make ci
```

This will:
1. ✅ Check code formatting
2. ✅ Run `go vet`
3. ✅ Run tests
4. ✅ Verify dependencies

## Fixing CI Issues

### Formatting Issues
```bash
# Fix formatting
make fmt

# Check if fixed
make fmt-check
```

### Go Vet Issues
```bash
# Run vet to see issues
make vet

# Fix issues manually, then verify
make vet
```

### Import Issues
```bash
# Install goimports (if needed)
make dev-tools

# Fix imports
make imports
```

## GitHub Actions

### CI Pipeline (`.github/workflows/ci.yml`)
- Runs on every push/PR to main/develop
- Tests multiple platforms
- No releases created

### Build & Release (`.github/workflows/build-and-release.yml`)
- **Manual trigger**: Go to Actions → Build and Release → Run workflow
- **Weekly**: Every Sunday at 2 AM UTC  
- **Tag push**: Automatic on version tags
- Creates multi-arch Docker images
- Creates GitHub releases with binaries

## Local Development Tips

### Quick Checks
```bash
make quick-check    # Fast formatting and vet check
```

### Version Testing
```bash
# Build with custom version
go build -ldflags="-X main.version=test-1.0.0" -o dbbackup-test main.go
./dbbackup-test --version
```

### Cross-compilation Testing
```bash
# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o dbbackup-linux main.go
GOOS=windows GOARCH=amd64 go build -o dbbackup-windows.exe main.go
```

## Docker Development

### Local Testing
```bash
# Build and test locally
make docker-build
docker run --rm dbbackup:dev --help

# Test with database
docker run --rm --network host dbbackup:dev dump postgres local --help
```

## Release Process

1. **Prepare release**:
   ```bash
   make release-prep
   ```

2. **Create tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Or trigger manually**:
   - Go to GitHub Actions
   - Run "Build and Release" workflow  
   - Specify version: `v1.0.0`
   - Enable "Create GitHub release"

## Troubleshooting

### "make: command not found"
Install make for your platform:
- macOS: `xcode-select --install`
- Ubuntu: `sudo apt-get install build-essential`
- Windows: Use WSL or install make via chocolatey

### Go formatting issues
```bash
# See what files need formatting
gofmt -s -l .

# Fix all files
make fmt
```

### Import issues
```bash
# Install goimports
go install golang.org/x/tools/cmd/goimports@latest

# Fix imports
make imports
```

### Docker build issues
```bash
# Clean Docker cache
docker system prune -f

# Rebuild without cache
docker build --no-cache -t dbbackup:dev .
```