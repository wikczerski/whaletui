# Installation Guide

This guide provides detailed installation instructions for whaletui, a terminal-based Docker management tool.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Building from Source](#building-from-source)
  - [Package Managers](#package-managers)
- [Platform-Specific Instructions](#platform-specific-instructions)
  - [Windows](#windows)
  - [macOS](#macos)
  - [Linux](#linux)
- [Docker Requirements](#docker-requirements)
- [Post-Installation](#post-installation)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, Linux (kernel 3.10+), or FreeBSD
- **Architecture**:
  - **x86_64 (AMD64)** - Primary support for Windows, Linux, macOS, and FreeBSD
  - **x86 (i386)** - 32-bit support for Linux and FreeBSD
  - **ARM64** - ARM64 support for Linux, Windows, and macOS
  - **PowerPC64 LE** - PowerPC64 Little Endian for Linux
  - **S390x** - IBM S390x for Linux
- **Memory**: Minimum 512MB RAM, recommended 2GB+
- **Terminal**: A modern terminal emulator that supports ANSI colors

### Required Software

- **Docker**: Docker Desktop or Docker Engine
  - **Local Mode**: Required when running whaletui locally to manage Docker on the same machine
  - **Connect Mode**: Not required locally when connecting to a remote Docker host via SSH
  - Docker Desktop: [Download for Windows](https://docs.docker.com/desktop/install/windows-install/) | [Download for macOS](https://docs.docker.com/desktop/install/mac-install/) | [Download for Linux](https://docs.docker.com/desktop/install/linux-install/)
  - Docker Engine: [Installation guide](https://docs.docker.com/engine/install/)
- **Go** (only required for building from source): Version 1.25.0 or higher

## Installation Methods

### Pre-built Binaries (Recommended)

The easiest way to install whaletui is to download pre-built binaries from the [Releases](https://github.com/wikczerski/whaletui/releases) page.

#### Download Steps

1. Visit the [Releases](https://github.com/wikczerski/whaletui/releases) page
2. Download the appropriate binary for your platform:
   - **Windows**:
     - `whaletui_windows_amd64.exe` (64-bit Intel/AMD)
     - `whaletui_windows_arm64.exe` (64-bit ARM)
   - **macOS**:
     - `whaletui_darwin_amd64` (64-bit Intel)
     - `whaletui_darwin_arm64` (Apple Silicon)
   - **Linux**:
     - `whaletui_linux_amd64` (64-bit Intel/AMD)
     - `whaletui_linux_386` (32-bit Intel/AMD)
     - `whaletui_linux_arm64` (64-bit ARM)
     - `whaletui_linux_ppc64le` (PowerPC64 LE)
     - `whaletui_linux_s390x` (IBM S390x)
   - **FreeBSD**:
     - `whaletui_freebsd_amd64` (64-bit Intel/AMD)
     - `whaletui_freebsd_386` (32-bit Intel/AMD)
     - `whaletui_freebsd_arm64` (64-bit ARM)

#### Installation Steps

**Windows:**
```cmd
# Download and rename to whaletui.exe
ren whaletui_windows_amd64.exe whaletui.exe

# Move to a directory in your PATH (optional)
move whaletui.exe C:\Windows\System32\
# OR add the directory containing whaletui.exe to your PATH environment variable
```

**macOS/Linux:**
```bash
# Download and make executable
chmod +x whaletui_linux_amd64

# Rename for convenience
mv whaletui_linux_amd64 whaletui

# Move to a directory in your PATH
sudo mv whaletui /usr/local/bin/
# OR add the directory to your PATH in ~/.bashrc or ~/.zshrc
export PATH="$PATH:$HOME/bin"
```

### Building from Source

If you prefer to build from source or need a custom build, follow these steps:

#### 1. Install Go

**Windows:**
- Download from [golang.org/dl](https://golang.org/dl/)
- Run the installer and follow the setup wizard
- Ensure Go is added to your PATH

**macOS:**
```bash
# Using Homebrew
brew install go

# Using MacPorts
sudo port install go
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL/Fedora
sudo yum install golang
# OR
sudo dnf install golang

# Arch Linux
sudo pacman -S go
```

#### 2. Clone the Repository

```bash
git clone https://github.com/wikczerski/whaletui.git
cd whaletui
```

#### 3. Build the Application

```bash
# Basic build
go build -o whaletui

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -o whaletui.exe
GOOS=windows GOARCH=arm64 go build -o whaletui_arm64.exe
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o whaletui
GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o whaletui_386
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o whaletui_arm64
GOOS=linux GOARCH=ppc64le go build -ldflags="-s -w" -o whaletui_ppc64le
GOOS=linux GOARCH=s390x go build -ldflags="-s -w" -o whaletui_s390x
GOOS=darwin GOARCH=amd64 go build -o whaletui
GOOS=darwin GOARCH=arm64 go build -o whaletui_arm64
GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o whaletui
GOOS=freebsd GOARCH=386 go build -ldflags="-s -w" -o whaletui_386
GOOS=freebsd GOARCH=arm64 go build -ldflags="-s -w" -o whaletui_arm64

# Build with version information
go build -ldflags="-X main.version=$(git describe --tags --always --dirty)" -o whaletui
```

#### 4. Install the Binary

```bash
# Move to a directory in your PATH
sudo mv whaletui /usr/local/bin/

# Verify installation
whaletui --version
```

### Package Managers

#### Using Homebrew (macOS/Linux)

> **Note**: Homebrew taps are custom repositories that contain formulas not available in the main Homebrew repository. You only need to add a tap if the formula isn't available with `brew install whaletui`.

```bash
# Install whaletui using Homebrew
# If available in main repository:
brew install whaletui

# If not in main repository, use custom tap:
brew tap wikczerski/tap
brew install whaletui

# Update to latest version
brew upgrade whaletui
```

#### Using Go Install

```bash
# Install directly from GitHub
go install github.com/wikczerski/whaletui@latest

# Install specific version
go install github.com/wikczerski/whaletui@v0.1.1
```

## Platform-Specific Instructions

### Windows

#### Prerequisites
- Windows 10 or later
- Docker Desktop for Windows
- PowerShell 5.0+ or Windows Terminal

#### Installation Steps

1. **Install Docker Desktop**
   - Download from [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)
   - Enable WSL 2 if prompted
   - Restart your computer after installation

2. **Install whaletui**
   ```cmd
   # Download pre-built binary
   # OR build from source
   go build -o whaletui.exe
   ```

3. **Add to PATH** (optional)
   - Copy `whaletui.exe` to a directory in your PATH
   - Or add the directory containing `whaletui.exe` to your PATH environment variable

#### Running on Windows

```cmd
# Start whaletui
whaletui.exe

# Connect to remote host
whaletui.exe connect --host 192.168.1.100 --user admin
```

### macOS

#### Prerequisites
- macOS 10.15 (Catalina) or later
- Docker Desktop for Mac
- Terminal.app or iTerm2

#### Installation Steps

1. **Install Docker Desktop**
   ```bash
   # Using Homebrew
   brew install --cask docker

   # Or download from Docker website
   # https://docs.docker.com/desktop/install/mac-install/
   ```

2. **Install whaletui**
   ```bash
   # Using Homebrew (recommended)
   brew install whaletui

   # Using Go install
   go install github.com/wikczerski/whaletui@latest

   # Or build from source
   git clone https://github.com/wikczerski/whaletui.git
   cd whaletui
   go build -o whaletui
   sudo mv whaletui /usr/local/bin/
   ```

#### Running on macOS

```bash
# Start whaletui
whaletui

# Connect to remote host
whaletui connect --host 192.168.1.100 --user admin
```

### Linux

#### Prerequisites
- Linux kernel 3.10 or later
- Docker Engine or Docker Desktop for Linux
- A modern terminal emulator

#### Installation Steps

**Ubuntu/Debian:**
```bash
# Install Docker
sudo apt update
sudo apt install docker.io docker-compose
sudo usermod -aG docker $USER

# Install whaletui (choose one method)
# Option 1: Using Homebrew (recommended)
brew install wikczerski/tap/whaletui

# Option 2: Using Go install
sudo apt install golang-go
go install github.com/wikczerski/whaletui@latest
```

**CentOS/RHEL/Fedora:**
```bash
# Install Docker
sudo yum install docker docker-compose
# OR
sudo dnf install docker docker-compose

sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Install Go
sudo yum install golang
# OR
sudo dnf install golang

# Install whaletui
go install github.com/wikczerski/whaletui@latest
```

**Arch Linux:**
```bash
# Install Docker
sudo pacman -S docker docker-compose
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Install Go
sudo pacman -S go

# Install whaletui
go install github.com/wikczerski/whaletui@latest
```

### FreeBSD

#### Prerequisites
- FreeBSD 12.0 or later
- Docker Engine (via ports or packages)
- A modern terminal emulator

#### Installation Steps

```bash
# Install Docker
pkg install docker

# Start Docker service
sysrc docker_enable="YES"
service docker start

# Install Go
pkg install go

# Install whaletui
go install github.com/wikczerski/whaletui@latest
```

#### Running on FreeBSD

```bash
# Start whaletui
whaletui

# Connect to remote host
whaletui connect --host 192.168.1.100 --user admin
```

## Docker Requirements

> **Note**: Docker is only required locally when running whaletui in local mode. When using `whaletui connect` to connect to a remote Docker host via SSH, Docker is not needed on your local machine.

### Docker Engine

- **Version**: Docker Engine 20.10.0 or later
- **API Version**: 1.41 or later
- **Features**: Must support Docker API v1.41+

### Docker Desktop

- **Windows**: Docker Desktop 4.0.0 or later
- **macOS**: Docker Desktop 4.0.0 or later
- **Linux**: Docker Desktop 4.0.0 or later

### Docker Compose

- **Version**: Docker Compose 2.0.0 or later
- **Required for**: Multi-container application management

### Permissions

Ensure your user has the necessary permissions to access the Docker daemon:

```bash
# Add user to docker group (Linux/macOS)
sudo usermod -aG docker $USER

# Verify access
docker ps

# If you get permission errors, restart your session or reboot
```

## Post-Installation

### Verify Installation

```bash
# Check version
whaletui --version

# Check help
whaletui --help

# Test basic functionality
whaletui
```

### First Run

1. **Start whaletui**
   ```bash
   whaletui
   ```

2. **Navigate the Interface**
   - Use arrow keys to navigate
   - Press Enter to select items
   - Use `:` to navigate views
   - Press 'q' to quit

3. **Test Docker Connection**
   - Ensure Docker is running
   - Verify containers are visible in the interface

### Configuration

Create a configuration directory for themes and settings:

```bash
# Create config directory
mkdir -p ~/.config/whaletui

# Copy default themes
cp -r config/* ~/.config/whaletui/
```

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Error: permission denied while trying to connect to the Docker daemon socket
sudo usermod -aG docker $USER
# Log out and log back in, or reboot
```

#### Docker Not Running
```bash
# Start Docker service
sudo systemctl start docker

# Check Docker status
sudo systemctl status docker
```

#### Binary Not Found
```bash
# Check if binary is in PATH
which whaletui

# Add to PATH if needed
export PATH="$PATH:$HOME/go/bin"
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
```

#### Build Errors
```bash
# Update Go modules
go mod tidy

# Clean build cache
go clean -cache

# Check Go version
go version
```

#### Remote Connection Issues
```bash
# Test SSH connection
ssh user@host

# Check Docker daemon accessibility
ssh user@host "docker ps"

# Verify firewall settings
# Ensure port 22 (SSH) and Docker daemon port are open
```

### Getting Help

- **Issues**: [GitHub Issues](https://github.com/wikczerski/whaletui/issues)
- **Discussions**: [GitHub Discussions](https://github.com/wikczerski/whaletui/discussions)
- **Documentation**: [Project Wiki](https://github.com/wikczerski/whaletui/wiki)

### Logs and Debugging

Enable debug logging for troubleshooting:

```bash
# Run with debug logging
whaletui --log-level DEBUG

# Check application logs
tail -f ~/.whaletui/logs/app.log
```

## Uninstallation

### Remove Binary

**Homebrew Installation:**
```bash
# Remove whaletui installed via Homebrew
brew uninstall whaletui

# Remove the tap only if you added it manually and no longer need it
# (Not needed if whaletui is in the main Homebrew repository)
brew untap wikczerski/tap
```

**Go Install:**
```bash
# Remove Go installation
go clean -i github.com/wikczerski/whaletui
```

**Manual Installation:**
```bash
# Remove from PATH directory
sudo rm /usr/local/bin/whaletui
```

### Remove Configuration
```bash
# Remove config directory
rm -rf ~/.config/whaletui

# Remove logs
rm -rf ~/.whaletui
```

### Remove Source Code
```bash
# Remove cloned repository
rm -rf ~/whaletui
```

---

For more information, see the [README.md](README.md) and [CONTRIBUTING.md](CONTRIBUTING.md) files.
