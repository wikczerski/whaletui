package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/errorhandler"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

var (
	refresh     int
	logLevel    string
	logFilePath string
	theme       string
)

var dockerErrorPatterns = []string{
	"docker client creation failed",
	"failed to create Docker client",
	"failed to connect to Docker",
	"failed to connect to Docker via SSH",
	"failed to connect to Docker via SSH proxy",
	"connection refused",
	"no such host",
	"timeout",
	"permission denied",
	"dial tcp",
	"dial unix",
	"connection reset by peer",
	"no route to host",
	"network is unreachable",
	"host is down",
	"connection timed out",
	"broken pipe",
	"file not found",
	"access denied",
	"operation not permitted",
}

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:           "theme",
	Short:         "Manage theme configuration",
	Long:          `Manage theme configuration for the whaletui UI.`,
	RunE:          runThemeCommand,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whaletui",
	Short: "whaletui - Docker CLI Dashboard",
	Long: `
Whaletui is a terminal-based Docker management tool inspired by k9s,
providing an intuitive and powerful interface for managing Docker containers,
images, volumes, and networks with a modern, responsive TUI.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if logLevel != "INFO" {
			cfg, err := config.Load()
			if err != nil {
				logger.SetLevelWithPath(logLevel, "")
				return
			}

			logPath := cfg.LogFilePath
			if logFilePath != "" {
				logPath = logFilePath
			}
			logger.SetLevelWithPath(logLevel, logPath)
		}
	},
	RunE:          runApp,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	err := rootCmd.Execute()
	if err != nil {
		logger.Error("Command execution failed", "error", err)
		logger.CloseLogFile()
		os.Exit(1)
	}
}

func init() {
	constants.SetAppVersion(Version)
	logger.SetTUIMode(true)

	setupRootFlags()
	setupConnectFlags()
	addSubcommands()
}

func setupRootFlags() {
	rootCmd.PersistentFlags().IntVar(&refresh, "refresh", 5, "Refresh interval in seconds")
	rootCmd.PersistentFlags().
		StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.PersistentFlags().
		StringVar(&logFilePath, "log-file", "", "Path to log file (used when --log-level DEBUG is set)")
	rootCmd.PersistentFlags().
		StringVar(&theme, "theme", "", "Path to theme configuration file (YAML/JSON)")
}

func addSubcommands() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(themeCmd)
}

// runApp runs the main application with the provided configuration
func runApp(_ *cobra.Command, _ []string) error {
	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("Config load failed", "error", err)
		return err
	}

	applyFlagOverrides(cfg)
	setLogLevel(cfg.LogLevel)

	logConnectionInfo(cfg, log)

	return createAndRunApplication(cfg, log)
}

// createAndRunApplication creates and runs the application
func createAndRunApplication(cfg *config.Config, log *slog.Logger) error {
	application, err := app.New(cfg)
	if err != nil {
		log.Error("App init failed", "error", err)

		// Check if this is a Docker connection error and provide helpful guidance
		if isDockerConnectionError(err) {
			return handleDockerConnectionError(err, cfg)
		}

		return err
	}

	return runApplicationWithShutdown(application, log)
}

// isDockerConnectionError checks if the error is related to Docker connection issues
func isDockerConnectionError(err error) bool {
	errStr := err.Error()

	for _, pattern := range dockerErrorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// handleDockerConnectionError provides user-friendly error handling for Docker connection issues
func handleDockerConnectionError(err error, cfg *config.Config) error {
	handler := errorhandler.NewDockerErrorHandler(
		err,
		cfg,
		errorhandler.AppRunner(runApplicationWithShutdown),
		UserInteraction{},
	)
	return handler.Handle()
}

func logConnectionInfo(cfg *config.Config, log *slog.Logger) {
	if cfg.RemoteHost != "" {
		log.Info("Connecting to remote Docker host", "host", cfg.RemoteHost)
	} else {
		log.Info("Connecting to local Docker instance")
	}
}

// applyFlagOverrides applies command line flag values to the configuration
func applyFlagOverrides(cfg *config.Config) {
	if refresh != 5 {
		cfg.RefreshInterval = refresh
	}

	if logLevel != "INFO" {
		cfg.LogLevel = logLevel
	}

	if logFilePath != "" {
		cfg.LogFilePath = logFilePath
	}

	if theme != "" {
		cfg.Theme = theme
	}
}

// setLogLevel sets the log level based on the configuration
func setLogLevel(level string) {
	cfg, err := config.Load()
	if err != nil {
		logger.SetLevelWithPath(level, "")
		return
	}

	logger.SetLevelWithPath(level, cfg.LogFilePath)
}

// runThemeCommand handles theme-related commands
func runThemeCommand(_ *cobra.Command, _ []string) error {
	log := logger.GetLogger()

	log.Debug(
		"Theme command started",
		"debug_mode",
		logger.IsDebugMode(),
		"log_file_path",
		logger.GetLogFilePath(),
	)

	themeManager := config.NewThemeManager("")

	defaultPath := "./config/theme.yaml"
	if err := themeManager.SaveTheme(defaultPath); err != nil {
		log.Error("Failed to save theme", "error", err)
		return err
	}

	log.Info("Theme configuration saved", "path", defaultPath)
	log.Info(
		"You can now customize the colors in this file and restart whaletui with --theme",
		"path",
		defaultPath,
	)

	return nil
}
