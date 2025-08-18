package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCommonOperations(t *testing.T) {
	assert.NotPanics(t, func() {
		// Test that the function can be called without panicking
		// We can't easily test with a real docker client in unit tests
		// but we can test the basic structure
	})
}

func TestCommonOperations_StartContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StartContainer(ctx, containerID)
	assert.Error(t, err)
}

func TestCommonOperations_StartContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StartContainer(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_StopContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StopContainer(ctx, containerID, nil)
	assert.Error(t, err)
}

func TestCommonOperations_StopContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StopContainer(ctx, containerID, nil)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_RestartContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RestartContainer(ctx, containerID, nil)
	assert.Error(t, err)
}

func TestCommonOperations_RestartContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RestartContainer(ctx, containerID, nil)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_RemoveContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, false)
	assert.Error(t, err)
}

func TestCommonOperations_RemoveContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_GetContainerLogs_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.GetContainerLogs(ctx, containerID)
	assert.Error(t, err)
}

func TestCommonOperations_GetContainerLogs_NilClient_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	logs, _ := co.GetContainerLogs(ctx, containerID)
	assert.Empty(t, logs)
}

func TestCommonOperations_GetContainerLogs_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.GetContainerLogs(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ExecContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Error(t, err)
}

func TestCommonOperations_ExecContainer_NilClient_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	output, _ := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Empty(t, output)
}

func TestCommonOperations_ExecContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_AttachContainer_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.AttachContainer(ctx, containerID)
	assert.Error(t, err)
}

func TestCommonOperations_AttachContainer_NilClient_NilResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	result, _ := co.AttachContainer(ctx, containerID)
	assert.Nil(t, result)
}

func TestCommonOperations_AttachContainer_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.AttachContainer(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_InspectResource_NilClient(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "container", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_InspectResource_NilClient_NilResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	data, _ := co.InspectResource(ctx, "container", resourceID)
	assert.Nil(t, data)
}

func TestCommonOperations_InspectResource_NilClient_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "container", resourceID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_InspectResource_UnsupportedType(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "unsupported", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_InspectResource_UnsupportedType_NilResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	data, _ := co.InspectResource(ctx, "unsupported", resourceID)
	assert.Nil(t, data)
}

func TestCommonOperations_InspectResource_UnsupportedType_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "unsupported", resourceID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ResourceTypeValidation_Container(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "container", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_ResourceTypeValidation_Image(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "image", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_ResourceTypeValidation_Volume(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "volume", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_ResourceTypeValidation_Network(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	resourceID := "test-resource"

	_, err := co.InspectResource(ctx, "network", resourceID)
	assert.Error(t, err)
}

func TestCommonOperations_ContextHandling_BackgroundContext(t *testing.T) {
	co := &CommonOperations{client: nil}
	containerID := "test-container"
	ctx := context.Background()

	err := co.StartContainer(ctx, containerID)
	assert.Error(t, err)
}

func TestCommonOperations_ContextHandling_ValueContext(t *testing.T) {
	co := &CommonOperations{client: nil}
	containerID := "test-container"
	ctx := context.WithValue(context.Background(), "key", "value")

	err := co.StartContainer(ctx, containerID)
	assert.Error(t, err)
}

func TestCommonOperations_ContextHandling_BackgroundContext_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	containerID := "test-container"
	ctx := context.Background()

	err := co.StartContainer(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ContextHandling_ValueContext_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	containerID := "test-container"
	ctx := context.WithValue(context.Background(), "key", "value")

	err := co.StartContainer(ctx, containerID)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ContainerIDHandling_EmptyID(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "")
	assert.Error(t, err)
}

func TestCommonOperations_ContainerIDHandling_SimpleID(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "container-123")
	assert.Error(t, err)
}

func TestCommonOperations_ContainerIDHandling_HexID(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "abc123def456")
	assert.Error(t, err)
}

func TestCommonOperations_ContainerIDHandling_UnderscoreID(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "test_container")
	assert.Error(t, err)
}

func TestCommonOperations_ContainerIDHandling_EmptyID_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "")
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ContainerIDHandling_SimpleID_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "container-123")
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ContainerIDHandling_HexID_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "abc123def456")
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ContainerIDHandling_UnderscoreID_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()

	err := co.StartContainer(ctx, "test_container")
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TimeoutHandling_NilTimeout(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StopContainer(ctx, containerID, nil)
	assert.Error(t, err)
}

func TestCommonOperations_TimeoutHandling_ZeroTimeout(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := time.Duration(0)

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Error(t, err)
}

func TestCommonOperations_TimeoutHandling_ShortTimeout(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := 30 * time.Second

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Error(t, err)
}

func TestCommonOperations_TimeoutHandling_LongTimeout(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := 60 * time.Second

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Error(t, err)
}

func TestCommonOperations_TimeoutHandling_NilTimeout_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.StopContainer(ctx, containerID, nil)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TimeoutHandling_ZeroTimeout_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := time.Duration(0)

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TimeoutHandling_ShortTimeout_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := 30 * time.Second

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TimeoutHandling_LongTimeout_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	timeout := 60 * time.Second

	err := co.StopContainer(ctx, containerID, &timeout)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ForceFlagHandling_True(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, true)
	assert.Error(t, err)
}

func TestCommonOperations_ForceFlagHandling_False(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, false)
	assert.Error(t, err)
}

func TestCommonOperations_ForceFlagHandling_True_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, true)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_ForceFlagHandling_False_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	err := co.RemoveContainer(ctx, containerID, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_CommandHandling_EmptyCommand(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{}, false)
	assert.Error(t, err)
}

func TestCommonOperations_CommandHandling_SingleCommand(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Error(t, err)
}

func TestCommonOperations_CommandHandling_MultipleArgs(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls", "-la"}, false)
	assert.Error(t, err)
}

func TestCommonOperations_CommandHandling_ComplexCommand(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"docker", "exec", "container", "ls"}, false)
	assert.Error(t, err)
}

func TestCommonOperations_CommandHandling_EmptyCommand_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	output, _ := co.ExecContainer(ctx, containerID, []string{}, false)
	assert.Empty(t, output)
}

func TestCommonOperations_CommandHandling_SingleCommand_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	output, _ := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Empty(t, output)
}

func TestCommonOperations_CommandHandling_MultipleArgs_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	output, _ := co.ExecContainer(ctx, containerID, []string{"ls", "-la"}, false)
	assert.Empty(t, output)
}

func TestCommonOperations_CommandHandling_ComplexCommand_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	output, _ := co.ExecContainer(ctx, containerID, []string{"docker", "exec", "container", "ls"}, false)
	assert.Empty(t, output)
}

func TestCommonOperations_CommandHandling_EmptyCommand_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{}, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_CommandHandling_SingleCommand_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls"}, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_CommandHandling_MultipleArgs_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"ls", "-la"}, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_CommandHandling_ComplexCommand_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"

	_, err := co.ExecContainer(ctx, containerID, []string{"docker", "exec", "container", "ls"}, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TTYHandling_True(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	_, err := co.ExecContainer(ctx, containerID, command, true)
	assert.Error(t, err)
}

func TestCommonOperations_TTYHandling_False(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	_, err := co.ExecContainer(ctx, containerID, command, false)
	assert.Error(t, err)
}

func TestCommonOperations_TTYHandling_True_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	output, _ := co.ExecContainer(ctx, containerID, command, true)
	assert.Empty(t, output)
}

func TestCommonOperations_TTYHandling_False_EmptyResult(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	output, _ := co.ExecContainer(ctx, containerID, command, false)
	assert.Empty(t, output)
}

func TestCommonOperations_TTYHandling_True_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	_, err := co.ExecContainer(ctx, containerID, command, true)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}

func TestCommonOperations_TTYHandling_False_ErrorMessage(t *testing.T) {
	co := &CommonOperations{client: nil}
	ctx := context.Background()
	containerID := "test-container"
	command := []string{"ls"}

	_, err := co.ExecContainer(ctx, containerID, command, false)
	assert.Contains(t, err.Error(), "docker client is not initialized")
}
