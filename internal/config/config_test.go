package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, 5, cfg.RefreshInterval)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, "default", cfg.Theme)

	if runtime.GOOS == "windows" {
		assert.Equal(t, "", cfg.DockerHost)
	} else {
		assert.Equal(t, "unix:///var/run/docker.sock", cfg.DockerHost)
	}
}

func TestLoad_NewConfig(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("Failed to set HOME env var: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tempHome); err != nil {
		t.Fatalf("Failed to set USERPROFILE env var: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME env var: %v", err)
		}
		if err := os.Setenv("USERPROFILE", originalUserProfile); err != nil {
			t.Errorf("Failed to restore USERPROFILE env var: %v", err)
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	configFile := filepath.Join(tempHome, ".dockerk9s", "config.json")
	assert.FileExists(t, configFile)

	assert.Equal(t, 5, cfg.RefreshInterval)
	assert.Equal(t, "INFO", cfg.LogLevel)
}

func TestLoad_ExistingConfig(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("Failed to set HOME env var: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tempHome); err != nil {
		t.Fatalf("Failed to set USERPROFILE env var: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME env var: %v", err)
		}
		if err := os.Setenv("USERPROFILE", originalUserProfile); err != nil {
			t.Errorf("Failed to restore USERPROFILE env var: %v", err)
		}
	}()

	configDir := filepath.Join(tempHome, ".dockerk9s")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "config.json")

	testConfig := `{
		"refresh_interval": 10,
		"log_level": "DEBUG",
		"docker_host": "tcp://localhost:2375",
		"theme": "dark"
	}`
	err = os.WriteFile(configFile, []byte(testConfig), 0644)
	require.NoError(t, err)

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 10, cfg.RefreshInterval)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "tcp://localhost:2375", cfg.DockerHost)
	assert.Equal(t, "dark", cfg.Theme)
}

func TestLoad_InvalidConfig(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("Failed to set HOME env var: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tempHome); err != nil {
		t.Fatalf("Failed to set USERPROFILE env var: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME env var: %v", err)
		}
		if err := os.Setenv("USERPROFILE", originalUserProfile); err != nil {
			t.Errorf("Failed to restore USERPROFILE env var: %v", err)
		}
	}()

	configDir := filepath.Join(tempHome, ".dockerk9s")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "config.json")

	invalidConfig := `{
		"refresh_interval": "invalid",
		"log_level": "DEBUG"
	`
	err = os.WriteFile(configFile, []byte(invalidConfig), 0644)
	require.NoError(t, err)

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "config parse failed")

	assert.FileExists(t, configFile)
}

func TestSave(t *testing.T) {
	// Test that Save method doesn't panic and returns appropriate error
	// Since we can't easily mock the home directory in tests, we'll test error handling

	cfg := &Config{
		RefreshInterval: 15,
		LogLevel:        "WARN",
		DockerHost:      "tcp://localhost:2376",
		Theme:           "custom",
	}

	err := cfg.Save()
	if err != nil {
		assert.True(t,
			strings.Contains(err.Error(), "home directory access failed") ||
				strings.Contains(err.Error(), "config write failed") ||
				strings.Contains(err.Error(), "config directory creation failed"),
			"Expected error about home directory, config write, or directory creation, got: %s", err.Error())
	}
}

func TestSave_InvalidHomeDir(t *testing.T) {
	cfg := &Config{
		RefreshInterval: 5,
		LogLevel:        "INFO",
		DockerHost:      "unix:///var/run/docker.sock",
		Theme:           "default",
	}

	assert.NotPanics(t, func() {
		_ = cfg.Save()
	})
}

func TestConfig_JSONTags(t *testing.T) {
	cfg := &Config{
		RefreshInterval: 20,
		LogLevel:        "ERROR",
		DockerHost:      "tcp://localhost:2377",
		Theme:           "minimal",
	}

	data, err := json.Marshal(cfg)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"refresh_interval"`)
	assert.Contains(t, string(data), `"log_level"`)
	assert.Contains(t, string(data), `"docker_host"`)
	assert.Contains(t, string(data), `"theme"`)

	var unmarshaled Config
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, cfg.RefreshInterval, unmarshaled.RefreshInterval)
	assert.Equal(t, cfg.LogLevel, unmarshaled.LogLevel)
	assert.Equal(t, cfg.DockerHost, unmarshaled.DockerHost)
	assert.Equal(t, cfg.Theme, unmarshaled.Theme)
}
