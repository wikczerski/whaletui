// Package errors provides error handling and custom error types for WhaleTUI.
package errors

import "fmt"

// DockerError represents a Docker-related error
type DockerError struct {
	Operation string
	Err       error
}

// NewDockerError creates a new Docker error
func NewDockerError(operation string, err error) *DockerError {
	return &DockerError{
		Operation: operation,
		Err:       err,
	}
}

func (e *DockerError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func (e *DockerError) Unwrap() error {
	return e.Err
}

// ConnectionError creates a connection error
func ConnectionError(host string, err error) error {
	return NewDockerError(fmt.Sprintf("connection to %s", host), err)
}

// ContainerError creates a container operation error
func ContainerError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s container %s", operation, id), err)
}

// ImageError creates an image operation error
func ImageError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s image %s", operation, id), err)
}

// VolumeError creates a volume operation error
func VolumeError(operation, name string, err error) error {
	return NewDockerError(fmt.Sprintf("%s volume %s", operation, name), err)
}

// NetworkError creates a network operation error
func NetworkError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s network %s", operation, id), err)
}

// ConfigError represents a configuration-related error
type ConfigError struct {
	Message string
}

// NewConfigError creates a new configuration error
func NewConfigError(message string) *ConfigError {
	return &ConfigError{Message: message}
}

func (e *ConfigError) Error() string {
	return e.Message
}

// UIError creates a UI-related error
func UIError(operation string, err error) error {
	return NewDockerError(fmt.Sprintf("UI %s", operation), err)
}
