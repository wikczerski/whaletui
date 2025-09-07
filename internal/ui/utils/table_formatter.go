package utils

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/ui/constants"
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
	text string, columnType string, viewName string,
) string {
	config := tf.GetColumnConfigForView(columnType, viewName)
	return TruncateText(text, config.Limit)
}

// FormatCellSmart formats a single cell with smart word boundary truncation (global configuration)
func (tf *TableFormatter) FormatCellSmart(text string, columnType string) string {
	return tf.FormatCellSmartForView(text, columnType, "")
}

// FormatCellSmartForView formats a single cell with view-specific smart word boundary truncation
func (tf *TableFormatter) FormatCellSmartForView(
	text string, columnType string, viewName string,
) string {
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
	if config := tf.getViewSpecificConfig(columnType, viewName); config != nil {
		defaultConfig.MergeWith(*config)
		return defaultConfig
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
func (tf *TableFormatter) GetColumnWidthForViewWithTerminalSize(
	columnType, viewName string, terminalWidth int,
) int {
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

// getViewSpecificConfig returns view-specific configuration if found
func (tf *TableFormatter) getViewSpecificConfig(columnType, viewName string) *config.ColumnConfig {
	if viewName == "" || tf.limits.Views == nil {
		return nil
	}
	viewConfig, exists := tf.limits.Views[viewName]
	if !exists {
		return nil
	}
	// Check in view-specific custom columns first
	if viewConfig.CustomColumns != nil {
		if customConfig, exists := viewConfig.CustomColumns[columnType]; exists {
			return &customConfig
		}
	}
	// Check in view-specific columns
	if viewConfig.Columns != nil {
		if columnConfig, exists := viewConfig.Columns[columnType]; exists {
			return &columnConfig
		}
	}
	return nil
}

// getDefaultAlignmentForColumn returns the default alignment for a column type
func (tf *TableFormatter) getDefaultAlignmentForColumn(columnType string) int {
	if alignment, exists := constants.DefaultColumnAlignments[strings.ToLower(columnType)]; exists {
		return alignment
	}
	// Default to left alignment for unknown column types
	return tview.AlignLeft
}

// getLimitForColumn returns the default character limit for a specific column type
func (tf *TableFormatter) getLimitForColumn(columnType string) int {
	if limit, exists := constants.DefaultColumnLimits[strings.ToLower(columnType)]; exists {
		return limit
	}
	// Default limit for unknown column types
	return 30
}

// getDefaultWidthForColumn returns sensible default widths for each column type
func (tf *TableFormatter) getDefaultWidthForColumn(columnType string) int {
	if width, exists := constants.DefaultColumnWidths[strings.ToLower(columnType)]; exists {
		return width
	}
	// Default width for unknown column types
	return 25
}
