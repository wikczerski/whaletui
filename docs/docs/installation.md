---
id: installation
title: Installation
sidebar_label: Installation
---

# Installation Guide

This guide will walk you through installing WhaleTUI on your system. WhaleTUI is available for multiple platforms and can be installed using various methods.

## Prerequisites

Before installing WhaleTUI, ensure you have the following following prerequisites:

- **Docker**: WhaleTUI requires Docker to be installed and running on your system
- **Go 1.21+**: Required for building from source
- **Git**: For cloning the repository

## Installation Methods

### Method 1: Download Pre-built Binary (Recommended)

The easiest way to get started with WhaleTUI is to download a pre-built binary for your platform.

#### Linux

```bash
# Download the latest release
wget https://github.com/wikczerski/whaletui/releases/latest/download/whaletui-linux-amd64

# Make it executable
chmod +x whaletui-linux-amd64

# Move to a directory in your PATH
sudo mv whaletui-linux-amd64 /usr/local/bin/whaletui
```

#### macOS

```bash
# Using Homebrew (if available)
brew install wikczerski/whaletui/whaletui

# Manual installation
curl -L https://github.com/wikczerski/whaletui/releases/latest/download/whaletui-darwin-amd64 -o whaletui
chmod +x whaletui
sudo mv whaletui /usr/local/bin/
```

#### Windows

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/wikczerski/whaletui/releases/latest/download/whaletui-windows-amd64.exe" -OutFile "whaletui.exe"

# Move to a directory in your PATH
Move-Item whaletui.exe "C:\Windows\System32\"
```

### Method 2: Build from Source

If you prefer to build WhaleTUI from source or need a custom build, follow these steps:

#### Clone the Repository

```bash
git clone https://github.com/wikczerski/whaletui.git
cd whaletui
```

#### Build the Application

```bash
# Build for your current platform
go build -o bin/whaletui cmd/root.go

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o bin/whaletui-linux-amd64 cmd/root.go
GOOS=darwin GOARCH=amd64 go build -o bin/whaletui-darwin-amd64 cmd/root.go
GOOS=windows GOARCH=amd64 go build -o bin/whaletui-windows-amd64.exe cmd/root.go
```

#### Install the Built Binary

```bash
# Move the binary to your PATH
sudo mv bin/whaletui /usr/local/bin/
```

### Method 3: Using Go Install

If you have Go installed, you can install WhaleTUI directly:

```bash
go install github.com/wikczerski/whaletui/cmd/root@latest
```

## Verification

After installation, verify that WhaleTUI is working correctly:

```bash
whaletui --version
```

You should see output similar to:
```
WhaleTUI version 1.0.0
```

## Configuration

WhaleTUI uses configuration files to customize its behavior. The default configuration file is located at:

- **Linux/macOS**: `~/.config/whaletui/config.yaml`
- **Windows**: `%APPDATA%\whaletui\config.yaml`

### Basic Configuration

Create a basic configuration file:

```yaml
# ~/.config/whaletui/config.yaml
docker:
  host: "unix:///var/run/docker.sock"
  timeout: 30s

ui:
  theme: "dark"
  refresh_rate: 1s

logging:
  level: "info"
  file: "~/.local/share/whaletui/logs/whaletui.log"
```

## Docker Integration

WhaleTUI integrates with Docker through the Docker Engine API. Ensure that:

1. Docker daemon is running
2. Your user has permission to access the Docker socket
3. Docker API is accessible (usually on `unix:///var/run/docker.sock`)

### Docker Permissions

On Linux, you may need to add your user to the `docker` group:

```bash
sudo usermod -aG docker $USER
newgrp docker
```

## Troubleshooting

### Common Issues

#### Permission Denied
If you encounter permission issues with Docker:

```bash
# Check Docker socket permissions
ls -la /var/run/docker.sock

# Fix permissions if needed
sudo chmod 666 /var/run/docker.sock
```

#### Connection Refused
If WhaleTUI can't connect to Docker:

```bash
# Check if Docker is running
sudo systemctl status docker

# Start Docker if it's not running
sudo systemctl start docker
```

#### Binary Not Found
If the `whaletui` command is not found:

```bash
# Check if the binary is in your PATH
which whaletui

# Add the directory to your PATH if needed
export PATH=$PATH:/path/to/whaletui
```

## Next Steps

Now that you have WhaleTUI installed, you can:

1. **Start WhaleTUI**: Run `whaletui` to launch the application
2. **Read the Quick Start Guide**: Learn the basics of using WhaleTUI
3. **Explore Features**: Discover all the capabilities of WhaleTUI
4. **Join the Community**: Get help and contribute to the project

## Support

If you encounter issues during installation:

- Check the [Troubleshooting section](#troubleshooting)
- Review the [GitHub Issues](https://github.com/wikczerski/whaletui/issues)

## Uninstallation

To remove WhaleTUI from your system:

```bash
# Remove the binary
sudo rm /usr/local/bin/whaletui

# Remove configuration files
rm -rf ~/.config/whaletui
rm -rf ~/.local/share/whaletui
```
