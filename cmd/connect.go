package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker/dockerssh"
	"github.com/wikczerski/whaletui/internal/logger"
	"golang.org/x/term"
)

var (
	sshKeyPath  string
	sshPassword bool
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
	RunE:          runConnectCommand,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func setupConnectFlags() {
	connectCmd.Flags().
		String("host", "", "Remote Docker host (e.g., 192.168.1.100, tcp://192.168.1.100, or ssh://user@host)")
	connectCmd.Flags().
		String("user", "", "SSH username for remote host connection (not needed if host includes ssh:// scheme)")
	connectCmd.Flags().Int("port", 2375, "Port for SSH fallback Docker proxy (default: 2375)")
	connectCmd.Flags().Bool("diagnose", false, "Run SSH connection diagnostics before connecting")
	connectCmd.Flags().
		StringVar(&sshKeyPath, "ssh-key-path", "", "Path to SSH private key file for authentication")
	connectCmd.Flags().
		BoolVar(&sshPassword, "password", false, "Prompt for SSH password authentication")
	_ = connectCmd.MarkFlagRequired("host")
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

	// Apply SSH authentication options
	applySSHOptions(cfg)

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

// applySSHOptions applies SSH authentication options to the configuration
func applySSHOptions(cfg *config.Config) {
	if sshKeyPath != "" {
		cfg.SSHKeyPath = sshKeyPath
	}

	if sshPassword {
		password, err := promptForPassword()
		if err != nil {
			logger.Warn("Failed to get password input", "error", err)
			return
		}
		cfg.SSHPassword = password
	}
}

// promptForPassword securely prompts the user for a password
func promptForPassword() (string, error) {
	fmt.Print("Enter SSH password: ")

	// Read password from stdin without echoing
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println() // Add newline after password input

	return string(passwordBytes), nil
}
