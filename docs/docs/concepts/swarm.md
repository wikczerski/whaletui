---
id: swarm
title: Docker Swarm
sidebar_label: Swarm
description: Understanding Docker Swarm and orchestration in WhaleTUI
---

# Docker Swarm

Docker Swarm is Docker's native clustering and orchestration solution that allows you to create and manage a cluster of Docker nodes as a single virtual system.

## Table of Contents

- [What is Docker Swarm?](#what-is-docker-swarm)
- [Swarm Architecture](#swarm-architecture)
- [Swarm Management in WhaleTUI](#swarm-management-in-whaletui)
- [Node Management](#node-management)
- [Service Management](#service-management)
- [Stack Deployment](#stack-deployment)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## What is Docker Swarm?

Docker Swarm is a container orchestration tool that enables you to manage a cluster of Docker hosts as a single system. It provides high availability, load balancing, scaling, and service discovery for containerized applications.

### Key Features

- **High Availability**: Automatic failover and recovery
- **Load Balancing**: Distribute traffic across multiple containers
- **Scaling**: Scale services up or down based on demand
- **Service Discovery**: Automatic DNS resolution between services
- **Rolling Updates**: Zero-downtime deployments
- **Security**: Encrypted communication and access control

### Use Cases

- **Production Deployments**: High-availability applications
- **Microservices**: Distributed application architectures
- **Load Balancing**: Traffic distribution across containers
- **Disaster Recovery**: Multi-node redundancy
- **Development Testing**: Local swarm for testing

## Swarm Architecture

### Swarm Components

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Docker Swarm Cluster                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   Manager Node  │    │   Manager Node  │    │   Manager Node  │        │
│  │                 │    │                 │    │                 │        │
│  │ • Swarm Manager│    │ • Swarm Manager│    │ • Swarm Manager│        │
│  │ • API Endpoint │    │ • API Endpoint │    │ • API Endpoint │        │
│  │ • Raft Consensus│    │ • Raft Consensus│    │ • Raft Consensus│        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
│           │                       │                       │                │
│           └───────────────────────┼───────────────────────┘                │
│                                   │                                        │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   Worker Node   │    │   Worker Node   │    │   Worker Node   │        │
│  │                 │    │                 │    │                 │        │
│  │ • Task Execution│    │ • Task Execution│    │ • Task Execution│        │
│  │ • Container Run │    │ • Container Run │    │ • Container Run │        │
│  │ • Resource Mgmt │    │ • Resource Mgmt │    │ • Resource Mgmt │        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Node Types

#### Manager Nodes

- **Swarm Management**: Control cluster operations
- **API Endpoint**: Handle client requests
- **Raft Consensus**: Maintain cluster state
- **Task Scheduling**: Distribute work across nodes

#### Worker Nodes

- **Task Execution**: Run containerized applications
- **Resource Management**: Allocate CPU, memory, storage
- **Health Monitoring**: Report node status
- **Load Distribution**: Accept assigned tasks

## Swarm Management in WhaleTUI

WhaleTUI provides comprehensive tools for managing Docker Swarm through an intuitive terminal interface.

### Swarm Overview

View cluster status and information:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WhaleTUI - Docker Swarm Overview                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ CLUSTER ID: abc123def456                                                   │
│ STATUS: Active                                                             │
│ NODES: 3 (2 managers, 1 worker)                                           │
│ SERVICES: 5                                                                │
│ TASKS: 15                                                                  │
│ STACKS: 2                                                                  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Node List View

Monitor all nodes in the cluster:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WhaleTUI - Swarm Nodes                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ ID           HOSTNAME      STATUS    AVAILABILITY    MANAGER STATUS        │
│ abc123...    manager-1     Ready     Active          Leader                │
│ def456...    manager-2     Ready     Active          Reachable             │
│ ghi789...    worker-1      Ready     Active          -                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Node Management

### Initializing Swarm

```bash
# Initialize swarm on first manager
docker swarm init --advertise-addr 192.168.1.100

# Join worker node to swarm
docker swarm join --token SWMTKN-1-... 192.168.1.100:2377

# Join manager node to swarm
docker swarm join --token SWMTKN-1-... 192.168.1.100:2377
```

### Managing Nodes

```bash
# Promote worker to manager
docker node promote worker-1

# Demote manager to worker
docker node demote manager-2

# Update node availability
docker node update --availability drain worker-1

# Remove node from swarm
docker swarm leave --force
```

### Node Inspection

```bash
# Get node information
docker node inspect node-name

# List all nodes
docker node ls

# View node tasks
docker node ps node-name
```

## Service Management

### Creating Services

```bash
# Basic service
docker service create --name web nginx:latest

# Service with replicas
docker service create --name api --replicas 3 myapp:latest

# Service with constraints
docker service create \
  --name db \
  --constraint 'node.role==manager' \
  postgres:15

# Service with resources
docker service create \
  --name compute \
  --limit-cpu 2 \
  --limit-memory 1g \
  --reserve-cpu 1 \
  --reserve-memory 512m \
  myapp:latest
```

### Service Configuration

```bash
# Service with environment variables
docker service create \
  --name app \
  --env NODE_ENV=production \
  --env DB_HOST=db \
  myapp:latest

# Service with secrets
docker service create \
  --name secure-app \
  --secret db-password \
  myapp:latest

# Service with configs
docker service create \
  --name config-app \
  --config source=app-config,target=/app/config.yml \
  myapp:latest
```

### Service Operations

```bash
# Scale service
docker service scale web=5

# Update service
docker service update --image nginx:1.21 web

# Remove service
docker service rm web

# View service logs
docker service logs web
```

## Stack Deployment

### Stack Definition

Deploy multiple services using Docker Compose:

```yaml
# docker-stack.yml
version: '3.8'

services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure

  api:
    image: myapp:latest
    environment:
      - NODE_ENV=production
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    secrets:
      - db-password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    deploy:
      placement:
        constraints:
          - node.role == manager

volumes:
  postgres-data:

secrets:
  db-password:
    external: true
```

### Stack Operations

```bash
# Deploy stack
docker stack deploy -c docker-stack.yml myapp

# List stacks
docker stack ls

# List stack services
docker stack services myapp

# List stack tasks
docker stack tasks myapp

# Remove stack
docker stack rm myapp
```

## Best Practices

### Cluster Design

- **Manager Nodes**: Use odd number (3, 5, 7) for consensus
- **Worker Distribution**: Distribute workload evenly
- **Resource Planning**: Plan for resource requirements
- **Network Design**: Use overlay networks for service communication

### Service Design

- **Stateless Services**: Design for horizontal scaling
- **Health Checks**: Implement proper health check endpoints
- **Resource Limits**: Set appropriate CPU and memory limits
- **Update Strategy**: Use rolling updates for zero downtime

### Security

- **TLS Encryption**: Enable TLS for node communication
- **Access Control**: Use secrets for sensitive data
- **Network Policies**: Implement network segmentation
- **Regular Updates**: Keep swarm and services updated

### Monitoring

- **Health Checks**: Monitor service health
- **Resource Usage**: Track CPU, memory, and network usage
- **Log Aggregation**: Centralize service logs
- **Alerting**: Set up alerts for critical issues

## Troubleshooting

### Common Issues

#### Node Communication Problems

```bash
# Check node status
docker node ls

# Inspect node details
docker node inspect node-name

# Check network connectivity
ping node-ip-address

# Verify swarm tokens
docker swarm join-token manager
docker swarm join-token worker
```

#### Service Deployment Issues

```bash
# Check service status
docker service ls

# Inspect service details
docker service inspect service-name

# View service logs
docker service logs service-name

# Check task status
docker service ps service-name
```

#### Scaling Problems

```bash
# Check resource availability
docker node inspect node-name | grep -A 10 "Status"

# Verify service constraints
docker service inspect service-name | grep -A 10 "Constraints"

# Check node resources
docker system df
```

### Swarm Diagnostics

```bash
# Check swarm status
docker info | grep -A 10 "Swarm"

# View swarm logs
sudo journalctl -u docker.service | grep swarm

# Inspect swarm configuration
docker swarm inspect

# Check cluster health
docker node ls --format "table {{.ID}}\t{{.Hostname}}\t{{.Status}}\t{{.Availability}}"
```

## Advanced Topics

### Rolling Updates

```bash
# Update service with rolling strategy
docker service update \
  --image nginx:1.21 \
  --update-parallelism 2 \
  --update-delay 10s \
  web

# Rollback service
docker service rollback web
```

### Service Dependencies

```bash
# Service with dependency
docker service create \
  --name app \
  --network app-network \
  --depends-on db \
  myapp:latest
```

### Load Balancing

```bash
# Service with load balancing
docker service create \
  --name web \
  --publish mode=host,target=80,published=8080 \
  --publish mode=host,target=443,published=8443 \
  nginx:latest
```

### Secrets Management

```bash
# Create secret
echo "mypassword" | docker secret create db-password -

# List secrets
docker secret ls

# Inspect secret
docker secret inspect secret-name

# Remove secret
docker secret rm secret-name
```

## Integration with WhaleTUI

WhaleTUI seamlessly integrates swarm management with other Docker operations:

- **Container Management**: Manage containers within swarm context
- **Network Configuration**: Use overlay networks for services
- **Volume Management**: Configure volumes for swarm services
- **Image Management**: Deploy images as swarm services

### Keyboard Shortcuts

- `:swarm` - Switch to Swarm Services view
- `:services` - Switch to Swarm Services view (alternative)
- `:nodes` - Switch to Swarm Nodes view
- `Enter` - Inspect selected item
- `i` - Inspect selected item
- `ESC` - Close modal or go back

> **Note**: Swarm management in WhaleTUI provides views for both services and nodes. For creating, modifying, and removing swarm resources, use the Docker CLI directly.

### Available Actions

Use keyboard shortcuts and navigation for quick access to:
- Service inspection and details
- Node inspection and details
- Viewing swarm information
- Container management (via Docker CLI)

## Next Steps

- [Containers](./containers.md) - Learn about Docker containers
- [Images](./images.md) - Manage Docker images
- [Networks](./networks.md) - Configure overlay networks
- [Volumes](./volumes.md) - Manage persistent storage
- [Nodes](./nodes.md) - Manage swarm nodes
- [Development Setup](../development/setup.md) - Set up your development environment
