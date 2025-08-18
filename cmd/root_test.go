package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/logger"
)

func TestRootCmd(t *testing.T) {
	// Test that root command is properly configured
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "d5r", rootCmd.Use)
	assert.Equal(t, "D5r - Docker CLI Dashboard", rootCmd.Short)
	assert.Contains(t, rootCmd.Long, "Container Management")
	assert.Contains(t, rootCmd.Long, "Image Management")
	assert.Contains(t, rootCmd.Long, "Volume Management")
	assert.Contains(t, rootCmd.Long, "Network Management")
}

func TestThemeCmd(t *testing.T) {
	// Test that theme command is properly configured
	assert.NotNil(t, themeCmd)
	assert.Equal(t, "theme", themeCmd.Use)
	assert.Equal(t, "Manage theme configuration", themeCmd.Short)
}

func TestRootCmdFlags(t *testing.T) {
	// Test that all expected flags are registered
	flags := rootCmd.PersistentFlags()

	hostFlag := flags.Lookup("host")
	assert.NotNil(t, hostFlag)
	assert.Equal(t, "host", hostFlag.Name)

	refreshFlag := flags.Lookup("refresh")
	assert.NotNil(t, refreshFlag)
	assert.Equal(t, "refresh", refreshFlag.Name)

	logLevelFlag := flags.Lookup("log-level")
	assert.NotNil(t, logLevelFlag)
	assert.Equal(t, "log-level", logLevelFlag.Name)

	themeFlag := flags.Lookup("theme")
	assert.NotNil(t, themeFlag)
	assert.Equal(t, "theme", themeFlag.Name)
}

func TestApplyFlagOverrides(t *testing.T) {
	// Test flag override functionality
	cfg := &config.Config{
		DockerHost:      "unix:///var/run/docker.sock",
		RemoteHost:      "",
		RefreshInterval: 5,
		LogLevel:        "INFO",
		Theme:           "default",
	}

	// Test remote host override
	remoteHost = "tcp://192.168.1.100:2375"
	applyFlagOverrides(cfg)
	assert.Equal(t, "tcp://192.168.1.100:2375", cfg.RemoteHost)
	assert.Equal(t, "tcp://192.168.1.100:2375", cfg.DockerHost)

	// Test refresh interval override
	refresh = 15
	applyFlagOverrides(cfg)
	assert.Equal(t, 15, cfg.RefreshInterval)

	// Test log level override
	logLevel = "DEBUG"
	applyFlagOverrides(cfg)
	assert.Equal(t, "DEBUG", cfg.LogLevel)

	// Test theme override
	theme = "dark"
	applyFlagOverrides(cfg)
	assert.Equal(t, "dark", cfg.Theme)

	// Reset global variables
	remoteHost = ""
	refresh = 5
	logLevel = "INFO"
	theme = ""
}

func TestSetLogLevel(t *testing.T) {
	// Test log level setting - verify functions don't panic
	log := logger.GetLogger()

	// Test that setting different log levels doesn't panic
	assert.NotPanics(t, func() {
		setLogLevel(log, "DEBUG")
	})

	assert.NotPanics(t, func() {
		setLogLevel(log, "WARN")
	})

	assert.NotPanics(t, func() {
		setLogLevel(log, "ERROR")
	})

	assert.NotPanics(t, func() {
		setLogLevel(log, "INFO")
	})

	assert.NotPanics(t, func() {
		setLogLevel(log, "INVALID")
	})
}

func TestRootCmdSubcommands(t *testing.T) {
	// Test that all expected subcommands are added
	subcommands := rootCmd.Commands()
	assert.Len(t, subcommands, 3) // config, theme, and version commands

	// Check that all expected commands are present
	commandNames := make([]string, len(subcommands))
	for i, cmd := range subcommands {
		commandNames[i] = cmd.Use
	}

	assert.Contains(t, commandNames, "config")
	assert.Contains(t, commandNames, "theme")
	assert.Contains(t, commandNames, "version")
}

func TestRootCmdIntegration(t *testing.T) {
	// Test that root command can be created and configured without panicking
	cmd := &cobra.Command{}
	cmd.AddCommand(rootCmd)

	// This is a basic test to ensure the command structure is valid
	assert.NotNil(t, cmd.Commands())
	assert.Len(t, cmd.Commands(), 1)
	assert.Equal(t, "d5r", cmd.Commands()[0].Use)
}
