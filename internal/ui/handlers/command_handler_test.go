package handlers

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wikczerski/D5r/internal/config"
	uimocks "github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func TestNewCommandHandler(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUI, handler.ui)
	assert.False(t, handler.isActive)
	assert.Nil(t, handler.commandInput)
}

func TestCreateCommandInput(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	themeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(themeManager)

	handler := NewCommandHandler(mockUI)
	input := handler.CreateCommandInput()

	assert.NotNil(t, input)
	assert.Equal(t, handler.commandInput, input)
	assert.Equal(t, ": ", input.GetLabel())

	mockUI.AssertExpectations(t)
}

func TestHideCommandInput(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	// Test with nil commandInput
	handler.hideCommandInput()

	// Test with initialized commandInput
	handler.commandInput = tview.NewInputField()
	handler.hideCommandInput()

	// Verify input is hidden (though we can't easily test the visual state)
	assert.NotNil(t, handler.commandInput)
}

func TestShowCommandInput(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	themeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(themeManager)

	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()

	handler.showCommandInput()

	mockUI.AssertExpectations(t)
}

func TestEnter(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockMainFlex := &tview.Flex{}
	mockApp := &tview.Application{}
	themeManager := config.NewThemeManager("")

	mockUI.On("GetMainFlex").Return(mockMainFlex)
	mockUI.On("GetApp").Return(mockApp)
	mockUI.On("GetThemeManager").Return(themeManager)

	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()

	handler.Enter()

	assert.True(t, handler.isActive)
	mockUI.AssertExpectations(t)
}

func TestExit(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockMainFlex := &tview.Flex{}
	mockApp := &tview.Application{}
	themeManager := config.NewThemeManager("")

	mockUI.On("GetMainFlex").Return(mockMainFlex)
	mockUI.On("GetViewContainer").Return(nil) // Return nil to avoid focus issues
	mockUI.On("GetApp").Return(mockApp).Maybe()
	mockUI.On("GetThemeManager").Return(themeManager).Maybe()

	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	handler.Exit()

	assert.False(t, handler.isActive)
	assert.Equal(t, "", handler.commandInput.GetText())
	mockUI.AssertExpectations(t)
}

func TestIsActive(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	assert.False(t, handler.IsActive())

	handler.isActive = true
	assert.True(t, handler.IsActive())
}

func TestGetInput(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	assert.Nil(t, handler.GetInput())

	input := tview.NewInputField()
	handler.commandInput = input
	assert.Equal(t, input, handler.GetInput())
}

func TestHandleInput(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockMainFlex := &tview.Flex{}
	themeManager := config.NewThemeManager("")

	mockUI.On("GetMainFlex").Return(mockMainFlex).Maybe()
	mockUI.On("GetViewContainer").Return(nil).Maybe()
	mockUI.On("GetThemeManager").Return(themeManager).Maybe()
	mockUI.On("ShowError", mock.AnythingOfType("*errors.errorString")).Maybe()

	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	// Test Escape key
	handler.HandleInput(tcell.KeyEscape)
	assert.False(t, handler.isActive)

	// Test Enter key with empty command
	handler.isActive = true
	handler.HandleInput(tcell.KeyEnter)
	// Empty commands now keep the input open (don't exit)
	assert.True(t, handler.isActive)

	mockUI.AssertExpectations(t)
}

func TestProcessCommand(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.isActive = true

	// Set up expectations for methods called by Exit()
	mockUI.On("GetViewContainer").Return(nil).Maybe()

	// Test view switching commands
	mockUI.On("SwitchView", "containers").Once()
	result := handler.processCommand("containers")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "images").Once()
	result = handler.processCommand("i")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "volumes").Once()
	result = handler.processCommand("v")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "networks").Once()
	result = handler.processCommand("networks")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	// Test quit command
	handler.isActive = true
	shutdownChan := make(chan struct{}, 1)
	mockUI.On("GetShutdownChan").Return(shutdownChan)
	result = handler.processCommand("quit")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	// Test help command
	handler.isActive = true
	mockUI.On("ShowHelp").Once()
	result = handler.processCommand("help")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	// Test unknown command
	handler.isActive = true
	// Note: Unknown commands now show error in command input instead of calling ShowError
	result = handler.processCommand("unknown")
	assert.False(t, result) // Should return false for unknown commands
	// The command input will stay active for 2 seconds to show the error message
	// We can't easily test the timing behavior in unit tests, so we just verify it doesn't call ShowError
	mockUI.AssertNotCalled(t, "ShowError")

	mockUI.AssertExpectations(t)
}

func TestGetAutocomplete(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	// Test empty input
	suggestions := handler.getAutocomplete("")
	assert.Len(t, suggestions, 9)
	assert.Contains(t, suggestions, "containers")
	assert.Contains(t, suggestions, "images")
	assert.Contains(t, suggestions, "volumes")
	assert.Contains(t, suggestions, "networks")
	assert.Contains(t, suggestions, "quit")
	assert.Contains(t, suggestions, "q")
	assert.Contains(t, suggestions, "exit")
	assert.Contains(t, suggestions, "help")
	assert.Contains(t, suggestions, "?")

	// Test partial input
	suggestions = handler.getAutocomplete("c")
	assert.Len(t, suggestions, 1)
	assert.Contains(t, suggestions, "containers")

	// Test case insensitive
	suggestions = handler.getAutocomplete("C")
	assert.Len(t, suggestions, 1)
	assert.Contains(t, suggestions, "containers")

	// Test no matches
	suggestions = handler.getAutocomplete("xyz")
	assert.Len(t, suggestions, 0)
}

func TestShowCommandError(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	// Set up theme manager mock with proper initialization
	themeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(themeManager)

	// Test showing error message
	handler.showCommandError("Test error message")

	// Verify the error message is displayed
	assert.Equal(t, "Test error message", handler.commandInput.GetText())

	// Verify the text color is set to red
	// Note: We can't easily test the color in unit tests, but we can verify the text is set
	assert.Equal(t, "Test error message", handler.commandInput.GetText())

	mockUI.AssertExpectations(t)
}

func TestClearError(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	// Set up theme manager mock
	themeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(themeManager).Maybe()

	// First show an error message
	handler.showCommandError("Test error message")
	assert.Equal(t, "Test error message", handler.commandInput.GetText())
	assert.NotNil(t, handler.errorTimer) // Timer should be set

	// Now clear the error
	handler.clearError()
	assert.Nil(t, handler.errorTimer) // Timer should be cleared

	mockUI.AssertExpectations(t)
}
