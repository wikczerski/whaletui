package e2e

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestVolumeList tests listing Docker volumes.
func TestVolumeList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create test volumes
	vol1 := dh.CreateTestVolume(framework.VolumeFixtures.Data1.Name)
	vol2 := dh.CreateTestVolume(framework.VolumeFixtures.Data2.Name)

	// List volumes
	volumes, err := client.VolumeList(ctx, volume.ListOptions{})
	require.NoError(t, err, "Failed to list volumes")

	// Verify volumes exist in list
	volumeNames := make([]string, 0)
	for _, v := range volumes.Volumes {
		volumeNames = append(volumeNames, v.Name)
	}

	assert.Contains(t, volumeNames, vol1, "Volume 1 should be in list")
	assert.Contains(t, volumeNames, vol2, "Volume 2 should be in list")
}

// TestVolumeCreate tests creating a volume.
func TestVolumeCreate(t *testing.T) {
	fw := framework.NewTestFramework(t)
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create volume
	volumeName := "e2e-test-volume-create"
	vol, err := client.VolumeCreate(ctx, volume.CreateOptions{Name: volumeName})
	require.NoError(t, err, "Failed to create volume")
	fw.RegisterTestVolume(vol.Name)

	// Verify volume exists
	inspect, err := client.VolumeInspect(ctx, volumeName)
	require.NoError(t, err, "Failed to inspect volume")
	assert.Equal(t, volumeName, inspect.Name, "Volume name should match")
}

// TestVolumeDelete tests deleting a volume.
func TestVolumeDelete(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create volume
	volumeName := dh.CreateTestVolume("e2e-test-volume-delete")

	// Delete volume
	dh.RemoveVolume(volumeName, false)

	// Verify volume is gone
	ctx := fw.GetContext()
	client := fw.GetDockerClient()
	_, err := client.VolumeInspect(ctx, volumeName)
	assert.Error(t, err, "Volume should be deleted")
}

// TestVolumeDeleteInUse tests deleting a volume that's in use by a container.
func TestVolumeDeleteInUse(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create volume
	volumeName := dh.CreateTestVolume("e2e-test-volume-in-use")

	// Create container using volume
	containerID := dh.CreateTestContainer(
		"e2e-test-volume-container",
		framework.ImageFixtures.Alpine,
		&container.Config{
			Image: framework.ImageFixtures.Alpine,
			Cmd:   []string{"sleep", "3600"},
		},
		&container.HostConfig{
			Binds: []string{volumeName + ":/data"},
		},
	)
	dh.StartContainer(containerID)

	// Try to delete volume (should fail)
	err := client.VolumeRemove(ctx, volumeName, false)
	assert.Error(t, err, "Should fail to delete volume in use")

	// Cleanup
	dh.StopContainer(containerID)
	dh.RemoveContainer(containerID, true)
	dh.RemoveVolume(volumeName, true)
}

// TestVolumeInspect tests inspecting a volume.
func TestVolumeInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create volume
	volumeName := dh.CreateTestVolume(framework.VolumeFixtures.Data3.Name)

	// Inspect volume
	inspect, err := client.VolumeInspect(ctx, volumeName)
	require.NoError(t, err, "Failed to inspect volume")

	// Verify inspect data
	assert.Equal(t, volumeName, inspect.Name, "Volume name should match")
	assert.Equal(t, "local", inspect.Driver, "Driver should be local")
	assert.NotEmpty(t, inspect.Mountpoint, "Mountpoint should not be empty")
}

// TestVolumeEmptyList tests handling of empty volume list.
func TestVolumeEmptyList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// List volumes
	volumes, err := client.VolumeList(ctx, volume.ListOptions{})
	require.NoError(t, err, "Failed to list volumes")

	// Verify operation works
	assert.NotNil(t, volumes, "Volume list should not be nil")
}
