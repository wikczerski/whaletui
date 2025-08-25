package handlers

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wikczerski/whaletui/internal/config"
	uimocks "github.com/wikczerski/whaletui/internal/mocks/ui"
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

	handler.hideCommandInput()

	handler.commandInput = tview.NewInputField()
	handler.hideCommandInput()

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
	mockUI.On("GetViewContainer").Return(nil)
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

	handler.HandleInput(tcell.KeyEscape)
	assert.False(t, handler.isActive)

	handler.isActive = true
	handler.HandleInput(tcell.KeyEnter)
	assert.True(t, handler.isActive)

	mockUI.AssertExpectations(t)
}

func TestProcessCommand(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.isActive = true

	mockUI.On("GetViewContainer").Return(nil).Maybe()

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

	handler.isActive = true
	mockUI.On("SwitchView", "swarmServices").Once()
	result = handler.processCommand("swarm services")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "swarmServices").Once()
	result = handler.processCommand("services")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "swarmNodes").Once()
	result = handler.processCommand("swarm nodes")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("SwitchView", "swarmNodes").Once()
	result = handler.processCommand("nodes")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	shutdownChan := make(chan struct{}, 1)
	mockUI.On("GetShutdownChan").Return(shutdownChan)
	result = handler.processCommand("quit")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	mockUI.On("ShowHelp").Once()
	result = handler.processCommand("help")
	assert.True(t, result)
	assert.False(t, handler.isActive)

	handler.isActive = true
	result = handler.processCommand("unknown")
	assert.False(t, result)
	mockUI.AssertNotCalled(t, "ShowError")

	mockUI.AssertExpectations(t)
}

func TestGetAutocomplete(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)

	suggestions := handler.getAutocomplete("")
	assert.Len(t, suggestions, 13)
	assert.Contains(t, suggestions, "containers")
	assert.Contains(t, suggestions, "images")
	assert.Contains(t, suggestions, "volumes")
	assert.Contains(t, suggestions, "networks")
	assert.Contains(t, suggestions, "swarm services")
	assert.Contains(t, suggestions, "swarm nodes")
	assert.Contains(t, suggestions, "services")
	assert.Contains(t, suggestions, "nodes")
	assert.Contains(t, suggestions, "quit")
	assert.Contains(t, suggestions, "q")
	assert.Contains(t, suggestions, "exit")
	assert.Contains(t, suggestions, "help")
	assert.Contains(t, suggestions, "?")

	suggestions = handler.getAutocomplete("c")
	assert.Len(t, suggestions, 1)
	assert.Contains(t, suggestions, "containers")

	suggestions = handler.getAutocomplete("C")
	assert.Len(t, suggestions, 1)
	assert.Contains(t, suggestions, "containers")

	suggestions = handler.getAutocomplete("s")
	assert.Len(t, suggestions, 3)
	assert.Contains(t, suggestions, "swarm services")
	assert.Contains(t, suggestions, "swarm nodes")
	assert.Contains(t, suggestions, "services")

	suggestions = handler.getAutocomplete("swarm")
	assert.Len(t, suggestions, 2)
	assert.Contains(t, suggestions, "swarm services")
	assert.Contains(t, suggestions, "swarm nodes")

	suggestions = handler.getAutocomplete("xyz")
	assert.Len(t, suggestions, 0)
}

func TestShowCommandError(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	themeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(themeManager)

	handler.showCommandError("Test error message")

	assert.Equal(t, "Test error message", handler.commandInput.GetText())

	assert.Equal(t, "Test error message", handler.commandInput.GetText())

	mockUI.AssertExpectations(t)
}

func TestClearError(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	handler := NewCommandHandler(mockUI)
	handler.commandInput = tview.NewInputField()
	handler.isActive = true

	themeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(themeManager).Maybe()

	handler.showCommandError("Test error message")
	assert.Equal(t, "Test error message", handler.commandInput.GetText())
	assert.NotNil(t, handler.errorTimer)

	handler.clearError()
	assert.Nil(t, handler.errorTimer)

	mockUI.AssertExpectations(t)
}
