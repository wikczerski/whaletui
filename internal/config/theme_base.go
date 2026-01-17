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
