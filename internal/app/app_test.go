package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/d5r/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		RefreshInterval: 5,
		LogLevel:        "INFO",
		DockerHost:      "unix:///var/run/docker.sock",
		Theme:           "default",
	}

	app, err := New(cfg)
	// This will likely fail in test environment since Docker is not available
	// but we can test the error handling
	if err != nil {
		assert.Contains(t, err.Error(), "docker client creation failed")
	} else {
		assert.NotNil(t, app)
	}
}

func TestNew_NilConfig(t *testing.T) {
	app, err := New(nil)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestNew_InvalidDockerHost(t *testing.T) {
	cfg := &config.Config{
		RefreshInterval: 5,
		LogLevel:        "INFO",
		DockerHost:      "invalid://host",
		Theme:           "default",
	}

	app, err := New(cfg)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestApp_Shutdown(t *testing.T) {
	cfg := &config.Config{
		RefreshInterval: 5,
		LogLevel:        "INFO",
		DockerHost:      "unix:///var/run/docker.sock",
		Theme:           "default",
	}

	app, err := New(cfg)
	if err == nil {
		assert.NotPanics(t, func() {
			app.Shutdown()
		})
	}
}

func TestApp_Run(t *testing.T) {
	t.Skip("Skipping full app test - requires proper mocking")
}

func TestApp_ConfigValidation(t *testing.T) {
	testCases := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "Valid Unix Socket",
			config: &config.Config{
				RefreshInterval: 5,
				LogLevel:        "INFO",
				DockerHost:      "unix:///var/run/docker.sock",
				Theme:           "default",
			},
			expectError: false,
		},
		{
			name: "Valid TCP Host",
			config: &config.Config{
				RefreshInterval: 5,
				LogLevel:        "INFO",
				DockerHost:      "tcp://localhost:2375",
				Theme:           "default",
			},
			expectError: false,
		},
		{
			name: "Valid Windows Named Pipe",
			config: &config.Config{
				RefreshInterval: 5,
				LogLevel:        "INFO",
				DockerHost:      "npipe:////./pipe/docker_engine",
				Theme:           "default",
			},
			expectError: false,
		},
		{
			name: "Invalid Host",
			config: &config.Config{
				RefreshInterval: 5,
				LogLevel:        "INFO",
				DockerHost:      "invalid://host",
				Theme:           "default",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := New(tc.config)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				if err != nil {
					assert.Contains(t, err.Error(), "docker client creation failed")
				} else {
					assert.NotNil(t, app)
				}
			}
		})
	}
}

func TestApp_LogLevelHandling(t *testing.T) {
	// Test different log levels
	logLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}

	for _, level := range logLevels {
		t.Run("LogLevel_"+level, func(t *testing.T) {
			cfg := &config.Config{
				RefreshInterval: 5,
				LogLevel:        level,
				DockerHost:      "unix:///var/run/docker.sock",
				Theme:           "default",
			}

			app, err := New(cfg)
			if err == nil {
				assert.NotNil(t, app)
			} else {
				assert.Contains(t, err.Error(), "docker client creation failed")
			}
		})
	}
}

func TestApp_ThemeHandling(t *testing.T) {
	themes := []string{"default", "dark", "light", "custom"}

	for _, theme := range themes {
		t.Run("Theme_"+theme, func(t *testing.T) {
			cfg := &config.Config{
				RefreshInterval: 5,
				LogLevel:        "INFO",
				DockerHost:      "unix:///var/run/docker.sock",
				Theme:           theme,
			}

			app, err := New(cfg)
			if err == nil {
				assert.NotNil(t, app)
			} else {
				assert.Contains(t, err.Error(), "docker client creation failed")
			}
		})
	}
}
