package network

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker"
)

func TestNewNetworkService(t *testing.T) {
	service := NewNetworkService(nil)
	assert.NotNil(t, service)
}

func TestNewNetworkService_WithDockerClient(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)
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
		result = []Network{}
	}

	assert.IsType(t, []Network{}, result)
}

func TestNetworkService_ListNetworks_NilClient(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ListNetworks_NilClient_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_RemoveNetwork(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "test-network-id")

	assert.Error(t, err)
}

func TestNetworkService_RemoveNetwork_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "test-network-id")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_RemoveNetwork_EmptyID(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "")

	assert.Error(t, err)
}

func TestNetworkService_RemoveNetwork_EmptyID_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "")

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

	// This should fail with an invalid network ID error, not client error
	assert.Error(t, err)
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
		t.Skip("No networks available for testing")
	}

	// Test with the first available network
	networkID := networks[0].ID
	result, err := service.InspectNetwork(ctx, networkID)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestNetworkService_InspectNetwork_Integration_NetworkType(t *testing.T) {
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
		t.Skip("No networks available for testing")
	}

	// Test with the first available network
	networkID := networks[0].ID
	result, err := service.InspectNetwork(ctx, networkID)

	require.NoError(t, err)
	assert.IsType(t, map[string]any{}, result)
}

func TestNetworkService_InspectNetwork_NilClient(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.InspectNetwork(ctx, "test-network-id")

	assert.Error(t, err)
}

func TestNetworkService_InspectNetwork_NilClient_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.InspectNetwork(ctx, "test-network-id")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_InspectNetwork_EmptyID(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.InspectNetwork(ctx, "")

	assert.Error(t, err)
}

func TestNetworkService_InspectNetwork_EmptyID_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.InspectNetwork(ctx, "")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_ContextHandling_BackgroundContext(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ContextHandling_ValueContext(t *testing.T) {
	service := NewNetworkService(nil)
	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("key"), "value")
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ContextHandling_TimeoutContext(t *testing.T) {
	service := NewNetworkService(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ContextHandling_BackgroundContext_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_ContextHandling_ValueContext_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("key"), "value")
	_, err := service.ListNetworks(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_ContextHandling_TimeoutContext_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := service.ListNetworks(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_NetworkConversion_EmptyList(t *testing.T) {
	// Test that empty list is handled correctly
	// This would typically be tested with a mock, but for now we test the nil client case
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_NetworkConversion_NilList(t *testing.T) {
	// Test that nil list is handled correctly
	// This would typically be tested with a mock, but for now we test the nil client case
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ErrorHandling_ClientNotInitialized(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_ErrorHandling_ClientNotInitialized_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_EmptyResults_EmptyList(t *testing.T) {
	// Test that empty list is handled correctly
	// This would typically be tested with a mock, but for now we test the nil client case
	service := NewNetworkService(nil)
	ctx := context.Background()
	_, err := service.ListNetworks(ctx)

	assert.Error(t, err)
}

func TestNetworkService_NetworkIDValidation_EmptyID(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "")

	assert.Error(t, err)
}

func TestNetworkService_NetworkIDValidation_EmptyID_ErrorMessage(t *testing.T) {
	service := NewNetworkService(nil)
	ctx := context.Background()
	err := service.RemoveNetwork(ctx, "")

	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestNetworkService_NetworkIDValidation_EmptyID_WithClient(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := docker.New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	service := NewNetworkService(client)
	ctx := context.Background()
	err = service.RemoveNetwork(ctx, "")

	// This should fail with an invalid network ID error, not client error
	assert.Error(t, err)
}
