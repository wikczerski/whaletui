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

WhaleTUI uses command mode for switching between different Docker resources:

- Press `:` to enter command mode
- Type one of the following commands:
  - `containers` - View and manage running and stopped containers
  - `images` - Browse available Docker images
  - `volumes` - Handle Docker volumes
  - `networks` - Manage Docker networks
  - `swarm` - Access Docker Swarm services
  - `services` - Access Docker Swarm services (alternative)
  - `nodes` - View Docker Swarm nodes
- Press `Enter` to execute the command

Use `↑/↓` arrows to navigate within each view and `Enter` to select items.

### 2. View Your Containers

Start by pressing `:` then typing `containers` to view the **Containers** view and see what's currently running on your system.

If you don't have any containers running, you can start one:

```bash
# In another terminal, start a simple container
docker run -d --name nginx-test nginx:alpine
```

### 3. Basic Container Operations

With a container running, you can:

- **Start/Stop**: Press `s` to start or `S` to stop containers
- **Restart**: Press `r` to restart a container
- **Delete**: Press `d` to delete a container
- **View Logs**: Press `l` to see real-time logs
- **Inspect**: Press `i` to view detailed container information

## Essential Keyboard Shortcuts

### Navigation
| Command | Action |
|---------|--------|
| `:containers` | Switch to Containers view |
| `:images` | Switch to Images view |
| `:volumes` | Switch to Volumes view |
| `:networks` | Switch to Networks view |
| `:swarm` | Switch to Swarm Services view |
| `:services` | Switch to Swarm Services view (alternative) |
| `:nodes` | Switch to Swarm Nodes view |
| `↑/↓` | Navigate items |
| `Enter` | Select item or view details |
| `ESC` | Close modal or go back |

### Global Actions
| Key | Action |
|-----|--------|
| `F5` | Refresh current view |
| `?` | Show help |
| `Ctrl+C` or `q` | Quit application |

### Container Actions
| Key | Action |
|-----|--------|
| `s` | Start container |
| `S` | Stop container |
| `r` | Restart container |
| `d` | Delete container |
| `l` | View logs |
| `i` | Inspect container |

## Common Operations

### Starting a Container

1. Press `:` then type `images` to navigate to the **Images** view
2. Select an image (e.g., `nginx:alpine`) using arrow keys
3. Press `Enter` to view image details
4. Use available actions to run the container

### Stopping a Container

1. Press `:` then type `containers` to go to the **Containers** view
2. Select a running container using arrow keys
3. Press `S` to stop the container
4. Confirm the action if prompted

### Viewing Container Logs

1. Select a container in the **Containers** view
2. Press `l` to open the logs view
3. Use `↑/↓` to scroll through logs
4. Press `ESC` to return to the main view

### Inspecting Containers

1. Select a container in the **Containers** view
2. Press `i` to inspect the container
3. View detailed information in JSON format
4. Press `ESC` to return to the main view

## Configuration

### Basic Settings

WhaleTUI can be configured through a configuration file. Create `~/.whaletui/config.json`:

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

### Themes

WhaleTUI supports multiple themes:

- **Default**: Default theme loaded when none provided
- **Dark**: Dark theme provided in repository
- **Custom**: Define your own color scheme using YAML or JSON

Change themes in the configuration file by setting the `theme` field to your desired theme name.

## Next Steps

Now that you're comfortable with the basics:

1. **Explore Concepts**: Dive deeper into [Docker concepts](concepts/containers.md)
2. **Customize Your Setup**: Learn about [Configuration Options](installation.md#configuration)
3. **Connect to Remote Hosts**: Set up SSH connections for remote Docker management
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
