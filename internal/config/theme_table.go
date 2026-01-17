package config

// ColumnConfig defines configuration for a single table column
type ColumnConfig struct {
	// Character limit for the column (0 = no limit)
	Limit int `json:"limit,omitempty" yaml:"limit,omitempty"`

	// Width of the column (0 = auto, >0 = fixed width in characters)
	Width int `json:"width,omitempty" yaml:"width,omitempty"`

	// Width as percentage of available terminal width (0-100, takes precedence over Width)
	WidthPercent int `json:"width_percent,omitempty" yaml:"width_percent,omitempty"`

	// Minimum width in characters (applied when using WidthPercent)
	MinWidth int `json:"min_width,omitempty" yaml:"min_width,omitempty"`

	// Maximum width in characters (applied when using WidthPercent)
	MaxWidth int `json:"max_width,omitempty" yaml:"max_width,omitempty"`

	// Whether the column is visible (true = visible, false = hidden)
	Visible bool `json:"visible,omitempty" yaml:"visible,omitempty"`

	// Alignment override for this column ("left", "right", "center")
	Alignment string `json:"alignment,omitempty" yaml:"alignment,omitempty"`

	// Display name for the column (if different from the key)
	DisplayName string `json:"display_name,omitempty" yaml:"display_name,omitempty"`
}

// MergeWith merges this ColumnConfig with another, copying non-zero/empty values
func (cc *ColumnConfig) MergeWith(other any) {
	otherConfig := extractColumnConfig(other)
	if otherConfig == nil {
		return
	}

	// Merge non-zero integer fields
	if otherConfig.Limit > 0 {
		cc.Limit = otherConfig.Limit
	}
	if otherConfig.Width > 0 {
		cc.Width = otherConfig.Width
	}
	if otherConfig.WidthPercent > 0 {
		cc.WidthPercent = otherConfig.WidthPercent
	}
	if otherConfig.MinWidth > 0 {
		cc.MinWidth = otherConfig.MinWidth
	}
	if otherConfig.MaxWidth > 0 {
		cc.MaxWidth = otherConfig.MaxWidth
	}

	// Merge non-empty string fields
	if otherConfig.Alignment != "" {
		cc.Alignment = otherConfig.Alignment
	}
	if otherConfig.DisplayName != "" {
		cc.DisplayName = otherConfig.DisplayName
	}

	// Always merge the Visible field (it's a boolean, so we need to handle it explicitly)
	// Only override if it's explicitly set to false in the config
	cc.Visible = otherConfig.Visible
}

// extractColumnConfig extracts ColumnConfig from various types
func extractColumnConfig(other any) *ColumnConfig {
	switch v := other.(type) {
	case ColumnConfig:
		return &v
	case *ColumnConfig:
		return v
	default:
		return nil
	}
}

// TableLimits defines column configuration for table columns
type TableLimits struct {
	// Global column configurations - allows fine-grained control over each column
	// Key is the column type (e.g., "id", "name", "size"), value is the configuration
	Columns map[string]ColumnConfig `json:"columns,omitempty" yaml:"columns,omitempty"`

	// Global custom columns - allows adding columns that are not part of the default set
	// Key is the column type, value is the configuration
	CustomColumns map[string]ColumnConfig `json:"custom_columns,omitempty" yaml:"custom_columns,omitempty"`

	// Per-view configurations - allows different column settings for each view
	// Key is the view name (e.g., "containers", "images", "volumes"), value is view-specific settings
	Views map[string]ViewConfig `json:"views,omitempty" yaml:"views,omitempty"`

	// Alignment configuration - allows overriding default alignment for specific column types
	// Valid values: "left", "right", "center" (deprecated, use Columns instead)
	AlignmentOverrides map[string]string `json:"alignment_overrides,omitempty" yaml:"alignment_overrides,omitempty"`
}

// ViewConfig defines column configuration for a specific view
type ViewConfig struct {
	// Column configurations specific to this view
	Columns map[string]ColumnConfig `json:"columns,omitempty" yaml:"columns,omitempty"`

	// Custom columns specific to this view
	CustomColumns map[string]ColumnConfig `json:"custom_columns,omitempty" yaml:"custom_columns,omitempty"`
}

// MergeWith merges this ViewConfig with another, copying non-zero/empty values
func (vc *ViewConfig) MergeWith(other any) {
	otherConfig := extractViewConfig(other)
	if otherConfig == nil {
		return
	}

	// Merge column configurations
	if otherConfig.Columns != nil {
		if vc.Columns == nil {
			vc.Columns = make(map[string]ColumnConfig)
		}
		for key, value := range otherConfig.Columns {
			if existing, exists := vc.Columns[key]; exists {
				existing.MergeWith(value)
				vc.Columns[key] = existing
			} else {
				vc.Columns[key] = value
			}
		}
	}

	// Merge custom columns
	if otherConfig.CustomColumns != nil {
		if vc.CustomColumns == nil {
			vc.CustomColumns = make(map[string]ColumnConfig)
		}
		for key, value := range otherConfig.CustomColumns {
			if existing, exists := vc.CustomColumns[key]; exists {
				existing.MergeWith(value)
				vc.CustomColumns[key] = existing
			} else {
				vc.CustomColumns[key] = value
			}
		}
	}
}

// extractViewConfig extracts ViewConfig from various types
func extractViewConfig(other any) *ViewConfig {
	switch v := other.(type) {
	case ViewConfig:
		return &v
	case *ViewConfig:
		return v
	default:
		return nil
	}
}

// MergeWith merges this TableLimits with another, copying non-zero/empty values
func (tl *TableLimits) MergeWith(other any) {
	otherLimits := extractTableLimits(other)
	if otherLimits == nil {
		return
	}

	tl.mergeConfigurationFields(otherLimits)
}

// extractTableLimits extracts TableLimits from various types
func extractTableLimits(other any) *TableLimits {
	switch v := other.(type) {
	case TableLimits:
		return &v
	case *TableLimits:
		return v
	default:
		return nil
	}
}

// mergeConfigurationFields merges column configurations and view configurations
func (tl *TableLimits) mergeConfigurationFields(other *TableLimits) {
	tl.mergeGlobalColumns(other)
	tl.mergeCustomColumns(other)
	tl.mergeViewConfigurations(other)
	tl.mergeAlignmentOverrides(other)
}

// mergeGlobalColumns merges global column configurations
func (tl *TableLimits) mergeGlobalColumns(other *TableLimits) {
	if other.Columns == nil {
		return
	}
	if tl.Columns == nil {
		tl.Columns = make(map[string]ColumnConfig)
	}
	for key, value := range other.Columns {
		if existing, exists := tl.Columns[key]; exists {
			existing.MergeWith(value)
			tl.Columns[key] = existing
		} else {
			tl.Columns[key] = value
		}
	}
}

// mergeCustomColumns merges global custom column configurations
func (tl *TableLimits) mergeCustomColumns(other *TableLimits) {
	if other.CustomColumns == nil {
		return
	}
	if tl.CustomColumns == nil {
		tl.CustomColumns = make(map[string]ColumnConfig)
	}
	for key, value := range other.CustomColumns {
		if existing, exists := tl.CustomColumns[key]; exists {
			existing.MergeWith(value)
			tl.CustomColumns[key] = existing
		} else {
			tl.CustomColumns[key] = value
		}
	}
}

// mergeViewConfigurations merges per-view configurations
func (tl *TableLimits) mergeViewConfigurations(other *TableLimits) {
	if other.Views == nil {
		return
	}
	if tl.Views == nil {
		tl.Views = make(map[string]ViewConfig)
	}
	for key, value := range other.Views {
		if existing, exists := tl.Views[key]; exists {
			existing.MergeWith(value)
			tl.Views[key] = existing
		} else {
			tl.Views[key] = value
		}
	}
}

// mergeAlignmentOverrides merges alignment overrides (deprecated)
func (tl *TableLimits) mergeAlignmentOverrides(other *TableLimits) {
	if other.AlignmentOverrides == nil {
		return
	}
	if tl.AlignmentOverrides == nil {
		tl.AlignmentOverrides = make(map[string]string)
	}
	for key, value := range other.AlignmentOverrides {
		if value != "" {
			tl.AlignmentOverrides[key] = value
		}
	}
}
