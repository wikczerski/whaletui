package services

import (
	"context"

	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/models"
)

type volumeService struct {
	*BaseService[models.Volume]
}

// NewVolumeService creates a new volume service
func NewVolumeService(client *docker.Client) VolumeService {
	base := NewBaseService[models.Volume](client, "volume")

	// Set up Docker-specific functions
	base.ListFunc = func(client *docker.Client, ctx context.Context) ([]models.Volume, error) {
		dockerVolumes, err := client.ListVolumes(ctx)
		if err != nil {
			return nil, err
		}

		// Convert docker types to models (now they're the same via type alias)
		result := make([]models.Volume, len(dockerVolumes))
		for i, vol := range dockerVolumes {
			result[i] = models.Volume(vol)
		}
		return result, nil
	}

	base.InspectFunc = func(client *docker.Client, ctx context.Context, name string) (map[string]any, error) {
		return client.InspectVolume(ctx, name)
	}

	base.RemoveFunc = func(client *docker.Client, ctx context.Context, name string, force bool) error {
		return client.RemoveVolume(ctx, name, force)
	}

	return &volumeService{BaseService: base}
}

func (s *volumeService) ListVolumes(ctx context.Context) ([]models.Volume, error) {
	return s.List(ctx)
}

func (s *volumeService) RemoveVolume(ctx context.Context, name string, force bool) error {
	return s.Remove(ctx, name, force)
}

func (s *volumeService) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	return s.Inspect(ctx, name)
}
