package shared

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
)

// DockerInfoService defines the interface for Docker system information
type DockerInfoService interface {
	GetDockerInfo(ctx context.Context) (*DockerInfo, error)
}

// dockerInfoService implements DockerInfoService
type dockerInfoService struct {
	client *docker.Client
}

// NewDockerInfoService creates a new Docker info service
func NewDockerInfoService(client *docker.Client) DockerInfoService {
	if client == nil {
		return nil
	}
	return &dockerInfoService{
		client: client,
	}
}

// GetDockerInfo retrieves Docker system information
func (s *dockerInfoService) GetDockerInfo(ctx context.Context) (*DockerInfo, error) {
	info, err := s.client.GetInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker info: %w", err)
	}

	// Helper functions to safely extract and convert values from map
	getInt := func(key string) int {
		if val, ok := info[key]; ok {
			if f, ok := val.(float64); ok {
				return int(f)
			}
		}
		return 0
	}

	getString := func(key string) string {
		if val, ok := info[key]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
		return ""
	}

	getBool := func(key string) bool {
		if val, ok := info[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	getStringSlice := func(key string) []string {
		if val, ok := info[key]; ok {
			if slice, ok := val.([]any); ok {
				result := make([]string, len(slice))
				for i, v := range slice {
					if s, ok := v.(string); ok {
						result[i] = s
					}
				}
				return result
			}
		}
		return nil
	}

	dockerInfo := &DockerInfo{
		Containers:         getInt("Containers"),
		Images:             getInt("Images"),
		Version:            getString("ServerVersion"),
		OS:                 getString("OperatingSystem"),
		Architecture:       getString("Architecture"),
		KernelVersion:      getString("KernelVersion"),
		Driver:             getString("Driver"),
		MemoryLimit:        getBool("MemoryLimit"),
		SwapLimit:          getBool("SwapLimit"),
		KernelMemory:       getBool("KernelMemory"),
		CPUCfsQuota:        getBool("CPUCfsQuota"),
		CPUCfsPeriod:       getBool("CPUCfsPeriod"),
		CPUShares:          getBool("CPUShares"),
		CPUSet:             getBool("CPUSet"),
		IPv4Forwarding:     getBool("IPv4Forwarding"),
		BridgeNfIptables:   getBool("BridgeNfIptables"),
		Debug:              getBool("Debug"),
		NFd:                getInt("NFd"),
		NGoroutines:        getInt("NGoroutines"),
		SystemTime:         getString("SystemTime"),
		LoggingDriver:      getString("LoggingDriver"),
		CgroupDriver:       getString("CgroupDriver"),
		OperatingSystem:    getString("OperatingSystem"),
		OSType:             getString("OSType"),
		IndexServerAddress: getString("IndexServerAddress"),
		ServerVersion:      getString("ServerVersion"),
		ClusterStore:       getString("ClusterStore"),
		ClusterAdvertise:   getString("ClusterAdvertise"),
		DefaultRuntime:     getString("DefaultRuntime"),
		LiveRestoreEnabled: getBool("LiveRestoreEnabled"),
		Isolation:          getString("Isolation"),
		InitBinary:         getString("InitBinary"),
		SecurityOptions:    getStringSlice("SecurityOptions"),
		ProductLicense:     getString("ProductLicense"),
		Warnings:           getStringSlice("Warnings"),
	}

	return dockerInfo, nil
}
