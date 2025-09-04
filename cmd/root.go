package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"

	"github.com/wikczerski/whaletui/internal/docker/dockerssh"
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

// connectCmd represents the connect command for remote Docker hosts
var connectCmd = &cobra.Command{
	Use:   "connect [flags]",
	Short: "Connect to a remote Docker host via SSH",
	Long: `Connect to a remote Docker host using SSH or TCP connections.
This command supports both direct TCP connections and SSH connections to remote Docker hosts.

Examples:
  # Connect using SSH (recommended for security)
  whaletui connect --host ssh://admin@192.168.1.100

  # Connect using separate host and user parameters
  whaletui connect --host 192.168.1.100 --user admin

  # Connect using TCP
  whaletui connect --host tcp://192.168.1.100:2375`,
	RunE:          runConnectCommand,
	SilenceUsage:  true,
	SilenceErrors: true,
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

func setupConnectFlags() {
	connectCmd.Flags().
		String("host", "", "Remote Docker host (e.g., 192.168.1.100, tcp://192.168.1.100, or ssh://user@host)")
	connectCmd.Flags().
		String("user", "", "SSH username for remote host connection (not needed if host includes ssh:// scheme)")
	connectCmd.Flags().Int("port", 2375, "Port for SSH fallback Docker proxy (default: 2375)")
	connectCmd.Flags().Bool("diagnose", false, "Run SSH connection diagnostics before connecting")
	_ = connectCmd.MarkFlagRequired("host")
}

func addSubcommands() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(themeCmd)
}

// runConnectCommand handles the connect subcommand for remote Docker hosts
func runConnectCommand(cmd *cobra.Command, _ []string) error {
	log := logger.GetLogger()

	host, user, port, diagnose := extractConnectFlags(cmd)

	if err := validateConnectParameters(host, user); err != nil {
		return err
	}

	cfg, err := loadAndConfigureConnection(host, user, port)
	if err != nil {
		return err
	}

	if diagnose {
		runDiagnosticsIfRequested(host, user, port, log)
	}

	return connectToRemoteHost(host, user, cfg, log)
}

// connectToRemoteHost connects to the remote host
func connectToRemoteHost(host, user string, cfg *config.Config, log *slog.Logger) error {
	log.Info("Connecting to remote Docker host", "host", host, "user", user)

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

func extractConnectFlags(cmd *cobra.Command) (host, user string, port int, diagnose bool) {
	host, _ = cmd.Flags().GetString("host")
	user, _ = cmd.Flags().GetString("user")
	port, _ = cmd.Flags().GetInt("port")
	diagnose, _ = cmd.Flags().GetBool("diagnose")
	return host, user, port, diagnose
}

func validateConnectParameters(host, user string) error {
	if strings.HasPrefix(host, "ssh://") {
		return nil
	}

	if user == "" {
		return errors.New("--user flag is required when not using ssh:// URL format")
	}

	return nil
}

func loadAndConfigureConnection(host, user string, port int) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	configureRemoteConnection(cfg, host, user, port)
	applyFlagOverrides(cfg)
	setLogLevel(cfg.LogLevel)

	return cfg, nil
}

func configureRemoteConnection(cfg *config.Config, host, user string, port int) {
	if strings.HasPrefix(host, "ssh://") {
		cfg.RemoteHost = host
		cfg.DockerHost = host
	} else {
		sshURL := fmt.Sprintf("ssh://%s@%s", user, host)
		cfg.RemoteHost = sshURL
		cfg.DockerHost = sshURL
	}
	cfg.RemoteUser = user
	cfg.RemotePort = port
}

func runDiagnosticsIfRequested(host, user string, port int, log *slog.Logger) {
	log.Info("Running SSH connection diagnostics...")
	if err := runSSHDiagnostics(host, user, port, log); err != nil {
		log.Error("SSH diagnostics failed", "error", err)
		log.Info("Continuing with connection attempt...")
	} else {
		log.Info("SSH diagnostics passed successfully")
	}
}

// runSSHDiagnostics runs SSH connection diagnostics
func runSSHDiagnostics(host, user string, _ int, log *slog.Logger) error {
	sshHost := extractSSHHost(host, user)

	// Parse the host to validate it
	_, _, _, err := dockerssh.ParseSSHHost(sshHost)
	if err != nil {
		return fmt.Errorf("failed to parse SSH host: %w", err)
	}

	// SSH host parsing successful - diagnostics passed
	log.Info("SSH host parsing successful")

	return nil
}

func extractSSHHost(host, user string) string {
	if strings.HasPrefix(host, "ssh://") {
		return host[6:]
	}

	if user == "" {
		return ""
	}

	return fmt.Sprintf("%s@%s", user, host)
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
	errorHandler := newDockerErrorHandler(err, cfg)
	return errorHandler.handle()
}

// isDockerRunning checks if Docker is running on the local system
func isDockerRunning() bool {
	// Try to run a simple Docker command to check if it's accessible
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run() == nil
}

func logConnectionInfo(cfg *config.Config, log *slog.Logger) {
	if cfg.RemoteHost != "" {
		log.Info("Connecting to remote Docker host", "host", cfg.RemoteHost)
	} else {
		log.Info("Connecting to local Docker instance")
	}
}

func runApplicationWithShutdown(application *app.App, log *slog.Logger) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	uiShutdownCh := application.GetUI().GetShutdownChan()

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed", "error", err)
		}
	}()

	waitForShutdownSignal(sigCh, uiShutdownCh, log)
	cleanupAndShutdown(application)

	return nil
}

// waitForShutdownSignal waits for shutdown signals
func waitForShutdownSignal(sigCh chan os.Signal, uiShutdownCh chan struct{}, log *slog.Logger) {
	select {
	case <-sigCh:
		log.Info("Received shutdown signal, shutting down gracefully...")
	case <-uiShutdownCh:
		log.Info("Received UI shutdown signal, shutting down gracefully...")
	}
}

// cleanupAndShutdown performs cleanup and shutdown
func cleanupAndShutdown(application *app.App) {
	defer cleanupTerminal()
	defer logger.CloseLogFile()
	defer logger.SetTUIMode(false)

	application.Shutdown()
}

// cleanupTerminal performs additional terminal cleanup operations
func cleanupTerminal() {
	if logger.IsTUIMode() {
		return
	}

	cleanupOperations := []func(){
		cleanupTerminalClearScreen,
		cleanupTerminalResetColors,
		cleanupTerminalShowCursor,
		cleanupTerminalMoveCursorToTop,
		cleanupTerminalSyncStdout,
	}

	for _, operation := range cleanupOperations {
		operation()
	}
}

// cleanupTerminalClearScreen clears the terminal screen
func cleanupTerminalClearScreen() {
	if _, err := fmt.Fprint(os.Stdout, "\033[2J"); err != nil {
		logger.Warn("Failed to clear screen", "error", err)
	}
}

// cleanupTerminalResetColors resets terminal colors
func cleanupTerminalResetColors() {
	if _, err := fmt.Fprint(os.Stdout, "\033[0m"); err != nil {
		logger.Warn("Failed to reset colors", "error", err)
	}
}

// cleanupTerminalShowCursor shows the terminal cursor
func cleanupTerminalShowCursor() {
	if _, err := fmt.Fprint(os.Stdout, "\033[?25h"); err != nil {
		logger.Warn("Failed to show cursor", "error", err)
	}
}

// cleanupTerminalMoveCursorToTop moves the cursor to the top of the terminal
func cleanupTerminalMoveCursorToTop() {
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		logger.Warn("Failed to move cursor", "error", err)
	}
}

// cleanupTerminalSyncStdout synchronizes stdout
func cleanupTerminalSyncStdout() {
	if err := os.Stdout.Sync(); err != nil {
		logger.Debug("Failed to sync stdout", "error", err)
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
