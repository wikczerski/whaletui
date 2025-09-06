---
id: images
title: Docker Images
sidebar_label: Images
description: Understanding Docker images and image management in WhaleTUI
---

# Docker Images

Docker images are the foundation of containerized applications. They contain everything needed to run an application: code, runtime, system tools, libraries, and settings.

## Table of Contents

- [What are Docker Images?](#what-are-docker-images)
- [Image Layers](#image-layers)
- [Image Management in WhaleTUI](#image-management-in-whaletui)
- [Common Image Operations](#common-image-operations)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## What are Docker Images?

A Docker image is a lightweight, standalone, executable package that includes everything needed to run a piece of software. Images are built from a set of instructions called a Dockerfile and can be shared, versioned, and reused across different environments.

### Key Characteristics

- **Immutable**: Once created, images cannot be modified
- **Layered**: Images are built in layers for efficiency
- **Portable**: Can run on any system with Docker
- **Versioned**: Can be tagged with specific versions
- **Reusable**: One image can create multiple containers

## Image Layers

Docker images are built using a layered filesystem. Each layer represents a set of changes to the filesystem:

```
┌─────────────────────────────────────┐
│           Application Layer         │ ← Your application code
├─────────────────────────────────────┤
│         Runtime Dependencies        │ ← Libraries and frameworks
├─────────────────────────────────────┤
│         System Dependencies         │ ← OS packages and tools
├─────────────────────────────────────┤
│           Base OS Layer             │ ← Operating system
└─────────────────────────────────────┘
```

### Benefits of Layering

- **Efficient Storage**: Shared layers between images
- **Faster Builds**: Only rebuild changed layers
- **Better Caching**: Reuse existing layers
- **Smaller Images**: Optimize each layer independently

## Image Management in WhaleTUI

WhaleTUI provides comprehensive tools for managing Docker images through an intuitive terminal interface.

### Image List View

The image list shows all available images with key information:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WhaleTUI - Docker Images                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ REPOSITORY          TAG       IMAGE ID      CREATED        SIZE             │
│ nginx              latest     a8758716bb6a   2 weeks ago    133MB           │
│ postgres           15         c43a65fa73b9   3 weeks ago    379MB           │
│ redis              alpine     9b7c6093c358   1 month ago    32.1MB          │
│ ubuntu             22.04      ba6acccedd29   2 months ago   77.8MB          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Image Details

View detailed information about specific images:

- **Image ID**: Unique identifier
- **Created**: Build timestamp
- **Size**: Total size on disk
- **Layers**: Number of filesystem layers
- **Architecture**: Target system architecture
- **OS**: Base operating system

## Common Image Operations

### Pulling Images

Download images from Docker Hub or other registries:

```bash
# Pull latest version
docker pull nginx

# Pull specific version
docker pull postgres:15

# Pull from private registry
docker pull registry.example.com/myapp:v1.0
```

### Building Images

Create custom images from Dockerfiles:

```dockerfile
# Example Dockerfile
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y nginx
COPY index.html /var/www/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

Build command:
```bash
docker build -t myapp:latest .
```

### Tagging Images

Organize images with meaningful tags:

```bash
# Tag with version
docker tag myapp:latest myapp:v1.0

# Tag for specific environment
docker tag myapp:latest myapp:production

# Tag for registry
docker tag myapp:latest registry.example.com/myapp:latest
```

### Removing Images

Clean up unused images:

```bash
# Remove specific image
docker rmi nginx:latest

# Remove all unused images
docker image prune

# Remove all images
docker rmi $(docker images -q)
```

## Best Practices

### Image Size Optimization

- **Use Multi-stage Builds**: Separate build and runtime stages
- **Choose Base Images Wisely**: Prefer Alpine Linux for smaller images
- **Combine RUN Commands**: Reduce layer count
- **Clean Up Package Managers**: Remove cache and temporary files

```dockerfile
# Multi-stage build example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Security Considerations

- **Scan for Vulnerabilities**: Regular security scans
- **Use Official Images**: Prefer official, maintained images
- **Update Regularly**: Keep base images current
- **Minimize Attack Surface**: Only install necessary packages

### Tagging Strategy

- **Semantic Versioning**: Use version numbers (v1.0.0, v1.1.0)
- **Environment Tags**: Production, staging, development
- **Date Tags**: Timestamp-based tags for debugging
- **Latest Tag**: Always point to current stable version

## Troubleshooting

### Common Issues

#### Image Pull Failures

```bash
# Check network connectivity
docker pull hello-world

# Verify registry credentials
docker login

# Check image name and tag
docker search nginx
```

#### Build Failures

```bash
# Check Dockerfile syntax
docker build --no-cache .

# View build context
docker build --progress=plain .

# Check available disk space
df -h
```

#### Image Size Issues

```bash
# Analyze image layers
docker history nginx:latest

# Check image details
docker inspect nginx:latest

# Use dive for layer analysis
dive nginx:latest
```

### Performance Optimization

- **Use .dockerignore**: Exclude unnecessary files
- **Leverage Build Cache**: Structure Dockerfile for better caching
- **Parallel Downloads**: Use multi-stage builds effectively
- **Registry Mirrors**: Use local registry mirrors when possible

## Advanced Topics

### Image Signing

Verify image authenticity with digital signatures:

```bash
# Enable content trust
export DOCKER_CONTENT_TRUST=1

# Pull signed images
docker pull nginx:latest

# Sign your images
docker trust sign myapp:latest
```

### Multi-architecture Images

Build images for different CPU architectures:

```bash
# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .

# Push multi-arch image
docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest --push .
```

### Registry Management

Work with different image registries:

- **Docker Hub**: Public registry
- **GitHub Container Registry**: Integrated with GitHub
- **AWS ECR**: Amazon's container registry
- **Google Container Registry**: Google Cloud Platform
- **Azure Container Registry**: Microsoft Azure

### Swarm Mode

In swarm mode, images are automatically distributed across nodes:

```bash
# Deploy service with image to multiple nodes
docker service create --name app --replicas 3 myapp:latest

# Update service image across all nodes
docker service update --image myapp:v2 app

# Check image distribution across nodes
docker service ps app
```

**Key Concepts:**
- **Node Distribution**: Images are pulled to nodes where containers will run
- **Image Caching**: Nodes cache frequently used images for faster deployment
- **Rolling Updates**: Services can be updated with new images across nodes
- **Resource Optimization**: Swarm optimizes image distribution based on node resources

## Integration with WhaleTUI

WhaleTUI seamlessly integrates image management with other Docker operations:

- **Container Creation**: Use images to create containers
- **Service Deployment**: Deploy images as services
- **Network Configuration**: Connect containers from different images
- **Volume Management**: Persist data across container restarts

### Keyboard Shortcuts

- `:images` - Switch to Images view
- `Enter` - View image details
- `i` - Inspect image details
- `ESC` - Close modal or go back

> **Note**: Image management in WhaleTUI currently focuses on viewing and inspecting images. For pulling, building, tagging, and removing images, use the Docker CLI directly.

### Available Actions

Use keyboard shortcuts and navigation for quick access to:
- Image inspection and details
- Viewing image information
- Container creation (via Docker CLI)

## Next Steps

- [Containers](./containers.md) - Learn about Docker containers
- [Networks](./networks.md) - Configure container networking
- [Volumes](./volumes.md) - Manage persistent storage
- [Swarm](./swarm.md) - Deploy images in a swarm
- [Nodes](./nodes.md) - Manage swarm nodes
- [Development Setup](../development/setup.md) - Set up your development environment
