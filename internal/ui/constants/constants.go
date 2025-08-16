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

// Default view
const DefaultView = ViewContainers

// UI layout constants
const (
	HeaderSectionHeight = 6
	StatusBarHeight     = 1
	TitleViewHeight     = 3
	BackButtonHeight    = 1
)

// Colors
const (
	HeaderColor     = tcell.ColorYellow
	BorderColor     = tcell.ColorWhite
	TextColor       = tcell.ColorWhite
	BackgroundColor = tcell.ColorDefault
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

// Docker info template
const DockerInfoTemplate = `🐳 Docker Info
✅ Connected
🐋 Version: %s
📊 Containers: %d`

// Status bar template
const StatusBarTemplate = "%s | %s for details | %s to quit"
