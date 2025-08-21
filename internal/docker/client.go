package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/logger"
)

// Client represents a Docker client wrapper
type Client struct {
	cli     *client.Client
	cfg     *config.Config
	log     *slog.Logger
	sshConn *SSHConnection
	sshCtx  *SSHContext
}

// detectWindowsDockerHost attempts to find the correct Docker host on Windows
func detectWindowsDockerHost(log *slog.Logger) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("not on Windows")
	}

	// Common Windows Docker Desktop pipe paths
	possibleHosts := []string{
		"npipe:////./pipe/dockerDesktopLinuxEngine", // Linux containers
		"npipe:////./pipe/docker_engine",            // Windows containers
		"npipe:////./pipe/dockerDesktopEngine",      // Legacy Docker Desktop
	}

	for _, host := range possibleHosts {
		if isHostWorking(host, log) {
			return host, nil
		}
	}

	return "", fmt.Errorf("no working Docker host found")
}

// isHostWorking tests if a Docker host is working
func isHostWorking(host string, log *slog.Logger) bool {
	opts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(5 * time.Second),
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return false
	}
	defer closeClientSafely(cli, log)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	return err == nil
}

// closeClientSafely closes a Docker client and logs any errors
func closeClientSafely(cli *client.Client, log *slog.Logger) {
	if err := cli.Close(); err != nil {
		log.Warn("Failed to close Docker client during host detection", "error", err)
	}
}

// New creates a new Docker client
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	log := logger.GetLogger()

	client, err := createDockerClient(cfg, log)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// createDockerClient creates a Docker client with the given configuration
func createDockerClient(cfg *config.Config, log *slog.Logger) (*Client, error) {
	// For SSH connections, try direct connection first
	if cfg.RemoteHost != "" && strings.HasPrefix(cfg.RemoteHost, "ssh://") {
		return createSSHDockerClient(cfg, log)
	}

	opts, err := buildClientOptions(cfg)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return handleClientCreationError(cfg, log, err)
	}

	return testAndCreateClient(cfg, log, cli)
}

// createSSHDockerClient creates a Docker client via SSH connection
func createSSHDockerClient(cfg *config.Config, log *slog.Logger) (*Client, error) {
	log.Info("Attempting SSH connection for Docker client", "host", cfg.RemoteHost)

	// Extract host from SSH URL
	host, err := extractHostFromURL(cfg.RemoteHost)
	if err != nil {
		return nil, fmt.Errorf("failed to extract host from SSH URL: %w", err)
	}

	// First try to connect via SSH without socat
	sshCtx, err := establishSSHConnection(host, cfg.RemoteUser, log)
	if err != nil {
		log.Warn("SSH connection without socat failed, trying with socat as fallback", "error", err)
		return trySSHFallbackWithSocat(cfg, log, host)
	}

	// Try to create Docker client directly via SSH
	cli, err := createDockerClientViaSSH(sshCtx, log)
	if err != nil {
		log.Warn("Direct SSH connection failed, trying socat fallback", "error", err)
		sshCtx.Close()
		return trySSHFallbackWithSocat(cfg, log, host)
	}

	client := &Client{
		cli:    cli,
		cfg:    cfg,
		log:    log,
		sshCtx: sshCtx,
	}

	log.Info("Successfully connected to remote Docker via SSH (direct)", "host", cfg.RemoteHost)
	return client, nil
}

// buildClientOptions builds the client options based on configuration
func buildClientOptions(cfg *config.Config) ([]client.Opt, error) {
	var opts []client.Opt

	if cfg.RemoteHost != "" {
		// Check if this is an SSH connection
		if strings.HasPrefix(cfg.RemoteHost, "ssh://") {
			// For SSH connections, we'll handle this specially
			// The Docker client will handle ssh:// URLs natively
			opts = append(opts, client.WithHost(cfg.RemoteHost), client.WithTimeout(30*time.Second))
		} else {
			dockerHost, err := formatRemoteHost(cfg.RemoteHost)
			if err != nil {
				return nil, fmt.Errorf("invalid remote host format: %w", err)
			}
			opts = append(opts, client.WithHost(dockerHost), client.WithTimeout(30*time.Second))
		}
	} else if cfg.DockerHost != "" {
		opts = append(opts, client.WithHost(cfg.DockerHost))
	}

	opts = append(opts,
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(10*time.Second),
	)

	return opts, nil
}

// formatRemoteHost formats and validates a remote host URL
func formatRemoteHost(host string) (string, error) {
	// Handle SSH URLs - they should be passed through as-is
	if strings.HasPrefix(host, "ssh://") {
		return host, nil
	}

	if err := validateRemoteHost(host); err != nil {
		return "", err
	}

	// Automatically add tcp:// prefix if user didn't provide a scheme
	if !strings.Contains(host, "://") {
		host = "tcp://" + host
	}

	return host, nil
}

// handleClientCreationError handles errors during client creation, including Windows auto-detection
func handleClientCreationError(cfg *config.Config, log *slog.Logger, err error) (*Client, error) {
	if runtime.GOOS == "windows" && cfg.RemoteHost == "" {
		return tryWindowsAutoDetection(cfg, log)
	}
	return nil, fmt.Errorf("failed to create Docker client: %w", err)
}

// tryWindowsAutoDetection attempts to auto-detect the correct Docker host on Windows
func tryWindowsAutoDetection(cfg *config.Config, log *slog.Logger) (*Client, error) {
	log.Warn("Docker client creation failed, attempting to auto-detect correct host...")

	detectedHost, detectErr := detectWindowsDockerHost(log)
	if detectErr != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", detectErr)
	}

	log.Info("Detected working Docker host", "host", detectedHost)
	return createClientWithHost(cfg, log, detectedHost)
}

// createClientWithHost creates a client with a specific host
func createClientWithHost(cfg *config.Config, log *slog.Logger, host string) (*Client, error) {
	detectedOpts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(10 * time.Second),
	}

	detectedCli, detectedErr := client.NewClientWithOpts(detectedOpts...)
	if detectedErr != nil {
		return nil, fmt.Errorf("failed to create Docker client with detected host: %w", detectedErr)
	}

	return testAndCreateClient(cfg, log, detectedCli)
}

// testAndCreateClient tests the connection and creates the client
func testAndCreateClient(cfg *config.Config, log *slog.Logger, cli *client.Client) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		if closeErr := cli.Close(); closeErr != nil {
			log.Warn("Failed to close Docker client", "error", closeErr)
		}

		// If this is a remote host and direct connection failed, try SSH fallback
		// Only try SSH fallback if it's not already an SSH connection
		if cfg.RemoteHost != "" && !strings.HasPrefix(cfg.RemoteHost, "ssh://") {
			return trySSHFallback(cfg, log)
		}

		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	client := &Client{
		cli: cli,
		cfg: cfg,
		log: log,
	}

	logConnectionSuccess(cfg, log)
	return client, nil
}

// logConnectionSuccess logs the successful connection
func logConnectionSuccess(cfg *config.Config, log *slog.Logger) {
	if cfg.RemoteHost != "" {
		log.Info("Successfully connected to remote Docker host", "host", cfg.RemoteHost)
	} else {
		log.Info("Successfully connected to local Docker instance")
	}
}

// Close closes the Docker client and cleans up SSH connections
func (c *Client) Close() error {
	var errors []string

	c.log.Info("Docker client closing, starting cleanup")

	if c.cli != nil {
		if err := c.cli.Close(); err != nil {
			c.log.Error("Failed to close Docker client", "error", err)
			errors = append(errors, fmt.Sprintf("Docker client close: %v", err))
		} else {
			c.log.Info("Docker client closed successfully")
		}
	}

	if c.sshConn != nil {
		c.log.Info("Closing SSH connection with socat")
		if err := c.sshConn.Close(); err != nil {
			c.log.Error("Failed to close SSH connection", "error", err)
			errors = append(errors, fmt.Sprintf("SSH connection close: %v", err))
		} else {
			c.log.Info("SSH connection closed successfully")
		}
	}

	if c.sshCtx != nil {
		c.log.Info("Closing SSH context")
		if err := c.sshCtx.Close(); err != nil {
			c.log.Error("Failed to close SSH context", "error", err)
			errors = append(errors, fmt.Sprintf("SSH context close: %v", err))
		} else {
			c.log.Info("SSH context closed successfully")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close client: %s", strings.Join(errors, "; "))
	}

	return nil
}

// trySSHFallback attempts to connect via SSH when direct connection fails
func trySSHFallback(cfg *config.Config, log *slog.Logger) (*Client, error) {
	log.Info("Direct connection failed, attempting SSH fallback...")

	host, err := extractHostFromURL(cfg.RemoteHost)
	if err != nil {
		return nil, fmt.Errorf("SSH connection failed: %w", err)
	}

	// First try to connect via SSH without socat
	sshCtx, err := establishSSHConnection(host, cfg.RemoteUser, log)
	if err != nil {
		log.Warn("SSH connection without socat failed, trying with socat as fallback", "error", err)
		return trySSHFallbackWithSocat(cfg, log, host)
	}

	// Try to create Docker client directly via SSH
	cli, err := createDockerClientViaSSH(sshCtx, log)
	if err != nil {
		log.Warn("Direct SSH connection failed, trying socat fallback", "error", err)
		sshCtx.Close()
		return trySSHFallbackWithSocat(cfg, log, host)
	}

	client := &Client{
		cli:    cli,
		cfg:    cfg,
		log:    log,
		sshCtx: sshCtx,
	}

	log.Info("Successfully connected to remote Docker via SSH (direct)", "host", cfg.RemoteHost)
	return client, nil
}

// trySSHFallbackWithSocat attempts to connect via SSH with socat as a last resort
func trySSHFallbackWithSocat(cfg *config.Config, log *slog.Logger, host string) (*Client, error) {
	log.Info("Attempting SSH connection with socat fallback...")

	sshConn, err := establishSSHConnectionWithSocat(host, cfg.RemoteUser, log)
	if err != nil {
		return nil, fmt.Errorf("SSH connection with socat failed: %w", err)
	}

	cli, err := createDockerClientViaSSHProxy(sshConn, log)
	if err != nil {
		sshConn.Close()
		return nil, fmt.Errorf("SSH connection with socat failed: %w", err)
	}

	client := &Client{
		cli:     cli,
		cfg:     cfg,
		log:     log,
		sshConn: sshConn,
	}

	log.Info("Successfully connected to remote Docker via SSH (socat fallback)", "host", cfg.RemoteHost)
	return client, nil
}

// establishSSHConnection establishes an SSH connection to the remote host without socat
func establishSSHConnection(host, username string, log *slog.Logger) (*SSHContext, error) {
	log.Info("Attempting SSH connection without socat", "host", host, "user", username)
	log.Info("Note: SSH key-based authentication required (no password will be provided)")

	// Format the host with username for SSH connection
	var sshHost string
	if username != "" {
		sshHost = fmt.Sprintf("%s@%s", username, host)
	} else {
		sshHost = host
	}

	sshClient, err := NewSSHClient(sshHost, 22) // Default SSH port
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	// Connect without socat
	sshCtx, err := sshClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	log.Info("SSH connection established without socat", "host", host)
	return sshCtx, nil
}

// establishSSHConnectionWithSocat establishes an SSH connection to the remote host with socat
func establishSSHConnectionWithSocat(host, username string, log *slog.Logger) (*SSHConnection, error) {
	log.Info("Attempting SSH connection with socat", "host", host, "user", username)
	log.Info("Note: SSH key-based authentication required (no password will be provided)")

	// Format the host with username for SSH connection
	var sshHost string
	if username != "" {
		sshHost = fmt.Sprintf("%s@%s", username, host)
	} else {
		sshHost = host
	}

	sshClient, err := NewSSHClient(sshHost, 22) // Default SSH port
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	// Pass port 0 to trigger dynamic port finding in the SSH client
	sshConn, err := sshClient.ConnectWithSocat(0)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	localProxyHost := sshConn.GetLocalProxyHost()
	log.Info("SSH connection established with socat proxy", "host", localProxyHost)

	return sshConn, nil
}

// createDockerClientViaSSH creates a Docker client using direct SSH connection
func createDockerClientViaSSH(sshCtx *SSHContext, _ *slog.Logger) (*client.Client, error) {
	// For direct SSH connection, we need to use the ssh:// URL format
	// This requires the Docker client to have SSH support built-in
	sshURL := fmt.Sprintf("ssh://%s", sshCtx.remoteHost)

	opts := []client.Opt{
		client.WithHost(sshURL),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(30 * time.Second),
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client via SSH: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		cli.Close()
		return nil, fmt.Errorf("failed to connect to Docker via SSH: %w", err)
	}

	return cli, nil
}

// createDockerClientViaSSHProxy creates a Docker client using the SSH proxy
func createDockerClientViaSSHProxy(sshConn *SSHConnection, _ *slog.Logger) (*client.Client, error) {
	localProxyHost := sshConn.GetLocalProxyHost()

	opts := []client.Opt{
		client.WithHost(localProxyHost),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(30 * time.Second),
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client via SSH proxy: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		cli.Close()
		return nil, fmt.Errorf("failed to connect to Docker via SSH proxy: %w", err)
	}

	return cli, nil
}

// extractHostFromURL extracts the hostname from a Docker host URL
func extractHostFromURL(hostURL string) (string, error) {
	// Handle SSH URLs specially
	if strings.HasPrefix(hostURL, "ssh://") {
		// For SSH URLs, extract the hostname from ssh://[user@]host[:port]
		sshHost := hostURL[6:] // Remove "ssh://" prefix
		return extractSSHHostname(sshHost)
	}

	if strings.Contains(hostURL, "://") {
		parts := strings.Split(hostURL, "://")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid host URL format: %s", hostURL)
		}
		hostURL = parts[1]
	}

	// Extract hostname (remove port if present)
	if strings.Contains(hostURL, ":") {
		parts := strings.Split(hostURL, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid host:port format: %s", hostURL)
		}
		hostURL = parts[0]
	}

	// Basic hostname validation
	if hostURL == "" {
		return "", fmt.Errorf("hostname cannot be empty")
	}

	// Check for common invalid hostname patterns
	if strings.HasPrefix(hostURL, ".") || strings.HasSuffix(hostURL, ".") {
		return "", fmt.Errorf("hostname '%s' cannot start or end with a dot", hostURL)
	}

	if strings.Contains(hostURL, "..") {
		return "", fmt.Errorf("hostname '%s' cannot contain consecutive dots", hostURL)
	}

	return hostURL, nil
}

// extractSSHHostname extracts the hostname from an SSH host string
func extractSSHHostname(sshHost string) (string, error) {
	if sshHost == "" {
		return "", fmt.Errorf("SSH host cannot be empty")
	}

	// SSH host format: [user@]host[:port]
	hostname := sshHost

	// Remove username part if present
	if strings.Contains(hostname, "@") {
		parts := strings.Split(hostname, "@")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid SSH host format: expected [user@]host[:port]")
		}
		if parts[0] == "" {
			return "", fmt.Errorf("username cannot be empty in SSH host format")
		}
		hostname = parts[1]
	}

	// Remove port part if present
	if strings.Contains(hostname, ":") {
		parts := strings.Split(hostname, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid SSH host format: expected [user@]host[:port]")
		}
		if parts[1] == "" {
			return "", fmt.Errorf("port cannot be empty in SSH host format")
		}
		hostname = parts[0]
	}

	// Validate hostname
	if hostname == "" {
		return "", fmt.Errorf("hostname cannot be empty in SSH host format")
	}

	// Basic hostname validation
	if strings.HasPrefix(hostname, ".") || strings.HasSuffix(hostname, ".") {
		return "", fmt.Errorf("hostname '%s' cannot start or end with a dot", hostname)
	}

	if strings.Contains(hostname, "..") {
		return "", fmt.Errorf("hostname '%s' cannot contain consecutive dots", hostname)
	}

	return hostname, nil
}

// validateRemoteHost validates the format of a remote Docker host
func validateRemoteHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Check if user provided a scheme
	if strings.Contains(host, "://") {
		// Support both tcp:// and ssh:// schemes
		switch {
		case strings.HasPrefix(host, "ssh://"):
			// SSH URLs are handled by the Docker client natively
			// Format: ssh://[user@]host[:port]
			return validateSSHHost(host[6:]) // Remove "ssh://" prefix
		case strings.HasPrefix(host, "tcp://"):
			// Remove scheme for further validation
			host = host[6:] // "tcp://" is 6 characters
		default:
			return fmt.Errorf("only tcp:// and ssh:// schemes are supported (e.g., tcp://192.168.1.100 or ssh://user@host)")
		}
	}

	// Validate hostname part
	if host == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// If port is specified in host, validate it
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid host:port format")
		}

		// Validate port is numeric
		if parts[1] == "" {
			return fmt.Errorf("port cannot be empty")
		}
	}

	return nil
}

// validateSSHHost validates the format of an SSH host string
func validateSSHHost(host string) error {
	if host == "" {
		return fmt.Errorf("SSH host cannot be empty")
	}

	// SSH host format: [user@]host[:port]
	// Extract hostname part (remove user@ and :port)
	hostname := host

	// Remove username part if present
	if strings.Contains(hostname, "@") {
		parts := strings.Split(hostname, "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid SSH host format: expected [user@]host[:port]")
		}
		if parts[0] == "" {
			return fmt.Errorf("username cannot be empty in SSH host format")
		}
		hostname = parts[1]
	}

	// Remove port part if present
	if strings.Contains(hostname, ":") {
		parts := strings.Split(hostname, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid SSH host format: expected [user@]host[:port]")
		}
		hostname = parts[0]
		if parts[1] == "" {
			return fmt.Errorf("port cannot be empty in SSH host format")
		}
		// Validate port is numeric
		if _, err := fmt.Sscanf(parts[1], "%d", new(int)); err != nil {
			return fmt.Errorf("invalid port in SSH host format: port must be numeric")
		}
	}

	// Validate hostname
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty in SSH host format")
	}

	// Basic hostname validation
	if strings.HasPrefix(hostname, ".") || strings.HasSuffix(hostname, ".") {
		return fmt.Errorf("hostname '%s' cannot start or end with a dot", hostname)
	}

	if strings.Contains(hostname, "..") {
		return fmt.Errorf("hostname '%s' cannot contain consecutive dots", hostname)
	}

	return nil
}

// GetInfo retrieves Docker system information
func (c *Client) GetInfo(ctx context.Context) (map[string]any, error) {
	info, err := c.cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("docker info failed: %w", err)
	}
	return marshalToMap(info)
}

// InspectContainer inspects a container
func (c *Client) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	container, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("container inspect failed %s: %w", id, err)
	}

	return marshalToMap(container)
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, id string) (string, error) {
	logs, err := c.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
		Tail:       "100",
	})
	if err != nil {
		return "", fmt.Errorf("container logs failed %s: %w", id, err)
	}
	defer logs.Close()

	return c.readAndFormatLogs(logs), nil
}

// readAndFormatLogs reads logs from the response and formats them
func (c *Client) readAndFormatLogs(logs io.ReadCloser) string {
	var logLines []string
	buffer := make([]byte, 1024)

	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			line := string(buffer[:n])
			formattedLine := c.formatLogLine(line)
			logLines = append(logLines, formattedLine)
		}
		if err != nil {
			break
		}
	}

	return strings.Join(logLines, "")
}

// formatLogLine formats a single log line by removing timestamp prefix if present
func (c *Client) formatLogLine(line string) string {
	if len(line) >= 8 {
		return line[8:] // Remove timestamp prefix
	}
	return line
}

// InspectImage inspects an image
func (c *Client) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	imageInfo, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("image inspect failed %s: %w", id, err)
	}
	return marshalToMap(imageInfo)
}

// InspectVolume inspects a volume
func (c *Client) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	volumeInfo, err := c.cli.VolumeInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("volume inspect failed %s: %w", name, err)
	}
	return marshalToMap(volumeInfo)
}

// RemoveVolume removes a volume
func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	if err := validateID(name, "volume name"); err != nil {
		return err
	}

	if err := c.cli.VolumeRemove(ctx, name, force); err != nil {
		return fmt.Errorf("failed to remove volume %s: %w", name, err)
	}

	return nil
}

// InspectNetwork inspects a network
func (c *Client) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	networkInfo, err := c.cli.NetworkInspect(ctx, id, network.InspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("network inspect failed %s: %w", id, err)
	}
	return marshalToMap(networkInfo)
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	opts := container.ListOptions{
		All: all,
	}

	containers, err := c.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]Container, 0, len(containers))
	for i := range containers {
		cont := &containers[i]
		ports := formatContainerPorts(cont.Ports)

		result = append(result, Container{
			ID:      cont.ID[:12],
			Name:    cont.Names[0][1:], // Remove leading slash
			Image:   cont.Image,
			Status:  cont.Status,
			State:   cont.State,
			Created: time.Unix(cont.Created, 0),
			Ports:   ports,
			Size:    "", // Size is not available in ContainerList
		})
	}

	sortContainersByCreationTime(result)
	return result, nil
}

// formatContainerPorts formats container ports into a readable string
func formatContainerPorts(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}

	var formattedPorts []string
	for _, p := range ports {
		if p.PublicPort > 0 {
			formattedPorts = append(formattedPorts, fmt.Sprintf("%d:%d/%s", p.PublicPort, p.PrivatePort, p.Type))
		} else if p.PrivatePort > 0 {
			formattedPorts = append(formattedPorts, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}

	return strings.Join(formattedPorts, ", ")
}

// sortContainersByCreationTime sorts containers by creation time (newest first)
func sortContainersByCreationTime(containers []Container) {
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Created.After(containers[j].Created)
	})
}

// GetContainerStats gets container stats
func (c *Client) GetContainerStats(ctx context.Context, id string) (map[string]any, error) {
	stats, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer stats.Body.Close()

	return c.decodeStatsResponse(stats.Body)
}

// decodeStatsResponse decodes the stats response body into a map
func (c *Client) decodeStatsResponse(body io.ReadCloser) (map[string]any, error) {
	var statsJSON map[string]any
	if err := json.NewDecoder(body).Decode(&statsJSON); err != nil {
		return nil, fmt.Errorf("failed to decode container stats: %w", err)
	}
	return statsJSON, nil
}

// ListImages lists all images
func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	opts := image.ListOptions{}

	images, err := c.cli.ImageList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]Image, 0, len(images))
	for i := range images {
		img := &images[i]
		repo, tag := parseImageRepository(img.RepoTags)
		size := formatSize(img.Size)

		result = append(result, Image{
			ID:         img.ID[7:19], // Remove "sha256:" prefix and truncate
			Repository: repo,
			Tag:        tag,
			Size:       size,
			Created:    time.Unix(img.Created, 0),
			Containers: int(img.Containers),
		})
	}

	sortImagesByCreationTime(result)
	return result, nil
}

// parseImageRepository parses repository and tag from image repoTags
func parseImageRepository(repoTags []string) (repository, tag string) {
	if len(repoTags) == 0 || repoTags[0] == "<none>:<none>" {
		return "<none>", "<none>"
	}

	parts := strings.Split(repoTags[0], ":")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return repoTags[0], "<none>"
}

// sortImagesByCreationTime sorts images by creation time (newest first)
func sortImagesByCreationTime(images []Image) {
	sort.Slice(images, func(i, j int) bool {
		return images[i].Created.After(images[j].Created)
	})
}

// ListVolumes lists all volumes
func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	vols, err := c.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]Volume, 0, len(vols.Volumes))
	for _, vol := range vols.Volumes {
		volume := c.createVolumeFromAPI(vol)
		result = append(result, volume)
	}

	sortVolumesByName(result)
	return result, nil
}

// createVolumeFromAPI creates a Volume from the API response
func (c *Client) createVolumeFromAPI(vol *volume.Volume) Volume {
	created := time.Time{}
	if vol.CreatedAt != "" {
		created, _ = time.Parse(time.RFC3339, vol.CreatedAt)
	}

	return Volume{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Created:    created,
		Size:       "", // Size is not available in VolumeList
	}
}

// sortVolumesByName sorts volumes by name
func sortVolumesByName(volumes []Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
}

// ListNetworks lists all networks
func (c *Client) ListNetworks(ctx context.Context) ([]Network, error) {
	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]Network, 0, len(networks))
	for i := range networks {
		net := &networks[i]
		result = append(result, Network{
			ID:         net.ID[:12],
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Created:    net.Created,
			Internal:   net.Internal,
			Attachable: net.Attachable,
			Ingress:    net.Ingress,
			IPv6:       net.EnableIPv6,
			EnableIPv6: net.EnableIPv6,
			Labels:     net.Labels,
			Containers: len(net.Containers),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := buildStopOptions(timeout)
	return c.cli.ContainerStop(ctx, id, opts)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := buildStopOptions(timeout)
	return c.cli.ContainerRestart(ctx, id, opts)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	opts := container.RemoveOptions{
		Force: force,
	}

	if err := validateID(id, "container ID"); err != nil {
		return err
	}

	return c.cli.ContainerRemove(ctx, id, opts)
}

// ExecContainer executes a command in a running container and returns the output
func (c *Client) ExecContainer(ctx context.Context, id string, command []string, _ bool) (string, error) {
	if err := c.validateClient(); err != nil {
		return "", err
	}

	execResp, err := c.createExecInstance(ctx, id, command)
	if err != nil {
		return "", err
	}

	output, err := c.executeAndCollectOutput(ctx, execResp.ID)
	if err != nil {
		return "", err
	}

	return output, nil
}

// createExecInstance creates an exec instance in the container
func (c *Client) createExecInstance(ctx context.Context, id string, command []string) (*container.ExecCreateResponse, error) {
	execConfig := container.ExecOptions{
		Cmd:          command,
		Tty:          false, // Set to false to capture output
		AttachStdin:  false, // We don't need stdin for command execution
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
	}

	execResp, err := c.cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec instance: %w", err)
	}

	return &execResp, nil
}

// executeAndCollectOutput executes the exec instance and collects the output
func (c *Client) executeAndCollectOutput(ctx context.Context, execID string) (string, error) {
	attachConfig := container.ExecStartOptions{
		Tty: false,
	}

	hijackedResp, err := c.cli.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer hijackedResp.Close()

	if err := c.cli.ContainerExecStart(ctx, execID, attachConfig); err != nil {
		return "", fmt.Errorf("failed to start exec instance: %w", err)
	}

	return c.readExecOutput(hijackedResp), nil
}

// readExecOutput reads the output from the hijacked response
func (c *Client) readExecOutput(hijackedResp types.HijackedResponse) string {
	var output strings.Builder
	buffer := make([]byte, 1024)

	for {
		n, err := hijackedResp.Reader.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}

// AttachContainer attaches to a running container
func (c *Client) AttachContainer(ctx context.Context, id string) (types.HijackedResponse, error) {
	if err := c.validateClient(); err != nil {
		return types.HijackedResponse{}, err
	}

	attachConfig := container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}

	if err := validateID(id, "container ID"); err != nil {
		return types.HijackedResponse{}, err
	}

	response, err := c.cli.ContainerAttach(ctx, id, attachConfig)
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("failed to attach to container: %w", err)
	}

	return response, nil
}

// RemoveImage removes an image
func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	if err := validateID(id, "image ID"); err != nil {
		return err
	}

	opts := image.RemoveOptions{
		Force:         force,
		PruneChildren: true, // Remove dependent images by default
	}

	_, err := c.cli.ImageRemove(ctx, id, opts)
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %w", id, err)
	}

	return nil
}

// RemoveNetwork removes a network
func (c *Client) RemoveNetwork(ctx context.Context, id string) error {
	if err := validateID(id, "network ID"); err != nil {
		return err
	}

	if err := c.cli.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network %s: %w", id, err)
	}

	return nil
}

// buildStopOptions builds stop options with optional timeout
func buildStopOptions(timeout *time.Duration) container.StopOptions {
	opts := container.StopOptions{}
	if timeout != nil {
		opts.Signal = "SIGTERM"
		seconds := int(timeout.Seconds())
		opts.Timeout = &seconds
	}
	return opts
}

// validateID validates that an ID is not empty
func validateID(id, idType string) error {
	if id == "" {
		return fmt.Errorf("%s cannot be empty", idType)
	}
	return nil
}

// validateClient validates that the Docker client is initialized
func (c *Client) validateClient() error {
	if c.cli == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return nil
}
