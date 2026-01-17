package config

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
