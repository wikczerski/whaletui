// Package volume provides Docker volume management functionality for WhaleTUI.
package volume

import "github.com/wikczerski/whaletui/internal/shared"

// Volume represents a Docker volume
type Volume = shared.Volume

// Details represents detailed volume information
type Details struct {
	Volume
	Status    map[string]any    `json:"status"`
	Options   map[string]string `json:"options"`
	UsageData UsageData         `json:"usage_data"`
}

// UsageData represents volume usage statistics
type UsageData struct {
	Size     int64 `json:"size"`
	RefCount int   `json:"ref_count"`
}
