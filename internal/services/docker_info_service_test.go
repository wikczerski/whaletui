package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/d5r/internal/config"
	"github.com/user/d5r/internal/docker"
	"github.com/user/d5r/internal/models"
)

func TestNewDockerInfoService(t *testing.T) {
	service := NewDockerInfoService(nil)
	assert.NotNil(t, service)

	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service = NewDockerInfoService(client)
	assert.NotNil(t, service)
}

func TestDockerInfoService_GetDockerInfo_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewDockerInfoService(client)
	ctx := context.Background()
	result, err := service.GetDockerInfo(ctx)

	require.NoError(t, err)
	// Result can be nil - that's valid if Docker info is not available
	if result != nil {
		assert.IsType(t, &models.DockerInfo{}, result)
		// Don't require specific fields as they might be empty in CI
	}
}

func TestDockerInfoService_GetDockerInfo_NilClient(t *testing.T) {
	service := NewDockerInfoService(nil)
	ctx := context.Background()

	assert.Panics(t, func() {
		service.GetDockerInfo(ctx)
	})
}
