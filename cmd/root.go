package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/logger"
)

var (
	refresh     int
	logLevel    string
	logFilePath string
	theme       string
)

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
	RunE: runConnectCommand,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
}

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage theme configuration",
	Long:  `Manage theme configuration for the whaletui UI.`,
	RunE:  runThemeCommand,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whaletui",
	Short: "whaletui - Docker CLI Dashboard",
	Long: `whaletui is a terminal-based Docker management tool inspired by k9s,
providing an intuitive and powerful interface for managing Docker containers,
images, volumes, and networks with a modern, responsive TUI.

Features:
  • Container Management: View, start, stop, restart, delete, and manage containers
  • Image Management: Browse, inspect, and remove Docker images
  • Volume Management: Manage Docker volumes with ease
  • Network Management: View and manage Docker networks
  • Remote Host Support: Connect to Docker hosts on different machines via SSH or TCP
  • SSH Connections: Secure SSH connections to remote Docker hosts using ssh:// URLs
  • TCP Connections: Direct TCP connections to remote Docker daemons
  • Theme Support: Customizable color schemes via YAML/JSON configuration
  • Advanced Logging: Multistream logging to both console and file when using --log-level DEBUG

Commands:
  connect  - Connect to a remote Docker host via SSH
  theme    - Manage theme configuration

Examples:
  whaletui                                    - Start with local Docker instance
  whaletui --log-level DEBUG                  - Start with debug logging to console and file
  whaletui --log-level DEBUG --log-file ./myapp.log  - Start with debug logging to custom file
  whaletui connect --host ssh://admin@host1  - Connect to remote host via SSH
  whaletui connect --host 192.168.1.100 --user admin  - Connect to remote host
  whaletui theme                             - Manage theme configuration`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		// Set log level before any command executes
		if logLevel != "INFO" {
			// Get the current config to access log file path
			cfg, err := config.Load()
			if err != nil {
				// If config load fails, use default path
				logger.SetLevelWithPath(logLevel, "")
			} else {
				// Use config log file path, but allow command line override
				logPath := cfg.LogFilePath
				if logFilePath != "" {
					logPath = logFilePath
				}
				logger.SetLevelWithPath(logLevel, logPath)
			}
		}
	},
	RunE: runApp,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Set up the command to not show help on errors
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	err := rootCmd.Execute()
	if err != nil {
		// Log the error and exit without showing help
		logger.Error("Command execution failed", "error", err)
		logger.CloseLogFile()
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&refresh, "refresh", 5, "Refresh interval in seconds")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log-file", "", "Path to log file (used when --log-level DEBUG is set)")
	rootCmd.PersistentFlags().StringVar(&theme, "theme", "", "Path to theme configuration file (YAML/JSON)")

	// Connect command flags
	connectCmd.Flags().String("host", "", "Remote Docker host (e.g., 192.168.1.100, tcp://192.168.1.100, or ssh://user@host)")
	connectCmd.Flags().String("user", "", "SSH username for remote host connection (not needed if host includes ssh:// scheme)")
	connectCmd.Flags().Int("port", 2375, "Port for SSH fallback Docker proxy (default: 2375)")
	connectCmd.Flags().Bool("diagnose", false, "Run SSH connection diagnostics before connecting")
	_ = connectCmd.MarkFlagRequired("host")
	// User flag is only required when not using ssh:// URL format
	// We'll validate this in the command execution

	// Add subcommands
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(themeCmd)
}

// runConnectCommand handles the connect subcommand for remote Docker hosts
func runConnectCommand(cmd *cobra.Command, _ []string) error {
	log := logger.GetLogger()

	// Get flags from the connect command
	host, _ := cmd.Flags().GetString("host")
	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	diagnose, _ := cmd.Flags().GetBool("diagnose")

	// Validate that we have the required parameters
	if strings.HasPrefix(host, "ssh://") {
		// If host is already an SSH URL, user is optional
		if user == "" {
			// Try to extract user from SSH URL
			if strings.Contains(host, "@") {
				parts := strings.Split(host[6:], "@") // Remove "ssh://" prefix
				if len(parts) == 2 {
					user = parts[0]
				}
			}
		}
	} else {
		// For non-SSH URLs, user is required
		if user == "" {
			return fmt.Errorf("--user flag is required when not using ssh:// URL format")
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error("Config load failed", "error", err)
		return err
	}

	// Set remote connection parameters
	// If host already has ssh:// scheme, use it as-is
	// Otherwise, construct the SSH URL
	if strings.HasPrefix(host, "ssh://") {
		cfg.RemoteHost = host
		cfg.DockerHost = host
	} else {
		// Construct SSH URL from separate host and user parameters
		sshURL := fmt.Sprintf("ssh://%s@%s", user, host)
		cfg.RemoteHost = sshURL
		cfg.DockerHost = sshURL
	}
	cfg.RemoteUser = user
	cfg.RemotePort = port

	// Apply flag overrides before setting log level
	applyFlagOverrides(cfg)
	setLogLevel(cfg.LogLevel)

	// Run diagnostics if requested
	if diagnose {
		log.Info("Running SSH connection diagnostics...")
		if err := runSSHDiagnostics(host, user, port); err != nil {
			log.Error("SSH diagnostics failed", "error", err)
			log.Info("Continuing with connection attempt...")
		} else {
			log.Info("SSH diagnostics passed successfully")
		}
	}

	log.Info("Connecting to remote Docker host", "host", host, "user", user)

	application, err := app.New(cfg)
	if err != nil {
		log.Error("App init failed", "error", err)
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	uiShutdownCh := application.GetUI().GetShutdownChan()

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed", "error", err)
		}
	}()

	select {
	case <-sigCh:
		log.Info("Received shutdown signal, shutting down gracefully...")
	case <-uiShutdownCh:
		log.Info("Received UI shutdown signal, shutting down gracefully...")
	}

	defer cleanupTerminal()
	defer logger.CloseLogFile()

	application.Shutdown()
	return nil
}

// runSSHDiagnostics runs SSH connection diagnostics
func runSSHDiagnostics(host, user string, _ int) error {
	var sshHost string

	// Check if host is already an SSH URL
	if strings.HasPrefix(host, "ssh://") {
		// Extract the host part from ssh://[user@]host[:port]
		hostPart := host[6:] // Remove "ssh://" prefix
		sshHost = hostPart
	} else {
		// Construct SSH host from separate host and user parameters
		if user == "" {
			return fmt.Errorf("user must be specified for diagnostics")
		}
		sshHost = fmt.Sprintf("%s@%s", user, host)
	}

	sshClient, err := docker.NewSSHClient(sshHost, 22) // Default SSH port
	if err != nil {
		return fmt.Errorf("failed to create SSH client for diagnostics: %w", err)
	}

	// Run diagnostics
	if err := sshClient.DiagnoseConnection(); err != nil {
		return fmt.Errorf("SSH diagnostics failed: %w", err)
	}

	return nil
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

	// Validate that --user is provided when --host is specified
	if cfg.RemoteHost != "" {
		log.Info("Connecting to remote Docker host", "host", cfg.RemoteHost)
	} else {
		log.Info("Connecting to local Docker instance")
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Error("App init failed", "error", err)
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	uiShutdownCh := application.GetUI().GetShutdownChan()

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed", "error", err)
		}
	}()

	select {
	case <-sigCh:
		log.Info("Received shutdown signal, shutting down gracefully...")
	case <-uiShutdownCh:
		log.Info("Received UI shutdown signal, shutting down gracefully...")
	}

	defer cleanupTerminal()
	defer logger.CloseLogFile()

	application.Shutdown()
	return nil
}

// cleanupTerminal performs additional terminal cleanup operations
func cleanupTerminal() {
	cleanupTerminalClearScreen()
	cleanupTerminalResetColors()
	cleanupTerminalShowCursor()
	cleanupTerminalMoveCursorToTop()
	cleanupTerminalSyncStdout()
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
		logger.Warn("Failed to sync stdout", "error", err)
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

	log.Debug("Theme command started", "debug_mode", logger.IsDebugMode(), "log_file_path", logger.GetLogFilePath())

	themeManager := config.NewThemeManager("")

	defaultPath := "./config/theme.yaml"
	if err := themeManager.SaveTheme(defaultPath); err != nil {
		log.Error("Failed to save theme", "error", err)
		return err
	}

	log.Info("Theme configuration saved", "path", defaultPath)
	log.Info("You can now customize the colors in this file and restart whaletui with --theme", "path", defaultPath)

	return nil
}
