package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			errorContains: []string{"lookup", "no such host"},
		},
		{
			name:          "InvalidHostnameStartingWithDot",
			host:          ".example.com",
			expectError:   true,
			errorContains: []string{"lookup", ".example.com", "no such host"},
		},
		{
			name:          "InvalidHostnameWithConsecutiveDots",
			host:          "example..com",
			expectError:   true,
			errorContains: []string{"lookup", "example..com", "no such host"},
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
