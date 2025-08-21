package image

import (
	"time"

	"github.com/wikczerski/whaletui/internal/docker"
)

// Image represents a Docker image
type Image = docker.Image

// Details represents detailed image information
type Details struct {
	Image
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	Author       string            `json:"author"`
	Comment      string            `json:"comment"`
	Config       Config            `json:"config"`
	History      []History         `json:"history"`
	Labels       map[string]string `json:"labels"`
}

// Config represents image configuration
type Config struct {
	User         string              `json:"user"`
	WorkingDir   string              `json:"working_dir"`
	Entrypoint   []string            `json:"entrypoint"`
	Cmd          []string            `json:"cmd"`
	Environment  []string            `json:"environment"`
	ExposedPorts map[string]struct{} `json:"exposed_ports"`
	Volumes      map[string]struct{} `json:"volumes"`
	Labels       map[string]string   `json:"labels"`
}

// History represents image layer history
type History struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Comment    string    `json:"comment"`
	EmptyLayer bool      `json:"empty_layer"`
}
