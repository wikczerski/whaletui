package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
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

	assert.Equal(t, mockUI, executor.ui)
}
