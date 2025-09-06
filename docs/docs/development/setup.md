---
id: setup
title: Development Setup
sidebar_label: Development Setup
---

# Development Setup

This guide will help you set up a development environment for contributing to WhaleTUI.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.25+**: [Download Go](https://golang.org/dl/)
- **Git**: [Download Git](https://git-scm.com/)
- **Docker**: [Download Docker](https://docker.com/)
- **pre-commit**: [Download pre-commit](https://pre-commit.com/)
- **Make** (optional): Usually pre-installed on Linux/macOS, [download for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

## Getting the Source Code

### 1. Fork the Repository

1. Go to [https://github.com/wikczerski/whaletui](https://github.com/wikczerski/whaletui)
2. Click the "Fork" button in the top-right corner
3. Clone your forked repository:

```bash
git clone https://github.com/wikczerski/whaletui.git
cd whaletui
```

### 2. Add Upstream Remote

```bash
git remote add upstream https://github.com/wikczerski/whaletui.git
git fetch upstream
```

## Building from Source

### 1. Build the Application

```bash
# Build for your current platform (outputs whaletui.exe by default)
make build

# Or build manually (Windows)
go build -o whaletui.exe .
# Unix based
go build -o whaletui .
```

### 2. Build for Multiple Platforms

```bash
# Build for all supported platforms using Goreleaser
make goreleaser-build

# Install Goreleaser first if needed
make install-goreleaser
```

### 3. Install Locally

```bash
# Install to your Go bin directory
go install .

# Or use the built executable directly
./whaletui.exe
```

## Development Workflow

### 1. Create a Feature Branch

```bash
# Make sure you're on your fork
git remote -v
# Should show your fork as origin, upstream as upstream

# Create and switch to a new feature branch
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

- Write your code following the [coding standards](coding-standards.md)
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests

```bash
# Run all tests
make test

# Or run manually
go test ./...

# Run specific test files
go test ./internal/app/...

# Run tests with verbose output
go test -v ./...
```

### 4. Run Code Quality Checks

```bash
# Run all quality checks
make test-all

# Run specific checks
make lint          # Run golangci-lint
make fmt           # Format code with gofumpt
make format-check  # Check formatting without changing
make format-fix    # Fix formatting and imports
make imports       # Fix import organization
make vet           # Run go vet
```

### 5. Build and Test

```bash
# Build and run quality checks
make test-all

# Build the application
make build
```

## Project Structure

```
whaletui/
├── cmd/                    # Application entry points
│   ├── root.go           # Main application
│   └── ...
├── internal/              # Private application code
│   ├── app/              # Core application logic
│   ├── config/           # Configuration management
│   ├── docker/           # Docker client integration
│   ├── domains/          # Business logic domains
│   ├── logger/           # Logging system
│   ├── mocks/            # Generated mocks
│   ├── shared/           # Shared utilities
│   └── ui/               # User interface components
├── docs/                  # Documentation (this site)
├── scripts/               # Build and utility scripts
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## Key Dependencies

### Core Dependencies

- **Docker SDK**: `github.com/docker/docker/client` and related types
- **Terminal UI**: `github.com/gdamore/tcell/v2` and `github.com/rivo/tview`
- **CLI Framework**: `github.com/spf13/cobra`
- **SSH Support**: `golang.org/x/crypto/ssh`
- **Configuration**: Built-in config management (no external config library)
- **Logging**: Standard `log/slog` package
- **YAML**: `gopkg.in/yaml.v3`

### Development Dependencies

- **Testing**: `github.com/stretchr/testify`
- **Mocking**: `github.com/vektra/mockery`
- **Linting**: `github.com/golangci/golangci-lint`
- **Formatting**: `mvdan.cc/gofumpt`

## Configuration

### Application Configuration

The application uses a JSON configuration file located at `~/.whaletui/config.json` with the following structure:

```json
{
  "refresh_interval": 5,
  "log_level": "INFO",
  "log_file_path": "./logs/whaletui.log",
  "docker_host": "unix:///var/run/docker.sock",
  "theme": "default",
  "remote_host": "",
  "remote_user": "",
  "remote_port": 2375
}
```

**Configuration Fields:**
- **`refresh_interval`**: UI refresh rate in seconds (default: 5)
- **`log_level`**: Logging level (default: "INFO")
- **`log_file_path`**: Path to log file (default: "./logs/whaletui.log")
- **`docker_host`**: Docker daemon socket (default: Unix socket on Linux/macOS, auto-detect on Windows)
- **`theme`**: Theme name (default: "default")
- **`remote_host`**: Remote Docker host for SSH connections
- **`remote_user`**: Username for SSH connections
- **`remote_port`**: SSH port (default: 2375)

### Theme Configuration

Themes are configured using YAML or JSON files. You can create custom themes by placing them in the config directory:

```yaml
# ~/.whaletui/themes/custom-theme.yaml
colors:
  header: "#00FF00"     # Bright green header
  border: "#444444"     # Dark gray borders
  text: "#E0E0E0"       # Light gray text
  background: "black"   # Black background
  success: "#00FF00"    # Bright green success
  warning: "#FFFF00"    # Bright yellow warning
  error: "#FF0000"      # Bright red error
  info: "#0080FF"       # Bright blue info

shell:
  border: "#444444"
  title: "#00FF00"
  text: "#E0E0E0"
  background: "black"
  cmd:
    label: "#00FF00"
    border: "#444444"
    text: "#E0E0E0"
    background: "black"
    placeholder: "#666666"
```

### Configuration File Location

- **Linux/macOS**: `~/.whaletui/config.json`
- **Windows**: `%USERPROFILE%\.whaletui\config.json`

The configuration directory is automatically created when you first run the application.

## Testing

### Unit Tests

```bash
# Run unit tests
make test

# Run with race detection
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Docker Testing Environment

The project includes a comprehensive Docker testing environment setup:

```bash
# Set up Docker testing environment (creates containers, networks, volumes, services)
make docker-test

# Clean up Docker testing environment
make docker-cleanup
```

This creates:
- Docker Swarm cluster
- Test networks and volumes
- Sample containers and services
- Perfect for testing all WhaleTUI features

### Mock Generation

```bash
# Generate mocks for interfaces
mockery --dir=internal/domains/container --name=ContainerService

# Generate all mocks in a directory
mockery --dir=internal/domains/container --all
```

## Debugging

### Using Delve

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug ./cmd/root.go

# Debug with arguments
dlv debug ./cmd/root.go -- --config=./config/local.yaml
```

### Using VS Code

1. Install the Go extension
2. Set breakpoints in your code
3. Press F5 to start debugging
4. Use the debug console and variables panel

### Logging

```bash
# View logs
tail -f logs/whaletui-dev.log
```

## Code Quality

### Pre-commit Workflow

```bash
# Run all quality checks before committing
make test-all

# This includes:
# - go vet (static analysis)
# - goimports (import organization)
# - gofumpt (code formatting)
# - golangci-lint (linting)
# - go test (unit tests)
```

### Code Formatting

```bash
# Format code
make fmt

# Check formatting without changing
make format-check

# Fix formatting and imports
make format-fix
```

### Linting

```bash
# Run linters
make lint

# Run go vet
make vet
```

## Documentation

### Building Documentation

The documentation is maintained on a separate `gh-pages` branch. To work with the documentation:

```bash
# Switch to the gh-pages branch
git checkout gh-pages

# Navigate to docs directory
cd docs

# Install dependencies
npm install

# Start development server
npm run start

# Build for production
npm run build
```

**Note**: The documentation is completely separate from the main application code. Make sure you're on the `gh-pages` branch when working with documentation files.

### Writing Documentation

- Use clear, concise language
- Include code examples
- Add screenshots for UI features
- Keep documentation up-to-date with code changes

## Contributing

### 1. Commit Your Changes

```bash
git add .
git commit -m "feat: add new container management feature

- Add container restart functionality
- Implement container health monitoring
- Update documentation with examples"
```

### 2. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 3. Create a Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill out the PR template
5. Submit the PR

## Getting Help

### Resources

- **Go Documentation**: [https://golang.org/doc/](https://golang.org/doc/)
- **Docker SDK**: [https://pkg.go.dev/github.com/docker/docker/client](https://pkg.go.dev/github.com/docker/docker/client)
- **Terminal UI**: [https://github.com/gdamore/tcell](https://github.com/gdamore/tcell) and [https://github.com/rivo/tview](https://github.com/rivo/tview)
- **CLI Framework**: [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
- **SSH Support**: [https://pkg.go.dev/golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh)

### Community

- **GitHub Issues**: [https://github.com/wikczerski/whaletui/issues](https://github.com/wikczerski/whaletui/issues)
- **Discussions**: [https://github.com/wikczerski/whaletui/discussions](https://github.com/wikczerski/whaletui/discussions)

## Next Steps

Now that you have your development environment set up:

1. **Explore the Codebase**: Familiarize yourself with the project structure
2. **Run the Tests**: Ensure everything is working correctly
3. **Set up Docker Testing**: Use `make docker-test` to create a test environment
4. **Pick an Issue**: Find a good first issue to work on
5. **Join the Community**: Connect with other contributors

Happy coding! 🚀
