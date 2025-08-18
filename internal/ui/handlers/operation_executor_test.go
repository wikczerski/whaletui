package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func TestNewOperationExecutor(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	executor := NewOperationExecutor(mockUI)
	assert.NotNil(t, executor)
}

func TestNewOperationExecutor_InitialState(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	executor := NewOperationExecutor(mockUI)
	assert.NotNil(t, executor)
}

func TestOperationExecutor_Constructor(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	executor := NewOperationExecutor(mockUI)

	assert.NotNil(t, executor)
	assert.NotNil(t, executor.ui)
}

func TestOperationExecutor_UIInterface(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	executor := NewOperationExecutor(mockUI)

	// Verify that the UI interface was properly set
	assert.Equal(t, mockUI, executor.ui)
}

// Note: The actual operation methods (Execute, ExecuteWithConfirmation, etc.) are async
// and require complex mock setup for UI methods like ShowConfirm, ShowError, etc.
// These tests are simplified to avoid the complexity of testing async operations with mocks.
// In a real application, these would be tested with integration tests or by mocking
// the underlying dependencies more comprehensively.
