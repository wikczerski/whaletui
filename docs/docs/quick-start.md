---
id: quick-start
title: Quick Start Guide
sidebar_label: Quick Start
---

# Quick Start Guide

Get up and running with WhaleTUI in minutes! This guide will walk you through the essential steps to start managing your Docker containers with WhaleTUI.

## Prerequisites

Before you begin, ensure you have:

- ✅ WhaleTUI installed (see [Installation Guide](installation.md))
- ✅ Docker running on your system
- ✅ At least one Docker container or image available

## Launch WhaleTUI

Start WhaleTUI by running the following command in your terminal:

```bash
whaletui
```

You should see the WhaleTUI interface appear with a clean, organized layout.

## First Steps

### 1. Navigate the Interface

WhaleTUI uses a command-based interface for different Docker resources:

- **Containers**: Press `:` then type `containers` to view and manage running and stopped containers
- **Images**: Press `:` then type `images` to browse available Docker images
- **Networks**: Press `:` then type `networks` to manage Docker networks
- **Volumes**: Press `:` then type `volumes` to handle Docker volumes
- **Swarm**: Press `:` then type `swarm` to access Docker Swarm features (if enabled)

Use `Tab` to switch between tabs and `↑/↓` arrows to navigate within each tab.

### 2. View Your Containers

Start by checking the **Containers** tab to see what's currently running on your system.

If you don't have any containers running, you can start one:

```bash
# In another terminal, start a simple container
docker run -d --name nginx-test nginx:alpine
```

### 3. Basic Container Operations

With a container running, you can:

- **Start/Stop**: Use the action buttons or keyboard shortcuts
- **View Logs**: Select a container using arrow keys and press `l` to see real-time logs
- **Execute Commands**: Press `e` to run a single command, or `a` to attach to interactive shell
- **Inspect**: Press `I` to view detailed container information

## Essential Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `:` | Open command view for view switching |
| `:containers` | Switch to Containers view |
| `:images` | Switch to Images view |
| `:networks` | Switch to Networks view |
| `:volumes` | Switch to Volumes view |
| `:swarm` | Switch to Swarm view |
| `↑/↓` | Navigate items |
| `Enter` | Select item or execute action |
| `Space` | Toggle selection |
| `l` | View logs |
| `e` | Execute single command |
| `a` | Attach to interactive shell |
| `i` | Inspect item |
| `d` | Delete item |
| `r` | Refresh view |
| `q` | Quit application |
| `?` | Show help |

## Common Operations

### Starting a Container

1. Navigate to the **Images** tab
2. Select an image (e.g., `nginx:alpine`)
3. Press `Enter` or use the "Run" action
4. Configure container options if prompted
5. Press `Enter` to start the container

### Stopping a Container

1. Go to the **Containers** tab
2. Select a running container
3. Press `S` or use the "Stop" action
4. Confirm the action

### Viewing Container Logs

1. Select a container in the **Containers** tab
2. Press `L` to open the logs view
3. Use `↑/↓` to scroll through logs
4. Press `Q` to return to the main view

### Executing Commands

1. Select a running container
2. Press `E` to open an interactive shell
3. Type your commands as usual
4. Press `Ctrl+D` or type `exit` to close the shell

## Configuration

### Basic Settings

WhaleTUI can be configured through a configuration file. Create `~/.config/whaletui/config.yaml`:

```yaml
docker:
  host: "unix:///var/run/docker.sock"
  timeout: 30s

ui:
  theme: "dark"
  refresh_rate: 1s
  show_help: true

logging:
  level: "info"
```

### Themes

WhaleTUI supports multiple themes:

- **Dark**: Dark theme provided in repository
- **Default**: Default theme loaded when none provided
- **Custom**: Define your own color scheme

Change themes in the configuration file or use the `--theme` flag:

```bash
whaletui --theme <theme_file_path>
```

## Next Steps

Now that you're comfortable with the basics:

1. **Explore Concepts**: Dive deeper into [Docker concepts](concepts/containers.md)
2. **Customize Your Setup**: Learn about [Configuration Options](installation.md#configuration)
3. **Connect to Remote Hosts**: Set up [SSH Connections](concepts/containers.md#executing-commands)
4. **Advanced Operations**: Master [Swarm management](concepts/swarm.md)

## Troubleshooting

If you encounter issues:

- **Check Docker Status**: Ensure Docker is running (`docker info`)
- **Verify Permissions**: Make sure you have access to Docker daemon
- **Review Logs**: Check WhaleTUI logs for error messages
- **Update WhaleTUI**: Ensure you're running the latest version

## Additional Resources

- **Documentation**: Explore all available guides and concepts
- **GitHub**: Visit our [GitHub repository](https://github.com/wikczerski/whaletui)
- **Issues**: Report bugs and request features on GitHub
