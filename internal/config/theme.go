package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"gopkg.in/yaml.v3"
)

// DefaultTheme provides hardcoded default colors
var DefaultTheme = ThemeConfig{
	Colors: ThemeColors{
		Header:     constants.ColorYellow,
		Border:     constants.ColorWhite,
		Text:       constants.ColorWhite,
		Background: constants.ColorDefault,
		Success:    constants.ColorGreen,
		Warning:    constants.ColorYellow,
		Error:      constants.ColorRed,
		Info:       constants.ColorBlue,
	},
	Shell: ShellTheme{
		Border:     constants.ShellThemeBorderColor,
		Title:      constants.ShellThemeTitleColor,
		Text:       constants.ShellThemeTextColor,
		Background: constants.ShellThemeBackgroundColor,
		Cmd: ShellCmdTheme{
			Label:       constants.ShellThemeCmdLabelColor,
			Border:      constants.ShellThemeCmdBorderColor,
			Text:        constants.ShellThemeCmdTextColor,
			Background:  constants.ShellThemeCmdBackgroundColor,
			Placeholder: constants.ShellThemeCmdPlaceholderColor,
		},
	},
	ContainerExec: ContainerExecTheme{
		Label:       constants.ContainerExecThemeLabelColor,
		Border:      constants.ContainerExecThemeBorderColor,
		Text:        constants.ContainerExecThemeTextColor,
		Background:  constants.ContainerExecThemeBackgroundColor,
		Placeholder: constants.ContainerExecThemePlaceholderColor,
		Title:       constants.ContainerExecThemeTitleColor,
	},
	CommandMode: CommandModeTheme{
		Label:       constants.CommandModeThemeLabelColor,
		Border:      constants.CommandModeThemeBorderColor,
		Text:        constants.CommandModeThemeTextColor,
		Background:  constants.CommandModeThemeBackgroundColor,
		Placeholder: constants.CommandModeThemePlaceholderColor,
		Title:       constants.CommandModeThemeTitleColor,
	},
}

// ThemeManager manages theme configuration
type ThemeManager struct {
	config *ThemeConfig
	path   string
}

// NewThemeManager creates a new theme manager
func NewThemeManager(configPath string) *ThemeManager {
	tm := &ThemeManager{
		config: &DefaultTheme,
		path:   configPath,
	}
	tm.LoadTheme()
	return tm
}

// LoadTheme loads the theme configuration from file
func (tm *ThemeManager) LoadTheme() {
	// Try to load from the specified path
	if tm.path != "" {
		if err := tm.loadFromFile(tm.path); err == nil {
			return
		}
	}

	// Try common config locations
	configDirs := []string{
		"./config",
		"./themes",
		"$HOME/.config/whaletui",
		"$HOME/.whaletui",
	}

	for _, dir := range configDirs {
		expandedDir := os.ExpandEnv(dir)
		paths := []string{
			filepath.Join(expandedDir, "theme.yaml"),
			filepath.Join(expandedDir, "theme.yml"),
			filepath.Join(expandedDir, "theme.json"),
		}

		for _, path := range paths {
			if err := tm.loadFromFile(path); err == nil {
				tm.path = path
				return
			}
		}
	}
}

// loadFromFile loads theme from a specific file
func (tm *ThemeManager) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read theme file: %w", err)
	}

	var config ThemeConfig
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse YAML theme file: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse JSON theme file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported theme file format: %s", ext)
	}

	// Validate and merge with defaults
	tm.config = tm.mergeWithDefaults(&config)
	return nil
}

// mergeWithDefaults merges loaded config with defaults using interface methods
func (tm *ThemeManager) mergeWithDefaults(loaded *ThemeConfig) *ThemeConfig {
	merged := DefaultTheme

	// Use interface methods to merge all fields automatically
	merged.MergeWith(loaded)

	return &merged
}

// GetColor converts a color name to tcell.Color
func (tm *ThemeManager) GetColor(colorName string) tcell.Color {
	switch colorName {
	case "black":
		return tcell.ColorBlack
	case "red":
		return tcell.ColorRed
	case "green":
		return tcell.ColorGreen
	case "yellow":
		return tcell.ColorYellow
	case "blue":
		return tcell.ColorBlue
	case "magenta":
		return tcell.ColorPurple
	case "cyan":
		return tcell.ColorTeal
	case "white":
		return tcell.ColorWhite
	case "default":
		return tcell.ColorDefault
	case "gray":
		return tcell.ColorGray
	case "darkgray":
		return tcell.ColorDarkGray
	default:
		// Try to parse as hex color
		if len(colorName) == 7 && colorName[0] == '#' {
			if color, err := parseHexColor(colorName); err == nil {
				return color
			}
		}
		return tcell.ColorDefault
	}
}

// parseHexColor parses a hex color string
func parseHexColor(hex string) (tcell.Color, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return tcell.ColorDefault, fmt.Errorf("invalid hex color format")
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return tcell.ColorDefault, err
	}

	return tcell.NewRGBColor(int32(r), int32(g), int32(b)), nil
}

// GetHeaderColor returns the header color
func (tm *ThemeManager) GetHeaderColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Header)
}

// GetBorderColor returns the border color
func (tm *ThemeManager) GetBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Border)
}

// GetTextColor returns the text color
func (tm *ThemeManager) GetTextColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Text)
}

// GetBackgroundColor returns the background color
func (tm *ThemeManager) GetBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Background)
}

// GetSuccessColor returns the success color
func (tm *ThemeManager) GetSuccessColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Success)
}

// GetWarningColor returns the warning color
func (tm *ThemeManager) GetWarningColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Warning)
}

// GetErrorColor returns the error color
func (tm *ThemeManager) GetErrorColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Error)
}

// GetInfoColor returns the info color
func (tm *ThemeManager) GetInfoColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Info)
}

// GetShellBorderColor returns the shell border color
func (tm *ThemeManager) GetShellBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Border)
}

// GetShellTitleColor returns the shell title color
func (tm *ThemeManager) GetShellTitleColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Title)
}

// GetShellTextColor returns the shell text color
func (tm *ThemeManager) GetShellTextColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Text)
}

// GetShellBackgroundColor returns the shell background color
func (tm *ThemeManager) GetShellBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Background)
}

// GetShellCmdLabelColor returns the shell command label color
func (tm *ThemeManager) GetShellCmdLabelColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Label)
}

// GetShellCmdBorderColor returns the shell command border color
func (tm *ThemeManager) GetShellCmdBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Border)
}

// GetShellCmdTextColor returns the shell command text color
func (tm *ThemeManager) GetShellCmdTextColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Text)
}

// GetShellCmdBackgroundColor returns the shell command background color
func (tm *ThemeManager) GetShellCmdBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Background)
}

// GetShellCmdPlaceholderColor returns the shell command placeholder color
func (tm *ThemeManager) GetShellCmdPlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Placeholder)
}

// GetContainerExecLabelColor returns the container exec label color
func (tm *ThemeManager) GetContainerExecLabelColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Label)
}

// GetContainerExecBorderColor returns the container exec border color
func (tm *ThemeManager) GetContainerExecBorderColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Border)
}

// GetContainerExecTextColor returns the container exec text color
func (tm *ThemeManager) GetContainerExecTextColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Text)
}

// GetContainerExecBackgroundColor returns the container exec background color
func (tm *ThemeManager) GetContainerExecBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Background)
}

// GetContainerExecPlaceholderColor returns the container exec placeholder color
func (tm *ThemeManager) GetContainerExecPlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Placeholder)
}

// GetContainerExecTitleColor returns the container exec title color
func (tm *ThemeManager) GetContainerExecTitleColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Title)
}

// GetCommandModeLabelColor returns the command mode label color
func (tm *ThemeManager) GetCommandModeLabelColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Label)
}

// GetCommandModeBorderColor returns the command mode border color
func (tm *ThemeManager) GetCommandModeBorderColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Border)
}

// GetCommandModeTextColor returns the command mode text color
func (tm *ThemeManager) GetCommandModeTextColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Text)
}

// GetCommandModeBackgroundColor returns the command mode background color
func (tm *ThemeManager) GetCommandModeBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Background)
}

// GetCommandModePlaceholderColor returns the command mode placeholder color
func (tm *ThemeManager) GetCommandModePlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Placeholder)
}

// GetCommandModeTitleColor returns the command mode title color
func (tm *ThemeManager) GetCommandModeTitleColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Title)
}

// SaveTheme saves the current theme configuration to file
func (tm *ThemeManager) SaveTheme(path string) error {
	if path == "" {
		path = tm.path
	}
	if path == "" {
		path = "./config/theme.yaml"
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	ext := filepath.Ext(path)
	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(tm.config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	case ".json":
		data, err = json.MarshalIndent(tm.config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	tm.path = path
	return nil
}

// GetConfig returns the current theme configuration
func (tm *ThemeManager) GetConfig() *ThemeConfig {
	return tm.config
}

// GetPath returns the current theme file path
func (tm *ThemeManager) GetPath() string {
	return tm.path
}
