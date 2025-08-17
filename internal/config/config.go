package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the application configuration
type Config struct {
	RefreshInterval int    `json:"refresh_interval"`
	LogLevel        string `json:"log_level"`
	DockerHost      string `json:"docker_host"`
	Theme           string `json:"theme"`
	RemoteHost      string `json:"remote_host,omitempty"` // Command line specified remote host
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	host := "unix:///var/run/docker.sock"
	if runtime.GOOS == "windows" {
		// For Docker Desktop on Windows, use the default (empty) to let Docker client auto-detect
		// This allows it to work with both Windows containers and WSL2 Linux containers
		host = ""
	}

	return &Config{
		RefreshInterval: 5,
		LogLevel:        "INFO",
		DockerHost:      host,
		Theme:           "default",
	}
}

// Load loads the configuration from file
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home directory access failed: %w", err)
	}

	configDir := filepath.Join(homeDir, ".dockerk9s")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("config directory creation failed: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := saveConfig(configFile, cfg); err != nil {
			return nil, fmt.Errorf("config save failed: %w", err)
		}
		return cfg, nil
	}

	cfg := &Config{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("config read failed: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory access failed: %w", err)
	}

	configDir := filepath.Join(homeDir, ".dockerk9s")
	configFile := filepath.Join(configDir, "config.json")

	return saveConfig(configFile, c)
}

func saveConfig(file string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("config marshal failed: %w", err)
	}

	if err := os.WriteFile(file, data, 0o600); err != nil {
		return fmt.Errorf("config write failed: %w", err)
	}

	return nil
}
