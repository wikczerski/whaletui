package image

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

// imageService implements ImageService
type imageService struct {
	client *docker.Client
}

// NewImageService creates a new image service
func NewImageService(client *docker.Client) interfaces.ImageService {
	if client == nil {
		return nil
	}
	return &imageService{
		client: client,
	}
}

// ListImages retrieves all images
func (s *imageService) ListImages(ctx context.Context) ([]Image, error) {
	images, err := s.client.ListImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]Image, 0, len(images))
	for _, img := range images {
		result = append(result, Image(img))
	}
	return result, nil
}

// RemoveImage removes an image
func (s *imageService) RemoveImage(ctx context.Context, id string, force bool) error {
	if s.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}

	if id == "" {
		return fmt.Errorf("image ID cannot be empty")
	}

	return s.client.RemoveImage(ctx, id, force)
}

// InspectImage inspects an image
func (s *imageService) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	return s.client.InspectImage(ctx, id)
}

// GetActions returns the available actions for images as a map
func (s *imageService) GetActions() map[rune]string {
	return map[rune]string{
		'r': "Remove",
		'h': "History",
		'f': "Filter",
		't': "Sort",
		'i': "Inspect",
	}
}

// GetActionsString returns the available actions for images as a formatted string
func (s *imageService) GetActionsString() string {
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}
