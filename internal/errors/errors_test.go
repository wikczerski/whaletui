package errors

import (
	"errors"
	"testing"
)

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

	err := ConnectionError(host, originalErr)

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("ConnectionError() should return *DockerError")
	}

	expectedOp := "connection to localhost:2375"
	if dockerErr.Operation != expectedOp {
		t.Errorf("ConnectionError() operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("ConnectionError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestContainerError(t *testing.T) {
	originalErr := errors.New("container not found")
	operation := "start"
	id := "abc123"

	err := ContainerError(operation, id, originalErr)

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("ContainerError() should return *DockerError")
	}

	expectedOp := "start container abc123"
	if dockerErr.Operation != expectedOp {
		t.Errorf("ContainerError() operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("ContainerError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestImageError(t *testing.T) {
	originalErr := errors.New("image not found")
	operation := "pull"
	id := "nginx:latest"

	err := ImageError(operation, id, originalErr)

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("ImageError() should return *DockerError")
	}

	expectedOp := "pull image nginx:latest"
	if dockerErr.Operation != expectedOp {
		t.Errorf("ImageError() operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("ImageError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestVolumeError(t *testing.T) {
	originalErr := errors.New("volume not found")
	operation := "create"
	name := "my-volume"

	err := VolumeError(operation, name, originalErr)

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("VolumeError() should return *DockerError")
	}

	expectedOp := "create volume my-volume"
	if dockerErr.Operation != expectedOp {
		t.Errorf("VolumeError() operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("VolumeError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("network not found")
	operation := "connect"
	id := "bridge"

	err := NetworkError(operation, id, originalErr)

	dockerErr, ok := err.(*DockerError)
	if !ok {
		t.Fatal("NetworkError() should return *DockerError")
	}

	expectedOp := "connect network bridge"
	if dockerErr.Operation != expectedOp {
		t.Errorf("NetworkError() operation = %v, want %v", dockerErr.Operation, expectedOp)
	}

	if dockerErr.Err != originalErr {
		t.Errorf("NetworkError() error = %v, want %v", dockerErr.Err, originalErr)
	}
}
