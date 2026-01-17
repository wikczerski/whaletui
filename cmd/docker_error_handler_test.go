package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
)

func TestDockerErrorHandlerCreation(t *testing.T) {
	// Test that dockerErrorHandler can be created
	cfg := &config.Config{}
	err := errors.New("test error")

	handler := newDockerErrorHandler(err, cfg, nil)
	assert.NotNil(t, handler)
	assert.Equal(t, err, handler.err)
	assert.Equal(t, cfg, handler.cfg)
}

func TestDockerErrorHandlerHelperFunctions(t *testing.T) {
	// Test helper functions that don't require user interaction
	handler := &dockerErrorHandler{
		err:         errors.New("test error"),
		cfg:         &config.Config{},
		interaction: UserInteraction{},
	}

	// Test hasLogFile
	assert.True(t, handler.hasLogFile("/path/to/log"))
	assert.False(t, handler.hasLogFile(""))

	// Test canReadLogFile
	assert.True(t, handler.canReadLogFile(nil))
	assert.False(t, handler.canReadLogFile(errors.New("read error")))

	// Test isValidStartIndex
	assert.True(t, handler.isValidStartIndex(5))
	assert.False(t, handler.isValidStartIndex(-1))

	// Test isValidLogLine
	assert.True(t, handler.isValidLogLine("some log line"))
	assert.False(t, handler.isValidLogLine(""))

	// Test calculateLogStartIndex
	assert.Equal(t, 0, handler.calculateLogStartIndex(10)) // 10 - 20 = -10, should return 0
	assert.Equal(t, 5, handler.calculateLogStartIndex(25)) // 25 - 20 = 5, should return 5
}

func TestDockerErrorHandlerConnectionTypeDetection(t *testing.T) {
	// Test remote vs local connection detection
	localCfg := &config.Config{}
	remoteCfg := &config.Config{RemoteHost: "192.168.1.100"}

	localHandler := &dockerErrorHandler{cfg: localCfg}
	remoteHandler := &dockerErrorHandler{cfg: remoteCfg}

	assert.False(t, localHandler.isRemoteConnection())
	assert.True(t, remoteHandler.isRemoteConnection())
}

func TestHandleDockerConnectionError(t *testing.T) {
	// Test the main error handling function
	cfg := &config.Config{}
	err := errors.New("docker client creation failed")

	// This function would normally interact with the user, so we just test it doesn't panic
	// In a real scenario, this would be tested with integration tests
	assert.NotPanics(t, func() {
		_ = handleDockerConnectionError(err, cfg)
	})
}
