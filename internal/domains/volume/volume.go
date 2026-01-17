// Package volume provides Docker volume management functionality for WhaleTUI.
package volume

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Volume represents a Docker volume
type Volume = shared.Volume

// Details represents detailed volume information
type Details struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Options    map[string]string `json:"options"`
	Scope      string            `json:"scope"`
}
