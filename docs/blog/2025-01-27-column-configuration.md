---
slug: column-configuration-system
title: Introducing the Column Configuration System
authors: [whaletui-team]
tags: [features, configuration, ui, customization]
---

# Introducing the Column Configuration System

We're excited to announce the release of WhaleTUI's comprehensive column configuration system! This powerful new feature gives you complete control over how data is displayed in all table views, making your Docker management experience more efficient and personalized.

<!--truncate-->

## What's New

The column configuration system introduces several key features that transform how you interact with WhaleTUI's table views:

### 🎯 **Percentage-Based Column Widths**
Set column widths as percentages of your terminal width for responsive layouts that adapt to any screen size.

```yaml
tableLimits:
  views:
    containers:
      columns:
        name:
          width_percent: 40
          min_width: 20
          max_width: 60
```

### 📊 **Intelligent Alignment**
Right-align numerical data for better readability and comparison:

```yaml
columns:
  size:
    alignment: "right"  # Perfect for comparing image sizes
  created:
    alignment: "right"  # Great for timestamp scanning
```

### 👁️ **Column Visibility Control**
Show or hide columns based on your workflow needs:

```yaml
views:
  containers:
    columns:
      id:
        visible: false  # Hide ID column for cleaner view
```

### 🏷️ **Custom Display Names**
Set meaningful column headers:

```yaml
columns:
  id:
    display_name: "Container ID"
  name:
    display_name: "Container Name"
```

## Per-View Configurations

Each view type now supports independent column configurations, allowing you to optimize the display for different data types:

- **Containers**: 7 configurable columns (ID, Name, Image, Status, State, Ports, Created)
- **Images**: 5 configurable columns (ID, Repository, Tag, Size, Created)
- **Volumes**: 5 configurable columns (Name, Driver, Mount Point, Created, Size)
- **Networks**: 7 configurable columns (ID, Name, Driver, Scope, Created, Subnet, Gateway)
- **Swarm Nodes**: 6 configurable columns (ID, Name, Status, Availability, Role, Engine)
- **Swarm Services**: 5 configurable columns (ID, Name, Mode, Replicas, Image)

## Real-World Benefits

### Better Readability
Right-aligned numerical data makes it easier to:
- Compare container sizes across multiple images
- Scan through port numbers consistently
- Analyze timestamp patterns in logs
- Review resource usage statistics

### Responsive Design
Percentage-based widths ensure your interface looks great on any terminal size:
- Small terminals: Columns remain readable with minimum width constraints
- Large terminals: Columns don't become excessively wide with maximum width limits
- Dynamic sizing: Layouts adapt automatically to terminal resizing

### Workflow Optimization
Customize views for your specific use cases:
- **Development**: Focus on container names and status
- **Production**: Emphasize resource usage and health status
- **Network Admin**: Highlight network details and configurations

## Configuration Examples

### Development Workflow
```yaml
tableLimits:
  views:
    containers:
      columns:
        id:
          visible: false  # Hide IDs for cleaner view
        name:
          width_percent: 50
          display_name: "Container"
        status:
          width_percent: 25
          alignment: "right"
        ports:
          width_percent: 25
```

### Production Monitoring
```yaml
tableLimits:
  views:
    containers:
      columns:
        name:
          width_percent: 35
          display_name: "Container"
        status:
          width_percent: 20
          alignment: "right"
        created:
          width_percent: 20
          alignment: "right"
        ports:
          width_percent: 25
```

## Getting Started

1. **Update to the latest version** of WhaleTUI
2. **Edit your theme file** (`dark-theme.yaml` or `dark-theme.json`)
3. **Add column configurations** under the `tableLimits` section
4. **Restart WhaleTUI** to apply changes

For detailed configuration options and examples, check out our [Column Configuration Guide](/docs/concepts/column-configuration) and [Configuration Examples](/docs/concepts/configuration-examples).

## Backward Compatibility

Existing configurations continue to work without any changes. The new system is designed to enhance your current setup while maintaining full backward compatibility.

## What's Next

This column configuration system is just the beginning. We're planning additional customization features including:

- **Custom column data sources** for additional information display
- **Column sorting and filtering** options
- **Saved configuration presets** for different workflows
- **Theme-specific column configurations**

## Feedback and Contributions

We'd love to hear your feedback on this new feature! Try it out and let us know:

- How you're using the column configurations
- What additional features would be helpful
- Any issues or suggestions for improvement

Visit our [GitHub repository](https://github.com/wikczerski/whaletui) to report issues, request features, or contribute to the project.

---

The column configuration system represents a significant step forward in making WhaleTUI more customizable and user-friendly. We're excited to see how you'll use these new features to optimize your Docker management workflow!

Happy container management! 🐳
