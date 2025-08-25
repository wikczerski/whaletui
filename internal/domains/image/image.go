// Package image provides Docker image management functionality for WhaleTUI.
package image

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Image represents a Docker image
type Image = shared.Image

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
