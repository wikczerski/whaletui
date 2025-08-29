// Package config provides configuration management for the WhaleTUI application.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Config represents the application configuration
type Config struct {
	RefreshInterval int    `json:"refresh_interval"`
	LogLevel        string `json:"log_level"`
	LogFilePath     string `json:"log_file_path,omitempty"`
	DockerHost      string `json:"docker_host"`
	Theme           string `json:"theme"`
	RemoteHost      string `json:"remote_host,omitempty"`
	RemoteUser      string `json:"remote_user,omitempty"`
	RemotePort      int    `json:"remote_port,omitempty"`
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
		LogFilePath:     "./logs/whaletui.log",
		DockerHost:      host,
		Theme:           "default",
		RemotePort:      2375,
	}
}

// Load loads the configuration from file
func Load() (*Config, error) {
	configDir, configFile, err := getConfigPaths()
	if err != nil {
		return nil, err
	}

	if isNewConfig(configFile) {
		return createNewConfig(configFile)
	}

	return loadExistingConfig(configFile, configDir)
}

// getConfigPaths gets the configuration directory and file paths
func getConfigPaths() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("home directory access failed: %w", err)
	}

	configDir := filepath.Join(homeDir, ".whaletui")
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return "", "", fmt.Errorf("config directory creation failed: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")
	return configDir, configFile, nil
}

// isNewConfig checks if the config file doesn't exist
func isNewConfig(configFile string) bool {
	_, err := os.Stat(configFile)
	return os.IsNotExist(err)
}

// createNewConfig creates a new default configuration
func createNewConfig(configFile string) (*Config, error) {
	cfg := DefaultConfig()
	if err := saveConfig(configFile, cfg); err != nil {
		return nil, fmt.Errorf("config save failed: %w", err)
	}
	return cfg, nil
}

// loadExistingConfig loads an existing configuration file
func loadExistingConfig(configFile, configDir string) (*Config, error) {
	cfg := DefaultConfig()

	if !isValidConfigPath(configFile, configDir) {
		return nil, errors.New("invalid config file path")
	}

	// nolint:gosec // Path is validated by isValidConfigPath before this function is called
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("config read failed: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	return cfg, nil
}

// isValidConfigPath validates the config file path to prevent directory traversal
func isValidConfigPath(configFile, configDir string) bool {
	// Clean both paths to remove any directory traversal attempts
	cleanConfigFile := filepath.Clean(configFile)
	cleanConfigDir := filepath.Clean(configDir)

	// Ensure both paths are absolute
	if !filepath.IsAbs(cleanConfigFile) || !filepath.IsAbs(cleanConfigDir) {
		return false
	}

	// Additional security: check for suspicious patterns
	// Check for directory traversal attempts
	if strings.Contains(cleanConfigFile, "..") {
		return false
	}

	// Check for home directory expansion attempts (but allow Windows short names like ~1)
	// Only reject paths that start with ~ or contain ~/ which could be home directory expansion
	if strings.HasPrefix(cleanConfigFile, "~") || strings.Contains(cleanConfigFile, "~/") {
		return false
	}

	// Ensure the config file is within the config directory
	// Use filepath.Rel to check if the config file is within the config directory
	relPath, err := filepath.Rel(cleanConfigDir, cleanConfigFile)
	if err != nil {
		// If filepath.Rel fails, fall back to string prefix check for compatibility
		// This handles edge cases in some CI environments
		return strings.HasPrefix(cleanConfigFile, cleanConfigDir)
	}

	// The relative path should not start with ".." (going up directories)
	return !strings.HasPrefix(relPath, "..")
}

// Save saves the configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory access failed: %w", err)
	}

	configDir := filepath.Join(homeDir, ".whaletui")
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
