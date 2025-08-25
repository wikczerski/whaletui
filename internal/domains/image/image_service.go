package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared"
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
func (s *imageService) ListImages(ctx context.Context) ([]shared.Image, error) {
	images, err := s.client.ListImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]shared.Image, 0, len(images))
	for _, img := range images {
		sharedImg := shared.Image{
			ID:       img.ID,
			RepoTags: []string{img.Repository + ":" + img.Tag},
			Created:  img.Created,
			Size:     img.Size,
			Labels:   make(map[string]string),
		}
		result = append(result, sharedImg)
	}
	return result, nil
}

// RemoveImage removes an image
func (s *imageService) RemoveImage(ctx context.Context, id string, force bool) error {
	if s.client == nil {
		return errors.New("docker client is not initialized")
	}

	if id == "" {
		return errors.New("image ID cannot be empty")
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
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n" +
		"<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}
