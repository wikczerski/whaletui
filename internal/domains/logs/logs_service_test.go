package logs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedmocks "github.com/wikczerski/whaletui/internal/mocks/shared"
)

func TestNewLogsService(t *testing.T) {
	mockContainerService := sharedmocks.NewMockContainerService(t)
	logsService := NewLogsService(mockContainerService)

	assert.NotNil(t, logsService)
}

func TestLogsService_GetLogs_Container(t *testing.T) {
	mockContainerService := sharedmocks.NewMockContainerService(t)
	mockContainerService.EXPECT().
		GetContainerLogs(context.Background(), "test-container-id").
		Return("test logs", nil)

	logsService := NewLogsService(mockContainerService)

	ctx := context.Background()

	// Test container logs
	logs, err := logsService.GetLogs(ctx, "container", "test-container-id")

	assert.NoError(t, err)
	assert.Equal(t, "test logs", logs)
}

func TestLogsService_GetLogs_UnsupportedResourceType(t *testing.T) {
	mockContainerService := sharedmocks.NewMockContainerService(t)
	logsService := NewLogsService(mockContainerService)

	ctx := context.Background()

	// Test unsupported resource type
	logs, err := logsService.GetLogs(ctx, "unsupported", "test-id")

	assert.Error(t, err)
	assert.Equal(t, "", logs)
	assert.Contains(t, err.Error(), "unsupported resource type")
}

func TestLogsService_GetActions(t *testing.T) {
	mockContainerService := sharedmocks.NewMockContainerService(t)
	logsService := NewLogsService(mockContainerService)

	actions := logsService.GetActions()

	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	assert.Equal(t, expectedActions, actions)
}

func TestLogsService_GetActionsString(t *testing.T) {
	mockContainerService := sharedmocks.NewMockContainerService(t)
	logsService := NewLogsService(mockContainerService)

	actionsString := logsService.GetActionsString()

	expectedString := "<f> Follow logs\n<t> Tail logs\n<s> Save logs\n<c> Clear logs\n<w> Wrap text"
	assert.Equal(t, expectedString, actionsString)
}
