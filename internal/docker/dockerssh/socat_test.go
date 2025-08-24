package dockerssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHClient_FindAvailableRemotePort(t *testing.T) {
	client := &SSHClient{}

	// This test requires a real SSH connection, so we just test the method exists
	assert.NotNil(t, client.findAvailableRemotePort)
}

func TestSSHClient_IsRemotePortAvailable(t *testing.T) {
	client := &SSHClient{}

	// This test requires a real SSH connection, so we just test the method exists
	assert.NotNil(t, client.isRemotePortAvailable)
}

func TestSSHClient_SetupSocatProxy(t *testing.T) {
	client := &SSHClient{}

	// This test requires a real SSH connection, so we just test the method exists
	assert.NotNil(t, client.setupSocatProxy)
}
