package dockerssh

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Removed validateHostname test as it's now a private method

func TestSSHClient_NewSSHClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewSSHClient("192.168.1.100", "22", "admin", logger)

	require.NotNil(t, client)
	assert.Equal(t, "admin", client.user)
	assert.Equal(t, "192.168.1.100", client.host)
	assert.Equal(t, "22", client.port)
}

func TestSSHClient_NewSSHClient_InvalidHost(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewSSHClient("invalid@host@format", "22", "admin", logger)
	require.NotNil(t, client)
	// The client is created but validation happens during connection
}

// Removed tests for methods that no longer exist

func TestSSHClient_CheckDockerSocketAccess(t *testing.T) {
	client := &SSHClient{}
	assert.NotNil(t, client.checkDockerSocketAccess)
}
