package dockerssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSSHClient_Integration is an integration test that requires a real SSH server
func TestSSHClient_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires real SSH server")

	host := "test@localhost"
	client, err := NewSSHClient(host, 22)
	require.NoError(t, err)

	conn, err := client.Connect()
	if err != nil {
		t.Logf("SSH connection failed as expected: %v", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("Failed to close SSH connection: %v", err)
		}
	}()

	assert.NotNil(t, conn)
	assert.Equal(t, "localhost", conn.remoteHost)
}
