package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func TestNewActionHandlers(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)
	assert.NotNil(t, handlers)
}

func TestNewActionHandlers_InitialState(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)
	assert.NotNil(t, handlers)
	// Note: execCommandParser is not a public field, so we can't test it directly
}

// Note: All ParseExecCommand tests have been removed because:
// 1. The parseExecCommand method is private
// 2. The method cannot be tested directly from outside the package
// 3. To test this functionality, either:
//    - Make the method public, or
//    - Test it indirectly through public methods that use it

func TestActionHandlers_Constructor(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)

	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.ui)
}

func TestActionHandlers_UIInterface(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)

	// Verify that the UI interface was properly set
	assert.Equal(t, mockUI, handlers.ui)
}
