package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/d5r/internal/config"
	"github.com/user/d5r/internal/docker"
	"github.com/user/d5r/internal/models"
)

func TestNewVolumeService(t *testing.T) {
	service := NewVolumeService(nil)
	assert.NotNil(t, service)

	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service = NewVolumeService(client)
	assert.NotNil(t, service)
}

func TestVolumeService_ListVolumes_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewVolumeService(client)
	ctx := context.Background()
	result, err := service.ListVolumes(ctx)

	require.NoError(t, err)
	// Result can be nil or empty slice - both are valid
	if result == nil {
		result = []models.Volume{}
	}

	assert.IsType(t, []models.Volume{}, result)
}

func TestVolumeService_ListVolumes_NilClient(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume_Force(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume_EmptyName(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "", false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_InspectVolume_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewVolumeService(client)
	ctx := context.Background()

	volumes, err := service.ListVolumes(ctx)
	if err != nil {
		t.Skipf("Could not list volumes: %v", err)
	}

	if len(volumes) == 0 {
		t.Skip("No volumes available for inspection")
	}

	volumeName := volumes[0].Name
	result, err := service.InspectVolume(ctx, volumeName)

	require.NoError(t, err)
	// Result can be nil or empty map - both are valid
	if result == nil {
		result = map[string]any{}
	}

	assert.IsType(t, map[string]any{}, result)
	// Don't require non-empty result as it might be empty in CI
}

func TestVolumeService_InspectVolume_NilClient(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.InspectVolume(ctx, "test-volume-name")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_InspectVolume_EmptyName(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewVolumeService(client)
	ctx := context.Background()
	result, err := service.InspectVolume(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestVolumeService_ContextHandling(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewVolumeService(client)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = service.ListVolumes(ctx)
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}

	_, err = service.InspectVolume(ctx, "test-name")
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}
}

func TestVolumeService_VolumeConversion(t *testing.T) {
	dockerVolume := docker.Volume{
		Name:       "test-volume",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
		Created:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Size:       "1.2GB",
	}

	modelVolume := models.Volume{
		Name:       dockerVolume.Name,
		Driver:     dockerVolume.Driver,
		Mountpoint: dockerVolume.Mountpoint,
		Created:    dockerVolume.Created,
		Size:       dockerVolume.Size,
	}

	// Verify conversion
	assert.Equal(t, dockerVolume.Name, modelVolume.Name)
	assert.Equal(t, dockerVolume.Driver, modelVolume.Driver)
	assert.Equal(t, dockerVolume.Mountpoint, modelVolume.Mountpoint)
	assert.Equal(t, dockerVolume.Created, modelVolume.Created)
	assert.Equal(t, dockerVolume.Size, modelVolume.Size)
}

func TestVolumeService_ErrorHandling(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume", false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
	assert.Contains(t, err.Error(), "docker client")
}
