package volume

import (
	"context"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

type volumeService struct {
	*shared.BaseService[shared.Volume]
}

// NewVolumeService creates a new volume service
func NewVolumeService(client *docker.Client) interfaces.VolumeService {
	base := shared.NewBaseService[shared.Volume](client, "volume")
	base.ListFunc = createListVolumesFunc()
	base.InspectFunc = createInspectVolumeFunc()
	base.RemoveFunc = createRemoveVolumeFunc()
	return &volumeService{BaseService: base}
}

// createListVolumesFunc creates the function for listing volumes
func createListVolumesFunc() func(*docker.Client, context.Context) ([]shared.Volume, error) {
	return func(client *docker.Client, ctx context.Context) ([]shared.Volume, error) {
		dockerVolumes, err := client.ListVolumes(ctx)
		if err != nil {
			return nil, err
		}
		return convertDockerVolumes(dockerVolumes), nil
	}
}

// convertDockerVolumes converts Docker volume types to shared types
func convertDockerVolumes(dockerVolumes []docker.Volume) []shared.Volume {
	result := make([]shared.Volume, len(dockerVolumes))
	for i, vol := range dockerVolumes {
		result[i] = shared.Volume{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  vol.Created,
			Status:     make(map[string]any),
			Labels:     vol.Labels,
			Scope:      vol.Scope,
			Size:       vol.Size,
		}
	}
	return result
}

// createInspectVolumeFunc creates the function for inspecting volumes
func createInspectVolumeFunc() func(*docker.Client, context.Context, string) (map[string]any, error) {
	return func(client *docker.Client, ctx context.Context, name string) (map[string]any, error) {
		return client.InspectVolume(ctx, name)
	}
}

// createRemoveVolumeFunc creates the function for removing volumes
func createRemoveVolumeFunc() func(*docker.Client, context.Context, string, bool) error {
	return func(client *docker.Client, ctx context.Context, name string, force bool) error {
		return client.RemoveVolume(ctx, name, force)
	}
}

func (s *volumeService) ListVolumes(ctx context.Context) ([]shared.Volume, error) {
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
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n" +
		"<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}

// GetNavigation returns the available navigation options for volumes as a map
func (s *volumeService) GetNavigation() map[rune]string {
	return map[rune]string{
		'↑': "Up",
		'↓': "Down",
		':': "Command",
		'/': "Filter",
	}
}

// GetNavigationString returns the available navigation options for volumes as a formatted string
func (s *volumeService) GetNavigationString() string {
	return "↑/↓: Navigate\n<:> Command mode\n/: Filter"
}
