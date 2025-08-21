package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/models"
)

// MockContainerService for testing
type mockContainerService struct{}

func (m *mockContainerService) ListContainers(_ context.Context) ([]models.Container, error) {
	return []models.Container{}, nil
}

func (m *mockContainerService) StartContainer(_ context.Context, _ string) error {
	return nil
}

func (m *mockContainerService) StopContainer(_ context.Context, _ string, _ *time.Duration) error {
	return nil
}

func (m *mockContainerService) RestartContainer(_ context.Context, _ string, _ *time.Duration) error {
	return nil
}

func (m *mockContainerService) RemoveContainer(_ context.Context, _ string, _ bool) error {
	return nil
}

func (m *mockContainerService) GetContainerLogs(_ context.Context, _ string) (string, error) {
	return "test logs", nil
}

func (m *mockContainerService) InspectContainer(_ context.Context, _ string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (m *mockContainerService) ExecContainer(_ context.Context, _ string, _ []string, _ bool) (string, error) {
	return "", nil
}

func (m *mockContainerService) AttachContainer(_ context.Context, _ string) (any, error) {
	return nil, nil
}

func (m *mockContainerService) GetActions() map[rune]string {
	return map[rune]string{}
}

func (m *mockContainerService) GetActionsString() string {
	return ""
}

func TestNewLogsService(t *testing.T) {
	mockContainerService := &mockContainerService{}
	logsService := NewLogsService(mockContainerService)

	assert.NotNil(t, logsService)
}

func TestLogsService_GetLogs_Container(t *testing.T) {
	mockContainerService := &mockContainerService{}
	logsService := NewLogsService(mockContainerService)

	ctx := context.Background()

	// Test container logs
	logs, err := logsService.GetLogs(ctx, "container", "test-container-id")

	assert.NoError(t, err)
	assert.Equal(t, "test logs", logs)
}

func TestLogsService_GetLogs_UnsupportedResourceType(t *testing.T) {
	mockContainerService := &mockContainerService{}
	logsService := NewLogsService(mockContainerService)

	ctx := context.Background()

	// Test unsupported resource type
	logs, err := logsService.GetLogs(ctx, "unsupported", "test-id")

	assert.Error(t, err)
	assert.Equal(t, "", logs)
	assert.Contains(t, err.Error(), "unsupported resource type")
}

func TestLogsService_GetActions(t *testing.T) {
	mockContainerService := &mockContainerService{}
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
	mockContainerService := &mockContainerService{}
	logsService := NewLogsService(mockContainerService)

	actionsString := logsService.GetActionsString()

	expectedString := "<f> Follow logs\n<t> Tail logs\n<s> Save logs\n<c> Clear logs\n<w> Wrap text"
	assert.Equal(t, expectedString, actionsString)
}
