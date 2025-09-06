package config

import "reflect"

// Mergeable defines types that can merge with defaults
type Mergeable interface {
	MergeWith(other any)
}

// mergeStringFields is a helper function that merges non-empty string fields from src to dst
func mergeStringFields(dst, src any) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Field(i)

		if srcField.Kind() == reflect.String && srcField.String() != "" {
			dstField.Set(srcField)
		}
	}
}

// ThemeConfig holds the color theme configuration
type ThemeConfig struct {
	Colors        ThemeColors        `json:"colors"        yaml:"colors"`
	Shell         ShellTheme         `json:"shell"         yaml:"shell"`
	ContainerExec ContainerExecTheme `json:"containerExec" yaml:"containerExec"`
	CommandMode   CommandModeTheme   `json:"commandMode"   yaml:"commandMode"`
	TableLimits   TableLimits        `json:"tableLimits"   yaml:"tableLimits"`
}

// MergeWith merges this config with another, copying non-empty values
func (tc *ThemeConfig) MergeWith(other any) {
	otherConfig := extractThemeConfig(other)
	if otherConfig == nil {
		return
	}

	tc.mergeAllSections(otherConfig)
}

// extractThemeConfig extracts ThemeConfig from various types
func extractThemeConfig(other any) *ThemeConfig {
	switch v := other.(type) {
	case ThemeConfig:
		return &v
	case *ThemeConfig:
		return v
	default:
		return nil
	}
}

// mergeAllSections merges all sections of the theme config
func (tc *ThemeConfig) mergeAllSections(other *ThemeConfig) {
	tc.Colors.MergeWith(other.Colors)
	tc.Shell.MergeWith(other.Shell)
	tc.ContainerExec.MergeWith(other.ContainerExec)
	tc.CommandMode.MergeWith(other.CommandMode)
	tc.TableLimits.MergeWith(other.TableLimits)
}

// ThemeColors defines the color scheme
type ThemeColors struct {
	Header     string `json:"header"     yaml:"header"`
	Border     string `json:"border"     yaml:"border"`
	Text       string `json:"text"       yaml:"text"`
	Background string `json:"background" yaml:"background"`
	Success    string `json:"success"    yaml:"success"`
	Warning    string `json:"warning"    yaml:"warning"`
	Error      string `json:"error"      yaml:"error"`
	Info       string `json:"info"       yaml:"info"`
}

// MergeWith merges this ThemeColors with another, copying non-empty values
func (tc *ThemeColors) MergeWith(other any) {
	// Handle both value and pointer types
	var otherColors ThemeColors
	switch v := other.(type) {
	case ThemeColors:
		otherColors = v
	case *ThemeColors:
		if v != nil {
			otherColors = *v
		} else {
			return
		}
	default:
		return
	}

	mergeStringFields(tc, otherColors)
}

// ShellTheme defines the shell-specific color scheme
type ShellTheme struct {
	Border     string        `json:"border"     yaml:"border"`
	Title      string        `json:"title"      yaml:"title"`
	Text       string        `json:"text"       yaml:"text"`
	Background string        `json:"background" yaml:"background"`
	Cmd        ShellCmdTheme `json:"cmd"        yaml:"cmd"`
}

// MergeWith merges this ShellTheme with another, copying non-empty values
func (st *ShellTheme) MergeWith(other any) {
	// Handle both value and pointer types
	var otherShell ShellTheme
	switch v := other.(type) {
	case ShellTheme:
		otherShell = v
	case *ShellTheme:
		if v != nil {
			otherShell = *v
		} else {
			return
		}
	default:
		return
	}

	mergeStringFields(st, otherShell)
	st.Cmd.MergeWith(otherShell.Cmd)
}

// ShellCmdTheme defines the shell command input field color scheme
type ShellCmdTheme struct {
	Label       string `json:"label"       yaml:"label"`
	Border      string `json:"border"      yaml:"border"`
	Text        string `json:"text"        yaml:"text"`
	Background  string `json:"background"  yaml:"background"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
}

// MergeWith merges this ShellCmdTheme with another, copying non-empty values
func (sct *ShellCmdTheme) MergeWith(other any) {
	// Handle both value and pointer types
	var otherCmd ShellCmdTheme
	switch v := other.(type) {
	case ShellCmdTheme:
		otherCmd = v
	case *ShellCmdTheme:
		if v != nil {
			otherCmd = *v
		} else {
			return
		}
	default:
		return
	}

	mergeStringFields(sct, otherCmd)
}

// ContainerExecTheme defines the container exec input field color scheme
type ContainerExecTheme struct {
	Label       string `json:"label"       yaml:"label"`
	Border      string `json:"border"      yaml:"border"`
	Text        string `json:"text"        yaml:"text"`
	Background  string `json:"background"  yaml:"background"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Title       string `json:"title"       yaml:"title"`
}

// MergeWith merges this ContainerExecTheme with another, copying non-empty values
func (cet *ContainerExecTheme) MergeWith(other any) {
	// Handle both value and pointer types
	var otherExec ContainerExecTheme
	switch v := other.(type) {
	case ContainerExecTheme:
		otherExec = v
	case *ContainerExecTheme:
		if v != nil {
			otherExec = *v
		} else {
			return
		}
	default:
		return
	}

	mergeStringFields(cet, otherExec)
}

// CommandModeTheme defines the command mode input field color scheme
type CommandModeTheme struct {
	Label       string `json:"label"       yaml:"label"`
	Border      string `json:"border"      yaml:"border"`
	Text        string `json:"text"        yaml:"text"`
	Background  string `json:"background"  yaml:"background"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Title       string `json:"title"       yaml:"title"`
}

// MergeWith merges this CommandModeTheme with another, copying non-empty values
func (cmt *CommandModeTheme) MergeWith(other any) {
	// Handle both value and pointer types
	var otherMode CommandModeTheme
	switch v := other.(type) {
	case CommandModeTheme:
		otherMode = v
	case *CommandModeTheme:
		if v != nil {
			otherMode = *v
		} else {
			return
		}
	default:
		return
	}

	mergeStringFields(cmt, otherMode)
}

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
	// Merge global column configurations
	if other.Columns != nil {
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

	// Merge global custom columns
	if other.CustomColumns != nil {
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

	// Merge per-view configurations
	if other.Views != nil {
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

	// Merge alignment overrides (deprecated, use Columns instead)
	if other.AlignmentOverrides != nil {
		if tl.AlignmentOverrides == nil {
			tl.AlignmentOverrides = make(map[string]string)
		}
		for key, value := range other.AlignmentOverrides {
			if value != "" {
				tl.AlignmentOverrides[key] = value
			}
		}
	}
}
