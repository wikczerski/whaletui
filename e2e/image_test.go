package e2e

import (
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestImageList tests listing Docker images.
func TestImageList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Ensure test images exist
	dh.EnsureImage(framework.ImageFixtures.Alpine)
	dh.EnsureImage(framework.ImageFixtures.Busybox)

	// List images
	images, err := client.ImageList(ctx, image.ListOptions{})
	require.NoError(t, err, "Failed to list images")

	assert.GreaterOrEqual(t, len(images), 2, "Should have at least 2 images")
}

// TestImageInspect tests inspecting an image.
func TestImageInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Ensure image exists
	imageName := framework.ImageFixtures.Nginx
	dh.EnsureImage(imageName)

	// Inspect image
	inspect, err := client.ImageInspect(ctx, imageName)
	require.NoError(t, err, "Failed to inspect image")

	// Verify inspect data
	assert.NotNil(t, inspect, "Inspect data should not be nil")
	assert.NotEmpty(t, inspect.ID, "Image ID should not be empty")
	assert.Contains(t, inspect.RepoTags, imageName, "Image should have correct tag")
}

// TestImageDelete tests deleting an image.
func TestImageDelete(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Pull a test image
	testImage := "hello-world:latest"
	dh.EnsureImage(testImage)

	// Delete image
	_, err := client.ImageRemove(ctx, testImage, image.RemoveOptions{Force: false})
	require.NoError(t, err, "Failed to delete image")

	// Verify image is gone
	_, err = client.ImageInspect(ctx, testImage)
	assert.Error(t, err, "Image should be deleted")
}

// TestImageDeleteInUse tests deleting an image that's in use by a container.
func TestImageDeleteInUse(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create container using image
	imageName := framework.ImageFixtures.Alpine
	containerID := dh.CreateTestContainer(
		"e2e-test-image-in-use",
		imageName,
		nil,
		nil,
	)

	// Try to delete image (should fail without force)
	_, err := client.ImageRemove(ctx, imageName, image.RemoveOptions{Force: false})
	assert.Error(t, err, "Should fail to delete image in use")

	// Cleanup container
	dh.RemoveContainer(containerID, true)
}

// TestImageEmptyList tests handling of empty image list scenario.
func TestImageEmptyList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// List images
	images, err := client.ImageList(ctx, image.ListOptions{})
	require.NoError(t, err, "Failed to list images")

	// We can't guarantee empty list (base images may exist)
	// but we can verify the operation works
	assert.NotNil(t, images, "Image list should not be nil")
}
