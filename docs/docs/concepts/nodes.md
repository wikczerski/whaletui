---
id: nodes
title: Docker Swarm Nodes
sidebar_label: Nodes
description: Understanding Docker Swarm nodes and node management in WhaleTUI
---

# Docker Swarm Nodes

Docker Swarm nodes are the individual Docker hosts that form a cluster. Each node can be either a manager node (responsible for cluster management) or a worker node (responsible for running containers).

## Table of Contents

- [What are Swarm Nodes?](#what-are-swarm-nodes)
- [Node Types](#node-types)
- [Node Management in WhaleTUI](#node-management-in-whaletui)
- [Common Node Operations](#common-node-operations)
- [Node Configuration](#node-configuration)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## What are Swarm Nodes?

Swarm nodes are individual Docker hosts that participate in a Docker Swarm cluster. Each node runs the Docker Engine and participates in the cluster's distributed system. Nodes can be physical machines, virtual machines, or cloud instances.

### Key Concepts

- **Cluster Membership**: Nodes join the swarm and become part of the cluster
- **Role Assignment**: Nodes are assigned either manager or worker roles
- **Resource Sharing**: Nodes contribute their resources to the cluster
- **Health Monitoring**: Nodes report their status and health to the swarm
- **Load Distribution**: Work is distributed across available nodes

## Node Types

### Manager Nodes

Manager nodes are responsible for cluster management and orchestration:

```bash
# Manager node responsibilities
- Maintain cluster state
- Schedule services
- Handle API requests
- Participate in Raft consensus
- Manage cluster membership
```

**Characteristics:**
- **Consensus**: Participate in Raft consensus for cluster state
- **Scheduling**: Make decisions about service placement
- **API**: Handle client requests and cluster operations
- **Leadership**: One manager serves as the leader
- **Reachability**: Other managers must be reachable

### Worker Nodes

Worker nodes execute tasks and run containers:

```bash
# Worker node responsibilities
- Execute assigned tasks
- Run containers
- Report status
- Provide resources
- Handle service requests
```

**Characteristics:**
- **Task Execution**: Run containers assigned by managers
- **Resource Provision**: Provide CPU, memory, and storage
- **Status Reporting**: Report health and status to managers
- **No Consensus**: Do not participate in cluster decisions
- **Scalability**: Can be easily added or removed

## Node Management in WhaleTUI

WhaleTUI provides comprehensive tools for managing Docker Swarm nodes through an intuitive terminal interface.

### Node List View

View all nodes in the cluster with key information:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WhaleTUI - Swarm Nodes                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ ID           HOSTNAME      STATUS    AVAILABILITY    MANAGER STATUS        │
│ abc123...    manager-1     Ready     Active          Leader                │
│ def456...    manager-2     Ready     Active          Reachable             │
│ ghi789...    worker-1      Ready     Active          -                     │
│ jkl012...    worker-2      Ready     Active          -                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Node Details

Inspect node configuration and status:

- **Node ID**: Unique identifier
- **Hostname**: Node's hostname
- **Status**: Current node status (Ready, Down, etc.)
- **Availability**: Node availability (Active, Pause, Drain)
- **Manager Status**: Role in manager consensus (Leader, Reachable, etc.)
- **Engine Version**: Docker Engine version
- **Resources**: CPU, memory, and storage capacity

## Common Node Operations

### Initializing Swarm

```bash
# Initialize swarm on first manager
docker swarm init --advertise-addr 192.168.1.100

# Join worker node to swarm
docker swarm join --token SWMTKN-1-... 192.168.1.100:2377

# Join manager node to swarm
docker swarm join --token SWMTKN-1-... 192.168.1.100:2377
```

### Managing Node Roles

```bash
# Promote worker to manager
docker node promote worker-1

# Demote manager to worker
docker node demote manager-2

# Check node role
docker node inspect node-name | grep -i role
```

### Updating Node Availability

```bash
# Set node to drain (stop accepting new tasks)
docker node update --availability drain worker-1

# Set node to pause (pause existing tasks)
docker node update --availability pause worker-1

# Set node to active (normal operation)
docker node update --availability active worker-1
```

### Removing Nodes

```bash
# Remove node from swarm
docker swarm leave --force

# Force remove node
docker node rm --force node-name
```

## Node Configuration

### Network Configuration

```bash
# Specify advertise address
docker swarm init --advertise-addr 192.168.1.100

# Use specific port
docker swarm init --advertise-addr 192.168.1.100:2377

# Listen on all interfaces
docker swarm init --advertise-addr 0.0.0.0
```

### Security Configuration

```bash
# Enable TLS encryption
docker swarm init \
  --advertise-addr 192.168.1.100 \
  --cert-expiry 2160h \
  --external-ca external-ca.pem

# Use custom CA
docker swarm init \
  --advertise-addr 192.168.1.100 \
  --ca-cert ca.pem \
  --ca-key ca-key.pem
```

### Resource Configuration

```bash
# Set resource limits
docker node update \
  --label-add memory=8g \
  --label-add cpu=4 \
  worker-1

# Add custom labels
docker node update \
  --label-add environment=production \
  --label-add team=backend \
  worker-1
```

## Best Practices

### Cluster Design

- **Manager Nodes**: Use odd number (3, 5, 7) for consensus
- **Worker Distribution**: Distribute workload evenly across workers
- **Resource Planning**: Plan for resource requirements and scaling
- **Network Design**: Ensure reliable network connectivity between nodes

### Node Management

- **Role Separation**: Keep manager and worker roles separate
- **Resource Monitoring**: Monitor node resource usage
- **Health Checks**: Implement node health monitoring
- **Backup Strategy**: Backup manager node configurations

### Security

- **TLS Encryption**: Enable TLS for node communication
- **Access Control**: Limit access to manager nodes
- **Regular Updates**: Keep Docker Engine updated
- **Network Security**: Implement network segmentation

### Monitoring

- **Resource Usage**: Track CPU, memory, and storage usage
- **Health Status**: Monitor node health and availability
- **Performance Metrics**: Track node performance
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

#### Node Availability Issues

```bash
# Check node availability
docker node inspect node-name | grep -i availability

# Update node availability
docker node update --availability active node-name

# Check for paused tasks
docker node ps node-name
```

#### Resource Problems

```bash
# Check node resources
docker node inspect node-name | grep -A 10 "Status"

# Check resource usage
docker system df

# Monitor node performance
docker stats --no-stream
```

### Node Diagnostics

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

### Node Labels and Constraints

```bash
# Add node labels
docker node update \
  --label-add environment=production \
  --label-add datacenter=us-east \
  worker-1

# Use labels in service constraints
docker service create \
  --name web \
  --constraint 'node.labels.environment==production' \
  nginx:latest
```

### Node Drain and Maintenance

```bash
# Drain node for maintenance
docker node update --availability drain worker-1

# Check drained tasks
docker node ps worker-1

# Re-enable node after maintenance
docker node update --availability active worker-1
```

### Multi-Architecture Support

```bash
# Check node architecture
docker node inspect node-name | grep -i architecture

# Deploy multi-arch services
docker service create \
  --name multiarch \
  --constraint 'node.arch==x86_64' \
  nginx:latest
```

## Integration with WhaleTUI

WhaleTUI seamlessly integrates node management with other Docker operations:

- **Service Management**: Deploy services with node constraints
- **Network Configuration**: Use overlay networks across nodes
- **Volume Management**: Configure volumes accessible to specific nodes
- **Image Management**: Ensure images are available on target nodes

### Keyboard Shortcuts

- `:nodes` - Switch to Nodes view
- `p` - Promote worker to manager (TBD - may not be implemented yet)
- `d` - Demote manager to worker (TBD - may not be implemented yet)
- `a` - Update node availability (TBD - may not be implemented yet)
- `Enter` - Inspect selected node

> **Note**: Some features marked as TBD (To Be Determined) may not be fully implemented in the current version of WhaleTUI. Check the latest release notes for current feature availability.

### Context Menus

Use keyboard shortcuts and navigation for quick access to:
- Node inspection
- Role management
- Availability updates
- Resource monitoring
- Configuration editing

## Next Steps

- [Swarm](./swarm.md) - Learn about Docker Swarm orchestration
- [Containers](./containers.md) - Deploy containers across nodes
- [Networks](./networks.md) - Configure overlay networks for nodes
- [Volumes](./volumes.md) - Manage persistent storage across nodes
- [Development Setup](../development/setup.md) - Set up your development environment
