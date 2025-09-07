# Configuration Examples

This guide provides practical examples of how to configure WhaleTUI's column system for different use cases.

## Basic Examples

### Hide ID Columns

Hide ID columns across all views for a cleaner interface:

```yaml
tableLimits:
  columns:
    id:
      visible: false
```

### Right-Align Numerical Data

Make numerical data easier to compare by right-aligning it:

```yaml
tableLimits:
  columns:
    size:
      alignment: "right"
    created:
      alignment: "right"
    containers:
      alignment: "right"
    replicas:
      alignment: "right"
```

### Customize Container View

Optimize the container view for your workflow:

```yaml
tableLimits:
  views:
    containers:
      columns:
        id:
          width_percent: 15
          alignment: "right"
          visible: true
        name:
          width_percent: 40
          min_width: 25
          display_name: "Container Name"
        status:
          width_percent: 20
          alignment: "right"
        ports:
          width_percent: 25
          limit: 30
```

## Advanced Examples

### Responsive Design

Create a responsive layout that works well on different terminal sizes:

```yaml
tableLimits:
  views:
    containers:
      columns:
        name:
          width_percent: 45
          min_width: 20
          max_width: 60
        status:
          width_percent: 20
          min_width: 15
          max_width: 25
        ports:
          width_percent: 35
          min_width: 20
          max_width: 50
```

### Image Repository Focus

Optimize the images view for repository management:

```yaml
tableLimits:
  views:
    images:
      columns:
        repository:
          width_percent: 60
          min_width: 30
          display_name: "Repository"
        tag:
          width_percent: 20
          alignment: "center"
          display_name: "Tag"
        size:
          width_percent: 15
          alignment: "right"
          display_name: "Size"
        created:
          width_percent: 15
          alignment: "right"
          display_name: "Created"
```

### Network Management

Configure the networks view for network administration:

```yaml
tableLimits:
  views:
    networks:
      columns:
        name:
          width_percent: 30
          min_width: 20
          display_name: "Network Name"
        driver:
          width_percent: 15
          alignment: "center"
          display_name: "Driver"
        scope:
          width_percent: 15
          alignment: "center"
          display_name: "Scope"
        subnet:
          width_percent: 20
          alignment: "right"
          display_name: "Subnet"
        gateway:
          width_percent: 20
          alignment: "right"
          display_name: "Gateway"
```

### Swarm Management

Optimize swarm views for cluster management:

```yaml
tableLimits:
  views:
    "Swarm Nodes":
      columns:
        name:
          width_percent: 35
          min_width: 20
          display_name: "Node Name"
        status:
          width_percent: 20
          alignment: "right"
          display_name: "Status"
        availability:
          width_percent: 20
          alignment: "center"
          display_name: "Availability"
        role:
          width_percent: 15
          alignment: "center"
          display_name: "Role"
        engine_version:
          width_percent: 20
          alignment: "right"
          display_name: "Engine Version"

    "Swarm Services":
      columns:
        name:
          width_percent: 40
          min_width: 25
          display_name: "Service Name"
        mode:
          width_percent: 20
          alignment: "center"
          display_name: "Mode"
        replicas:
          width_percent: 15
          alignment: "right"
          display_name: "Replicas"
        image:
          width_percent: 25
          min_width: 20
          display_name: "Image"
```

## Use Case Examples

### Development Workflow

Optimize for development with focus on container names and status:

```yaml
tableLimits:
  views:
    containers:
      columns:
        id:
          visible: false  # Hide IDs for cleaner view
        name:
          width_percent: 50
          min_width: 30
          display_name: "Container"
        status:
          width_percent: 25
          alignment: "right"
          display_name: "Status"
        ports:
          width_percent: 25
          limit: 35
          display_name: "Ports"
```

### Production Monitoring

Focus on status and resource usage for production monitoring:

```yaml
tableLimits:
  views:
    containers:
      columns:
        name:
          width_percent: 35
          min_width: 20
          display_name: "Container"
        status:
          width_percent: 20
          alignment: "right"
          display_name: "Status"
        created:
          width_percent: 20
          alignment: "right"
          display_name: "Created"
        ports:
          width_percent: 25
          limit: 30
          display_name: "Ports"

    images:
      columns:
        repository:
          width_percent: 50
          min_width: 30
          display_name: "Repository"
        size:
          width_percent: 25
          alignment: "right"
          display_name: "Size"
        created:
          width_percent: 25
          alignment: "right"
          display_name: "Created"
```

### Network Administration

Optimize for network management tasks:

```yaml
tableLimits:
  views:
    networks:
      columns:
        name:
          width_percent: 25
          min_width: 15
          display_name: "Network"
        driver:
          width_percent: 15
          alignment: "center"
          display_name: "Driver"
        scope:
          width_percent: 15
          alignment: "center"
          display_name: "Scope"
        subnet:
          width_percent: 25
          alignment: "right"
          display_name: "Subnet"
        gateway:
          width_percent: 20
          alignment: "right"
          display_name: "Gateway"
```

## Custom Column Examples

### Add Uptime Column to Containers

```yaml
tableLimits:
  views:
    containers:
      custom_columns:
        uptime:
          width_percent: 15
          alignment: "right"
          display_name: "Uptime"
          visible: true
          limit: 20
```

### Add Health Status Column

```yaml
tableLimits:
  views:
    containers:
      custom_columns:
        health:
          width_percent: 15
          alignment: "center"
          display_name: "Health"
          visible: true
          limit: 15
```

## Complete Configuration Example

Here's a complete configuration that demonstrates all features:

```yaml
tableLimits:
  # Global settings applied to all views
  columns:
    id:
      limit: 12
      width_percent: 15
      min_width: 10
      max_width: 20
      visible: true
      alignment: "right"
      display_name: "ID"

    name:
      limit: 30
      width_percent: 35
      min_width: 20
      max_width: 50
      visible: true
      alignment: "left"
      display_name: "Name"

    status:
      limit: 20
      width_percent: 20
      min_width: 15
      max_width: 25
      visible: true
      alignment: "right"
      display_name: "Status"

    created:
      limit: 20
      width_percent: 15
      min_width: 12
      max_width: 20
      visible: true
      alignment: "right"
      display_name: "Created"

  # View-specific configurations
  views:
    containers:
      columns:
        id:
          width_percent: 20
          alignment: "right"
        name:
          width_percent: 40
          min_width: 25
          display_name: "Container Name"
        status:
          width_percent: 20
          alignment: "right"
        ports:
          width_percent: 20
          limit: 25
          display_name: "Ports"

      custom_columns:
        uptime:
          width_percent: 15
          alignment: "right"
          display_name: "Uptime"
          visible: true
          limit: 20

    images:
      columns:
        repository:
          width_percent: 50
          min_width: 30
          display_name: "Repository"
        tag:
          width_percent: 20
          alignment: "center"
          display_name: "Tag"
        size:
          width_percent: 20
          alignment: "right"
          display_name: "Size"
        created:
          width_percent: 15
          alignment: "right"
          display_name: "Created"

    volumes:
      columns:
        name:
          width_percent: 40
          min_width: 25
          display_name: "Volume Name"
        driver:
          width_percent: 20
          alignment: "center"
          display_name: "Driver"
        mountpoint:
          width_percent: 40
          min_width: 30
          display_name: "Mount Point"

    networks:
      columns:
        name:
          width_percent: 30
          min_width: 20
          display_name: "Network Name"
        driver:
          width_percent: 15
          alignment: "center"
          display_name: "Driver"
        scope:
          width_percent: 15
          alignment: "center"
          display_name: "Scope"
        subnet:
          width_percent: 20
          alignment: "right"
          display_name: "Subnet"
        gateway:
          width_percent: 20
          alignment: "right"
          display_name: "Gateway"

    "Swarm Nodes":
      columns:
        name:
          width_percent: 35
          min_width: 20
          display_name: "Node Name"
        status:
          width_percent: 20
          alignment: "right"
          display_name: "Status"
        availability:
          width_percent: 20
          alignment: "center"
          display_name: "Availability"
        role:
          width_percent: 15
          alignment: "center"
          display_name: "Role"
        engine_version:
          width_percent: 20
          alignment: "right"
          display_name: "Engine Version"

    "Swarm Services":
      columns:
        name:
          width_percent: 40
          min_width: 25
          display_name: "Service Name"
        mode:
          width_percent: 20
          alignment: "center"
          display_name: "Mode"
        replicas:
          width_percent: 15
          alignment: "right"
          display_name: "Replicas"
        image:
          width_percent: 25
          min_width: 20
          display_name: "Image"
```

## Tips and Tricks

### Quick Width Adjustments

Use these percentage guidelines for quick setup:

- **Small columns** (IDs, counts): 10-15%
- **Medium columns** (names, status): 20-30%
- **Large columns** (descriptions, paths): 40-50%

### Alignment Patterns

- **Left-align**: Text content, names, descriptions
- **Right-align**: Numbers, IDs, timestamps, sizes
- **Center-align**: Short values, states, modes

### Responsive Design

Always use `min_width` and `max_width` for responsive layouts:

```yaml
name:
  width_percent: 40
  min_width: 20  # Readable on small terminals
  max_width: 60  # Not too wide on large terminals
```

### Configuration Testing

Test your configuration by:

1. Starting with minimal changes
2. Testing on different terminal sizes
3. Verifying all views work correctly
4. Checking alignment and width behavior

This approach ensures your configuration works well across different environments and use cases.
