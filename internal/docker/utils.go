package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func marshalToMap(v any) (map[string]any, error) {
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

func formatSize(size int64) string {
	const (
		KB int64 = 1024
		MB int64 = KB * 1024
		GB int64 = MB * 1024
		TB int64 = GB * 1024
	)

	var (
		unit  string
		value float64
	)

	switch {
	case size >= TB:
		unit = "TB"
		value = float64(size) / float64(TB)
	case size >= GB:
		unit = "GB"
		value = float64(size) / float64(GB)
	case size >= MB:
		unit = "MB"
		value = float64(size) / float64(MB)
	case size >= KB:
		unit = "KB"
		value = float64(size) / float64(KB)
	default:
		unit = "B"
		value = float64(size)
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

// SuggestConfigUpdate suggests updating the configuration file with a working Docker host
func SuggestConfigUpdate(detectedHost string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("config update suggestion only available on Windows")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory access failed: %w", err)
	}

	configDir := filepath.Join(homeDir, ".dockerk9s")
	configFile := filepath.Join(configDir, "config.json")

	// Read existing config
	var config map[string]any
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err == nil {
			if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
				fmt.Printf("Warning: failed to parse config file: %v\n", unmarshalErr)
			}
		}
	}

	// Update the Docker host
	if config == nil {
		config = make(map[string]any)
	}
	config["docker_host"] = detectedHost

	// Ensure other required fields exist
	if _, exists := config["refresh_interval"]; !exists {
		config["refresh_interval"] = 5
	}
	if _, exists := config["log_level"]; !exists {
		config["log_level"] = "INFO"
	}
	if _, exists := config["theme"]; !exists {
		config["theme"] = "default"
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("config directory creation failed: %w", err)
	}

	// Write updated config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("config marshal failed: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0o600); err != nil {
		return fmt.Errorf("config write failed: %w", err)
	}

	return nil
}
