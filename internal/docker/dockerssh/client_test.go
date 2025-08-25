package dockerssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHClient_validateHostname(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		expectError   bool
		errorContains []string
	}{
		{
			name:        "ValidIPAddress",
			host:        "192.168.1.100",
			expectError: false,
		},
		{
			name:        "ValidIPv6Address",
			host:        "::1",
			expectError: false,
		},
		{
			name:        "ValidHostname",
			host:        "localhost",
			expectError: false,
		},
		{
			name:        "ValidDomainName",
			host:        "example.com",
			expectError: false,
		},
		{
			name:          "EmptyHostname",
			host:          "",
			expectError:   true,
			errorContains: []string{"hostname cannot be empty"},
		},
		{
			name:          "InvalidHostnameStartingWithDot",
			host:          ".example.com",
			expectError:   true,
			errorContains: []string{"cannot start or end with a dot"},
		},
		{
			name:          "InvalidHostnameWithConsecutiveDots",
			host:          "example..com",
			expectError:   true,
			errorContains: []string{"cannot contain consecutive dots"},
		},
		{
			name:        "NonExistentHostname",
			host:        "this-hostname-definitely-does-not-exist-12345ssss.com",
			expectError: true,
			errorContains: []string{
				"lookup",
				"this-hostname-definitely-does-not-exist-12345ssss.com",
				"no such host",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &SSHClient{host: tt.host}
			err := client.validateHostname()

			if tt.expectError {
				assert.Error(t, err)
				if len(tt.errorContains) > 0 {
					for _, expectedText := range tt.errorContains {
						assert.Contains(t, err.Error(), expectedText)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSSHClient_NewSSHClient(t *testing.T) {
	client, err := NewSSHClient("admin@192.168.1.100", 22)

	if err != nil {
		assert.Contains(t, err.Error(), "SSH key")
	} else {
		require.NotNil(t, client)
		assert.Equal(t, "admin", client.user)
		assert.Equal(t, "192.168.1.100", client.host)
		assert.Equal(t, "22", client.port)
	}
}

func TestSSHClient_NewSSHClient_InvalidHost(t *testing.T) {
	_, err := NewSSHClient("invalid@host@format", 22)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid SSH host format")
}

func TestSSHClient_GetConnectionInfo(t *testing.T) {
	client := &SSHClient{
		user: "admin",
		host: "192.168.1.100",
		port: "22",
	}

	info := client.GetConnectionInfo()
	assert.Equal(t, "SSH Client: admin@192.168.1.100:22", info)
}

func TestSSHClient_DiagnoseConnection(t *testing.T) {
	client := &SSHClient{
		user: "admin",
		host: "invalid-hostname-that-does-not-exist-12345.com",
		port: "22",
	}

	err := client.DiagnoseConnection()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Hostname resolution")
}

func TestSSHClient_CheckSocatAvailability(t *testing.T) {
	client := &SSHClient{}
	assert.NotNil(t, client.checkSocatAvailability)
}

func TestSSHClient_CheckDockerSocketAccess(t *testing.T) {
	client := &SSHClient{}
	assert.NotNil(t, client.checkDockerSocketAccess)
}

func TestSSHClient_VerifySocatProcess(t *testing.T) {
	client := &SSHClient{}
	assert.NotNil(t, client.verifySocatProcess)
}
