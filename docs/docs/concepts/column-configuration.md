# Column Configuration

WhaleTUI provides a comprehensive column configuration system that allows you to customize how data is displayed in all table views. This system gives you complete control over column widths, alignment, visibility, and display names.

## Overview

The column configuration system supports:

- **Percentage-based column widths** with responsive sizing
- **Per-view configurations** for different data types
- **Column visibility control** to show/hide specific columns
- **Custom alignment** (left, right, center) for better readability
- **Display name customization** for column headers
- **Character limits** for content truncation
- **Custom columns** for additional data display

## Configuration Structure

Column configurations are defined in your theme files (`dark-theme.yaml` or `dark-theme.json`) under the `tableLimits` section:

```yaml
tableLimits:
  # Global column configurations (applied to all views)
  columns:
    id:
      limit: 12
      width_percent: 15
      min_width: 10
      max_width: 20
      visible: true
      alignment: "left"
      display_name: "ID"

  # Per-view configurations (override global settings)
  views:
    containers:
      columns:
        id:
          width_percent: 20
          alignment: "right"
        name:
          width_percent: 40
          min_width: 20
          max_width: 60
          display_name: "Container Name"
```

## Column Properties

### Width Configuration

#### `width_percent`
Sets the column width as a percentage of the available terminal width (0-100).

```yaml
name:
  width_percent: 40  # Takes 40% of terminal width
```

#### `min_width`
Minimum width in characters. Prevents columns from becoming too narrow.

```yaml
name:
  width_percent: 30
  min_width: 20  # Never smaller than 20 characters
```

#### `max_width`
Maximum width in characters. Prevents columns from becoming too wide.

```yaml
name:
  width_percent: 50
  max_width: 60  # Never larger than 60 characters
```

### Visibility Control

#### `visible`
Controls whether a column is displayed (true/false).

```yaml
id:
  visible: false  # Hide the ID column
```

### Alignment

#### `alignment`
Sets the text alignment within the column:
- `"left"` - Left-aligned (default for text)
- `"right"` - Right-aligned (recommended for numbers, IDs, timestamps)
- `"center"` - Center-aligned

```yaml
size:
  alignment: "right"  # Right-align numerical data
created:
  alignment: "right"  # Right-align timestamps
```

### Display Customization

#### `display_name`
Custom header text for the column.

```yaml
id:
  display_name: "Container ID"  # Instead of just "ID"
```

#### `limit`
Maximum number of characters to display before truncation.

```yaml
name:
  limit: 30  # Truncate names longer than 30 characters
```

## Per-View Configurations

Each view type has its own configuration section, allowing you to customize columns differently for different data types.

### Available Views

- **`containers`** - Container list view
- **`images`** - Image list view
- **`volumes`** - Volume list view
- **`networks`** - Network list view
- **`Swarm Nodes`** - Swarm node list view
- **`Swarm Services`** - Swarm service list view

### View-Specific Example

```yaml
views:
  containers:
    columns:
      id:
        width_percent: 15
        alignment: "right"
        visible: true
      name:
        width_percent: 35
        min_width: 20
        display_name: "Container Name"
      status:
        width_percent: 20
        alignment: "right"
      ports:
        width_percent: 30
        limit: 25

  images:
    columns:
      repository:
        width_percent: 50
        min_width: 30
      tag:
        width_percent: 20
        alignment: "center"
      size:
        width_percent: 15
        alignment: "right"
      created:
        width_percent: 15
        alignment: "right"
```

## Custom Columns

You can add custom columns to specific views using the `custom_columns` section:

```yaml
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

## Configuration Priority

The system applies configurations in the following priority order (highest to lowest):

1. **View-specific custom columns** (`views.{view}.custom_columns`)
2. **View-specific regular columns** (`views.{view}.columns`)
3. **Global custom columns** (`custom_columns`)
4. **Global regular columns** (`columns`)
5. **Default values** (built-in defaults)

## Best Practices

### Alignment Guidelines

- **Right-align numerical data**: IDs, sizes, ports, timestamps, counts
- **Left-align text data**: Names, descriptions, status messages
- **Center-align short values**: Tags, states, modes

### Width Guidelines

- **ID columns**: 15-20% width
- **Name columns**: 30-40% width
- **Status/State columns**: 15-25% width
- **Timestamp columns**: 15-20% width
- **Size columns**: 10-15% width

### Responsive Design

Use `min_width` and `max_width` to ensure columns remain readable across different terminal sizes:

```yaml
name:
  width_percent: 40
  min_width: 20  # Readable on small terminals
  max_width: 60  # Not too wide on large terminals
```

## Example Configurations

### Minimal Configuration

```yaml
tableLimits:
  views:
    containers:
      columns:
        id:
          visible: false  # Hide ID column
        name:
          width_percent: 50  # Make name column wider
```

### Comprehensive Configuration

```yaml
tableLimits:
  columns:
    id:
      limit: 12
      width_percent: 15
      min_width: 10
      max_width: 20
      visible: true
      alignment: "right"
      display_name: "ID"

  views:
    containers:
      columns:
        id:
          width_percent: 20
          alignment: "right"
        name:
          width_percent: 35
          min_width: 20
          max_width: 50
          display_name: "Container Name"
        status:
          width_percent: 20
          alignment: "right"
        ports:
          width_percent: 25
          limit: 30

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
          width_percent: 45
          min_width: 25
        tag:
          width_percent: 20
          alignment: "center"
        size:
          width_percent: 20
          alignment: "right"
        created:
          width_percent: 15
          alignment: "right"
```

## Migration from Legacy Configuration

If you have existing configurations using the old `HeaderLayout` system, they will continue to work. However, we recommend migrating to the new system for better control and features.

### Legacy vs New Configuration

**Old (deprecated):**
```yaml
headerLayout:
  id: 15
  name: 30
```

**New (recommended):**
```yaml
tableLimits:
  views:
    containers:
      columns:
        id:
          width_percent: 15
        name:
          width_percent: 30
```

## Troubleshooting

### Configuration Not Applied

1. **Check view name case**: Use exact case (`"Swarm Nodes"`, not `"swarm nodes"`)
2. **Verify YAML syntax**: Ensure proper indentation and structure
3. **Check file location**: Configuration should be in your active theme file
4. **Restart application**: Changes require application restart

### Column Width Issues

1. **Total width exceeds 100%**: Ensure all `width_percent` values don't exceed 100%
2. **Columns too narrow**: Increase `min_width` values
3. **Columns too wide**: Decrease `max_width` values or `width_percent`

### Alignment Not Working

1. **Check alignment values**: Use `"left"`, `"right"`, or `"center"`
2. **Verify column type**: Ensure the column exists in the view
3. **Check configuration priority**: View-specific settings override global settings

## Advanced Features

### Dynamic Width Calculation

The system automatically calculates column widths based on:
- Terminal width
- Percentage settings
- Min/max constraints
- Available space

### Content Padding

For fixed-width columns, content is automatically padded to maintain alignment:
- **Right-aligned**: Left-padded with spaces
- **Left-aligned**: Right-padded with spaces
- **Center-aligned**: Center-padded with spaces

### Configuration Merging

User configurations are intelligently merged with defaults:
- Non-zero values override defaults
- Empty strings are ignored
- Boolean values always override defaults
- Nested configurations are merged recursively

This ensures that partial configurations work correctly while maintaining sensible defaults for unspecified properties.
