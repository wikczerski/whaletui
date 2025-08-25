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

// valueExtractors holds helper functions for extracting values from Docker info
type valueExtractors struct {
	getInt         func(key string) int
	getString      func(key string) string
	getBool        func(key string) bool
	getStringSlice func(key string) []string
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

	extractors := s.createExtractors(info)
	dockerInfo := s.buildDockerInfo(extractors)

	return dockerInfo, nil
}

// createExtractors creates helper functions to safely extract values from the info map
func (s *dockerInfoService) createExtractors(info map[string]any) *valueExtractors {
	return &valueExtractors{
		getInt: func(key string) int {
			if val, ok := info[key]; ok {
				if f, ok := val.(float64); ok {
					return int(f)
				}
			}
			return 0
		},
		getString: func(key string) string {
			if val, ok := info[key]; ok {
				if s, ok := val.(string); ok {
					return s
				}
			}
			return ""
		},
		getBool: func(key string) bool {
			if val, ok := info[key]; ok {
				if b, ok := val.(bool); ok {
					return b
				}
			}
			return false
		},
		getStringSlice: func(key string) []string {
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
		},
	}
}

// buildDockerInfo builds the DockerInfo struct using the provided extractors
func (s *dockerInfoService) buildDockerInfo(extractors *valueExtractors) *DockerInfo {
	return &DockerInfo{
		Containers:         extractors.getInt("Containers"),
		Images:             extractors.getInt("Images"),
		Version:            extractors.getString("ServerVersion"),
		OS:                 extractors.getString("OperatingSystem"),
		Architecture:       extractors.getString("Architecture"),
		KernelVersion:      extractors.getString("KernelVersion"),
		Driver:             extractors.getString("Driver"),
		MemoryLimit:        extractors.getBool("MemoryLimit"),
		SwapLimit:          extractors.getBool("SwapLimit"),
		KernelMemory:       extractors.getBool("KernelMemory"),
		CPUCfsQuota:        extractors.getBool("CPUCfsQuota"),
		CPUCfsPeriod:       extractors.getBool("CPUCfsPeriod"),
		CPUShares:          extractors.getBool("CPUShares"),
		CPUSet:             extractors.getBool("CPUSet"),
		IPv4Forwarding:     extractors.getBool("IPv4Forwarding"),
		BridgeNfIptables:   extractors.getBool("BridgeNfIptables"),
		Debug:              extractors.getBool("Debug"),
		NFd:                extractors.getInt("NFd"),
		NGoroutines:        extractors.getInt("NGoroutines"),
		SystemTime:         extractors.getString("SystemTime"),
		LoggingDriver:      extractors.getString("LoggingDriver"),
		CgroupDriver:       extractors.getString("CgroupDriver"),
		OperatingSystem:    extractors.getString("OperatingSystem"),
		OSType:             extractors.getString("OSType"),
		IndexServerAddress: extractors.getString("IndexServerAddress"),
		ServerVersion:      extractors.getString("ServerVersion"),
		ClusterStore:       extractors.getString("ClusterStore"),
		ClusterAdvertise:   extractors.getString("ClusterAdvertise"),
		DefaultRuntime:     extractors.getString("DefaultRuntime"),
		LiveRestoreEnabled: extractors.getBool("LiveRestoreEnabled"),
		Isolation:          extractors.getString("Isolation"),
		InitBinary:         extractors.getString("InitBinary"),
		SecurityOptions:    extractors.getStringSlice("SecurityOptions"),
		ProductLicense:     extractors.getString("ProductLicense"),
		Warnings:           extractors.getStringSlice("Warnings"),
	}
}
