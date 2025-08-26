# WhaleTUI - Docker CLI Dashboard

[![Go Version](https://img.shields.io/badge/Go-1.25.0+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-GNU%20AGPL%20v3-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Cross--Platform-blue.svg)]()
[![Docker](https://img.shields.io/badge/Docker-Required-blue.svg)](https://docker.com)

A terminal-based Docker management tool inspired by k9s, providing an intuitive and powerful interface for managing Docker containers, images, volumes, and networks with a modern, responsive TUI.

## ✨ Features

- **Container Management** - View, start, stop, restart, and manage containers
- **Image Management** - Browse, pull, remove, and inspect Docker images
- **Volume Management** - Manage Docker volumes and their data
- **Network Management** - Configure and manage Docker networks
- **Real-time Monitoring** - Live updates of container status and resource usage
- **Theme Support** - Customizable color schemes and UI appearance
- **Remote Host Support** - Connect to remote Docker hosts via SSH
- **Log Viewing** - Real-time container logs with search and filtering

## 🚀 Quick Start

### Prerequisites

- **Docker Desktop** or **Docker Engine** - [Download Docker](https://docker.com/products/docker-desktop/)
- **Cross-platform support**: Windows, Linux, macOS

> **Note:** You only need Go 1.25.0+ if you want to build from source. Pre-built binaries do not require Go to be installed.

### Installation

For detailed installation instructions, visit our [Installation Guide](https://wikczerski.github.io/whaletui/docs/installation).

**Quick commands:**
```bash
# Clone and build
git clone https://github.com/wikczerski/whaletui.git
cd whaletui
go build -o whaletui

# Run
./whaletui
```

**Pre-built binaries:** Available on the [Releases](https://github.com/wikczerski/whaletui/releases) page.

## Usage

### Basic Commands

1. **Launch whaletui** - Run `./whaletui` (or `whaletui.exe` on Windows)

2. **Navigate the Interface** - Use arrow keys, Tab, and Enter to navigate

3. **Container Operations** - Select containers and use keyboard shortcuts for actions

### Remote Host Connection

whaletui supports connecting to remote Docker hosts using the `connect` subcommand:

```bash
# Basic SSH connection
./whaletui connect --host 192.168.1.100 --user admin

# With custom port
./whaletui connect --host 192.168.1.100 --user admin --port 2376

# With additional options
./whaletui connect --host 192.168.1.100 --user admin --refresh 10 --log-level DEBUG

# Using TCP protocol
./whaletui connect --host tcp://192.168.1.100 --user admin

# With port in host string
./whaletui connect --host 192.168.1.100:2375 --user admin
```

### Command Line Options

- `whaletui` - Start with local Docker instance (default)
- `whaletui connect` - Connect to a remote Docker host via SSH
- `whaletui theme` - Manage theme configuration
- `whaletui --help` - Show help and available options

### 🎨 Theme Configuration

whaletui supports multiple theme formats for customizing the UI appearance:

**JSON Theme Example:**
```json
{
  "colors": {
    "primary": "#00ff00",
    "secondary": "#ff00ff",
    "background": "#000000",
    "text": "#ffffff"
  }
}
```

**YAML Theme Example:**
```yaml
colors:
  primary: "#00ff00"
  secondary: "#ff00ff"
  background: "#000000"
  text: "#ffffff"
```

**Apply a theme:**
```bash
whaletui --theme config/custom-theme.yaml
whaletui --theme config/theme.json
```

## 🏗️ Architecture

whaletui follows a modular architecture with clear separation of concerns:

```
whaletui/
├── cmd/           # Command line interface
├── internal/      # Internal application logic
│   ├── app/       # Application core
│   ├── config/    # Configuration management
│   ├── docker/    # Docker client operations
│   ├── services/  # Business logic services
│   └── ui/        # Terminal UI components
├── config/        # Configuration files
└── docs/          # Documentation
```

## 🛠️ Development

### Building from Source

```bash
# Basic build
go build -o whaletui

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -o whaletui.exe
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o whaletui
```

### Running Tests

```bash
go test ./...
```

### Local Development

```bash
# Run locally
./whaletui

# Connect to remote host for testing
./whaletui connect --host 192.168.1.100 --user admin

# Run with debug logging
./whaletui --log-level DEBUG

# Run with custom refresh rate
./whaletui --refresh 10

# Full remote connection with options
./whaletui connect --host 192.168.1.100 --user admin --port 2376 --refresh 10 --log-level DEBUG
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
## Acknowledgments

- Inspired by [k9s](https://k9scli.io/) for Kubernetes management
- Built with [tview](https://github.com/rivo/tview) for the terminal UI
- Uses [tcell](https://github.com/gdamore/tcell) for terminal handling
