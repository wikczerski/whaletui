package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/logger"
)

type dockerErrorHandler struct {
	err         error
	cfg         *config.Config
	interaction UserInteraction
}

func newDockerErrorHandler(err error, cfg *config.Config) *dockerErrorHandler {
	return &dockerErrorHandler{
		err:         err,
		cfg:         cfg,
		interaction: UserInteraction{},
	}
}

func (h *dockerErrorHandler) handle() error {
	h.prepareTerminal()
	h.showErrorHeader()

	h.handleErrorDetails()
	h.handleConnectionType()

	h.showGeneralHelp()
	h.showLogOption()
	h.waitForExit()
	h.showRemoteOption()

	return fmt.Errorf("docker connection failed: %w", h.err)
}

// handleErrorDetails handles the error details display
func (h *dockerErrorHandler) handleErrorDetails() {
	if h.askForDetails() {
		h.showDetailedError()
	}
}

// handleConnectionType handles different connection types
func (h *dockerErrorHandler) handleConnectionType() {
	if h.isRemoteConnection() {
		h.showRemoteHelp()
	} else {
		h.showLocalHelp()
		if h.askForRetry() {
			if err := h.attemptRetry(); err != nil {
				fmt.Fprintf(os.Stderr, "Retry attempt failed: %v\n", err)
			}
		}
	}
}

func (h *dockerErrorHandler) prepareTerminal() {
	logger.SetTUIMode(false)
	if _, err := fmt.Fprint(os.Stdout, "\033[2J\033[H"); err != nil {
		// Log the error but continue since this is just terminal clearing
		fmt.Fprintf(os.Stderr, "Failed to clear terminal: %v\n", err)
	}
}

func (h *dockerErrorHandler) showErrorHeader() {
	fmt.Println("❌ Docker Connection Failed")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Printf("Error details: %v\n", h.err)
	fmt.Println()
}

func (h *dockerErrorHandler) askForDetails() bool {
	return h.interaction.askYesNo("Show detailed error information?")
}

func (h *dockerErrorHandler) askForRetry() bool {
	return h.interaction.askYesNo("Would you like to retry?")
}

func (h *dockerErrorHandler) isRemoteConnection() bool {
	return h.cfg.RemoteHost != ""
}

func (h *dockerErrorHandler) showDetailedError() {
	fmt.Println()
	fmt.Println("Detailed error information:")
	fmt.Println("==========================")
	fmt.Printf("Error type: %T\n", h.err)
	fmt.Printf("Error message: %s\n", h.err.Error())

	h.showExitErrorDetails()
	h.showSpecificErrorGuidance()
	fmt.Println()
}

func (h *dockerErrorHandler) showExitErrorDetails() {
	if dockerErr, ok := h.err.(*exec.ExitError); ok {
		fmt.Printf("Exit code: %d\n", dockerErr.ExitCode())
		if len(dockerErr.Stderr) > 0 {
			fmt.Printf("Stderr: %s\n", string(dockerErr.Stderr))
		}
	}
}

func (h *dockerErrorHandler) showSpecificErrorGuidance() {
	errStr := h.err.Error()
	switch {
	case strings.Contains(errStr, "permission denied"):
		h.showPermissionErrorHelp()
	case strings.Contains(errStr, "connection refused"):
		h.showConnectionRefusedHelp()
	case strings.Contains(errStr, "timeout"):
		h.showTimeoutErrorHelp()
	}
}

func (h *dockerErrorHandler) showPermissionErrorHelp() {
	fmt.Println()
	fmt.Println("Permission denied error detected:")
	fmt.Println("• Check if your user has access to Docker")
	fmt.Println("• Try running: sudo usermod -aG docker $USER")
	fmt.Println("• Log out and log back in after adding to docker group")
}

func (h *dockerErrorHandler) showConnectionRefusedHelp() {
	fmt.Println()
	fmt.Println("Connection refused error detected:")
	fmt.Println("• Docker daemon is not listening on the expected port/socket")
	fmt.Println("• Check if Docker is running")
	fmt.Println("• Verify Docker socket permissions")
}

func (h *dockerErrorHandler) showTimeoutErrorHelp() {
	fmt.Println()
	fmt.Println("Timeout error detected:")
	fmt.Println("• Docker daemon might be overloaded")
	fmt.Println("• Check system resources")
	fmt.Println("• Try increasing Docker daemon timeout")
}

func (h *dockerErrorHandler) showRemoteHelp() {
	h.showRemoteHelpHeader()
	h.showRemoteHelpChecklist()
	h.showRemoteHelpSuggestions()
}

// showRemoteHelpHeader shows the remote help header
func (h *dockerErrorHandler) showRemoteHelpHeader() {
	fmt.Printf("Unable to connect to remote Docker host: %s\n", h.cfg.RemoteHost)
	fmt.Println()
}

// showRemoteHelpChecklist shows the remote help checklist
func (h *dockerErrorHandler) showRemoteHelpChecklist() {
	fmt.Println("Please check:")
	fmt.Println("• The remote host is accessible")
	fmt.Println("• Docker daemon is running on the remote host")
	fmt.Println("• SSH connection is working (if using SSH)")
	fmt.Println("• Firewall settings allow the connection")
	fmt.Println("• Port 2375/2376 is open (for TCP connections)")
	fmt.Println()
}

// showRemoteHelpSuggestions shows the remote help suggestions
func (h *dockerErrorHandler) showRemoteHelpSuggestions() {
	fmt.Println("You can try:")
	fmt.Printf("  whaletui connect --host %s --user <username>\n", h.cfg.RemoteHost)
	fmt.Println("  • Test SSH connection: ssh <username>@<host>")
	fmt.Println("  • Test Docker connection: docker -H <host> ps")
}

func (h *dockerErrorHandler) showLocalHelp() {
	fmt.Println("Unable to connect to local Docker daemon")
	fmt.Println()

	h.showDockerStatus()
	h.showGeneralChecks()
	h.showOSSpecificChecks()
	h.showLocalSuggestions()
}

func (h *dockerErrorHandler) showDockerStatus() {
	if isDockerRunning() {
		fmt.Println("⚠️  Docker appears to be running but connection failed")
		fmt.Println("This might be a permission issue or Docker socket problem")
	} else {
		fmt.Println("🚫 Docker daemon is not running")
	}
}

func (h *dockerErrorHandler) showGeneralChecks() {
	fmt.Println()
	fmt.Println("Please check:")
	fmt.Println("• Docker Desktop is running (Windows/macOS)")
	fmt.Println("• Docker daemon is running (Linux)")
	fmt.Println("• You have permission to access Docker")
	fmt.Println("• Docker socket is accessible")
}

func (h *dockerErrorHandler) showOSSpecificChecks() {
	switch runtime.GOOS {
	case "windows":
		h.showWindowsChecks()
	case "linux":
		h.showLinuxChecks()
	case "darwin":
		h.showMacOSChecks()
	}
}

func (h *dockerErrorHandler) showWindowsChecks() {
	fmt.Println()
	fmt.Println("Windows-specific checks:")
	fmt.Println("• Docker Desktop is installed and running")
	fmt.Println("• WSL2 is enabled and running (if using Linux containers)")
	fmt.Println("• Docker Desktop has finished starting up")
	fmt.Println("• Check Docker Desktop system tray icon")
	fmt.Println("• Try restarting Docker Desktop")
}

func (h *dockerErrorHandler) showLinuxChecks() {
	fmt.Println()
	fmt.Println("Linux-specific checks:")
	fmt.Println("• Docker daemon is running: sudo systemctl status docker")
	fmt.Println("• Docker socket permissions: ls -la /var/run/docker.sock")
	fmt.Println("• User is in docker group: groups $USER")
	fmt.Println("• Try: sudo systemctl start docker")
}

func (h *dockerErrorHandler) showMacOSChecks() {
	fmt.Println()
	fmt.Println("macOS-specific checks:")
	fmt.Println("• Docker Desktop is installed and running")
	fmt.Println("• Docker Desktop has finished starting up")
	fmt.Println("• Check Docker Desktop menu bar icon")
	fmt.Println("• Try restarting Docker Desktop")
}

func (h *dockerErrorHandler) showLocalSuggestions() {
	fmt.Println()
	fmt.Println("You can try:")
	fmt.Println("  • Starting Docker Desktop")
	fmt.Println("  • Running 'docker ps' to test connection")
	fmt.Println("  • Checking Docker service: sudo systemctl status docker")
	fmt.Println("  • Connecting to a remote host: whaletui connect --host <host> --user <username>")
}

func (h *dockerErrorHandler) attemptRetry() error {
	fmt.Println("Retrying connection...")
	time.Sleep(2 * time.Second)

	application, retryErr := app.New(h.cfg)
	if retryErr == nil {
		fmt.Println("✅ Connection successful! Starting whaletui...")
		time.Sleep(1 * time.Second)
		return runApplicationWithShutdown(application, logger.GetLogger())
	}

	fmt.Printf("❌ Retry failed: %v\n", retryErr)
	h.interaction.waitForEnter()
	return retryErr
}

func (h *dockerErrorHandler) showGeneralHelp() {
	fmt.Println()
	fmt.Println("For more help:")
	fmt.Println("  • Check the logs: whaletui --log-level DEBUG")
	fmt.Println("  • View documentation: https://github.com/wikczerski/whaletui")
	fmt.Println("  • Report issues: https://github.com/wikczerski/whaletui/issues")
	fmt.Println()
}

func (h *dockerErrorHandler) showLogOption() {
	logFilePath := logger.GetLogFilePath()
	if h.hasLogFile(logFilePath) {
		fmt.Printf("Recent logs are available at: %s\n", logFilePath)
		if h.interaction.askYesNo("Would you like to view recent logs?") {
			h.showRecentLogs(logFilePath)
		}
	}
}

func (h *dockerErrorHandler) hasLogFile(logFilePath string) bool {
	return logFilePath != ""
}

func (h *dockerErrorHandler) showRecentLogs(logFilePath string) {
	h.showRecentLogsHeader()
	h.validateAndReadLogFile(logFilePath)
}

// showRecentLogsHeader shows the recent logs header
func (h *dockerErrorHandler) showRecentLogsHeader() {
	fmt.Println()
	fmt.Println("Recent logs:")
	fmt.Println("============")
}

// validateAndReadLogFile validates and reads the log file
func (h *dockerErrorHandler) validateAndReadLogFile(logFilePath string) {
	if !h.isValidLogFilePath(logFilePath) {
		h.showInvalidLogPathMessage()
		return
	}

	h.readAndDisplayLogFile(logFilePath)
	fmt.Println()
}

// isValidLogFilePath checks if the log file path is valid
func (h *dockerErrorHandler) isValidLogFilePath(logFilePath string) bool {
	// Clean the path to remove any directory traversal attempts
	cleanPath := filepath.Clean(logFilePath)

	// Ensure it's an absolute path
	if !filepath.IsAbs(cleanPath) {
		return false
	}

	// Additional security: check for suspicious patterns
	if strings.Contains(cleanPath, "..") || strings.Contains(cleanPath, "~") {
		return false
	}

	// Ensure the cleaned path matches the original (after cleaning)
	return cleanPath == filepath.Clean(logFilePath)
}

// showInvalidLogPathMessage shows the invalid log path message
func (h *dockerErrorHandler) showInvalidLogPathMessage() {
	fmt.Println("Invalid log file path")
	fmt.Println()
}

// readAndDisplayLogFile reads and displays the log file
func (h *dockerErrorHandler) readAndDisplayLogFile(logFilePath string) {
	// nolint:gosec // Path is validated by isValidLogFilePath before this function is called
	logContent, readErr := os.ReadFile(logFilePath)
	if h.canReadLogFile(readErr) {
		h.displayLastLogLines(string(logContent))
	} else {
		fmt.Printf("Could not read log file: %v\n", readErr)
	}
}

func (h *dockerErrorHandler) canReadLogFile(err error) bool {
	return err == nil
}

func (h *dockerErrorHandler) displayLastLogLines(content string) {
	lines := strings.Split(content, "\n")
	start := h.calculateLogStartIndex(len(lines))

	for i := start; i < len(lines); i++ {
		if h.isValidLogLine(lines[i]) {
			fmt.Println(lines[i])
		}
	}
}

func (h *dockerErrorHandler) calculateLogStartIndex(totalLines int) int {
	const maxLogLines = 20
	start := totalLines - maxLogLines
	if h.isValidStartIndex(start) {
		return start
	}
	return 0
}

func (h *dockerErrorHandler) isValidStartIndex(start int) bool {
	return start >= 0
}

func (h *dockerErrorHandler) isValidLogLine(line string) bool {
	return line != ""
}

func (h *dockerErrorHandler) waitForExit() {
	h.interaction.waitForEnter()
}

func (h *dockerErrorHandler) showRemoteOption() {
	fmt.Println()
	if h.interaction.askYesNo("Would you like to try connecting to a remote host instead?") {
		h.showRemoteConnectionExamples()
	}
}

func (h *dockerErrorHandler) showRemoteConnectionExamples() {
	fmt.Println()
	fmt.Println("To connect to a remote Docker host:")
	fmt.Println("  whaletui connect --host <host> --user <username>")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  whaletui connect --host 192.168.1.100 --user admin")
	fmt.Println("  whaletui connect --host ssh://admin@192.168.1.100")
	fmt.Println("  whaletui connect --host tcp://192.168.1.100:2375")
	fmt.Println()
	h.interaction.waitForEnter()
}
