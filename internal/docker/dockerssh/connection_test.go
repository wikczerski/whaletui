package dockerssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/logger"
)

func TestSSHConnection_GetLocalProxyHost(t *testing.T) {
	conn := &SSHConnection{
		localPort: 2376,
	}

	proxyHost := conn.GetLocalProxyHost()
	assert.Equal(t, "tcp://127.0.0.1:2376", proxyHost)
}

func TestSSHConnection_Close(t *testing.T) {
	conn := &SSHConnection{
		log: logger.GetLogger(),
	}
	err := conn.Close()
	assert.NoError(t, err)

	conn = &SSHConnection{
		client:  nil,
		session: nil,
		log:     logger.GetLogger(),
	}
	err = conn.Close()
	assert.NoError(t, err)
}

func TestSSHContext_Close(t *testing.T) {
	conn := &SSHContext{
		client: nil,
		log:    logger.GetLogger(),
	}
	err := conn.Close()
	assert.NoError(t, err)
}

func TestSSHConnection_GetLocalProxyHost_Nil(t *testing.T) {
	var conn *SSHConnection
	proxyHost := conn.GetLocalProxyHost()
	assert.Equal(t, "tcp://127.0.0.1:0", proxyHost)
}

func TestSSHConnection_GetLocalProxyHost_InvalidPort(t *testing.T) {
	conn := &SSHConnection{
		localPort: -1,
	}

	proxyHost := conn.GetLocalProxyHost()
	assert.Equal(t, "tcp://127.0.0.1:2375", proxyHost)
}
