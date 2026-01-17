package container

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Container represents a Docker container
type Container = shared.Container

// Details represents detailed container information
type Details struct {
	ID      string            `json:"id"`
	Names   []string          `json:"names"`
	Image   string            `json:"image"`
	ImageID string            `json:"image_id"`
	Command string            `json:"command"`
	Created int64             `json:"created"`
	Ports   []shared.Port     `json:"ports"`
	Labels  map[string]string `json:"labels"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
}
