package utils

import (
	"strings"

	"github.com/wikczerski/whaletui/internal/config"
)

// TableFormatter handles text truncation for table columns based on configuration
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

// FormatCell formats a single cell based on the column type and limit
func (tf *TableFormatter) FormatCell(text string, columnType string) string {
	limit := tf.getLimitForColumn(columnType)
	return TruncateText(text, limit)
}

// FormatCellSmart formats a single cell with smart word boundary truncation
func (tf *TableFormatter) FormatCellSmart(text string, columnType string) string {
	limit := tf.getLimitForColumn(columnType)
	return TruncateTextSmart(text, limit)
}

// UpdateLimits updates the character limits for the formatter
func (tf *TableFormatter) UpdateLimits(limits config.TableLimits) {
	tf.limits = limits
}

// GetLimits returns the current character limits
func (tf *TableFormatter) GetLimits() config.TableLimits {
	return tf.limits
}

// getLimitForColumn returns the character limit for a specific column type
func (tf *TableFormatter) getLimitForColumn(columnType string) int {
	switch strings.ToLower(columnType) {
	case "id":
		return tf.limits.ID
	case "name":
		return tf.limits.Name
	case "image":
		return tf.limits.Image
	case "status":
		return tf.limits.Status
	case "state":
		return tf.limits.State
	case "ports":
		return tf.limits.Ports
	case "created":
		return tf.limits.Created
	case "size":
		return tf.limits.Size
	case "driver":
		return tf.limits.Driver
	case "mountpoint":
		return tf.limits.Mountpoint
	case "repository":
		return tf.limits.Repository
	case "tag":
		return tf.limits.Tag
	case "network":
		return tf.limits.Network
	case "scope":
		return tf.limits.Scope
	case "description":
		return tf.limits.Description
	default:
		// Default limit for unknown column types
		return 30
	}
}
