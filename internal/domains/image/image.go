package image

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Image represents a Docker image
type Image = shared.Image

// Details represents detailed image information
type Details struct {
	ID           string            `json:"id"`
	RepoTags     []string          `json:"repo_tags"`
	Created      int64             `json:"created"`
	Size         int64             `json:"size"`
	Labels       map[string]string `json:"labels"`
	Config       map[string]any    `json:"config"`
	Architecture string            `json:"architecture"`
	Os           string            `json:"os"`
}
