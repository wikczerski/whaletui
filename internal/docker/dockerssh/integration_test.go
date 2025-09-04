package dockerssh

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSSHClient_Integration is an integration test that requires a real SSH server
func TestSSHClient_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires real SSH server")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewSSHClient("localhost", "22", "test", logger)

	// Test that client was created successfully
	assert.NotNil(t, client)
	assert.Equal(t, "test", client.user)
	assert.Equal(t, "localhost", client.host)
	assert.Equal(t, "22", client.port)
}
