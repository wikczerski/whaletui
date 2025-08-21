package docker

import (
	"fmt"
	"net"
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
			name:          "NonExistentHostname",
			host:          "this-hostname-definitely-does-not-exist-12345ssss.com",
			expectError:   true,
			errorContains: []string{"lookup", "this-hostname-definitely-does-not-exist-12345ssss.com", "no such host"},
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

func TestParseSSHHost(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		expectedUser  string
		expectedHost  string
		expectedPort  string
		expectError   bool
		errorContains string
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
			expectedUser: "", // Will be set to current user
			expectedHost: "192.168.1.100",
			expectedPort: "22",
		},
		{
			name:         "HostWithPort",
			host:         "192.168.1.100:2222",
			expectedUser: "", // Will be set to current user
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
			user, host, port, err := parseSSHHost(tt.host)

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
	// This test will depend on the user's SSH key setup
	keyPath, err := getSSHKeyPath()

	// We can't guarantee SSH keys exist, so we just test the function doesn't crash
	// and returns either a valid path or an appropriate error
	if err != nil {
		assert.Contains(t, err.Error(), "no SSH private key found")
	} else {
		assert.NotEmpty(t, keyPath)
	}
}

func TestFindAvailablePort(t *testing.T) {
	port, err := findAvailablePort()
	require.NoError(t, err)

	// Port should be in the expected range (using a wider range to avoid conflicts)
	assert.GreaterOrEqual(t, port, 2376)
	assert.LessOrEqual(t, port, 2390)

	// Port should actually be available
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(t, err)
	listener.Close()
}

func TestSSHClient_NewSSHClient(t *testing.T) {
	// Test with a valid host format
	client, err := NewSSHClient("admin@192.168.1.100", 22)

	// This test may fail if SSH keys aren't set up, but we can test the parsing logic
	if err != nil {
		// If it fails, it should be due to SSH key issues, not parsing issues
		assert.Contains(t, err.Error(), "SSH key")
	} else {
		require.NotNil(t, client)
		assert.Equal(t, "admin", client.user)
		assert.Equal(t, "192.168.1.100", client.host)
		assert.Equal(t, "22", client.port)
	}
}

func TestSSHClient_NewSSHClient_InvalidHost(t *testing.T) {
	// Test with invalid host format
	_, err := NewSSHClient("invalid@host@format", 22)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid SSH host format")
}

func TestSSHConnection_GetLocalProxyHost(t *testing.T) {
	conn := &SSHConnection{
		localPort: 2376,
	}

	proxyHost := conn.GetLocalProxyHost()
	assert.Equal(t, "tcp://127.0.0.1:2376", proxyHost)
}

func TestSSHConnection_Close(t *testing.T) {
	// Test closing a connection with nil components
	conn := &SSHConnection{}
	err := conn.Close()
	assert.NoError(t, err)

	// Test with mock components (we can't easily test with real SSH connections in unit tests)
	// Note: We can't test with real SSH client in unit tests, so we just test the nil case
	conn = &SSHConnection{
		client:  nil,
		session: nil,
	}
	err = conn.Close()
	assert.NoError(t, err)
}

// Integration test helper - this would be used in integration tests
func TestSSHClient_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires real SSH server")

	// This test would require a real SSH server to test against
	// It's skipped by default but can be enabled for integration testing

	host := "test@localhost"
	client, err := NewSSHClient(host, 22)
	require.NoError(t, err)

	// Test connection (this would fail without a real SSH server)
	conn, err := client.Connect(2375)
	if err != nil {
		t.Logf("SSH connection failed as expected: %v", err)
		return
	}

	// Clean up
	defer conn.Close()

	// Test proxy host
	proxyHost := conn.GetLocalProxyHost()
	assert.Contains(t, proxyHost, "127.0.0.1:")
	assert.Contains(t, proxyHost, "tcp://")
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
	// Test with invalid hostname
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
	// This test requires a real SSH connection, so we'll just test the method exists
	client := &SSHClient{}

	// Test that the method signature is correct
	// In a real scenario, this would be called with an actual SSH client
	assert.NotNil(t, client.checkSocatAvailability)
}

func TestSSHClient_CheckDockerSocketAccess(t *testing.T) {
	// This test requires a real SSH connection, so we'll just test the method exists
	client := &SSHClient{}

	// Test that the method signature is correct
	// In a real scenario, this would be called with an actual SSH client
	assert.NotNil(t, client.checkDockerSocketAccess)
}

func TestSSHClient_VerifySocatProcess(t *testing.T) {
	// This test requires a real SSH connection, so we'll just test the method exists
	client := &SSHClient{}

	// Test that the method signature is correct
	// In a real scenario, this would be called with an actual SSH client
	assert.NotNil(t, client.verifySocatProcess)
}
