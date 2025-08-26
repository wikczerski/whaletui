---
id: networks
title: Docker Networks
sidebar_label: Networks
description: Understanding Docker networking and network management in WhaleTUI
---

# Docker Networks

Docker networks enable containers to communicate with each other and with external systems. They provide isolation, security, and connectivity for your containerized applications.

## Table of Contents

- [What are Docker Networks?](#what-are-docker-networks)
- [Network Types](#network-types)
- [Network Management in WhaleTUI](#network-management-in-whaletui)
- [Common Network Operations](#common-network-operations)
- [Network Configuration](#network-configuration)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## What are Docker Networks?

Docker networks are virtual networks that allow containers to communicate with each other and with external networks. They provide network isolation, security, and connectivity management for containerized applications.

### Key Concepts

- **Network Isolation**: Containers in different networks cannot communicate directly
- **Service Discovery**: Containers can find each other by name within the same network
- **Port Mapping**: Expose container ports to the host or external networks
- **Load Balancing**: Distribute traffic across multiple containers
- **Security**: Control which containers can communicate with each other

## Network Types

### Bridge Networks

The default network type, providing automatic DNS resolution between containers:

```bash
# Create a custom bridge network
docker network create my-network

# Run containers on the network
docker run --network my-network --name web nginx
docker run --network my-network --name db postgres
```

**Characteristics:**
- Automatic DNS resolution
- Port publishing to host
- Isolated from other networks
- Good for single-host deployments

### Host Networks

Containers share the host's network stack directly:

```bash
# Run container with host networking
docker run --network host nginx
```

**Characteristics:**
- No network isolation
- Direct access to host ports
- Better performance
- Less secure

### Overlay Networks

Enable communication between containers across multiple Docker hosts:

```bash
# Create overlay network for swarm
docker network create --driver overlay my-overlay
```

**Characteristics:**
- Multi-host communication
- Swarm service integration
- Encrypted communication
- Complex configuration

### Macvlan Networks

Assign MAC addresses to containers for direct network access:

```bash
# Create macvlan network
docker network create -d macvlan \
  --subnet=192.168.1.0/24 \
  --gateway=192.168.1.1 \
  -o parent=eth0 my-macvlan
```

**Characteristics:**
- Direct network access
- MAC address assignment
- VLAN support
- Network performance

### None Networks

Complete network isolation:

```bash
# Run container without network access
docker run --network none alpine
```

**Characteristics:**
- No network connectivity
- Maximum isolation
- Useful for security testing
- Limited functionality

## Network Management in WhaleTUI

WhaleTUI provides an intuitive interface for managing Docker networks through the terminal.

### Network List View

View all available networks with key information:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WhaleTUI - Docker Networks                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ NETWORK ID     NAME           DRIVER    SCOPE    CONTAINERS                 │
│ abc123def456   bridge         bridge    local    3                          │
│ def456ghi789   host           host      local    0                          │
│ ghi789jkl012   none           null      local    1                          │
│ jkl012mno345   my-network     bridge    local    2                          │
│ mno345pqr678   my-overlay     overlay   swarm    5                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Network Details

Inspect network configuration and connected containers:

- **Network ID**: Unique identifier
- **Driver**: Network driver type
- **Scope**: Local or swarm
- **Subnet**: IP address range
- **Gateway**: Default gateway
- **Connected Containers**: List of attached containers

## Common Network Operations

### Creating Networks

```bash
# Create bridge network with custom subnet
docker network create \
  --driver bridge \
  --subnet 172.18.0.0/16 \
  --gateway 172.18.0.1 \
  my-custom-network

# Create network with labels
docker network create \
  --label environment=production \
  --label team=backend \
  prod-network
```

### Connecting Containers

```bash
# Connect existing container to network
docker network connect my-network container-name

# Run container on specific network
docker run --network my-network --name web nginx

# Connect to multiple networks
docker run --network frontend --network backend --name app myapp
```

### Disconnecting Containers

```bash
# Disconnect container from network
docker network disconnect my-network container-name

# Force disconnect (even if container is running)
docker network disconnect -f my-network container-name
```

### Removing Networks

```bash
# Remove unused network
docker network rm my-network

# Remove all unused networks
docker network prune

# Force remove network with connected containers
docker network rm -f my-network
```

## Network Configuration

### IP Address Assignment

```bash
# Specify static IP for container
docker run --network my-network \
  --ip 172.18.0.10 \
  --name web nginx

# Create network with specific IP range
docker network create \
  --subnet 10.0.0.0/24 \
  --ip-range 10.0.0.0/25 \
  my-subnet
```

### DNS Configuration

```bash
# Set custom DNS servers
docker run --dns 8.8.8.8 --dns 8.8.4.4 nginx

# Use host DNS
docker run --dns-search example.com nginx

# Custom DNS options
docker run --dns-opt use-vc nginx
```

### Port Publishing

```bash
# Publish specific port
docker run -p 8080:80 nginx

# Publish all ports
docker run -P nginx

# Bind to specific interface
docker run -p 127.0.0.1:8080:80 nginx

# Publish range of ports
docker run -p 8080-8090:80 nginx
```

## Best Practices

### Network Design

- **Use Custom Networks**: Avoid default bridge for production
- **Network Segmentation**: Separate frontend, backend, and database
- **Service Discovery**: Leverage Docker's built-in DNS
- **Security Groups**: Use network isolation for security

### Performance Optimization

- **Choose Appropriate Driver**: Bridge for single-host, overlay for multi-host
- **Minimize Network Hops**: Keep related containers on same network
- **Use Host Networking**: For performance-critical applications
- **Monitor Network Usage**: Track bandwidth and latency

### Security Considerations

- **Network Isolation**: Separate networks for different environments
- **Port Exposure**: Only publish necessary ports
- **Access Control**: Limit container-to-container communication
- **Encryption**: Use overlay networks with encryption for sensitive data

## Troubleshooting

### Common Issues

#### Container Communication Problems

```bash
# Check network connectivity
docker exec container-name ping other-container

# Verify network configuration
docker network inspect network-name

# Check container network settings
docker inspect container-name | grep -A 20 "NetworkSettings"
```

#### DNS Resolution Issues

```bash
# Test DNS resolution
docker exec container-name nslookup other-container

# Check DNS configuration
docker exec container-name cat /etc/resolv.conf

# Restart network
docker network disconnect network-name container-name
docker network connect network-name container-name
```

#### Port Binding Problems

```bash
# Check port usage
netstat -tulpn | grep :8080

# Verify port mapping
docker port container-name

# Check firewall settings
sudo ufw status
```

### Network Diagnostics

```bash
# Inspect network details
docker network inspect network-name

# View network logs
docker logs container-name

# Check network statistics
docker stats container-name

# Monitor network traffic
docker exec container-name tcpdump -i eth0
```

## Advanced Topics

### Network Policies

Control traffic flow between containers:

```bash
# Create network with custom policies
docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=br0 \
  --opt com.docker.network.bridge.enable_icc=false \
  isolated-network
```

### Service Mesh Integration

Integrate with service mesh solutions:

- **Istio**: Advanced traffic management
- **Linkerd**: Lightweight service mesh
- **Consul Connect**: Service-to-service communication
- **Traefik**: Reverse proxy and load balancer

### Network Monitoring

Monitor network performance and health:

```bash
# Install network monitoring tools
docker run -d --name netdata \
  -p 19999:19999 \
  -v /proc:/host/proc:ro \
  -v /sys:/host/sys:ro \
  -v /var/run/docker.sock:/var/run/docker.sock \
  netdata/netdata
```

### Swarm Mode

In swarm mode, overlay networks span across multiple nodes:

```bash
# Create overlay network for swarm
docker network create --driver overlay --subnet 10.0.9.0/24 my-network

# Deploy service using overlay network
docker service create --name web --network my-network nginx:latest

# Check network connectivity across nodes
docker network inspect my-network
```

**Key Concepts:**
- **Node Spanning**: Overlay networks automatically span across all swarm nodes
- **Service Discovery**: Services can discover each other across nodes
- **Load Balancing**: Built-in load balancing across network endpoints
- **Security**: Encrypted communication between nodes

## Integration with WhaleTUI

WhaleTUI seamlessly integrates network management with other Docker operations:

- **Container Management**: Create and manage containers on networks
- **Service Deployment**: Deploy services with network configuration
- **Volume Management**: Persist data across network changes
- **Image Management**: Use images with network requirements

### Keyboard Shortcuts

- `:networks` - Switch to Networks view
- `c` - Create new network (TBD - may not be implemented yet)
- `i` - Inspect network details (TBD - may not be implemented yet)
- `r` - Remove network (TBD - may not be implemented yet)
- `Enter` - View network details

> **Note**: Some features marked as TBD (To Be Determined) may not be fully implemented in the current version of WhaleTUI. Check the latest release notes for current feature availability.

### Context Menus

Use keyboard shortcuts and navigation for quick access to:
- Network inspection
- Container management
- Configuration editing
- Removal options

## Next Steps

- [Containers](./containers.md) - Learn about Docker containers
- [Images](./images.md) - Manage Docker images
- [Volumes](./volumes.md) - Manage persistent storage
- [Swarm](./swarm.md) - Deploy overlay networks in a swarm
- [Nodes](./nodes.md) - Manage swarm nodes
- [Development Setup](../development/setup.md) - Set up your development environment
