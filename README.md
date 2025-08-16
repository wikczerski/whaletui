# D5r - Docker CLI Dashboard

[![Go Version](https://img.shields.io/badge/Go-1.25.0+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Cross--Platform-blue.svg)]()
[![Docker](https://img.shields.io/badge/Docker-Required-blue.svg)](https://docker.com)

A modern, terminal-based Docker management tool inspired by k9s for kubernetes, providing an intuitive interface for managing Docker containers, images, volumes, and networks with a responsive TUI built in Go.

## ✨ Features

### 🐳 Core Docker Management
- **Container Management**: View, start, stop, restart, delete, and manage containers
- **Image Management**: Browse, inspect, and remove Docker images  
- **Volume Management**: Manage Docker volumes with ease
- **Network Management**: View and manage Docker networks

### 🎨 Modern Terminal Interface
- **Responsive Design**: Automatically adapts to terminal dimensions
- **Dynamic Legend**: Context-sensitive shortcuts and actions
- **Real-time Updates**: Live updates of Docker resources

### 🔍 Advanced Functionality
- **Container Logs**: View real-time container logs with full scrolling support
- **Resource Inspection**: Detailed JSON inspection of Docker resources
- **Command Mode**: k9s-inspired `:` command input for view switching
- **Smart Navigation**: Context-aware navigation hints

### ⌨️ Enhanced Keyboard Controls
- **Full Keyboard Navigation**: Complete keyboard-driven interface
- **Context-Sensitive Shortcuts**: Different actions based on current view

## 🚀 Quick Start

### Prerequisites

- **Docker Desktop** or **Docker Engine** - [Download Docker](https://docker.com/products/docker-desktop/)
- **Cross-platform support**: Windows, Linux, macOS

> **Note:** You only need Go 1.25.0+ if you want to build from source. Pre-built binaries do not require Go to be installed.

### Installation

#### Option 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/user/d5r.git
cd d5r

# Build the application
go build -o d5r

# Run D5r
./d5r
```

#### Option 2: Download Pre-built Binary

Visit the [Releases](https://github.com/user/d5r/releases) page to download the latest pre-built binary for your platform.

#### Option 3: Cross-platform Build

The project includes a build system that generates binaries and packages for multiple platforms:

```bash
# Build for all platforms (Linux, macOS, FreeBSD, Windows)
./scripts/build.sh          # Linux/macOS
.\scripts\build.ps1         # Windows PowerShell
.\scripts\build.bat         # Windows Batch

# Or use Makefile (Unix-like systems)
make build-all
make build-packages
make release   # Builds all binaries/packages, generates checksums, and prepares the full release in the dist/ directory
```

**Supported Platforms:**
- **Linux**: amd64, arm64, armv7, ppc64le, s390x
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **FreeBSD**: amd64, arm64
- **Windows**: amd64, arm64

**Package Formats:**
- **Linux**: .tar.gz, .deb, .rpm, .apk
- **macOS**: .tar.gz
- **FreeBSD**: .tar.gz
- **Windows**: .zip

### First Run

1. **Ensure Docker is running** - Start Docker Desktop or Docker Engine
2. **Launch D5r** - Run `./d5r` (or `d5r.exe` on Windows)
3. **Navigate** - Use `:` to enter command mode
4. **Explore** - Press `Enter` to inspect items, `l` for logs

### Remote Docker Hosts

D5r supports connecting to remote Docker hosts via command line arguments:

```bash
# Connect to a remote Docker host
./d5r --host tcp://192.168.1.100:2375

# Connect with custom refresh interval
./d5r --host tcp://192.168.1.100:2375 --refresh 10

# Connect with custom log level
./d5r --host tcp://192.168.1.100:2375 --log-level DEBUG
```

**Available flags:**
- `--host`: Remote Docker host (e.g., `tcp://192.168.1.100:2375`, `http://docker.example.com:2375`) - **⚠️ Not tested yet**
- `--refresh`: Refresh interval in seconds (default: 5)
- `--log-level`: Log level (DEBUG, INFO, WARN, ERROR, default: INFO)
- `--theme`: UI theme (default: default)

**Available commands:**
- `d5r` - Start the application (default command)
- `d5r version` - Show version information
- `d5r config` - Show configuration information
- `d5r --help` - Show help and available options

**Note:** Remote Docker hosts must have the Docker daemon configured to accept remote connections. For security, consider using TLS certificates for production environments.

**⚠️ Warning:** Remote Docker host functionality (`--host` flag) has not been tested yet. It might not work at all also use at your own risk and only in development/testing environments.

## 📖 Usage Guide

### 🎯 Command Mode (k9s-style)

Press `:` to enter command mode, then type:
- `containers` or `c` - Switch to containers view
- `images` or `i` - Switch to images view  
- `volumes` or `v` - Switch to volumes view
- `networks` or `n` - Switch to networks view
- `help` or `?` - Show help
- `quit`, `q`, or `exit` - Exit application

**Tip:** Command mode supports autocomplete - start typing and press Tab for suggestions.

### 📋 View Navigation

- **Command Mode**: `:` + view name
- **Context Awareness**: Navigation hints update based on current view

### 🔧 Actions by View

#### Containers View
- **s**: Start container
- **S**: Stop container  
- **r**: Restart container
- **d**: Delete container
- **l**: View container logs
- **i**: Inspect container

#### Images View
- **d**: Delete image
- **i**: Inspect image

#### Volumes View
- **d**: Delete volume
- **i**: Inspect volume

#### Networks View
- **d**: Delete network
- **i**: Inspect network

### 📖 Log View Navigation

When viewing container logs:
- **↑/↓**: Scroll line by line
- **PgUp/PgDn**: Page scrolling
- **Home/End**: Jump to top/bottom
- **Spacebar**: Half-page scrolling
- **ESC/Enter**: Return to container table

### 🔍 Inspect View Navigation

When viewing resource details:
- **↑/↓**: Scroll line by line
- **PgUp/PgDn**: Page scrolling
- **Home/End**: Jump to top/bottom
- **Spacebar**: Half-page scrolling
- **ESC/Enter**: Return to resource table

### 🎨 Table Navigation

- **↑/↓**: Navigate between rows
- **Enter**: View details and available actions
- **ESC**: Close details view
- **Tab**: Navigate between interactive elements

### 🌐 Global Shortcuts

- **Q**: Quit application
- **Ctrl+C**: Exit application
- **ESC**: Close modal or return to previous view
- **F5**: Refresh current view
- **?**: Show help dialog

### 💡 Tips & Tricks

- **Help Dialog**: Press `?` anytime to see all available shortcuts
- **Command Mode**: Use `:` for quick view switching and commands
- **Confirmation Dialogs**: Delete operations show confirmation prompts
- **Real-time Updates**: Views automatically refresh every 5 seconds (configurable)
- **Error Handling**: Errors are displayed in modal dialogs with clear messages

## ⚙️ Configuration

The application can be configured through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DOCKER_HOST` | Docker daemon address | Platform-specific default |
| `LOG_LEVEL` | Logging level | `INFO` |
| `REFRESH_INTERVAL` | UI refresh interval (seconds) | `5` |

### Example Configuration

```bash
# Set custom Docker host
export DOCKER_HOST=tcp://localhost:2375

# Enable debug logging
export LOG_LEVEL=DEBUG

# Set refresh interval to 10 seconds
export REFRESH_INTERVAL=10
```

## 🏗️ Architecture

D5r follows the architecture below:

```
d5r/
├── cmd/                   # Command-line interface
├── internal/              # Internal packages
│   ├── app/              # Main application logic
│   ├── config/           # Configuration management
│   ├── docker/           # Docker client wrapper
│   ├── errors/           # Error handling
│   ├── logger/           # Structured logging system
│   ├── models/           # Data models and types
│   ├── services/         # Business logic layer
│   └── ui/               # Terminal UI components
│       ├── builders/     # UI component builders
│       ├── constants/    # UI constants and colors
│       ├── core/         # Core UI framework
│       ├── handlers/     # Event handlers
│       ├── interfaces/   # UI interfaces
│       ├── managers/     # UI managers
│       └── views/        # Individual view implementations
├── scripts/               # Build and utility scripts
├── go.mod                 # Go module file
└── README.md              # This file
```

### 🔧 Key Components

- **Service Layer**: Business logic separated from UI
- **Docker Client**: High-level Docker API wrapper
- **UI Framework**: tview-based terminal interface
- **Builder Pattern**: Consistent UI component creation
- **View Manager**: Centralized view switching logic

## 🛠️ Development

### Building

```bash
# Build for current platform
go build -o d5r

# Build for specific platform
GOOS=windows GOARCH=amd64 go build -o d5r.exe

# Build with debug info
go build -ldflags="-s -w" -o d5r
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run vet
go vet ./...
```

### Dependencies

```bash
# Add new dependency
go get github.com/example/package

# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## 🚧 Features Roadmap

### Planned Enhancements
- **Container Stats**: Real-time resource usage monitoring
- **Multi-Container Operations**: Bulk start/stop/restart
- **Custom Filters**: Advanced filtering and search
- **Export Functionality**: Save logs, configs, and data
- **Theme Support**: Customizable color schemes

### Future Ideas
- **Docker Compose Support**: Manage multi-container applications
- **Plugin System**: Extensible architecture for custom views
- **Metrics Dashboard**: Performance and health monitoring
- **Container Shell Integration**: Direct terminal access
- **Docker Swarm**: Manage Services and stacks in Docker Swarms

## 🤝 Contributing

Contributions are Welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. **Follow Go conventions**: Use `gofmt`, `golint`, and `go vet`
2. **Add tests**: Include tests for new functionality
3. **Update documentation**: Keep README and code comments current
4. **Use clean architecture**: Maintain separation of concerns
5. **Follow existing patterns**: Match the established code style

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **k9s**: Inspiration for the UI design and interaction patterns
- **tview**: Terminal UI library for Go
- **Docker**: Containerization platform
- **Go Community**: Amazing ecosystem and tooling

## 📊 Project Status

- **Version**: 0.1.0a
- **Status**: Alpha Development
- **Platform**: Cross-platform (Windows, Linux, macOS)
- **Go Version**: 1.25.0+

---

**Made with ❤️ using Go and tview**