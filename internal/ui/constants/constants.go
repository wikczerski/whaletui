package constants

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// View names
const (
	ViewContainers = "containers"
	ViewImages     = "images"
	ViewVolumes    = "volumes"
	ViewNetworks   = "networks"
	ViewLogs       = "logs"
)

// DefaultView is the default view to show when the application starts
const DefaultView = ViewContainers

// UI layout constants
const (
	HeaderSectionHeight = 6
	StatusBarHeight     = 1
	TitleViewHeight     = 3
	BackButtonHeight    = 1
)

// Colors - These are now managed by the theme manager
// Default values are kept for backward compatibility
// Note: These constants are deprecated and will be removed in future versions
// Use the theme manager instead
const (
	HeaderColor     = tcell.ColorYellow
	BorderColor     = tcell.ColorWhite
	TextColor       = tcell.ColorWhite
	BackgroundColor = tcell.ColorDefault
)

// Basic color constants for theme configuration
const (
	ColorYellow   = "yellow"
	ColorWhite    = "white"
	ColorDefault  = "default"
	ColorGreen    = "green"
	ColorRed      = "red"
	ColorBlue     = "blue"
	ColorBlack    = "black"
	ColorGray     = "gray"
	ColorDarkGray = "darkgray"
)

// Main theme color constants
const (
	ThemeHeaderColor     = ColorYellow
	ThemeBorderColor     = ColorWhite
	ThemeTextColor       = ColorWhite
	ThemeBackgroundColor = ColorDefault
	ThemeSuccessColor    = ColorGreen
	ThemeWarningColor    = ColorYellow
	ThemeErrorColor      = ColorRed
	ThemeInfoColor       = ColorBlue
)

// Shell color constants
const (
	ShellBorderColor     = tcell.ColorYellow
	ShellTitleColor      = tcell.ColorYellow
	ShellTextColor       = tcell.ColorWhite
	ShellBackgroundColor = tcell.ColorBlack

	ShellCmdLabelColor       = tcell.ColorGreen
	ShellCmdBorderColor      = tcell.ColorGreen
	ShellCmdTextColor        = tcell.ColorWhite
	ShellCmdBackgroundColor  = tcell.ColorBlack
	ShellCmdPlaceholderColor = tcell.ColorGray
)

// Shell theme string constants
const (
	ShellThemeBorderColor         = ColorYellow
	ShellThemeTitleColor          = ColorYellow
	ShellThemeTextColor           = ColorWhite
	ShellThemeBackgroundColor     = ColorBlack
	ShellThemeCmdLabelColor       = ColorGreen
	ShellThemeCmdBorderColor      = ColorGreen
	ShellThemeCmdTextColor        = ColorWhite
	ShellThemeCmdBackgroundColor  = ColorBlack
	ShellThemeCmdPlaceholderColor = ColorGray
)

// ContainerExec color constants
const (
	ContainerExecLabelColor       = tcell.ColorYellow
	ContainerExecBorderColor      = tcell.ColorYellow
	ContainerExecTextColor        = tcell.ColorWhite
	ContainerExecBackgroundColor  = tcell.ColorDarkGray
	ContainerExecPlaceholderColor = tcell.ColorGray
	ContainerExecTitleColor       = tcell.ColorYellow
)

// ContainerExec theme string constants
const (
	ContainerExecThemeLabelColor       = ColorYellow
	ContainerExecThemeBorderColor      = ColorYellow
	ContainerExecThemeTextColor        = ColorWhite
	ContainerExecThemeBackgroundColor  = ColorDarkGray
	ContainerExecThemePlaceholderColor = ColorGray
	ContainerExecThemeTitleColor       = ColorYellow
)

// CommandMode color constants
const (
	CommandModeLabelColor       = tcell.ColorYellow
	CommandModeBorderColor      = tcell.ColorYellow
	CommandModeTextColor        = tcell.ColorWhite
	CommandModeBackgroundColor  = tcell.ColorDarkGray
	CommandModePlaceholderColor = tcell.ColorGray
	CommandModeTitleColor       = tcell.ColorYellow
)

// CommandMode theme string constants
const (
	CommandModeThemeLabelColor       = ColorYellow
	CommandModeThemeBorderColor      = ColorYellow
	CommandModeThemeTextColor        = ColorWhite
	CommandModeThemeBackgroundColor  = ColorDarkGray
	CommandModeThemePlaceholderColor = ColorGray
	CommandModeThemeTitleColor       = ColorYellow
)

// Table and row color constants
const (
	TableDefaultRowColor = tcell.ColorWhite
	TableSuccessColor    = tcell.ColorGreen
	TableWarningColor    = tcell.ColorYellow
	TableErrorColor      = tcell.ColorRed
	TableInfoColor       = tcell.ColorBlue
)

// UI element color constants
const (
	UIInvisibleColor = tcell.ColorDefault
)

// Key bindings
const (
	KeyStartContainer   = 's'
	KeyStopContainer    = 'S'
	KeyRestartContainer = 'r'
	KeyDeleteContainer  = 'd'
	KeyViewLogs         = 'l'
	KeyInspect          = 'i'
	KeyDeleteImage      = 'd'
	KeyDeleteVolume     = 'd'
	KeyDeleteNetwork    = 'd'
)

// Time formatting
const (
	TimeFormatRelative = "ago"
	TimeFormatAbsolute = "Jan 2, 2006 15:04:05"
	TimeThreshold24h   = 24 * time.Hour
)

// DockerInfoTemplate is the template for displaying Docker system information
const DockerInfoTemplate = `🐳 Docker Info
✅ Connected
🐋 Version: %s
📊 Containers: %d
🖼️  Images: %d
💾 Volumes: %d
🌐 Networks: %d
💻 OS: %s
🏗️  Architecture: %s
🔧 Driver: %s
📝 Logging: %s`

// StatusBarTemplate is the template for the status bar display
const StatusBarTemplate = "[%s] [Enter] Details [Q] Quit"
