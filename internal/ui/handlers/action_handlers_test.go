package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	uimocks "github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func TestNewActionHandlers(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)
	assert.NotNil(t, handlers)
}

func TestNewActionHandlers_InitialState(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)
	assert.NotNil(t, handlers)
}

func TestActionHandlers_Constructor(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)

	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.ui)
}

func TestActionHandlers_UIInterface(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handlers := NewActionHandlers(mockUI)

	assert.Equal(t, mockUI, handlers.ui)
}
