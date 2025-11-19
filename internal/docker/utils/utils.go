package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/docker/types"
)

// MarshalToMap converts any value to a map[string]any using JSON marshaling
func MarshalToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return result, nil
}

// FormatSize formats a size in bytes to a human-readable string
func FormatSize(size int64) string {
	sizeInfo := getSizeInfo(size)
	return fmt.Sprintf("%.2f %s", sizeInfo.value, sizeInfo.unit)
}

// sizeInfo holds size formatting information
type sizeInfo struct {
	unit  string
	value float64
}

// getSizeInfo determines the appropriate unit and value for size formatting
func getSizeInfo(size int64) sizeInfo {
	const (
		KB int64 = 1024
		MB int64 = KB * 1024
		GB int64 = MB * 1024
		TB int64 = GB * 1024
	)

	switch {
	case size >= TB:
		return sizeInfo{unit: "TB", value: float64(size) / float64(TB)}
	case size >= GB:
		return sizeInfo{unit: "GB", value: float64(size) / float64(GB)}
	case size >= MB:
		return sizeInfo{unit: "MB", value: float64(size) / float64(MB)}
	case size >= KB:
		return sizeInfo{unit: "KB", value: float64(size) / float64(KB)}
	default:
		return sizeInfo{unit: "B", value: float64(size)}
	}
}

// DetectWindowsDockerHost attempts to find the correct Docker host on Windows
func DetectWindowsDockerHost(log *slog.Logger) (string, error) {
	if runtime.GOOS != "windows" {
		return "", errors.New("not on Windows")
	}

	possibleHosts := getWindowsDockerHosts()
	for _, host := range possibleHosts {
		if isHostWorking(host, log) {
			return host, nil
		}
	}

	return "", errors.New("no working Docker host found")
}

// getWindowsDockerHosts returns the list of possible Windows Docker host paths
func getWindowsDockerHosts() []string {
	return []string{
		"npipe:////./pipe/dockerDesktopLinuxEngine", // Linux containers
		"npipe:////./pipe/docker_engine",            // Windows containers
		"npipe:////./pipe/dockerDesktopEngine",      // Legacy Docker Desktop
	}
}

// isHostWorking tests if a Docker host is working
func isHostWorking(host string, log *slog.Logger) bool {
	cli, err := createTestClient(host)
	if err != nil {
		return false
	}
	defer closeClientSafely(cli, log)

	return testClientConnection(cli)
}

// createTestClient creates a Docker client for testing host connectivity
func createTestClient(host string) (*client.Client, error) {
	opts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(5 * time.Second),
	}

	return client.NewClientWithOpts(opts...)
}

// testClientConnection tests if the Docker client can connect successfully
func testClientConnection(cli *client.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cli.Ping(ctx)
	return err == nil
}

// closeClientSafely closes a Docker client and logs any errors
func closeClientSafely(cli *client.Client, log *slog.Logger) {
	if err := cli.Close(); err != nil {
		log.Warn("Failed to close Docker client during host detection", "error", err)
	}
}

// FormatContainerPorts formats container ports into a readable string
func FormatContainerPorts(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}

	var formattedPorts []string
	for _, p := range ports {
		if p.PublicPort > 0 {
			formattedPorts = append(
				formattedPorts,
				fmt.Sprintf("%d:%d/%s", p.PublicPort, p.PrivatePort, p.Type),
			)
		} else if p.PrivatePort > 0 {
			formattedPorts = append(formattedPorts, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}

	return strings.Join(formattedPorts, ", ")
}

// SortContainersByCreationTime sorts containers by creation time (newest first)
func SortContainersByCreationTime(containers []types.Container) {
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Created.After(containers[j].Created)
	})
}

// ParseImageRepository parses repository and tag from image repoTags
func ParseImageRepository(repoTags []string) (repository, tag string) {
	if len(repoTags) == 0 || repoTags[0] == "<none>:<none>" {
		return "<none>", "<none>"
	}

	parts := strings.Split(repoTags[0], ":")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return repoTags[0], "<none>"
}

// SortImagesByCreationTime sorts images by creation time (newest first)
func SortImagesByCreationTime(images []types.Image) {
	sort.Slice(images, func(i, j int) bool {
		return images[i].Created.After(images[j].Created)
	})
}

// SortVolumesByName sorts volumes by name
func SortVolumesByName(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
}

// BuildStopOptions builds stop options with optional timeout
func BuildStopOptions(timeout *time.Duration) container.StopOptions {
	opts := container.StopOptions{}
	if timeout != nil {
		opts.Signal = "SIGTERM"
		seconds := int(timeout.Seconds())
		opts.Timeout = &seconds
	}
	return opts
}

// ValidateID validates that an ID is not empty
func ValidateID(id, idType string) error {
	if id == "" {
		return fmt.Errorf("%s cannot be empty", idType)
	}
	return nil
}
