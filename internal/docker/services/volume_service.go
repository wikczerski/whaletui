package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	domaintypes "github.com/wikczerski/whaletui/internal/docker/types"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// VolumeService handles volume-related operations
type VolumeService struct {
	cli *client.Client
	log *slog.Logger
}

// NewVolumeService creates a new VolumeService
func NewVolumeService(cli *client.Client, log *slog.Logger) *VolumeService {
	return &VolumeService{
		cli: cli,
		log: log,
	}
}

// ListVolumes lists all volumes
func (s *VolumeService) ListVolumes(ctx context.Context) ([]domaintypes.Volume, error) {
	vols, err := s.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]domaintypes.Volume, 0, len(vols.Volumes))
	for _, vol := range vols.Volumes {
		volume := s.createVolumeFromAPI(vol)
		result = append(result, volume)
	}

	utils.SortVolumesByName(result)
	return result, nil
}

// InspectVolume inspects a volume
func (s *VolumeService) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	volumeInfo, err := s.cli.VolumeInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("volume inspect failed %s: %w", name, err)
	}
	return utils.MarshalToMap(volumeInfo)
}

// RemoveVolume removes a volume
func (s *VolumeService) RemoveVolume(ctx context.Context, name string, force bool) error {
	if err := utils.ValidateID(name, "volume name"); err != nil {
		return err
	}

	if err := s.cli.VolumeRemove(ctx, name, force); err != nil {
		return fmt.Errorf("failed to remove volume %s: %w", name, err)
	}

	return nil
}

// Helper methods

func (s *VolumeService) createVolumeFromAPI(vol *volume.Volume) domaintypes.Volume {
	created := time.Time{}
	if vol.CreatedAt != "" {
		created, _ = time.Parse(time.RFC3339, vol.CreatedAt)
	}

	return domaintypes.Volume{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Created:    created,
		Size:       "", // Size is not available in VolumeList
	}
}
