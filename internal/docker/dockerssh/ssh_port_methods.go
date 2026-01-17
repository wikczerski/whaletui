package dockerssh

import (
	"errors"
	"fmt"
	"time"
)

// createLocalTCPPortOnRemote creates a local TCP port on remote machine using socat/nc/netcat
func (s *SSHTunnelClient) createLocalTCPPortOnRemote() error {
	s.log.Info("Creating local TCP port on remote machine")

	// Find an available port on remote machine
	remotePort, err := s.findAvailableRemotePort()
	if err != nil {
		return fmt.Errorf("failed to find available remote port: %w", err)
	}

	s.remoteLocalPort = remotePort

	// Try to create TCP port using different methods
	return s.tryCreateTCPPortWithMethods(remotePort)
}

// tryCreateTCPPortWithMethods tries to create TCP port using different methods
func (s *SSHTunnelClient) tryCreateTCPPortWithMethods(remotePort int) error {
	methods := s.getPrimaryMethods(remotePort)

	// Try the primary methods first
	for _, method := range methods {
		if err := s.tryCreateTCPPort(method.name, method.cmd); err != nil {
			s.log.Warn("Failed to create TCP port with method",
				"method", method.name, "error", err)
			continue
		}

		s.log.Info("Successfully created local TCP port on remote",
			"method", method.name, "port", remotePort)
		s.connectionMethod = fmt.Sprintf("SSH Tunnel (%s)", method.name)
		return nil
	}

	// If primary methods fail, try fallback methods using built-in tools
	s.log.Info("Primary methods failed, trying fallback methods using built-in tools")
	return s.tryCreateTCPPortWithFallbackMethods(remotePort)
}

// getPrimaryMethods returns the primary methods for creating TCP ports
func (s *SSHTunnelClient) getPrimaryMethods(remotePort int) []struct {
	name string
	cmd  string
} {
	return []struct {
		name string
		cmd  string
	}{
		{
			"socat",
			fmt.Sprintf(
				"nohup socat TCP-LISTEN:%d,reuseaddr,fork "+
					"UNIX-CONNECT:/var/run/docker.sock >/dev/null 2>&1 &",
				remotePort),
		},
		{
			"nc",
			fmt.Sprintf(
				"nohup nc -l -p %d -e 'cat /var/run/docker.sock' >/dev/null 2>&1 &",
				remotePort),
		},
		{
			"netcat",
			fmt.Sprintf(
				"nohup netcat -l -p %d -e 'cat /var/run/docker.sock' >/dev/null 2>&1 &",
				remotePort),
		},
	}
}

// tryCreateTCPPortWithFallbackMethods tries to create TCP port using built-in tools
func (s *SSHTunnelClient) tryCreateTCPPortWithFallbackMethods(remotePort int) error {
	fallbackMethods := s.getFallbackMethods(remotePort)

	for _, method := range fallbackMethods {
		if err := s.tryCreateTCPPort(method.name, method.cmd); err != nil {
			s.log.Warn("Failed to create TCP port with fallback method",
				"method", method.name, "error", err)
			continue
		}

		s.log.Info("Successfully created local TCP port on remote using fallback method",
			"method", method.name,
			"port", remotePort)
		s.connectionMethod = fmt.Sprintf("SSH Tunnel (fallback %s)", method.name)
		return nil
	}

	return errors.New("failed to create TCP port with any fallback method")
}

// getFallbackMethods returns the fallback methods for creating TCP ports
//
//nolint:gocognit,function-length
func (s *SSHTunnelClient) getFallbackMethods(remotePort int) []struct { //nolint:gocognit,function-length
	name string
	cmd  string
} {
	return []struct {
		name string
		cmd  string
	}{
		{
			"python3",
			fmt.Sprintf(
				`nohup python3 -c "
import socket, os, sys, threading, signal, time
import sys

def handle_client(client_socket, addr):
    try:
        with open('/var/run/docker.sock', 'rb') as docker_sock:
            while True:
                data = client_socket.recv(4096)
                if not data:
                    break
                docker_sock.write(data)
                docker_sock.flush()
                response = docker_sock.read(4096)
                if response:
                    client_socket.send(response)
    except Exception:
        pass
    finally:
        client_socket.close()

def signal_handler(signum, frame):
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(('127.0.0.1', %d))
server.listen(1)

while True:
    try:
        client_socket, addr = server.accept()
        thread = threading.Thread(target=handle_client, args=(client_socket, addr))
        thread.daemon = True
        thread.start()
    except Exception:
        break
" >/dev/null 2>&1 &`, remotePort),
		},
		{
			"python",
			fmt.Sprintf(
				`nohup python -c "
import socket, os, sys, threading, signal, time
import sys

def handle_client(client_socket, addr):
    try:
        with open('/var/run/docker.sock', 'rb') as docker_sock:
            while True:
                data = client_socket.recv(4096)
                if not data:
                    break
                docker_sock.write(data)
                docker_sock.flush()
                response = docker_sock.read(4096)
                if response:
                    client_socket.send(response)
    except Exception:
        pass
    finally:
        client_socket.close()

def signal_handler(signum, frame):
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(('127.0.0.1', %d))
server.listen(1)

while True:
    try:
        client_socket, addr = server.accept()
        thread = threading.Thread(target=handle_client, args=(client_socket, addr))
        thread.daemon = True
        thread.start()
    except Exception:
        break
" >/dev/null 2>&1 &`, remotePort),
		},
		{
			"bash",
			fmt.Sprintf(
				`nohup bash -c '
set -e
exec 3<>/var/run/docker.sock
while true; do
    timeout 1 bash -c "exec 4<>/dev/tcp/127.0.0.1/%d" 2>/dev/null || break
    {
        while IFS= read -r line <&4; do
            echo "$line" >&3
            read -r response <&3
            echo "$response" >&4
        done
    } &
    wait
    exec 4<&-
    exec 4>&-
done
' >/dev/null 2>&1 &`, remotePort),
		},
	}
}

// tryCreateTCPPort attempts to create a TCP port using the specified method
func (s *SSHTunnelClient) tryCreateTCPPort(method, cmd string) error {
	if err := s.executeTCPPortCommand(method, cmd); err != nil {
		return err
	}

	// Give the process a moment to start
	time.Sleep(2 * time.Second)

	return s.verifyTCPPortCreation()
}

// executeTCPPortCommand executes the command to create the TCP port
func (s *SSHTunnelClient) executeTCPPortCommand(method, cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	// Execute the command to create the TCP port
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to execute %s command: %w", method, err)
	}

	return nil
}
