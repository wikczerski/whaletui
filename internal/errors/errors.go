package errors

import "fmt"

// DockerError represents a Docker-related error
type DockerError struct {
	Operation string
	Err       error
}

func (e *DockerError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func (e *DockerError) Unwrap() error {
	return e.Err
}

// NewDockerError creates a new Docker error
func NewDockerError(operation string, err error) *DockerError {
	return &DockerError{
		Operation: operation,
		Err:       err,
	}
}

// Connection errors
func ConnectionError(host string, err error) error {
	return NewDockerError(fmt.Sprintf("connection to %s", host), err)
}

// Operation errors
func ContainerError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s container %s", operation, id), err)
}

func ImageError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s image %s", operation, id), err)
}

func VolumeError(operation, name string, err error) error {
	return NewDockerError(fmt.Sprintf("%s volume %s", operation, name), err)
}

func NetworkError(operation, id string, err error) error {
	return NewDockerError(fmt.Sprintf("%s network %s", operation, id), err)
}
