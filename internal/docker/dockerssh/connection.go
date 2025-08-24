package dockerssh

import (
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

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

// Close closes the SSH connection and cleans up
func (s *SSHConnection) Close() error {
	if s == nil {
		return nil
	}

	var errors []string

	// Close the port forwarding listener if it exists
	if s.listener != nil {
		s.log.Info("Closing port forwarding listener")
		if err := s.listener.Close(); err != nil {
			s.log.Error("Failed to close port forwarding listener", "error", err)
			errors = append(errors, fmt.Sprintf("listener close: %v", err))
		} else {
			s.log.Info("Port forwarding listener closed successfully")
		}
	}

	if s.client != nil {
		s.log.Info("SSH connection closing")
		if err := s.client.Close(); err != nil {
			s.log.Error("Failed to close SSH client", "error", err)
			errors = append(errors, fmt.Sprintf("SSH client close: %v", err))
		} else {
			s.log.Info("SSH client closed successfully")
		}
	}

	if s.session != nil {
		if err := s.session.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("session close: %v", err))
		}
	}

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

	if err := s.killSocatByPID(); err != nil {
		s.log.Error("Failed to kill socat process by PID", "pid", s.socatPID, "error", err)
	}

	if err := s.checkPortStatus(); err != nil {
		s.log.Warn("Port status check failed", "error", err)
	}

	s.log.Info("Socat cleanup completed")
	return nil
}

// killSocatByPID kills the specific socat process by its PID and all its children
func (s *SSHConnection) killSocatByPID() error {
	if err := s.killChildProcesses(); err != nil {
		s.log.Warn("Failed to kill child processes", "error", err)
	}

	if err := s.killParentProcess(); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

// killChildProcesses kills all child processes of the socat process
func (s *SSHConnection) killChildProcesses() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	killTreeCmd := fmt.Sprintf("pkill -P %s", s.socatPID)
	s.log.Debug("Killing child processes", "command", killTreeCmd, "parent_pid", s.socatPID)

	_ = session.Run(killTreeCmd)
	time.Sleep(200 * time.Millisecond)
	return nil
}

// killParentProcess attempts graceful kill first, then force kill if needed
func (s *SSHConnection) killParentProcess() error {
	if err := s.attemptGracefulKill(); err != nil {
		s.log.Debug("Graceful kill failed, trying force kill", "error", err)
		return s.attemptForceKill()
	}

	s.log.Info("Graceful kill executed successfully", "pid", s.socatPID)
	return nil
}

// attemptGracefulKill attempts to gracefully kill the socat process
func (s *SSHConnection) attemptGracefulKill() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	killCmd := fmt.Sprintf("kill %s", s.socatPID)
	s.log.Debug("Executing graceful kill of socat parent process", "command", killCmd, "pid", s.socatPID)

	return session.Run(killCmd)
}

// attemptForceKill forcefully kills the socat process
func (s *SSHConnection) attemptForceKill() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create force kill session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	forceKillCmd := fmt.Sprintf("kill -9 %s", s.socatPID)
	s.log.Debug("Executing force kill of socat process", "command", forceKillCmd, "pid", s.socatPID)

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
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	checkCmd := fmt.Sprintf("lsof -i:%d || netstat -tlnp | grep ':%d ' || ss -tlnp | grep ':%d '", s.remotePort, s.remotePort, s.remotePort)
	s.log.Debug("Checking port status", "command", checkCmd, "port", s.remotePort)

	output, err := session.Output(checkCmd)
	if err != nil {
		s.log.Info("Port appears to be free after cleanup", "port", s.remotePort)
		return nil
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		s.log.Info("Port appears to be free after cleanup", "port", s.remotePort)
	} else {
		s.log.Warn("Port still appears to be in use after cleanup", "port", s.remotePort, "output", outputStr)
	}

	return nil
}
