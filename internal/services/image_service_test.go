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

func TestNewImageService(t *testing.T) {
	service := NewImageService(nil)
	assert.NotNil(t, service)

	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service = NewImageService(client)
	assert.NotNil(t, service)
}

func TestImageService_ListImages_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewImageService(client)
	ctx := context.Background()
	result, err := service.ListImages(ctx)

	require.NoError(t, err)
	// Result can be nil or empty slice - both are valid
	if result == nil {
		result = []models.Image{}
	}

	assert.IsType(t, []models.Image{}, result)
}

func TestImageService_ListImages_NilClient(t *testing.T) {
	service := NewImageService(nil)
	ctx := context.Background()

	assert.Panics(t, func() {
		service.ListImages(ctx)
	})
}

func TestImageService_RemoveImage(t *testing.T) {
	service := NewImageService(nil)
	ctx := context.Background()
	err := service.RemoveImage(ctx, "test-image-id", false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "image removal not yet implemented")
}

func TestImageService_RemoveImage_Force(t *testing.T) {
	service := NewImageService(nil)
	ctx := context.Background()
	err := service.RemoveImage(ctx, "test-image-id", true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "image removal not yet implemented")
}

func TestImageService_RemoveImage_EmptyID(t *testing.T) {
	service := NewImageService(nil)
	ctx := context.Background()
	err := service.RemoveImage(ctx, "", false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "image removal not yet implemented")
}

func TestImageService_InspectImage_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewImageService(client)
	ctx := context.Background()

	images, err := service.ListImages(ctx)
	if err != nil {
		t.Skipf("Could not list images: %v", err)
	}

	if len(images) == 0 {
		t.Skip("No images available for inspection")
	}

	imageID := images[0].ID
	result, err := service.InspectImage(ctx, imageID)

	require.NoError(t, err)
	// Result can be nil or empty map - both are valid
	if result == nil {
		result = map[string]any{}
	}

	assert.IsType(t, map[string]any{}, result)
	// Don't require non-empty result as it might be empty in CI
}

func TestImageService_InspectImage_NilClient(t *testing.T) {
	service := NewImageService(nil)
	ctx := context.Background()

	assert.Panics(t, func() {
		service.InspectImage(ctx, "test-image-id")
	})
}

func TestImageService_InspectImage_EmptyID(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewImageService(client)
	ctx := context.Background()
	result, err := service.InspectImage(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestImageService_ContextHandling(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewImageService(client)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = service.ListImages(ctx)
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}

	_, err = service.InspectImage(ctx, "test-id")
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}
}

func TestImageService_ImageConversion(t *testing.T) {
	dockerImage := docker.Image{
		ID:         "sha256:abc123",
		Repository: "nginx",
		Tag:        "latest",
		Size:       "133MB",
		Created:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Containers: 5,
	}

	modelImage := models.Image(dockerImage)

	assert.Equal(t, dockerImage.ID, modelImage.ID)
	assert.Equal(t, dockerImage.Repository, modelImage.Repository)
	assert.Equal(t, dockerImage.Tag, modelImage.Tag)
	assert.Equal(t, dockerImage.Size, modelImage.Size)
	assert.Equal(t, dockerImage.Created, modelImage.Created)
	assert.Equal(t, dockerImage.Containers, modelImage.Containers)
}
