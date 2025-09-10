package dockerssh

import (
	"log/slog"
	"os"
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
			password: "testpass",
			wantErr:  false,
		},
		{
			name:     "valid client with both auth methods",
			host:     "localhost",
			port:     "22",
			user:     "testuser",
			keyPath:  "/path/to/key",
			password: "testpass",
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
			name:          "no custom key path",
			customKeyPath: "",
			wantErr:       false, // May succeed if default keys exist
		},
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
		})
	}
}
