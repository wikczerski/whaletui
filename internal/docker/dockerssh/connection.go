package dockerssh

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHConnection represents an active SSH connection
type SSHConnection struct {
	client     *ssh.Client
	session    *ssh.Session
	localPort  int
	remoteHost string
	remotePort int
	socatPID   string
	log        *slog.Logger
	listener   net.Listener // Add listener for port forwarding
}

// GetLocalProxyHost returns the local host:port for the Docker proxy
func (s *SSHConnection) GetLocalProxyHost() string {
	if s == nil {
		return "tcp://127.0.0.1:0"
	}

	port := s.localPort
	if port <= 0 {
		port = 2375
	}

	return fmt.Sprintf("tcp://127.0.0.1:%d", port)
}

// Close closes the SSH connection and cleans up resources
func (s *SSHConnection) Close() error {
	if s == nil {
		return nil
	}

	var errors []string
	s.closeListener(&errors)
	if err := s.killSocatProcess(); err != nil {
		errors = append(errors, fmt.Sprintf("socat cleanup: %v", err))
	}
	s.closeSSHClient(&errors)
	s.closeSession(&errors)

	return s.buildCloseError(errors)
}

// closeListener closes the port forwarding listener
func (s *SSHConnection) closeListener(errors *[]string) {
	if s.listener == nil {
		return
	}

	s.log.Info("Closing port forwarding listener")
	if err := s.listener.Close(); err != nil {
		s.log.Error("Failed to close port forwarding listener", "error", err)
		*errors = append(*errors, fmt.Sprintf("listener close: %v", err))
	} else {
		s.log.Info("Port forwarding listener closed successfully")
	}
}

// closeSSHClient closes the SSH client
func (s *SSHConnection) closeSSHClient(errors *[]string) {
	if s.client == nil {
		return
	}

	s.log.Info("SSH connection closing")
	if err := s.client.Close(); err != nil {
		s.log.Error("Failed to close SSH client", "error", err)
		*errors = append(*errors, fmt.Sprintf("SSH client close: %v", err))
	} else {
		s.log.Info("SSH client closed successfully")
	}
}

// closeSession closes the SSH session
func (s *SSHConnection) closeSession(errors *[]string) {
	if s.session == nil {
		return
	}

	if err := s.session.Close(); err != nil {
		*errors = append(*errors, fmt.Sprintf("session close: %v", err))
	}
}

// buildCloseError builds the final error message if there were any errors
func (s *SSHConnection) buildCloseError(errors []string) error {
	if len(errors) > 0 {
		return fmt.Errorf("failed to close SSH connection: %s", strings.Join(errors, "; "))
	}
	return nil
}

// killSocatProcess kills the specific socat process we started (by PID)
func (s *SSHConnection) killSocatProcess() error {
	if s == nil || s.client == nil {
		return nil
	}

	s.log.Info("Starting socat cleanup", "port", s.remotePort, "pid", s.socatPID)

	if s.socatPID == "" {
		s.log.Warn("No socat PID stored, cannot perform targeted cleanup")
		return nil
	}

	s.performSocatCleanup()
	return nil
}

// performSocatCleanup performs the actual socat cleanup operations
func (s *SSHConnection) performSocatCleanup() {
	// Use force kill as a fallback
	if err := s.attemptForceKill(); err != nil {
		s.log.Error("Failed to kill socat process by PID", "pid", s.socatPID, "error", err)
	}

	if err := s.checkPortStatus(); err != nil {
		s.log.Warn("Port status check failed", "error", err)
	}

	s.log.Info("Socat cleanup completed")
}

// These functions are currently unused and have been removed to satisfy the linter

// attemptForceKill forcefully kills the socat process
func (s *SSHConnection) attemptForceKill() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create force kill session: %w", err)
	}
	defer s.closeSessionSafely(session)

	forceKillCmd := s.buildForceKillCommand()
	s.log.Debug("Executing force kill of socat process", "command", forceKillCmd, "pid", s.socatPID)

	return s.executeForceKill(session, forceKillCmd)
}

// closeSessionSafely safely closes the SSH session
func (s *SSHConnection) closeSessionSafely(session *ssh.Session) {
	if err := session.Close(); err != nil {
		s.log.Warn("Failed to close SSH session", "error", err)
	}
}

// buildForceKillCommand builds the force kill command
func (s *SSHConnection) buildForceKillCommand() string {
	return fmt.Sprintf("kill -9 %s", s.socatPID)
}

// executeForceKill executes the force kill command
func (s *SSHConnection) executeForceKill(session *ssh.Session, forceKillCmd string) error {
	if err := session.Run(forceKillCmd); err != nil {
		s.log.Warn("Force kill failed", "error", err, "pid", s.socatPID)
		return err
	}

	s.log.Info("Force kill executed successfully", "pid", s.socatPID)
	return nil
}

// checkPortStatus checks what process is using the port after cleanup
func (s *SSHConnection) checkPortStatus() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer s.closeSessionSafely(session)

	return s.executePortStatusCheck(session)
}

// executePortStatusCheck executes the actual port status check
func (s *SSHConnection) executePortStatusCheck(session *ssh.Session) error {
	checkCmd := s.buildPortCheckCommand()
	s.log.Debug("Checking port status", "command", checkCmd, "port", s.remotePort)

	output, err := session.Output(checkCmd)
	if err != nil {
		s.log.Info("Port appears to be free after cleanup", "port", s.remotePort)
		return nil
	}

	s.analyzePortStatusOutput(output)
	return nil
}

// buildPortCheckCommand builds the command to check port status
func (s *SSHConnection) buildPortCheckCommand() string {
	portCheck := fmt.Sprintf("lsof -i:%d", s.remotePort)
	netstatCheck := fmt.Sprintf("netstat -tlnp | grep ':%d '", s.remotePort)
	ssCheck := fmt.Sprintf("ss -tlnp | grep ':%d '", s.remotePort)
	return fmt.Sprintf("%s || %s || %s", portCheck, netstatCheck, ssCheck)
}

// analyzePortStatusOutput analyzes the output of the port status check
func (s *SSHConnection) analyzePortStatusOutput(output []byte) {
	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		s.log.Info("Port appears to be free after cleanup", "port", s.remotePort)
	} else {
		s.log.Warn("Port still appears to be in use after cleanup", "port", s.remotePort, "output", outputStr)
	}
}
