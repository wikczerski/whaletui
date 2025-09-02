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
	HeaderLayout  HeaderLayout       `json:"headerLayout"  yaml:"headerLayout"`
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
	tc.HeaderLayout.MergeWith(other.HeaderLayout)
	tc.TableLimits.MergeWith(other.TableLimits)
}

// HeaderLayout defines the header column width configuration
type HeaderLayout struct {
	DockerInfoWidth int `json:"dockerInfoWidth" yaml:"dockerInfoWidth"`
	SpacerWidth     int `json:"spacerWidth"     yaml:"spacerWidth"`
	NavigationWidth int `json:"navigationWidth" yaml:"navigationWidth"`
	ActionsWidth    int `json:"actionsWidth"    yaml:"actionsWidth"`
	ContentWidth    int `json:"contentWidth"    yaml:"contentWidth"`
	LogoWidth       int `json:"logoWidth"       yaml:"logoWidth"`
}

// MergeWith merges this HeaderLayout with another, copying non-zero values
func (hl *HeaderLayout) MergeWith(other any) {
	otherLayout := extractHeaderLayout(other)
	if otherLayout == nil {
		return
	}

	hl.mergeIntegerFields(otherLayout)
}

// extractHeaderLayout extracts HeaderLayout from various types
func extractHeaderLayout(other any) *HeaderLayout {
	switch v := other.(type) {
	case HeaderLayout:
		return &v
	case *HeaderLayout:
		return v
	default:
		return nil
	}
}

// mergeIntegerFields merges non-zero integer fields
func (hl *HeaderLayout) mergeIntegerFields(other *HeaderLayout) {
	hl.mergeWidthField(&hl.DockerInfoWidth, other.DockerInfoWidth)
	hl.mergeWidthField(&hl.SpacerWidth, other.SpacerWidth)
	hl.mergeWidthField(&hl.NavigationWidth, other.NavigationWidth)
	hl.mergeWidthField(&hl.ActionsWidth, other.ActionsWidth)
	hl.mergeWidthField(&hl.ContentWidth, other.ContentWidth)
	hl.mergeWidthField(&hl.LogoWidth, other.LogoWidth)
}

// mergeWidthField merges a single width field if the other value is greater than 0
func (hl *HeaderLayout) mergeWidthField(field *int, otherValue int) {
	if otherValue > 0 {
		*field = otherValue
	}
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

// TableLimits defines character limits for table columns
type TableLimits struct {
	ID          int `json:"id,omitempty"          yaml:"id,omitempty"`
	Name        int `json:"name,omitempty"        yaml:"name,omitempty"`
	Image       int `json:"image,omitempty"       yaml:"image,omitempty"`
	Status      int `json:"status,omitempty"      yaml:"status,omitempty"`
	State       int `json:"state,omitempty"       yaml:"state,omitempty"`
	Ports       int `json:"ports,omitempty"       yaml:"ports,omitempty"`
	Created     int `json:"created,omitempty"     yaml:"created,omitempty"`
	Size        int `json:"size,omitempty"        yaml:"size,omitempty"`
	Driver      int `json:"driver,omitempty"      yaml:"driver,omitempty"`
	Mountpoint  int `json:"mountpoint,omitempty"  yaml:"mountpoint,omitempty"`
	Repository  int `json:"repository,omitempty"  yaml:"repository,omitempty"`
	Tag         int `json:"tag,omitempty"         yaml:"tag,omitempty"`
	Network     int `json:"network,omitempty"     yaml:"network,omitempty"`
	Scope       int `json:"scope,omitempty"       yaml:"scope,omitempty"`
	Description int `json:"description,omitempty" yaml:"description,omitempty"`
}

// MergeWith merges this TableLimits with another, copying non-zero values
func (tl *TableLimits) MergeWith(other any) {
	otherLimits := extractTableLimits(other)
	if otherLimits == nil {
		return
	}

	tl.mergeIntegerFields(otherLimits)
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

// mergeIntegerFields merges non-zero integer fields
func (tl *TableLimits) mergeIntegerFields(other *TableLimits) {
	tl.mergeLimitField(&tl.ID, other.ID)
	tl.mergeLimitField(&tl.Name, other.Name)
	tl.mergeLimitField(&tl.Image, other.Image)
	tl.mergeLimitField(&tl.Status, other.Status)
	tl.mergeLimitField(&tl.State, other.State)
	tl.mergeLimitField(&tl.Ports, other.Ports)
	tl.mergeLimitField(&tl.Created, other.Created)
	tl.mergeLimitField(&tl.Size, other.Size)
	tl.mergeLimitField(&tl.Driver, other.Driver)
	tl.mergeLimitField(&tl.Mountpoint, other.Mountpoint)
	tl.mergeLimitField(&tl.Repository, other.Repository)
	tl.mergeLimitField(&tl.Tag, other.Tag)
	tl.mergeLimitField(&tl.Network, other.Network)
	tl.mergeLimitField(&tl.Scope, other.Scope)
	tl.mergeLimitField(&tl.Description, other.Description)
}

// mergeLimitField merges a single limit field if the other value is greater than 0
func (tl *TableLimits) mergeLimitField(field *int, otherValue int) {
	if otherValue > 0 {
		*field = otherValue
	}
}
