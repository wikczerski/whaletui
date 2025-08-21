package volume

import "github.com/wikczerski/whaletui/internal/docker"

// Volume represents a Docker volume
type Volume = docker.Volume

// VolumeDetails represents detailed volume information
type VolumeDetails struct {
	Volume
	Status    map[string]any    `json:"status"`
	Options   map[string]string `json:"options"`
	UsageData VolumeUsageData   `json:"usage_data"`
}

// VolumeUsageData represents volume usage statistics
type VolumeUsageData struct {
	Size     int64 `json:"size"`
	RefCount int   `json:"ref_count"`
}
