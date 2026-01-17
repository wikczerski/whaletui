package cmd

import (
	"fmt"
	"runtime"
	"strings"
)

func (h *dockerErrorHandler) showGeneralHelp() {
	fmt.Println()
	fmt.Println("For more help:")
	fmt.Println("  • Check the logs: whaletui --log-level DEBUG")
	fmt.Println("  • View documentation: https://github.com/wikczerski/whaletui")
	fmt.Println("  • Report issues: https://github.com/wikczerski/whaletui/issues")
	fmt.Println()
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
