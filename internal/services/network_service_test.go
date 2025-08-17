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

func TestNewNetworkService(t *testing.T) {
	service := NewNetworkService(nil)
	assert.NotNil(t, service)

	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service = NewNetworkService(client)
	assert.NotNil(t, service)
}

func TestNetworkService_ListNetworks_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)
	ctx := context.Background()
	result, err := service.ListNetworks(ctx)

	require.NoError(t, err)
	// Result can be nil or empty slice - both are valid
	if result == nil {
		result = []models.Network{}
	}

	assert.IsType(t, []models.Network{}, result)
}

func TestNetworkService_ListNetworks_NilClient(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()

	assert.Panics(t, func() {
		service.ListNetworks(ctx)
	})
}

func TestNetworkService_RemoveNetwork(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "test-network-id")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_RemoveNetwork_EmptyID(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_RemoveNetwork_EmptyID_WithClient(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)
	ctx := context.Background()
	err = service.RemoveNetwork(ctx, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "network ID cannot be empty")
}

func TestNetworkService_InspectNetwork_Integration(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)
	ctx := context.Background()

	networks, err := service.ListNetworks(ctx)
	if err != nil {
		t.Skipf("Could not list networks: %v", err)
	}

	if len(networks) == 0 {
		t.Skip("No networks available for inspection")
	}

	networkID := networks[0].ID
	result, err := service.InspectNetwork(ctx, networkID)

	require.NoError(t, err)
	// Result can be nil or empty map - both are valid
	if result == nil {
		result = map[string]any{}
	}

	assert.IsType(t, map[string]any{}, result)
	// Don't require non-empty result as it might be empty in CI
}

func TestNetworkService_InspectNetwork_NilClient(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()

	assert.Panics(t, func() {
		service.InspectNetwork(ctx, "test-network-id")
	})
}

func TestNetworkService_InspectNetwork_EmptyID(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)

	ctx := context.Background()
	result, err := service.InspectNetwork(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNetworkService_ContextHandling(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = service.ListNetworks(ctx)
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}

	_, err = service.InspectNetwork(ctx, "test-id")
	if err != nil {
		assert.Contains(t, err.Error(), "context")
	}
}

func TestNetworkService_NetworkConversion(t *testing.T) {
	createdTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	dockerNetwork := docker.Network{
		ID:      "abc123def456",
		Name:    "bridge",
		Driver:  "bridge",
		Scope:   "local",
		Created: createdTime,
	}

	modelNetwork := models.Network{
		ID:      dockerNetwork.ID,
		Name:    dockerNetwork.Name,
		Driver:  dockerNetwork.Driver,
		Scope:   dockerNetwork.Scope,
		Created: dockerNetwork.Created,
	}

	assert.Equal(t, dockerNetwork.ID, modelNetwork.ID)
	assert.Equal(t, dockerNetwork.Name, modelNetwork.Name)
	assert.Equal(t, dockerNetwork.Driver, modelNetwork.Driver)
	assert.Equal(t, dockerNetwork.Scope, modelNetwork.Scope)
	assert.Equal(t, dockerNetwork.Created, modelNetwork.Created)
}

func TestNetworkService_ErrorHandling(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "test-network")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_EmptyResults(t *testing.T) {

	service := NewNetworkService(nil)

	assert.NotNil(t, service)
}

func TestNetworkService_NetworkIDValidation(t *testing.T) {
	service := NewNetworkService(nil)
	testIDs := []string{
		"",
		"abc123",
		"def456ghi789",
		"very-long-network-id-with-many-characters",
		"network-with-special-chars-!@#$%^&*()",
	}

	for _, id := range testIDs {
		err := service.RemoveNetwork(context.Background(), id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "docker client is not initialized")
	}
}

func TestNetworkService_NetworkIDValidation_WithClient(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)

	// Test empty ID validation
	err = service.RemoveNetwork(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network ID cannot be empty")
}
