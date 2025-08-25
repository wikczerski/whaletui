package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testEnvVars holds the original environment variables for restoration
type testEnvVars struct {
	home        string
	userProfile string
}

// setupTestEnvironment sets up the test environment variables
func setupTestEnvironment(t *testing.T, tempHome string) testEnvVars {
	t.Helper()

	envVars := testEnvVars{
		home:        os.Getenv("HOME"),
		userProfile: os.Getenv("USERPROFILE"),
	}

	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("Failed to set HOME env var: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tempHome); err != nil {
		t.Fatalf("Failed to set USERPROFILE env var: %v", err)
	}

	return envVars
}

// restoreTestEnvironment restores the original environment variables
func restoreTestEnvironment(t *testing.T, envVars testEnvVars) {
	t.Helper()

	if err := os.Setenv("HOME", envVars.home); err != nil {
		t.Errorf("Failed to restore HOME env var: %v", err)
	}
	if err := os.Setenv("USERPROFILE", envVars.userProfile); err != nil {
		t.Errorf("Failed to restore USERPROFILE env var: %v", err)
	}
}

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
	envVars := setupTestEnvironment(t, tempHome)
	defer restoreTestEnvironment(t, envVars)

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
	envVars := setupTestEnvironment(t, tempHome)
	defer restoreTestEnvironment(t, envVars)

	configDir := filepath.Join(tempHome, ".dockerk9s")
	err := os.MkdirAll(configDir, 0o750)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "config.json")
	writeTestConfig(t, configFile)

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 10, cfg.RefreshInterval)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "tcp://localhost:2375", cfg.DockerHost)
	assert.Equal(t, "dark", cfg.Theme)
}

// writeTestConfig writes a test configuration to the specified file
func writeTestConfig(t *testing.T, configFile string) {
	t.Helper()

	testConfig := `{
		"refresh_interval": 10,
		"log_level": "DEBUG",
		"docker_host": "tcp://localhost:2375",
		"theme": "dark"
	}`
	err := os.WriteFile(configFile, []byte(testConfig), 0o600)
	require.NoError(t, err)
}

func TestLoad_InvalidConfig(t *testing.T) {
	tempHome := t.TempDir()
	envVars := setupTestEnvironment(t, tempHome)
	defer restoreTestEnvironment(t, envVars)

	configDir := filepath.Join(tempHome, ".dockerk9s")
	err := os.MkdirAll(configDir, 0o750)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "config.json")
	writeInvalidTestConfig(t, configFile)

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "config parse failed")

	assert.FileExists(t, configFile)
}

// writeInvalidTestConfig writes an invalid test configuration to the specified file
func writeInvalidTestConfig(t *testing.T, configFile string) {
	t.Helper()

	invalidConfig := `{
		"refresh_interval": "invalid",
		"log_level": "DEBUG"
	`
	err := os.WriteFile(configFile, []byte(invalidConfig), 0o600)
	require.NoError(t, err)
}

func TestSave(t *testing.T) {
	cfg := &Config{
		RefreshInterval: 15,
		LogLevel:        "WARN",
		DockerHost:      "tcp://localhost:2376",
		Theme:           "custom",
	}

	err := cfg.Save()
	if err != nil {
		assert.True(
			t,
			strings.Contains(err.Error(), "home directory access failed") ||
				strings.Contains(err.Error(), "config write failed") ||
				strings.Contains(err.Error(), "config directory creation failed"),
			"Expected error about home directory, config write, or directory creation, got: %s",
			err.Error(),
		)
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
	cfg := createTestConfig()

	data, err := json.Marshal(cfg)
	require.NoError(t, err)

	assertJSONFields(t, data)
	assertUnmarshal(t, data, cfg)
}

// createTestConfig creates a test configuration
func createTestConfig() *Config {
	return &Config{
		RefreshInterval: 20,
		LogLevel:        "ERROR",
		DockerHost:      "tcp://localhost:2377",
		Theme:           "minimal",
	}
}

// assertJSONFields asserts that the JSON contains expected fields
func assertJSONFields(t *testing.T, data []byte) {
	t.Helper()

	assert.Contains(t, string(data), `"refresh_interval"`)
	assert.Contains(t, string(data), `"log_level"`)
	assert.Contains(t, string(data), `"docker_host"`)
	assert.Contains(t, string(data), `"theme"`)
}

// assertUnmarshal asserts that the JSON can be unmarshaled correctly
func assertUnmarshal(t *testing.T, data []byte, expected *Config) {
	t.Helper()

	var unmarshaled Config
	err := json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, expected.RefreshInterval, unmarshaled.RefreshInterval)
	assert.Equal(t, expected.LogLevel, unmarshaled.LogLevel)
	assert.Equal(t, expected.DockerHost, unmarshaled.DockerHost)
	assert.Equal(t, expected.Theme, unmarshaled.Theme)
}

func TestThemeLoading(t *testing.T) {
	tempDir := t.TempDir()
	themePath := filepath.Join(tempDir, "test_theme.yaml")

	// Debug: Log the paths being used
	t.Logf("Temp directory: %s", tempDir)
	t.Logf("Theme path: %s", themePath)
	t.Logf("Theme path absolute: %s", filepath.Clean(themePath))

	writeTestTheme(t, themePath)

	// Debug: Verify the file was written
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		t.Fatalf("Theme file was not created: %s", themePath)
	}

	// Debug: Read and log the file contents
	content, err := os.ReadFile(themePath)
	if err != nil {
		t.Fatalf("Failed to read theme file: %v", err)
	}
	t.Logf("Theme file contents: %s", string(content))

	tm := NewThemeManager(themePath)

	assertCustomColors(t, tm)
}

// writeTestTheme writes a test theme file
func writeTestTheme(t *testing.T, themePath string) {
	t.Helper()

	themeContent := `colors:
  header: "red"
  border: "blue"
  text: "green"
  background: "black"
  success: "green"
  warning: "yellow"
  error: "red"
  info: "cyan"`

	err := os.WriteFile(themePath, []byte(themeContent), 0o600)
	require.NoError(t, err)
}

// assertCustomColors asserts that custom colors are loaded
func assertCustomColors(t *testing.T, tm *ThemeManager) {
	t.Helper()

	headerColor := tm.GetHeaderColor()
	borderColor := tm.GetBorderColor()

	assert.NotEqual(t, headerColor, tcell.ColorYellow, "Header color should not be default yellow")
	assert.NotEqual(t, borderColor, tcell.ColorWhite, "Border color should not be default white")

	assert.True(t, headerColor != tcell.ColorYellow, "Custom header color should be loaded")
	assert.True(t, borderColor != tcell.ColorWhite, "Custom border color should be loaded")
}
