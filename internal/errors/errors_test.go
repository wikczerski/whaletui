package errors

import (
	"errors"
	"testing"
)

// assertDockerError is a helper function to assert DockerError properties
func assertDockerError(t *testing.T, err error, expectedOp string, originalErr error) {
	t.Helper()

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("Error should return *DockerError")
	}

	if dockerErr.Operation != expectedOp {
		t.Errorf("operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestDockerError_Error(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		err       error
		expected  string
	}{
		{
			name:      "basic error",
			operation: "start",
			err:       errors.New("container not found"),
			expected:  "start failed: container not found",
		},
		{
			name:      "empty operation",
			operation: "",
			err:       errors.New("test error"),
			expected:  " failed: test error",
		},
		{
			name:      "nil error",
			operation: "stop",
			err:       nil,
			expected:  "stop failed: <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerErr := &DockerError{
				Operation: tt.operation,
				Err:       tt.err,
			}
			result := dockerErr.Error()
			if result != tt.expected {
				t.Errorf("DockerError.Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDockerError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	dockerErr := &DockerError{
		Operation: "test",
		Err:       originalErr,
	}

	unwrapped := dockerErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("DockerError.Unwrap() = %v, want %v", unwrapped, originalErr)
	}
}

func TestNewDockerError(t *testing.T) {
	originalErr := errors.New("test error")
	dockerErr := NewDockerError("test operation", originalErr)

	if dockerErr.Operation != "test operation" {
		t.Errorf("NewDockerError() operation = %v, want %v", dockerErr.Operation, "test operation")
	}
	if dockerErr.Err != originalErr {
		t.Errorf("NewDockerError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestConnectionError(t *testing.T) {
	originalErr := errors.New("connection refused")
	host := "localhost:2375"
	expectedOp := "connection to localhost:2375"

	err := ConnectionError(host, originalErr)
	assertDockerError(t, err, expectedOp, originalErr)
}

func TestContainerError(t *testing.T) {
	originalErr := errors.New("container not found")
	operation := "start"
	id := "abc123"
	expectedOp := "start container abc123"

	err := ContainerError(operation, id, originalErr)
	assertDockerError(t, err, expectedOp, originalErr)
}

func TestImageError(t *testing.T) {
	originalErr := errors.New("image not found")
	operation := "pull"
	id := "nginx:latest"
	expectedOp := "pull image nginx:latest"

	err := ImageError(operation, id, originalErr)
	assertDockerError(t, err, expectedOp, originalErr)
}

func TestVolumeError(t *testing.T) {
	originalErr := errors.New("volume not found")
	operation := "create"
	name := "my-volume"
	expectedOp := "create volume my-volume"

	err := VolumeError(operation, name, originalErr)
	assertDockerError(t, err, expectedOp, originalErr)
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("network not found")
	operation := "connect"
	id := "bridge"
	expectedOp := "connect network bridge"

	err := NetworkError(operation, id, originalErr)
	assertDockerError(t, err, expectedOp, originalErr)
}
