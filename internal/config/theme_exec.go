package config

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
