package errorhandler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
)

type mockInteractor struct{}

func (m *mockInteractor) AskYesNo(question string) bool { return false }
func (m *mockInteractor) WaitForEnter()                 {}

func TestDockerErrorHandlerCreation(t *testing.T) {
	// Test that dockerErrorHandler can be created
	cfg := &config.Config{}
	err := errors.New("test error")

	handler := NewDockerErrorHandler(err, cfg, nil, &mockInteractor{})
	assert.NotNil(t, handler)
	assert.Equal(t, err, handler.err)
	assert.Equal(t, cfg, handler.cfg)
}

func TestDockerErrorHandlerHelperFunctions(t *testing.T) {
	// Test helper functions that don't require user interaction
	handler := &DockerErrorHandler{
		err:         errors.New("test error"),
		cfg:         &config.Config{},
		interaction: &mockInteractor{},
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

	localHandler := &DockerErrorHandler{cfg: localCfg}
	remoteHandler := &DockerErrorHandler{cfg: remoteCfg}

	assert.False(t, localHandler.isRemoteConnection())
	assert.True(t, remoteHandler.isRemoteConnection())
}
