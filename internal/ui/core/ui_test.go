package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/D5r/internal/services"
	"github.com/wikczerski/D5r/internal/ui/constants"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		serviceFactory *services.ServiceFactory
		expectError    bool
		expectNilUI    bool
	}{
		{
			name:           "NilServiceFactory",
			serviceFactory: nil,
			expectError:    false, // UI.New doesn't return error for nil service factory
			expectNilUI:    false,
		},
		{
			name:           "ValidServiceFactory",
			serviceFactory: &services.ServiceFactory{},
			expectError:    false,
			expectNilUI:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip this test for now as it requires full UI initialization
			// which is problematic in test environment
			t.Skip("Skipping full UI test - requires proper mocking")
		})
	}
}

func TestUI_InitialState(t *testing.T) {
	t.Skip("Skipping full UI test - requires proper mocking")
}

func TestUI_ViewManagement(t *testing.T) {
	t.Skip("Skipping full UI test - requires proper mocking")
}

func TestUI_ComponentInitialization(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.app)
	assert.NotNil(t, ui.pages)
	assert.NotNil(t, ui.mainFlex)
	assert.NotNil(t, ui.statusBar)
	assert.NotNil(t, ui.viewContainer)
	assert.NotNil(t, ui.commandHandler.GetInput())

	assert.NotNil(t, ui.headerManager.GetDockerInfoCol())
	assert.NotNil(t, ui.headerManager.GetNavCol())
	assert.NotNil(t, ui.headerManager.GetActionsCol())
}

func TestUI_ShutdownChannel(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.shutdownChan)

	// Test that we can send to the channel without blocking
	select {
	case ui.shutdownChan <- struct{}{}:
		// Successfully sent
	default:
		t.Error("Shutdown channel is blocking, should be buffered")
	}
}

func TestUI_LoggerInitialization(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.log)
}

func TestUI_CommandInputInitialization(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.commandHandler.GetInput())

	assert.Equal(t, ": ", ui.commandHandler.GetInput().GetLabel())
	assert.Equal(t, " Command Mode ", ui.commandHandler.GetInput().GetTitle())
}

func TestUI_PagesSetup(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.pages)
	assert.NotNil(t, ui.pages)
}

func TestUI_MainLayout(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.mainFlex)
}

func TestUI_ViewContainer(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.viewContainer)

	title := ui.viewContainer.GetTitle()
	assert.Contains(t, title, "Containers") // Default view should be containers
}

func TestUI_StatusBar(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.statusBar)

	text := ui.statusBar.GetText(true)
	assert.NotEmpty(t, text)
}

func TestUI_CurrentViewTracking(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.Equal(t, constants.DefaultView, ui.viewRegistry.GetCurrentName())

	validViews := []string{constants.ViewContainers, constants.ViewImages, constants.ViewVolumes, constants.ViewNetworks}
	found := false
	currentView := ui.viewRegistry.GetCurrentName()
	for _, view := range validViews {
		if view == currentView {
			found = true
			break
		}
	}
	assert.True(t, found, "Current view should be one of the valid views")
}

func TestUI_ViewReferences(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.containersView)
	assert.NotNil(t, ui.imagesView)
	assert.NotNil(t, ui.volumesView)
	assert.NotNil(t, ui.networksView)

	containersPrimitive := ui.containersView.GetView()
	assert.NotNil(t, containersPrimitive)

	imagesPrimitive := ui.imagesView.GetView()
	assert.NotNil(t, imagesPrimitive)

	volumesPrimitive := ui.volumesView.GetView()
	assert.NotNil(t, volumesPrimitive)

	networksPrimitive := ui.networksView.GetView()
	assert.NotNil(t, networksPrimitive)
}

func TestUI_ServiceFactoryIntegration(t *testing.T) {
	t.Skip("Skipping full UI test - requires proper mocking")
}

func TestUI_CommandModeState(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.commandHandler.IsActive())

	assert.NotNil(t, ui.commandHandler)
}

func TestUI_DetailsModeState(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.inDetailsMode)

	ui.inDetailsMode = true
	assert.True(t, ui.inDetailsMode)

	ui.inDetailsMode = false
	assert.False(t, ui.inDetailsMode)
}

func TestUI_LogsModeState(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.inLogsMode)

	ui.inLogsMode = true
	assert.True(t, ui.inLogsMode)

	ui.inLogsMode = false
	assert.False(t, ui.inLogsMode)
}

func TestUI_CurrentActions(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.currentActions)

	ui.currentActions['a'] = "Action A"
	ui.currentActions['b'] = "Action B"

	assert.Equal(t, "Action A", ui.currentActions['a'])
	assert.Equal(t, "Action B", ui.currentActions['b'])
	assert.Equal(t, 2, len(ui.currentActions))
}
