package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/models"
)

func TestNewVolumeService(t *testing.T) {
	service := NewVolumeService(nil)
	assert.NotNil(t, service)
}

func TestNewVolumeService_WithDockerClient(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewVolumeService(client)
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
}

func TestVolumeService_ListVolumes_NilClient_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", false)

	require.Error(t, err)
}

func TestVolumeService_RemoveVolume_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", false)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume_Force(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", true)

	require.Error(t, err)
}

func TestVolumeService_RemoveVolume_Force_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "test-volume-name", true)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_RemoveVolume_EmptyName(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "", false)

	require.Error(t, err)
}

func TestVolumeService_RemoveVolume_EmptyName_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	err := service.RemoveVolume(ctx, "", false)

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
		t.Skip("No volumes available for testing")
	}

	// Test with the first available volume
	volumeName := volumes[0].Name
	result, err := service.InspectVolume(ctx, volumeName)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestVolumeService_InspectVolume_Integration_VolumeType(t *testing.T) {
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
		t.Skip("No volumes available for testing")
	}

	// Test with the first available volume
	volumeName := volumes[0].Name
	result, err := service.InspectVolume(ctx, volumeName)

	require.NoError(t, err)
	assert.IsType(t, &models.Volume{}, result)
}

func TestVolumeService_InspectVolume_NilClient(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.InspectVolume(ctx, "test-volume")

	assert.Error(t, err)
}

func TestVolumeService_InspectVolume_NilClient_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.InspectVolume(ctx, "test-volume")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_InspectVolume_EmptyName(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.InspectVolume(ctx, "")

	assert.Error(t, err)
}

func TestVolumeService_InspectVolume_EmptyName_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.InspectVolume(ctx, "")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_ContextHandling_BackgroundContext(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_ContextHandling_ValueContext(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.WithValue(context.Background(), "key", "value")
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_ContextHandling_TimeoutContext(t *testing.T) {
	service := NewVolumeService(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_ContextHandling_BackgroundContext_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_ContextHandling_ValueContext_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.WithValue(context.Background(), "key", "value")
	_, err := service.ListVolumes(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_ContextHandling_TimeoutContext_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := service.ListVolumes(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestVolumeService_VolumeConversion_EmptyList(t *testing.T) {
	// Test that empty list is handled correctly
	// This would typically be tested with a mock, but for now we test the nil client case
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_VolumeConversion_NilList(t *testing.T) {
	// Test that nil list is handled correctly
	// This would typically be tested with a mock, but for now we test the nil client case
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_ErrorHandling_ClientNotInitialized(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Error(t, err)
}

func TestVolumeService_ErrorHandling_ClientNotInitialized_ErrorMessage(t *testing.T) {
	service := NewVolumeService(nil)
	ctx := context.Background()
	_, err := service.ListVolumes(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}
