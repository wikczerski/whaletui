package docker

import (
	"fmt"
	"strings"
)

// closeDockerClient closes the Docker client and logs the result
func (c *Client) closeDockerClient(errors *[]string) {
	if c.cli != nil {
		if err := c.cli.Close(); err != nil {
			c.log.Error("Failed to close Docker client", "error", err)
			*errors = append(*errors, fmt.Sprintf("Docker client close: %v", err))
		} else {
			c.log.Info("Docker client closed successfully")
		}
	}
}

// closeSSHConnection closes the SSH connection and logs the result
func (c *Client) closeSSHConnection(errors *[]string) {
	if c.sshConn != nil {
		c.log.Info("Closing SSH connection with socat")
		if err := c.sshConn.Close(); err != nil {
			c.log.Error("Failed to close SSH connection", "error", err)
			*errors = append(*errors, fmt.Sprintf("SSH connection close: %v", err))
		} else {
			c.log.Info("SSH connection closed successfully")
		}
	}
}

// closeSSHContext closes the SSH context and logs the result
func (c *Client) closeSSHContext(errors *[]string) {
	if c.sshCtx != nil {
		c.log.Info("Closing SSH context")
		if err := c.sshCtx.Close(); err != nil {
			c.log.Error("Failed to close SSH context", "error", err)
			*errors = append(*errors, fmt.Sprintf("SSH context close: %v", err))
		} else {
			c.log.Info("SSH context closed successfully")
		}
	}
}

// buildCloseError builds the final error message if any errors occurred
func (c *Client) buildCloseError(errors []string) error {
	if len(errors) > 0 {
		return fmt.Errorf("failed to close client: %s", strings.Join(errors, "; "))
	}
	return nil
}
