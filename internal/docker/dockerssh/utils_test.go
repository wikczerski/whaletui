package dockerssh

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSSHHost(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		expectedUser  string
		expectedHost  string
		expectedPort  string
		errorContains string
		expectError   bool
	}{
		{
			name:         "UserAtHost",
			host:         "admin@192.168.1.100",
			expectedUser: "admin",
			expectedHost: "192.168.1.100",
			expectedPort: "22",
		},
		{
			name:         "UserAtHostWithPort",
			host:         "admin@192.168.1.100:2222",
			expectedUser: "admin",
			expectedHost: "192.168.1.100",
			expectedPort: "2222",
		},
		{
			name:         "HostOnly",
			host:         "192.168.1.100",
			expectedUser: "",
			expectedHost: "192.168.1.100",
			expectedPort: "22",
		},
		{
			name:         "HostWithPort",
			host:         "192.168.1.100:2222",
			expectedUser: "",
			expectedHost: "192.168.1.100",
			expectedPort: "2222",
		},
		{
			name:          "InvalidFormat",
			host:          "user@host@invalid",
			expectError:   true,
			errorContains: "invalid SSH host format",
		},
		{
			name:          "InvalidPortFormat",
			host:          "user@host:port:extra",
			expectError:   true,
			errorContains: "invalid SSH host format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, host, port, err := ParseSSHHost(tt.host)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				if tt.expectedUser != "" {
					assert.Equal(t, tt.expectedUser, user)
				}
				assert.Equal(t, tt.expectedHost, host)
				assert.Equal(t, tt.expectedPort, port)
			}
		})
	}
}

func TestGetSSHKeyPath(t *testing.T) {
	keyPath, err := getSSHKeyPath()

	if err != nil {
		assert.Contains(t, err.Error(), "no SSH private key found")
	} else {
		assert.NotEmpty(t, keyPath)
	}
}

func TestFindAvailablePort(t *testing.T) {
	port, err := findAvailablePort()
	require.NoError(t, err)

	assert.GreaterOrEqual(t, port, 2376)
	assert.LessOrEqual(t, port, 2390)

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(t, err)
	if closeErr := listener.Close(); closeErr != nil {
		t.Logf("Failed to close test listener: %v", closeErr)
	}
}
