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
	Colors        ThemeColors        `json:"colors" yaml:"colors"`
	Shell         ShellTheme         `json:"shell" yaml:"shell"`
	ContainerExec ContainerExecTheme `json:"containerExec" yaml:"containerExec"`
	CommandMode   CommandModeTheme   `json:"commandMode" yaml:"commandMode"`
	HeaderLayout  HeaderLayout       `json:"headerLayout" yaml:"headerLayout"`
}

// MergeWith merges this config with another, copying non-empty values
func (tc *ThemeConfig) MergeWith(other any) {
	// Handle both value and pointer types
	var otherConfig ThemeConfig
	switch v := other.(type) {
	case ThemeConfig:
		otherConfig = v
	case *ThemeConfig:
		if v != nil {
			otherConfig = *v
		} else {
			return
		}
	default:
		return
	}

	tc.Colors.MergeWith(otherConfig.Colors)
	tc.Shell.MergeWith(otherConfig.Shell)
	tc.ContainerExec.MergeWith(otherConfig.ContainerExec)
	tc.CommandMode.MergeWith(otherConfig.CommandMode)
	tc.HeaderLayout.MergeWith(otherConfig.HeaderLayout)
}

// HeaderLayout defines the header column width configuration
type HeaderLayout struct {
	DockerInfoWidth int `json:"dockerInfoWidth" yaml:"dockerInfoWidth"`
	SpacerWidth     int `json:"spacerWidth" yaml:"spacerWidth"`
	NavigationWidth int `json:"navigationWidth" yaml:"navigationWidth"`
	ActionsWidth    int `json:"actionsWidth" yaml:"actionsWidth"`
	ContentWidth    int `json:"contentWidth" yaml:"contentWidth"`
	LogoWidth       int `json:"logoWidth" yaml:"logoWidth"`
}

// MergeWith merges this HeaderLayout with another, copying non-zero values
func (hl *HeaderLayout) MergeWith(other any) {
	// Handle both value and pointer types
	var otherLayout HeaderLayout
	switch v := other.(type) {
	case HeaderLayout:
		otherLayout = v
	case *HeaderLayout:
		if v != nil {
			otherLayout = *v
		} else {
			return
		}
	default:
		return
	}

	// Merge non-zero integer fields
	if otherLayout.DockerInfoWidth > 0 {
		hl.DockerInfoWidth = otherLayout.DockerInfoWidth
	}
	if otherLayout.SpacerWidth > 0 {
		hl.SpacerWidth = otherLayout.SpacerWidth
	}
	if otherLayout.NavigationWidth > 0 {
		hl.NavigationWidth = otherLayout.NavigationWidth
	}
	if otherLayout.ActionsWidth > 0 {
		hl.ActionsWidth = otherLayout.ActionsWidth
	}
	if otherLayout.ContentWidth > 0 {
		hl.ContentWidth = otherLayout.ContentWidth
	}
	if otherLayout.LogoWidth > 0 {
		hl.LogoWidth = otherLayout.LogoWidth
	}
}

// ThemeColors defines the color scheme
type ThemeColors struct {
	Header     string `json:"header" yaml:"header"`
	Border     string `json:"border" yaml:"border"`
	Text       string `json:"text" yaml:"text"`
	Background string `json:"background" yaml:"background"`
	Success    string `json:"success" yaml:"success"`
	Warning    string `json:"warning" yaml:"warning"`
	Error      string `json:"error" yaml:"error"`
	Info       string `json:"info" yaml:"info"`
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
	Border     string        `json:"border" yaml:"border"`
	Title      string        `json:"title" yaml:"title"`
	Text       string        `json:"text" yaml:"text"`
	Background string        `json:"background" yaml:"background"`
	Cmd        ShellCmdTheme `json:"cmd" yaml:"cmd"`
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
	Label       string `json:"label" yaml:"label"`
	Border      string `json:"border" yaml:"border"`
	Text        string `json:"text" yaml:"text"`
	Background  string `json:"background" yaml:"background"`
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
	Label       string `json:"label" yaml:"label"`
	Border      string `json:"border" yaml:"border"`
	Text        string `json:"text" yaml:"text"`
	Background  string `json:"background" yaml:"background"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Title       string `json:"title" yaml:"title"`
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
	Label       string `json:"label" yaml:"label"`
	Border      string `json:"border" yaml:"border"`
	Text        string `json:"text" yaml:"text"`
	Background  string `json:"background" yaml:"background"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Title       string `json:"title" yaml:"title"`
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
