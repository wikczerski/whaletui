// Package constants provides UI constants and configuration values for the WhaleTUI application.
package constants

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// View names
const (
	ViewContainers    = "containers"
	ViewImages        = "images"
	ViewVolumes       = "volumes"
	ViewNetworks      = "networks"
	ViewLogs          = "logs"
	ViewSwarmServices = "swarmServices"
	ViewSwarmNodes    = "swarmNodes"
	ViewDockerInfo    = "dockerInfo"
)

// DefaultView is the default view to show when the application starts
const DefaultView = ViewContainers

// UI layout constants
const (
	HeaderSectionHeight = 9 // Increased from 6 to 9 (3 extra rows)
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

// WhaleTuiLogo is the ASCII art logo for the WhaleTUI application
const WhaleTuiLogo = ` _    _ _           _      _______    _
| |  | | |         | |    |__   __|  (_)
| |  | | |__   __ _| | ___   | |_   _ _
| |/\| | '_ \ / _` + "`" + ` | |/ _ \  | | | | | |
\  /\  / | | | (_| | |  __/  |_| |_| | |
 \/  \/|_| |_|\__,_|_|\___|  |_|\__,_|_|`

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
	KeyScaleService     = 's'
	KeyRemoveService    = 'd'
	KeyRemoveNode       = 'd'
	KeyDrainNode        = 'D'
	KeyActivateNode     = 'a'
)

// Time formatting
const (
	TimeFormatRelative = "ago"
	TimeFormatAbsolute = "Jan 2, 2006 15:04:05"
	TimeThreshold24h   = 24 * time.Hour
)

// DockerInfoTemplate is the template for displaying Docker system information
const DockerInfoTemplate = `🐳 Docker Info
%s
🐋 Docker: %s
💻 OS: %s
📝 Logging: %s
🔗 Method: %s
🚀 WhaleTui: %s`

// StatusBarTemplate is the template for the status bar display
const StatusBarTemplate = "[%s] [Enter] Details [Q] Quit"

// AppVersion holds the application version for display in UI
// This will be set by the build process or defaults to "dev"
var AppVersion = "dev"

// SetAppVersion sets the application version from the cmd package
func SetAppVersion(version string) {
	AppVersion = version
}

// Header column width constants
const (
	HeaderDockerInfoWidth = 30 // Width for Docker info column
	HeaderSpacerWidth     = 1  // Width for spacer column between sections (minimal spacing)
	HeaderNavigationWidth = 25 // Width for navigation column
	HeaderActionsWidth    = 25 // Width for actions column
	HeaderContentWidth    = 25 // Width for general content columns (fallback)
	HeaderLogoWidth       = 20 // Width for logo column
)

// Default column character limits
var DefaultColumnLimits = map[string]int{
	// Container/Image IDs
	"id": 12,
	// Container/Volume names
	"name": 30,
	// Image names
	"image": 35,
	// Status messages
	"status": 18,
	// State values
	"state": 15,
	// Port mappings
	"ports": 25,
	// Created timestamps
	"created": 25,
	// Size values
	"size": 15,
	// Volume drivers
	"driver": 15,
	// Mount points
	"mountpoint": 50,
	// Image repositories
	"repository": 40,
	// Image tags
	"tag": 20,
	// Network names
	"network": 25,
	// Network scope
	"scope": 10,
	// General descriptions
	"description": 50,
	// Container count
	"containers": 12,
	// Replica count
	"replicas": 10,
	// Engine version
	"engine_version": 25,
	// Hostname
	"hostname": 25,
	// IP addresses
	"address": 20,
	// Swarm mode
	"mode": 15,
	// Swarm role
	"role": 10,
	// Swarm availability
	"availability": 15,
	// Manager status
	"manager_status": 15,
}

// Default column widths
var DefaultColumnWidths = map[string]int{
	// Container/Image IDs - slightly wider than limit for readability
	"id": 15,
	// Container/Volume names - wider for better readability
	"name": 35,
	// Image names - often long repository paths
	"image": 35,
	// Status messages - need space for status text
	"status": 18,
	// State values - running, stopped, etc.
	"state": 15,
	// Port mappings - 0.0.0.0:8080->80/tcp format
	"ports": 30,
	// Created timestamps - ISO format dates
	"created": 25,
	// Size values - 1.2GB, 500MB format
	"size": 15,
	// Volume drivers - local, nfs, etc.
	"driver": 15,
	// Mount points - often long paths
	"mountpoint": 55,
	// Image repositories - docker.io/library/nginx
	"repository": 40,
	// Image tags - latest, v1.0.0, etc.
	"tag": 20,
	// Network names - bridge, host, custom networks
	"network": 25,
	// Network scope - local, global
	"scope": 10,
	// General descriptions - can be long
	"description": 55,
	// Container count - just numbers
	"containers": 10,
	// Replica count - just numbers
	"replicas": 10,
	// Engine version - 20.10.0 format
	"engine_version": 25,
	// Hostname - server names
	"hostname": 25,
	// IP addresses - 192.168.1.100 format
	"address": 20,
	// Swarm mode - replicated, global
	"mode": 15,
	// Swarm role - manager, worker
	"role": 10,
	// Swarm availability - active, pause, drain
	"availability": 15,
	// Manager status - leader, reachable, etc.
	"manager_status": 15,
}

// Default column alignments
var DefaultColumnAlignments = map[string]int{
	// Right-aligned columns for numerical data and IDs
	"id":             tview.AlignRight,
	"size":           tview.AlignRight,
	"containers":     tview.AlignRight,
	"created":        tview.AlignRight,
	"replicas":       tview.AlignRight,
	"engine_version": tview.AlignRight,
	"status":         tview.AlignRight,
	"state":          tview.AlignRight,
	// Left-aligned columns for text-based content
	"name":           tview.AlignLeft,
	"image":          tview.AlignLeft,
	"repository":     tview.AlignLeft,
	"tag":            tview.AlignLeft,
	"driver":         tview.AlignLeft,
	"mountpoint":     tview.AlignLeft,
	"network":        tview.AlignLeft,
	"scope":          tview.AlignLeft,
	"description":    tview.AlignLeft,
	"hostname":       tview.AlignLeft,
	"address":        tview.AlignLeft,
	"ports":          tview.AlignLeft,
	"mode":           tview.AlignLeft,
	"role":           tview.AlignLeft,
	"availability":   tview.AlignLeft,
	"manager_status": tview.AlignLeft,
}
