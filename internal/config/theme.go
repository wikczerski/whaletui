package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	TableLimits: TableLimits{
		// Default column configurations will be handled by TableFormatter defaults
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

	// Load theme if path is provided
	if configPath != "" {
		tm.LoadTheme()
	}

	return tm
}

// LoadTheme loads the theme configuration from file
func (tm *ThemeManager) LoadTheme() {
	// If a specific path is provided, try to load from it and don't fall back
	if tm.path != "" {
		if tm.tryLoadFromSpecifiedPath() {
			return
		}
		// If loading from specified path fails, don't fall back to other locations
		// This ensures tests fail when they expect a specific theme file
		return
	}

	// Only try fallback locations if no specific path was provided
	tm.tryLoadFromFallbackLocations()
}

// GetColor converts a color name to tcell.Color
func (tm *ThemeManager) GetColor(colorName string) tcell.Color {
	if color := getNamedColor(colorName); color != tcell.ColorDefault {
		return color
	}

	return getHexColorOrDefault(colorName)
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
	path = getSavePath(path, tm.path)

	if err := ensureSaveDirectory(path); err != nil {
		return err
	}

	data, err := marshalThemeData(path, tm.config)
	if err != nil {
		return err
	}

	if err := writeThemeFile(path, data); err != nil {
		return err
	}

	tm.path = path
	return nil
}

// GetConfig returns the current theme configuration
func (tm *ThemeManager) GetConfig() *ThemeConfig {
	return tm.config
}

// GetTableLimits returns the current table limits configuration
func (tm *ThemeManager) GetTableLimits() TableLimits {
	return tm.config.TableLimits
}

// ReloadTheme reloads the theme configuration from file
func (tm *ThemeManager) ReloadTheme() error {
	if tm.path == "" {
		return errors.New("no theme path specified")
	}

	// Load the theme from file
	err := tm.loadFromFile(tm.path)
	if err != nil {
		return fmt.Errorf("failed to reload theme: %w", err)
	}

	return nil
}

// GetPath returns the current theme file path
func (tm *ThemeManager) GetPath() string {
	return tm.path
}

// GetCurrentThemeInfo returns debug information about the current theme
func (tm *ThemeManager) GetCurrentThemeInfo() map[string]any {
	return map[string]any{
		"path":            tm.path,
		"headerColor":     tm.config.Colors.Header,
		"borderColor":     tm.config.Colors.Border,
		"textColor":       tm.config.Colors.Text,
		"backgroundColor": tm.config.Colors.Background,
	}
}

// tryLoadFromSpecifiedPath tries to load from the specified path first
func (tm *ThemeManager) tryLoadFromSpecifiedPath() bool {
	if tm.path != "" {
		err := tm.loadFromFile(tm.path)
		if err != nil {
			// Debug: Log the error for troubleshooting
			fmt.Printf("DEBUG: Failed to load theme from specified path %s: %v\n", tm.path, err)
		}
		return err == nil
	}
	return false
}

// tryLoadFromFallbackLocations tries to load from common config locations
func (tm *ThemeManager) tryLoadFromFallbackLocations() {
	configDirs := getFallbackConfigDirs()

	for _, dir := range configDirs {
		expandedDir := os.ExpandEnv(dir)
		if tm.tryLoadFromDirectory(expandedDir) {
			return
		}
	}
}

// getFallbackConfigDirs returns the list of fallback configuration directories
func getFallbackConfigDirs() []string {
	return []string{
		"./config",
		"./themes",
		"$HOME/.config/whaletui",
		"$HOME/.whaletui",
	}
}

// tryLoadFromDirectory tries to load a theme from a specific directory
func (tm *ThemeManager) tryLoadFromDirectory(dir string) bool {
	paths := getThemeFilePaths(dir)

	for _, path := range paths {
		if err := tm.loadFromFile(path); err == nil {
			tm.path = path
			return true
		}
	}
	return false
}

// getThemeFilePaths returns the list of possible theme file paths in a directory
func getThemeFilePaths(dir string) []string {
	return []string{
		filepath.Join(dir, "theme.yaml"),
		filepath.Join(dir, "theme.yml"),
		filepath.Join(dir, "theme.json"),
	}
}

// getNamedColor returns a named color or ColorDefault if not found
func getNamedColor(colorName string) tcell.Color {
	colorMap := getColorMap()
	if color, exists := colorMap[colorName]; exists {
		return color
	}
	return tcell.ColorDefault
}

// getColorMap returns the mapping of color names to tcell colors
func getColorMap() map[string]tcell.Color {
	return map[string]tcell.Color{
		"black":    tcell.ColorBlack,
		"red":      tcell.ColorRed,
		"green":    tcell.ColorGreen,
		"yellow":   tcell.ColorYellow,
		"blue":     tcell.ColorBlue,
		"magenta":  tcell.ColorPurple,
		"cyan":     tcell.ColorTeal,
		"white":    tcell.ColorWhite,
		"default":  tcell.ColorDefault,
		"gray":     tcell.ColorGray,
		"darkgray": tcell.ColorDarkGray,
	}
}

// getHexColorOrDefault tries to parse a hex color or returns default
func getHexColorOrDefault(colorName string) tcell.Color {
	if len(colorName) == 7 && colorName[0] == '#' {
		if color, err := parseHexColor(colorName); err == nil {
			return color
		}
	}
	return tcell.ColorDefault
}

// getSavePath determines the path to save the theme to
func getSavePath(path, defaultPath string) string {
	if path == "" {
		path = defaultPath
	}
	if path == "" {
		path = "./config/theme.yaml"
	}
	return path
}

// ensureSaveDirectory ensures the directory for saving exists
func ensureSaveDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o750)
}

// marshalThemeData marshals the theme data based on file extension
func marshalThemeData(path string, config *ThemeConfig) ([]byte, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		return yaml.Marshal(config)
	case ".json":
		return json.MarshalIndent(config, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// writeThemeFile writes the theme data to the specified path
func writeThemeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}

// loadFromFile loads theme from a specific file
func (tm *ThemeManager) loadFromFile(path string) error {
	resolvedPath, err := tm.validateThemePath(path)
	if err != nil {
		return err
	}

	// nolint:gosec // Path is validated by validateThemePath before this function is called
	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to read theme file: %w", err)
	}

	config, err := tm.parseThemeData(resolvedPath, data)
	if err != nil {
		return err
	}

	tm.config = tm.mergeWithDefaults(config)
	return nil
}

// validateThemePath validates the theme file path to prevent directory traversal
// and returns the resolved absolute path
func (tm *ThemeManager) validateThemePath(path string) (string, error) {
	// Clean the path to remove any directory traversal attempts
	cleanPath := filepath.Clean(path)

	// Convert relative paths to absolute paths based on current working directory
	if !filepath.IsAbs(cleanPath) {
		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve relative path: %w", err)
		}
		cleanPath = absPath
	}

	// Additional security: check for suspicious patterns
	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", errors.New("theme file path contains directory traversal attempts")
	}

	// Check for home directory expansion attempts (but allow Windows short names like ~1)
	// Only reject paths that start with ~ or contain ~/ which could be home directory expansion
	if strings.HasPrefix(cleanPath, "~") || strings.Contains(cleanPath, "~/") {
		return "", errors.New("theme file path contains home directory expansion attempts")
	}

	return cleanPath, nil
}

// parseThemeData parses theme data based on file extension
func (tm *ThemeManager) parseThemeData(path string, data []byte) (*ThemeConfig, error) {
	var config ThemeConfig
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML theme file: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON theme file: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported theme file format: %s", ext)
	}

	return &config, nil
}

// mergeWithDefaults merges loaded config with defaults using interface methods
func (tm *ThemeManager) mergeWithDefaults(loaded *ThemeConfig) *ThemeConfig {
	merged := DefaultTheme
	merged.MergeWith(loaded)
	return &merged
}

// parseHexColor parses a hex color string
func parseHexColor(hex string) (tcell.Color, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return tcell.ColorDefault, errors.New("invalid hex color format")
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return tcell.ColorDefault, err
	}

	return tcell.NewRGBColor(int32(r), int32(g), int32(b)), nil
}
