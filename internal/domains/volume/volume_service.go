package volume

import (
	"context"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

type volumeService struct {
	*shared.BaseService[Volume]
}

// NewVolumeService creates a new volume service
func NewVolumeService(client *docker.Client) interfaces.VolumeService {
	base := shared.NewBaseService[Volume](client, "volume")

	// Set up Docker-specific functions
	base.ListFunc = func(client *docker.Client, ctx context.Context) ([]Volume, error) {
		dockerVolumes, err := client.ListVolumes(ctx)
		if err != nil {
			return nil, err
		}

		// Convert docker types to models (now they're the same via type alias)
		result := make([]Volume, len(dockerVolumes))
		for i, vol := range dockerVolumes {
			result[i] = Volume(vol)
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

func (s *volumeService) ListVolumes(ctx context.Context) ([]Volume, error) {
	return s.List(ctx)
}

func (s *volumeService) RemoveVolume(ctx context.Context, name string, force bool) error {
	return s.Remove(ctx, name, force)
}

func (s *volumeService) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	return s.Inspect(ctx, name)
}

// GetActions returns the available actions for volumes as a map
func (s *volumeService) GetActions() map[rune]string {
	return map[rune]string{
		'r': "Remove",
		'h': "History",
		'f': "Filter",
		't': "Sort",
		'i': "Inspect",
	}
}

// GetActionsString returns the available actions for volumes as a formatted string
func (s *volumeService) GetActionsString() string {
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}
