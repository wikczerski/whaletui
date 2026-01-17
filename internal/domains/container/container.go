package container

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Container represents a Docker container
type Container = shared.Container

// Details represents detailed container information
type Details struct {
	Labels  map[string]string `json:"labels"`
	ID      string            `json:"id"`
	Image   string            `json:"image"`
	ImageID string            `json:"image_id"`
	Command string            `json:"command"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Names   []string          `json:"names"`
	Ports   []shared.Port     `json:"ports"`
	Created int64             `json:"created"`
}
