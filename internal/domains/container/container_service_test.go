package container

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerService(t *testing.T) {
	service := NewContainerService(nil)
	assert.NotNil(t, service)

	service2 := NewContainerService(nil)
	assert.NotNil(t, service2)

	_ = service
}

func TestContainerService_ListContainers(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	containers, err := service.ListContainers(ctx)
	assert.Error(t, err)
	assert.Nil(t, containers)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_StartContainer(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	err := service.StartContainer(ctx, "test-container-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_StopContainer(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()
	timeout := 30 * time.Second

	err := service.StopContainer(ctx, "test-container-id", &timeout)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_StopContainer_NoTimeout(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	err := service.StopContainer(ctx, "test-container-id", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_RestartContainer(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()
	timeout := 30 * time.Second

	err := service.RestartContainer(ctx, "test-container-id", &timeout)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_RestartContainer_NoTimeout(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	err := service.RestartContainer(ctx, "test-container-id", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_RemoveContainer(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	err := service.RemoveContainer(ctx, "test-container-id", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")

	err = service.RemoveContainer(ctx, "test-container-id", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_GetContainerLogs(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	logs, err := service.GetContainerLogs(ctx, "test-container-id")
	assert.Error(t, err)
	assert.Empty(t, logs)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestContainerService_InspectContainer(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	info, err := service.InspectContainer(ctx, "test-container-id")
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestContainerService_ErrorMessages(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()
	containerID := "test-container-123"

	err := service.StartContainer(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")

	err = service.StopContainer(ctx, containerID, nil)
	assert.Contains(t, err.Error(), "docker client is not initialized")

	err = service.RestartContainer(ctx, containerID, nil)
	assert.Contains(t, err.Error(), "docker client is not initialized")

	err = service.RemoveContainer(ctx, containerID, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")

	logs, err := service.GetContainerLogs(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
	assert.Empty(t, logs)
}

func TestContainerService_ContextHandling(t *testing.T) {
	service := NewContainerService(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := service.StartContainer(ctx, "test-container-id")
	assert.Error(t, err)

	err = service.StopContainer(ctx, "test-container-id", nil)
	assert.Error(t, err)

	err = service.RestartContainer(ctx, "test-container-id", nil)
	assert.Error(t, err)

	err = service.RemoveContainer(ctx, "test-container-id", false)
	assert.Error(t, err)

	logs, err := service.GetContainerLogs(ctx, "test-container-id")
	assert.Error(t, err)
	assert.Empty(t, logs)

	info, err := service.InspectContainer(ctx, "test-container-id")
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestContainerService_EmptyContainerID(t *testing.T) {
	service := NewContainerService(nil)
	ctx := context.Background()

	err := service.StartContainer(ctx, "")
	assert.Error(t, err)

	err = service.StopContainer(ctx, "", nil)
	assert.Error(t, err)

	err = service.RestartContainer(ctx, "", nil)
	assert.Error(t, err)

	err = service.RemoveContainer(ctx, "", false)
	assert.Error(t, err)

	logs, err := service.GetContainerLogs(ctx, "")
	assert.Error(t, err)
	assert.Empty(t, logs)

	info, err := service.InspectContainer(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, info)
}
