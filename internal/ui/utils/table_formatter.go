package utils

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
)

// TableFormatter handles text truncation and alignment for table columns based on configuration
type TableFormatter struct {
	limits config.TableLimits
}

// NewTableFormatter creates a new table formatter with the given limits
func NewTableFormatter(limits config.TableLimits) *TableFormatter {
	return &TableFormatter{
		limits: limits,
	}
}

// NewTableFormatterFromTheme creates a new table formatter from theme manager
func NewTableFormatterFromTheme(themeManager *config.ThemeManager) *TableFormatter {
	return &TableFormatter{
		limits: themeManager.GetTableLimits(),
	}
}

// FormatCell formats a single cell based on the column type and limit (global configuration)
func (tf *TableFormatter) FormatCell(text string, columnType string) string {
	return tf.FormatCellForView(text, columnType, "")
}

// FormatCellForView formats a single cell with view-specific character limits applied
func (tf *TableFormatter) FormatCellForView(
	text string, columnType string, viewName string) string {
	config := tf.GetColumnConfigForView(columnType, viewName)
	return TruncateText(text, config.Limit)
}

// FormatCellSmart formats a single cell with smart word boundary truncation (global configuration)
func (tf *TableFormatter) FormatCellSmart(text string, columnType string) string {
	return tf.FormatCellSmartForView(text, columnType, "")
}

// FormatCellSmartForView formats a single cell with view-specific smart word boundary truncation
func (tf *TableFormatter) FormatCellSmartForView(text string, columnType string, viewName string) string {
	config := tf.GetColumnConfigForView(columnType, viewName)
	return TruncateTextSmart(text, config.Limit)
}

// UpdateLimits updates the character limits for the formatter
func (tf *TableFormatter) UpdateLimits(limits config.TableLimits) {
	tf.limits = limits
}

// GetLimits returns the current character limits
func (tf *TableFormatter) GetLimits() config.TableLimits {
	return tf.limits
}

// GetAlignmentForColumn returns the appropriate alignment for a specific column type (global configuration)
func (tf *TableFormatter) GetAlignmentForColumn(columnType string) int {
	return tf.GetAlignmentForColumnForView(columnType, "")
}

// GetAlignmentForColumnForView returns the appropriate alignment for a specific column type in a specific view
func (tf *TableFormatter) GetAlignmentForColumnForView(columnType string, viewName string) int {
	config := tf.GetColumnConfigForView(columnType, viewName)
	if config.Alignment != "" {
		return tf.parseAlignmentString(config.Alignment)
	}

	// Fallback to default alignment
	return tf.getDefaultAlignmentForColumn(columnType)
}

// GetColumnConfig returns the configuration for a specific column type (global configuration)
func (tf *TableFormatter) GetColumnConfig(columnType string) config.ColumnConfig {
	return tf.GetColumnConfigForView(columnType, "")
}

// GetColumnConfigForView returns the configuration for a specific column type in a specific view
func (tf *TableFormatter) GetColumnConfigForView(columnType, viewName string) config.ColumnConfig {
	columnType = strings.ToLower(columnType)
	// Keep viewName as-is to match the exact keys in configuration

	// Start with default configuration
	defaultConfig := config.ColumnConfig{
		Limit:     tf.getLimitForColumn(columnType),
		Width:     tf.getDefaultWidthForColumn(columnType),
		Visible:   true, // Visible by default
		Alignment: tf.getAlignmentString(tf.getDefaultAlignmentForColumn(columnType)),
	}

	// Check per-view configuration first (highest priority)
	if viewName != "" && tf.limits.Views != nil {
		if viewConfig, exists := tf.limits.Views[viewName]; exists {
			// Check in view-specific custom columns first
			if viewConfig.CustomColumns != nil {
				if customConfig, exists := viewConfig.CustomColumns[columnType]; exists {
					defaultConfig.MergeWith(customConfig)
					return defaultConfig
				}
			}

			// Check in view-specific columns
			if viewConfig.Columns != nil {
				if columnConfig, exists := viewConfig.Columns[columnType]; exists {
					defaultConfig.MergeWith(columnConfig)
					return defaultConfig
				}
			}
		}
	}

	// Check in global custom columns and merge
	if tf.limits.CustomColumns != nil {
		if customConfig, exists := tf.limits.CustomColumns[columnType]; exists {
			defaultConfig.MergeWith(customConfig)
			return defaultConfig
		}
	}

	// Check in global regular columns and merge
	if tf.limits.Columns != nil {
		if columnConfig, exists := tf.limits.Columns[columnType]; exists {
			defaultConfig.MergeWith(columnConfig)
			return defaultConfig
		}
	}

	// Return default configuration (no overrides found)
	return defaultConfig
}

// GetColumnWidth returns the width for a specific column type
func (tf *TableFormatter) GetColumnWidth(columnType string) int {
	return tf.GetColumnWidthForView(columnType, "")
}

// GetColumnWidthForView returns the width for a specific column type in a specific view
func (tf *TableFormatter) GetColumnWidthForView(columnType, viewName string) int {
	return tf.GetColumnWidthForViewWithTerminalSize(columnType, viewName, 0)
}

// GetColumnWidthForViewWithTerminalSize returns the width for a specific column type in a specific view
// with a given terminal width (0 means use default behavior)
func (tf *TableFormatter) GetColumnWidthForViewWithTerminalSize(columnType, viewName string, terminalWidth int) int {
	config := tf.GetColumnConfigForView(columnType, viewName)

	// If percentage width is specified, calculate from terminal width
	if config.WidthPercent > 0 && terminalWidth > 0 {
		calculatedWidth := (terminalWidth * config.WidthPercent) / 100

		// Apply min/max constraints
		if config.MinWidth > 0 && calculatedWidth < config.MinWidth {
			calculatedWidth = config.MinWidth
		}
		if config.MaxWidth > 0 && calculatedWidth > config.MaxWidth {
			calculatedWidth = config.MaxWidth
		}

		return calculatedWidth
	}

	// Fallback to fixed width
	if config.Width > 0 {
		return config.Width
	}

	// Default to character limit if no width specified
	return config.Limit
}

// IsColumnVisible returns whether a column should be visible
func (tf *TableFormatter) IsColumnVisible(columnType string) bool {
	return tf.IsColumnVisibleForView(columnType, "")
}

// IsColumnVisibleForView returns whether a column should be visible in a specific view
func (tf *TableFormatter) IsColumnVisibleForView(columnType, viewName string) bool {
	config := tf.GetColumnConfigForView(columnType, viewName)
	return config.Visible
}

// GetColumnDisplayName returns the display name for a column
func (tf *TableFormatter) GetColumnDisplayName(columnType string) string {
	return tf.GetColumnDisplayNameForView(columnType, "")
}

// GetColumnDisplayNameForView returns the display name for a column in a specific view
func (tf *TableFormatter) GetColumnDisplayNameForView(columnType, viewName string) string {
	config := tf.GetColumnConfigForView(columnType, viewName)
	if config.DisplayName != "" {
		return config.DisplayName
	}
	// Return capitalized version of the column type
	// Simple capitalization: first letter uppercase, rest lowercase
	if len(columnType) == 0 {
		return columnType
	}
	return strings.ToUpper(columnType[:1]) + strings.ToLower(columnType[1:])
}

// getAlignmentString converts tview alignment constant to string
func (tf *TableFormatter) getAlignmentString(alignment int) string {
	switch alignment {
	case tview.AlignLeft:
		return "left"
	case tview.AlignRight:
		return "right"
	case tview.AlignCenter:
		return "center"
	default:
		return "left"
	}
}

// parseAlignmentString converts string alignment to tview alignment constant
func (tf *TableFormatter) parseAlignmentString(alignment string) int {
	switch strings.ToLower(alignment) {
	case "left":
		return tview.AlignLeft
	case "right":
		return tview.AlignRight
	case "center":
		return tview.AlignCenter
	default:
		// Default to left if invalid alignment string
		return tview.AlignLeft
	}
}

// getDefaultAlignmentForColumn returns the default alignment for a column type
func (tf *TableFormatter) getDefaultAlignmentForColumn(columnType string) int {
	switch strings.ToLower(columnType) {
	// Right-aligned columns for numerical data and IDs
	case "id":
		return tview.AlignRight
	case "size":
		return tview.AlignRight
	case "containers":
		return tview.AlignRight
	case "created":
		return tview.AlignRight
	case "replicas":
		return tview.AlignRight
	case "engine_version":
		return tview.AlignRight
	// Left-aligned columns for text-based content
	case "name":
		return tview.AlignLeft
	case "image":
		return tview.AlignLeft
	case "repository":
		return tview.AlignLeft
	case "tag":
		return tview.AlignLeft
	case "driver":
		return tview.AlignLeft
	case "mountpoint":
		return tview.AlignLeft
	case "network":
		return tview.AlignLeft
	case "scope":
		return tview.AlignLeft
	case "description":
		return tview.AlignLeft
	case "hostname":
		return tview.AlignLeft
	case "address":
		return tview.AlignLeft
	case "ports":
		return tview.AlignLeft
	// Status and state can be either - default to left for consistency
	case "status":
		return tview.AlignRight
	case "state":
		return tview.AlignRight
	case "mode":
		return tview.AlignLeft
	case "role":
		return tview.AlignLeft
	case "availability":
		return tview.AlignLeft
	case "manager_status":
		return tview.AlignLeft
	default:
		// Default to left alignment for unknown column types
		return tview.AlignLeft
	}
}

// getLimitForColumn returns the default character limit for a specific column type
func (tf *TableFormatter) getLimitForColumn(columnType string) int {
	switch strings.ToLower(columnType) {
	case "id":
		return 12 // Container/Image IDs
	case "name":
		return 30 // Container/Volume names
	case "image":
		return 40 // Image names
	case "status":
		return 20 // Status messages
	case "state":
		return 15 // State values
	case "ports":
		return 25 // Port mappings
	case "created":
		return 20 // Created timestamps
	case "size":
		return 15 // Size values
	case "driver":
		return 15 // Volume drivers
	case "mountpoint":
		return 50 // Mount points
	case "repository":
		return 35 // Image repositories
	case "tag":
		return 20 // Image tags
	case "network":
		return 25 // Network names
	case "scope":
		return 15 // Network scope
	case "description":
		return 50 // General descriptions
	case "containers":
		return 10 // Container count
	case "replicas":
		return 10 // Replica count
	case "engine_version":
		return 15 // Engine version
	case "hostname":
		return 20 // Hostname
	case "address":
		return 20 // IP addresses
	case "mode":
		return 15 // Swarm mode
	case "role":
		return 12 // Swarm role
	case "availability":
		return 15 // Swarm availability
	case "manager_status":
		return 20 // Manager status
	default:
		// Default limit for unknown column types
		return 30
	}
}

// getDefaultWidthForColumn returns sensible default widths for each column type
func (tf *TableFormatter) getDefaultWidthForColumn(columnType string) int {
	switch strings.ToLower(columnType) {
	case "id":
		return 15 // Container/Image IDs - slightly wider than limit for readability
	case "name":
		return 35 // Container/Volume names - wider for better readability
	case "image":
		return 80 // Image names - often long repository paths
	case "status":
		return 25 // Status messages - need space for status text
	case "state":
		return 18 // State values - running, stopped, etc.
	case "ports":
		return 30 // Port mappings - 0.0.0.0:8080->80/tcp format
	case "created":
		return 22 // Created timestamps - ISO format dates
	case "size":
		return 18 // Size values - 1.2GB, 500MB format
	case "driver":
		return 18 // Volume drivers - local, nfs, etc.
	case "mountpoint":
		return 55 // Mount points - often long paths
	case "repository":
		return 40 // Image repositories - docker.io/library/nginx
	case "tag":
		return 25 // Image tags - latest, v1.0.0, etc.
	case "network":
		return 30 // Network names - bridge, host, custom networks
	case "scope":
		return 18 // Network scope - local, global
	case "description":
		return 55 // General descriptions - can be long
	case "containers":
		return 12 // Container count - just numbers
	case "replicas":
		return 10 // Replica count - just numbers
	case "engine_version":
		return 20 // Engine version - 20.10.0 format
	case "hostname":
		return 25 // Hostname - server names
	case "address":
		return 20 // IP addresses - 192.168.1.100 format
	case "mode":
		return 15 // Swarm mode - replicated, global
	case "role":
		return 12 // Swarm role - manager, worker
	case "availability":
		return 15 // Swarm availability - active, pause, drain
	case "manager_status":
		return 25 // Manager status - leader, reachable, etc.
	default:
		// Default width for unknown column types
		return 25
	}
}
