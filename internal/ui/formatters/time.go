package formatters

import (
	"fmt"
	"time"

	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// FormatTime formats a time.Time to a human-readable string
func FormatTime(t time.Time) string {
	if time.Since(t) < constants.TimeThreshold24h {
		return fmt.Sprintf("%s %s", formatDuration(time.Since(t)), constants.TimeFormatRelative)
	}
	return t.Format(constants.TimeFormatAbsolute)
}

// formatDuration formats a duration to a human-readable string
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())

	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	case seconds < 86400:
		return fmt.Sprintf("%dh %dm", seconds/3600, (seconds%3600)/60)
	default:
		return fmt.Sprintf("%dd %dh", seconds/86400, (seconds%86400)/3600)
	}
}
