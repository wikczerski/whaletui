---
id: containers
title: Docker Containers
sidebar_label: Containers
---

# Docker Containers

Containers are the fundamental building blocks of Docker applications. In WhaleTUI, you can manage containers efficiently through an intuitive interface that simplifies common operations.

## What are Containers?

Docker containers are lightweight, standalone, and executable packages that include everything needed to run an application: code, runtime, system tools, system libraries, and settings.

## Container Lifecycle in WhaleTUI

### 1. Viewing Containers

The **Containers** tab in WhaleTUI displays all containers on your system:

- **Running Containers**: Currently active containers
- **Stopped Containers**: Containers that have been stopped
- **Container Status**: Health, uptime, and resource usage

### 2. Container Information

Each container displays key information:

- **Name**: Human-readable container identifier
- **Image**: Base image used to create the container
- **Status**: Current state (running, stopped, paused)
- **Ports**: Exposed ports and mappings
- **Size**: Disk space used by the container

## Basic Container Operations

### Starting a Container

```bash
# Using WhaleTUI UI
1. Navigate to Containers tab
2. Select a stopped container
3. Press 's' to start a container
4. Container will begin running

# Using Docker CLI (alternative)
docker start <container_id>
```

### Stopping a Container

```bash
# Using WhaleTUI UI
1. Select a running container
2. Press 'S' to stop a container
3. Container will stop gracefully

# Using Docker CLI (alternative)
docker stop <container_id>
```

### Restarting a Container

```bash
# Using WhaleTUI UI
1. Select a container (running or stopped)
2. Press 'r' to restart a container
3. Container will restart with the same configuration

# Using Docker CLI (alternative)
docker restart <container_id>
```

### Removing a Container

```bash
# Using WhaleTUI UI

1. Select a stopped container using arrow keys
2. Press 'd' to delete the container
3. Confirm the action with 'Enter'
4. Container will be permanently removed

# Using Docker CLI (alternative)
docker rm <container_id>
```

## Advanced Container Operations

### Viewing Container Logs

```bash
# Using WhaleTUI UI
1. Select a container
2. Press 'l' to open logs view
3. Navigate through log entries
4. Press 'Enter' to return to main view

# Using Docker CLI (alternative)
docker logs <container_id>
```

### Executing Commands

```bash
# Using WhaleTUI UI - Single Command Execution
1. Select a running container using arrow keys
2. Press 'e' to execute a single command
3. Enter the command to run
4. View command output
5. Return to container list

# Using WhaleTUI UI - Interactive Shell (Attach)
1. Select a running container using arrow keys
2. Press 'a' to attach to container's main process
3. Interact with the running application
4. Use Ctrl+C or Ctrl+D to detach

# Using Docker CLI - Single Command
docker exec <container_id> <command>
docker exec my-container ls -la

# Using Docker CLI - Interactive Shell
docker exec -it <container_id> /bin/bash
docker attach <container_id>
```

### Key Differences

- **Exec (`e`)**: Runs a single command and returns to WhaleTUI
- **Attach (`a`)**: Connects to the container's main process for interactive use

### Inspecting Containers

```bash
# Using WhaleTUI UI
1. Select a container
2. Press 'i' to view detailed information
3. Review configuration, networking, and resources
4. Press 'BackSpace' to return to main view

# Using Docker CLI (alternative)
docker inspect <container_id>
```

## Container Configuration

### Environment Variables

Containers can be configured with environment variables:

```bash
# Set environment variables when running
docker run -e VAR_NAME=value -e ANOTHER_VAR=value image_name

# In WhaleTUI, use the Run dialog to configure environment variables
# TBD
```

### Port Mappings

Expose container ports to the host:

```bash
# Map host port 8080 to container port 80
docker run -p 8080:80 nginx

# In WhaleTUI, configure port mappings in the Run dialog
# TBD
```

### Volume Mounts

Persist data between container runs:

```bash
# Mount host directory to container
docker run -v /host/path:/container/path image_name

# In WhaleTUI, configure volume mounts in the Run dialog
# TBD
```

## Container Resource Management

### Resource Limits

Set limits on container resource usage:

```bash
# Limit memory and CPU
docker run --memory=512m --cpus=1.0 image_name

# In WhaleTUI, configure resource limits in the Run dialog
# TBD
```

### Resource Monitoring

WhaleTUI provides real-time resource monitoring:

- **CPU Usage**: Current CPU consumption
- **Memory Usage**: RAM usage and limits
- **Network I/O**: Network traffic statistics
- **Disk I/O**: Storage access patterns

## Best Practices

### Container Naming

Use descriptive names for containers:

```bash
# Good naming examples
web-server-prod
database-staging
api-gateway-dev

# Avoid generic names
container1
test
temp
```

### Resource Cleanup

Regularly clean up unused containers:

```bash
# Remove stopped containers
docker container prune

# In WhaleTUI, use the Cleanup action to remove stopped containers
# TBD
```

### Health Checks

Implement health checks for production containers:

```bash
# Add health check to Dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost/ || exit 1
```

## Troubleshooting

### Common Issues

#### Container Won't Start
```bash
# Check container logs
docker logs <container_id>

# Verify image exists
docker images

# Check resource availability
docker system df
```

#### Container Performance Issues
```bash
# Monitor resource usage
docker stats <container_id>

# Check for resource limits
docker inspect <container_id> | grep -i memory
```

#### Network Connectivity
```bash
# Test container networking
docker exec <container_id> ping google.com

# Check port bindings
docker port <container_id>
```

## Next Steps

- [Images](./images.md) - Learn about Docker images
- [Networks](./networks.md) - Configure container networking
- [Volumes](./volumes.md) - Manage persistent storage
- [Swarm](./swarm.md) - Deploy containers in a swarm
- [Nodes](./nodes.md) - Manage swarm nodes
- [Development Setup](../development/setup.md) - Set up your development environment

## Related Topics

- [Docker Images](images.md)
- [Container Networks](networks.md)
- [Data Volumes](volumes.md)
- [Container Orchestration](swarm.md)

### Swarm Mode

When running in swarm mode, containers are distributed across multiple nodes:

```bash
# Deploy service across multiple nodes
docker service create --name web --replicas 3 nginx:latest

# Scale service across available nodes
docker service scale web=5

# Check service distribution across nodes
docker service ps web
```

**Key Concepts:**
- **Node Distribution**: Containers are automatically distributed across available nodes
- **Load Balancing**: Swarm provides built-in load balancing for services
- **High Availability**: Services continue running even if individual nodes fail
- **Resource Management**: Swarm considers node resources when scheduling containers
