package dockerssh

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestNewSSHClientWithAuth(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		host     string
		port     string
		user     string
		keyPath  string
		password string
		wantErr  bool
	}{
		{
			name:     "valid client with key path",
			host:     "localhost",
			port:     "22",
			user:     "testuser",
			keyPath:  "/path/to/key",
			password: "",
			wantErr:  false,
		},
		{
			name:     "valid client with password",
			host:     "localhost",
			port:     "22",
			user:     "testuser",
			keyPath:  "",
			password: getTestPassword(),
			wantErr:  false,
		},
		{
			name:     "valid client with both auth methods",
			host:     "localhost",
			port:     "22",
			user:     "testuser",
			keyPath:  "/path/to/key",
			password: getTestPassword(),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSSHClientWithAuth(tt.host, tt.port, tt.user, tt.keyPath, tt.password, log)

			if client == nil {
				t.Fatal("Expected client to be created, got nil")
			}

			if client.host != tt.host {
				t.Errorf("Expected host %s, got %s", tt.host, client.host)
			}

			if client.port != tt.port {
				t.Errorf("Expected port %s, got %s", tt.port, client.port)
			}

			if client.user != tt.user {
				t.Errorf("Expected user %s, got %s", tt.user, client.user)
			}

			if client.keyPath != tt.keyPath {
				t.Errorf("Expected keyPath %s, got %s", tt.keyPath, client.keyPath)
			}

			if client.password != tt.password {
				t.Errorf("Expected password %s, got %s", tt.password, client.password)
			}
		})
	}
}

func TestAddSSHKeyAuth(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name          string
		customKeyPath string
		wantErr       bool
	}{
		{
			name:          "custom key path that doesn't exist",
			customKeyPath: "/nonexistent/key",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ssh.ClientConfig{
				Auth: []ssh.AuthMethod{},
			}

			err := addSSHKeyAuth(config, log, tt.customKeyPath)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify that the error message is appropriate
			if err != nil {
				if !strings.Contains(err.Error(), "SSH key not found at specified path") {
					t.Errorf("Expected 'SSH key not found at specified path' error, got: %v", err)
				}
			}
		})
	}
}

func TestAddSSHKeyAuth_NoCustomKeyPath(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	config := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{},
	}

	err := addSSHKeyAuth(config, log, "")
	// This test may pass or fail depending on whether SSH keys exist on the system
	// In CI environments, it will typically fail with "no SSH keys found"
	// On local development machines with SSH keys, it may succeed
	// We just verify that the function doesn't panic and returns a reasonable result
	if err != nil {
		// If there's an error, it should be about missing SSH keys
		if !strings.Contains(err.Error(), "no SSH keys found") &&
			!strings.Contains(err.Error(), "SSH key not found") {
			t.Errorf("Unexpected error type: %v", err)
		}
	}
	// If there's no error, that's also acceptable (SSH keys were found)
}

func getTestPassword() string {
	if pwd := os.Getenv("SSH_TEST_PASSWORD"); pwd != "" {
		return pwd
	}
	return "dummy-test-password-for-unit-tests"
}
