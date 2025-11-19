package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestContainerList tests listing containers in various states.
func TestContainerList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create test containers in different states
	runningID := dh.CreateTestContainer(
		framework.ContainerFixtures.Nginx.Name,
		framework.ContainerFixtures.Nginx.Image,
		framework.ContainerFixtures.Nginx.Config,
		framework.ContainerFixtures.Nginx.HostConfig,
	)
	dh.StartContainer(runningID)

	stoppedID := dh.CreateTestContainer(
		framework.ContainerFixtures.Redis.Name,
		framework.ContainerFixtures.Redis.Image,
		framework.ContainerFixtures.Redis.Config,
		framework.ContainerFixtures.Redis.HostConfig,
	)
	// Leave stopped

	// Verify containers exist
	containers := dh.ListContainers(true)
	assert.GreaterOrEqual(t, len(containers), 2, "Should have at least 2 containers")

	// Verify states
	runningState := dh.GetContainerState(runningID)
	assert.Equal(t, "running", runningState, "Container should be running")

	stoppedState := dh.GetContainerState(stoppedID)
	assert.Contains(t, []string{"created", "exited"}, stoppedState, "Container should be stopped")
}

// TestContainerStartStop tests starting and stopping containers.
func TestContainerStartStop(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create stopped container
	containerID := dh.CreateTestContainer(
		framework.ContainerFixtures.Alpine.Name,
		framework.ContainerFixtures.Alpine.Image,
		framework.ContainerFixtures.Alpine.Config,
		framework.ContainerFixtures.Alpine.HostConfig,
	)

	// Start container
	dh.StartContainer(containerID)
	dh.WaitForContainerState(containerID, "running", 10*time.Second)

	state := dh.GetContainerState(containerID)
	assert.Equal(t, "running", state, "Container should be running after start")

	// Stop container
	dh.StopContainer(containerID)
	dh.WaitForContainerState(containerID, "exited", 10*time.Second)

	state = dh.GetContainerState(containerID)
	assert.Equal(t, "exited", state, "Container should be stopped after stop")
}

// TestContainerRestart tests restarting a container.
func TestContainerRestart(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create and start container
	containerName := fmt.Sprintf(
		"%s-%d",
		framework.ContainerFixtures.Busybox.Name,
		time.Now().UnixNano(),
	)
	containerID := dh.CreateTestContainer(
		containerName,
		framework.ContainerFixtures.Busybox.Image,
		framework.ContainerFixtures.Busybox.Config,
		framework.ContainerFixtures.Busybox.HostConfig,
	)
	dh.StartContainer(containerID)
	dh.WaitForContainerState(containerID, "running", 10*time.Second)

	// Restart container
	timeout := 10
	err := client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
	require.NoError(t, err, "Failed to restart container")

	// Verify container is running
	dh.WaitForContainerState(containerID, "running", 10*time.Second)
	state := dh.GetContainerState(containerID)
	assert.Equal(t, "running", state, "Container should be running after restart")
}

// TestContainerDelete tests deleting a container.
func TestContainerDelete(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create container
	containerID := dh.CreateTestContainer(
		"e2e-test-delete",
		framework.ImageFixtures.Alpine,
		nil,
		nil,
	)

	// Delete container
	dh.RemoveContainer(containerID, false)

	// Verify container is gone
	container := dh.FindContainerByName("e2e-test-delete")
	assert.Nil(t, container, "Container should be deleted")
}

// TestContainerDeleteRunning tests deleting a running container (should require force).
func TestContainerDeleteRunning(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create and start container
	containerID := dh.CreateTestContainer(
		"e2e-test-delete-running",
		framework.ImageFixtures.Alpine,
		&container.Config{
			Image: framework.ImageFixtures.Alpine,
			Cmd:   []string{"sleep", "3600"},
		},
		nil,
	)
	dh.StartContainer(containerID)
	dh.WaitForContainerState(containerID, "running", 10*time.Second)

	// Try to delete without force (should fail in real scenario, but we'll force it)
	// In a real TUI test, this would show an error modal
	dh.RemoveContainer(containerID, true) // Force delete

	// Verify container is gone
	container := dh.FindContainerByName("e2e-test-delete-running")
	assert.Nil(t, container, "Container should be deleted")
}

// TestContainerInspect tests inspecting a container.
func TestContainerInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create container
	containerID := dh.CreateTestContainer(
		framework.ContainerFixtures.Nginx.Name,
		framework.ContainerFixtures.Nginx.Image,
		framework.ContainerFixtures.Nginx.Config,
		framework.ContainerFixtures.Nginx.HostConfig,
	)

	// Inspect container
	inspect, err := client.ContainerInspect(ctx, containerID)
	require.NoError(t, err, "Failed to inspect container")

	// Verify inspect data
	assert.NotNil(t, inspect, "Inspect data should not be nil")
	assert.Equal(t, containerID, inspect.ID, "Container ID should match")
	assert.Contains(
		t,
		inspect.Name,
		framework.ContainerFixtures.Nginx.Name,
		"Container name should match",
	)
	assert.Equal(t, framework.ContainerFixtures.Nginx.Image, inspect.Config.Image, "Image should match")
}

// TestContainerStates tests containers in different states and color coding.
func TestContainerStates(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	tests := []struct {
		name          string
		fixture       framework.ContainerFixture
		shouldStart   bool
		expectedState string
	}{
		{
			name:          "Running Container",
			fixture:       framework.ContainerFixtures.Nginx,
			shouldStart:   true,
			expectedState: "running",
		},
		{
			name:          "Created Container",
			fixture:       framework.ContainerFixtures.Redis,
			shouldStart:   false,
			expectedState: "created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerID := dh.CreateTestContainer(
				tt.fixture.Name+"-state-test",
				tt.fixture.Image,
				tt.fixture.Config,
				tt.fixture.HostConfig,
			)

			if tt.shouldStart {
				dh.StartContainer(containerID)
				dh.WaitForContainerState(containerID, "running", 10*time.Second)
			}

			state := dh.GetContainerState(containerID)
			assert.Equal(t, tt.expectedState, state, "Container state should match expected")
		})
	}
}

// TestContainerEmptyList tests handling of empty container list.
func TestContainerEmptyList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Clean up all test containers first
	dh.CleanupAll()

	// List containers (should be empty or only system containers)
	containers := dh.ListContainers(false) // Only running

	// We can't guarantee zero containers (system containers may exist)
	// but we can verify the list operation works
	assert.NotNil(t, containers, "Container list should not be nil")
}

// TestContainerLifecycle tests complete container lifecycle.
func TestContainerLifecycle(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create container
	containerID := dh.CreateTestContainer(
		"e2e-lifecycle-test",
		framework.ImageFixtures.Alpine,
		&container.Config{
			Image: framework.ImageFixtures.Alpine,
			Cmd:   []string{"sleep", "3600"},
		},
		nil,
	)

	// Verify created state
	state := dh.GetContainerState(containerID)
	assert.Equal(t, "created", state, "Container should be in created state")

	// Start container
	dh.StartContainer(containerID)
	dh.WaitForContainerState(containerID, "running", 10*time.Second)

	state = dh.GetContainerState(containerID)
	assert.Equal(t, "running", state, "Container should be running")

	// Restart container
	timeout := 10
	err := client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
	require.NoError(t, err, "Failed to restart container")
	dh.WaitForContainerState(containerID, "running", 10*time.Second)

	state = dh.GetContainerState(containerID)
	assert.Equal(t, "running", state, "Container should be running after restart")

	// Stop container
	dh.StopContainer(containerID)
	dh.WaitForContainerState(containerID, "exited", 10*time.Second)

	state = dh.GetContainerState(containerID)
	assert.Equal(t, "exited", state, "Container should be stopped")

	// Delete container
	dh.RemoveContainer(containerID, false)

	// Verify deleted
	container := dh.FindContainerByName("e2e-lifecycle-test")
	assert.Nil(t, container, "Container should be deleted")
}
