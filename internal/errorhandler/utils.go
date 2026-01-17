package errorhandler

import (
	"os/exec"
)

// IsDockerRunning checks if Docker is running on the local system
func IsDockerRunning() bool {
	// Try to run a simple Docker command to check if it's accessible
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run() == nil
}
