package docker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		expectError bool
	}{
		{
			name:        "NilConfig",
			cfg:         nil,
			expectError: true,
		},
		{
			name: "ValidConfig",
			cfg: &config.Config{
				DockerHost: "unix:///var/run/docker.sock",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.cfg)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				if err != nil {
					t.Skipf("Docker not available: %v", err)
				}
				assert.NoError(t, err)
				assert.NotNil(t, client)
				if client != nil {
					defer func() {
						if err := client.Close(); err != nil {
							t.Logf("Warning: failed to close client: %v", err)
						}
					}()
				}
			}
		})
	}
}

func TestClient_GetInfo(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()
	info, err := client.GetInfo(ctx)
	if err != nil {
		t.Skipf("Docker info not available: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, info)

	id, hasID := info["ID"]
	assert.True(t, hasID, "Docker info should contain ID field")
	if hasID {
		assert.NotEmpty(t, id)
	}

	driver, hasDriver := info["Driver"]
	assert.True(t, hasDriver, "Docker info should contain Driver field")
	if hasDriver {
		assert.NotEmpty(t, driver)
	}
}

func TestClient_InspectContainer(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.InspectContainer(ctx, "invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "container inspect failed")

	_, err = client.InspectContainer(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "container inspect failed")
}

func TestClient_GetContainerLogs(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.GetContainerLogs(ctx, "invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "container logs failed")

	_, err = client.GetContainerLogs(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "container logs failed")
}

func TestClient_InspectImage(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.InspectImage(ctx, "invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image inspect failed")

	_, err = client.InspectImage(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image inspect failed")
}

func TestClient_InspectVolume(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.InspectVolume(ctx, "invalid-volume")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "volume inspect failed")

	_, err = client.InspectVolume(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "volume inspect failed")
}

func TestClient_InspectNetwork(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.InspectNetwork(ctx, "invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network inspect failed")

	_, err = client.InspectNetwork(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network inspect failed")
}

func TestClient_Close(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	assert.NotPanics(t, func() {
		err := client.Close()
		assert.NoError(t, err)
	})

	assert.NotPanics(t, func() {
		err := client.Close()
		assert.NoError(t, err)
	})
}

func TestClient_ListContainers(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	containers, err := client.ListContainers(ctx, true)
	if err != nil {
		t.Skipf("Docker containers not available: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, containers)

	containers, err = client.ListContainers(ctx, false)
	assert.NoError(t, err)
	assert.NotNil(t, containers)
}

func TestClient_GetContainerStats(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	_, err = client.GetContainerStats(ctx, "invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get container stats")

	_, err = client.GetContainerStats(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get container stats")
}

func TestClient_ListImages(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	images, err := client.ListImages(ctx)
	if err != nil {
		t.Skipf("Docker images not available: %v", err)
	}

	assert.NoError(t, err)
	if images == nil {
		images = []Image{}
	}
	assert.IsType(t, []Image{}, images)
}

func TestClient_ListVolumes(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	volumes, err := client.ListVolumes(ctx)
	if err != nil {
		t.Skipf("Docker volumes not available: %v", err)
	}

	assert.NoError(t, err)
	if volumes == nil {
		volumes = []Volume{}
	}
	assert.IsType(t, []Volume{}, volumes)
}

func TestClient_ListNetworks(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	networks, err := client.ListNetworks(ctx)
	if err != nil {
		t.Skipf("Docker networks not available: %v", err)
	}

	assert.NoError(t, err)
	if networks == nil {
		networks = []Network{}
	}
	assert.IsType(t, []Network{}, networks)
}

func TestClient_ListNetworks_CreatedField(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	networks, err := client.ListNetworks(ctx)
	if err != nil {
		t.Skipf("Docker networks not available: %v", err)
	}

	for _, net := range networks {
		assert.NotZero(t, net.Created, "Network %s should have a creation time", net.Name)
		assert.True(t, net.Created.Before(time.Now().Add(24*time.Hour)), "Network %s creation time should not be in the future", net.Name)
		assert.True(t, net.Created.After(time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)), "Network %s creation time should be reasonable", net.Name)
	}
}

func TestClient_StartContainer(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.StartContainer(ctx, "invalid-id")
	assert.Error(t, err)

	err = client.StartContainer(ctx, "")
	assert.Error(t, err)
}

func TestClient_StopContainer(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.StopContainer(ctx, "invalid-id", nil)
	assert.Error(t, err)

	err = client.StopContainer(ctx, "", nil)
	assert.Error(t, err)

	timeout := 10 * time.Second
	err = client.StopContainer(ctx, "invalid-id", &timeout)
	assert.Error(t, err)
}

func TestClient_RestartContainer(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.RestartContainer(ctx, "invalid-id", nil)
	assert.Error(t, err)

	err = client.RestartContainer(ctx, "", nil)
	assert.Error(t, err)

	timeout := 10 * time.Second
	err = client.RestartContainer(ctx, "invalid-id", &timeout)
	assert.Error(t, err)
}

func TestClient_RemoveContainer(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.RemoveContainer(ctx, "invalid-id", false)
	assert.Error(t, err)

	err = client.RemoveContainer(ctx, "", false)
	assert.Error(t, err)

	err = client.RemoveContainer(ctx, "invalid-id", true)
	assert.Error(t, err)
}

func TestClient_RemoveImage(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.RemoveImage(ctx, "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image ID cannot be empty")

	err = client.RemoveImage(ctx, "invalid-image-id", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove image")
}

func TestClient_RemoveNetwork(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning: failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	err = client.RemoveNetwork(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network ID cannot be empty")

	err = client.RemoveNetwork(ctx, "invalid-network-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove network")
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"ZeroBytes", 0, "0.00 B"},
		{"Bytes", 123, "123.00 B"},
		{"Kilobytes", 1024, "1.00 KB"},
		{"KilobytesDecimal", 1536, "1.50 KB"},
		{"Megabytes", 1048576, "1.00 MB"},
		{"MegabytesDecimal", 1572864, "1.50 MB"},
		{"Gigabytes", 1073741824, "1.00 GB"},
		{"GigabytesDecimal", 1610612736, "1.50 GB"},
		{"Terabytes", 1099511627776, "1.00 TB"},
		{"TerabytesDecimal", 1649267441664, "1.50 TB"},
		{"LargeValue", 9223372036854775807, "8388608.00 TB"}, // Max int64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSize(tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainer_Fields(t *testing.T) {
	container := Container{
		ID:      "abc123",
		Name:    "test-container",
		Image:   "test-image:latest",
		Status:  "running",
		State:   "running",
		Created: time.Now(),
		Ports:   "8080:80/tcp",
		Size:    "1.5 MB",
	}

	assert.Equal(t, "abc123", container.ID)
	assert.Equal(t, "test-container", container.Name)
	assert.Equal(t, "test-image:latest", container.Image)
	assert.Equal(t, "running", container.Status)
	assert.Equal(t, "running", container.State)
	assert.NotZero(t, container.Created)
	assert.Equal(t, "8080:80/tcp", container.Ports)
	assert.Equal(t, "1.5 MB", container.Size)
}

func TestImage_Fields(t *testing.T) {
	image := Image{
		ID:         "def456",
		Repository: "test-repo",
		Tag:        "latest",
		Size:       "2.1 GB",
		Created:    time.Now(),
		Containers: 3,
	}

	assert.Equal(t, "def456", image.ID)
	assert.Equal(t, "test-repo", image.Repository)
	assert.Equal(t, "latest", image.Tag)
	assert.Equal(t, "2.1 GB", image.Size)
	assert.NotZero(t, image.Created)
	assert.Equal(t, 3, image.Containers)
}

func TestVolume_Fields(t *testing.T) {
	volume := Volume{
		Name:       "test-volume",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/test-volume",
		Created:    time.Now(),
		Size:       "100 MB",
	}

	assert.Equal(t, "test-volume", volume.Name)
	assert.Equal(t, "local", volume.Driver)
	assert.Equal(t, "/var/lib/docker/volumes/test-volume", volume.Mountpoint)
	assert.NotZero(t, volume.Created)
	assert.Equal(t, "100 MB", volume.Size)
}

func TestNetwork_Fields(t *testing.T) {
	createdTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	network := Network{
		ID:         "ghi789",
		Name:       "test-network",
		Driver:     "bridge",
		Scope:      "local",
		Created:    createdTime,
		Containers: 2,
	}

	assert.Equal(t, "ghi789", network.ID)
	assert.Equal(t, "test-network", network.Name)
	assert.Equal(t, "bridge", network.Driver)
	assert.Equal(t, "local", network.Scope)
	assert.Equal(t, createdTime, network.Created)
	assert.Equal(t, 2, network.Containers)
}

func TestClient_ContextHandling(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.GetInfo(ctx)
	assert.Error(t, err)
}

func TestClient_ConcurrentAccess(t *testing.T) {
	cfg := &config.Config{DockerHost: "unix:///var/run/docker.sock"}
	client, err := New(cfg)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_, _ = client.ListContainers(ctx, true)
			_, _ = client.ListImages(ctx)
			_, _ = client.ListVolumes(ctx)
			_, _ = client.ListNetworks(ctx)
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestExtractHostFromURL(t *testing.T) {
	tests := []struct {
		name        string
		hostURL     string
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "ValidTCPURL",
			hostURL:     "tcp://192.168.1.100:2375",
			expected:    "192.168.1.100",
			expectError: false,
		},
		{
			name:        "ValidTCPURLNoPort",
			hostURL:     "tcp://192.168.1.100",
			expected:    "192.168.1.100",
			expectError: false,
		},
		{
			name:        "ValidHostnameURL",
			hostURL:     "tcp://myserver.example.com:2375",
			expected:    "myserver.example.com",
			expectError: false,
		},
		{
			name:        "ValidHostnameURLNoPort",
			hostURL:     "tcp://myserver.example.com",
			expected:    "myserver.example.com",
			expectError: false,
		},
		{
			name:        "NoScheme",
			hostURL:     "192.168.1.100:2375",
			expected:    "192.168.1.100",
			expectError: false,
		},
		{
			name:        "NoSchemeNoPort",
			hostURL:     "192.168.1.100",
			expected:    "192.168.1.100",
			expectError: false,
		},
		{
			name:        "InvalidURLFormat",
			hostURL:     "tcp://192.168.1.100:2375:extra",
			expectError: true,
			errorMsg:    "invalid host:port format",
		},
		{
			name:        "EmptyHostname",
			hostURL:     "tcp://:2375",
			expectError: true,
			errorMsg:    "hostname cannot be empty",
		},
		{
			name:        "HostnameStartingWithDot",
			hostURL:     "tcp://.example.com:2375",
			expectError: true,
			errorMsg:    "hostname '.example.com' cannot start or end with a dot",
		},
		{
			name:        "HostnameEndingWithDot",
			hostURL:     "tcp://example.com.:2375",
			expectError: true,
			errorMsg:    "hostname 'example.com.' cannot start or end with a dot",
		},
		{
			name:        "HostnameWithConsecutiveDots",
			hostURL:     "tcp://example..com:2375",
			expectError: true,
			errorMsg:    "hostname 'example..com' cannot contain consecutive dots",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractHostFromURL(tt.hostURL)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
