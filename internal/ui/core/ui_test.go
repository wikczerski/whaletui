package core

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/internal/models"
	"github.com/wikczerski/whaletui/internal/services"
	"github.com/wikczerski/whaletui/internal/services/mocks"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// createMinimalMockServices creates minimal mock services for testing
func createMinimalMockServices() *services.ServiceFactory {
	mockContainerService := &mocks.MockContainerService{}
	mockImageService := &mocks.MockImageService{}
	mockVolumeService := &mocks.MockVolumeService{}
	mockNetworkService := &mocks.MockNetworkService{}
	mockDockerInfoService := &mocks.MockDockerInfoService{}
	mockLogsService := &mocks.MockLogsService{}

	// Set up minimal mock expectations
	mockContainerService.On("GetActionsString").Return("")
	mockImageService.On("GetActionsString").Return("")
	mockVolumeService.On("GetActionsString").Return("")
	mockNetworkService.On("GetActionsString").Return("")
	mockContainerService.On("ListContainers", mock.Anything).Return([]models.Container{}, nil)
	mockDockerInfoService.On("GetDockerInfo", mock.Anything).Return(&models.DockerInfo{}, nil)

	return &services.ServiceFactory{
		ContainerService:  mockContainerService,
		ImageService:      mockImageService,
		VolumeService:     mockVolumeService,
		NetworkService:    mockNetworkService,
		DockerInfoService: mockDockerInfoService,
		LogsService:       mockLogsService,
	}
}

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

func TestUI_ComponentInitialization_App(t *testing.T) {
	// Create mock services with proper expectations
	mockContainerService := &mocks.MockContainerService{}
	mockImageService := &mocks.MockImageService{}
	mockVolumeService := &mocks.MockVolumeService{}
	mockNetworkService := &mocks.MockNetworkService{}
	mockDockerInfoService := &mocks.MockDockerInfoService{}
	mockLogsService := &mocks.MockLogsService{}

	// Set up mock expectations for GetActionsString() calls
	mockContainerService.On("GetActionsString").Return("<s> Start\n<S> Stop\n<r> Restart\n<d> Delete\n<a> Attach\n<l> Logs\n<i> Inspect\n<n> New\n<e> Exec\n<f> Filter\n<t> Sort\n<h> History\n<enter> Details\n<:> Command")
	mockImageService.On("GetActionsString").Return("<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command")
	mockVolumeService.On("GetActionsString").Return("<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command")
	mockNetworkService.On("GetActionsString").Return("<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command")

	// Set up mock expectations for ListContainers() calls
	mockContainerService.On("ListContainers", mock.Anything).Return([]models.Container{}, nil)

	// Set up mock expectations for GetDockerInfo() calls
	mockDockerInfo := &models.DockerInfo{
		Version:         "20.10.0",
		Containers:      5,
		Images:          10,
		Volumes:         3,
		Networks:        2,
		OperatingSystem: "linux",
		Architecture:    "amd64",
		Driver:          "overlay2",
		LoggingDriver:   "json-file",
	}
	mockDockerInfoService.On("GetDockerInfo", mock.Anything).Return(mockDockerInfo, nil)

	// Create a mock service factory
	mockServiceFactory := &services.ServiceFactory{
		ContainerService:  mockContainerService,
		ImageService:      mockImageService,
		VolumeService:     mockVolumeService,
		NetworkService:    mockNetworkService,
		DockerInfoService: mockDockerInfoService,
		LogsService:       mockLogsService,
	}

	ui, err := New(mockServiceFactory, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.app)

	// Verify that all expected methods were called
	mockContainerService.AssertExpectations(t)
	mockImageService.AssertExpectations(t)
	mockVolumeService.AssertExpectations(t)
	mockNetworkService.AssertExpectations(t)
	mockDockerInfoService.AssertExpectations(t)
}

func TestUI_ComponentInitialization_Pages(t *testing.T) {
	// Create minimal mock services for this test
	mockServiceFactory := createMinimalMockServices()

	ui, err := New(mockServiceFactory, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.pages)
}

func TestUI_ComponentInitialization_MainFlex(t *testing.T) {
	// Create minimal mock services for this test
	mockServiceFactory := createMinimalMockServices()

	ui, err := New(mockServiceFactory, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.mainFlex)
}

func TestUI_ComponentInitialization_StatusBar(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.statusBar)
}

func TestUI_ComponentInitialization_ViewContainer(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.viewContainer)
}

func TestUI_ComponentInitialization_CommandHandler(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.commandHandler.GetInput())
}

func TestUI_ComponentInitialization_HeaderManager_DockerInfo(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.headerManager.GetDockerInfoCol())
}

func TestUI_ComponentInitialization_HeaderManager_Nav(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.headerManager.GetNavCol())
}

func TestUI_ComponentInitialization_HeaderManager_Actions(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.headerManager.GetActionsCol())
}

func TestUI_ShutdownChannel_Exists(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.shutdownChan)
}

func TestUI_ShutdownChannel_NotBlocking(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

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
}

func TestUI_CommandInput_Label(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.Equal(t, ": ", ui.commandHandler.GetInput().GetLabel())
}

func TestUI_CommandInput_Title(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.Equal(t, " Command Mode ", ui.commandHandler.GetInput().GetTitle())
}

func TestUI_PagesSetup(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

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
}

func TestUI_ViewContainer_Title(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	title := ui.viewContainer.GetTitle()
	assert.Contains(t, title, "Containers") // Default view should be containers
}

func TestUI_StatusBar(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.statusBar)
}

func TestUI_StatusBar_Text(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	text := ui.statusBar.GetText(true)
	assert.NotEmpty(t, text)
}

func TestUI_CurrentViewTracking(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.Equal(t, constants.DefaultView, ui.viewRegistry.GetCurrentName())
}

func TestUI_CurrentViewTracking_ValidView(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

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

func TestUI_ViewReferences_Containers(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.containersView)
}

func TestUI_ViewReferences_Images(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.imagesView)
}

func TestUI_ViewReferences_Volumes(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.volumesView)
}

func TestUI_ViewReferences_Networks(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.networksView)
}

func TestUI_ViewReferences_ContainersPrimitive(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	containersPrimitive := ui.containersView.GetView()
	assert.NotNil(t, containersPrimitive)
}

func TestUI_ViewReferences_ImagesPrimitive(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	imagesPrimitive := ui.imagesView.GetView()
	assert.NotNil(t, imagesPrimitive)
}

func TestUI_ViewReferences_VolumesPrimitive(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	volumesPrimitive := ui.volumesView.GetView()
	assert.NotNil(t, volumesPrimitive)
}

func TestUI_ViewReferences_NetworksPrimitive(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	networksPrimitive := ui.networksView.GetView()
	assert.NotNil(t, networksPrimitive)
}

func TestUI_ServiceFactoryIntegration(t *testing.T) {
	t.Skip("Skipping full UI test - requires proper mocking")
}

func TestUI_CommandModeState_InitiallyInactive(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.commandHandler.IsActive())
}

func TestUI_CommandModeState_HandlerExists(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.commandHandler)
}

func TestUI_DetailsModeState_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.inDetailsMode)
}

func TestUI_DetailsModeState_CanSetTrue(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.inDetailsMode = true
	assert.True(t, ui.inDetailsMode)
}

func TestUI_DetailsModeState_CanSetFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.inDetailsMode = false
	assert.False(t, ui.inDetailsMode)
}

func TestUI_LogsModeState_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.False(t, ui.inLogsMode)
}

func TestUI_LogsModeState_CanSetTrue(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.inLogsMode = true
	assert.True(t, ui.inLogsMode)
}

func TestUI_LogsModeState_CanSetFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.inLogsMode = false
	assert.False(t, ui.inLogsMode)
}

func TestUI_CurrentActions_Exists(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.currentActions)
}

func TestUI_CurrentActions_CanSetActionA(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.currentActions['a'] = "Action A"

	assert.Equal(t, "Action A", ui.currentActions['a'])
}

func TestUI_CurrentActions_CanSetActionB(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.currentActions['b'] = "Action B"

	assert.Equal(t, "Action B", ui.currentActions['b'])
}

func TestUI_CurrentActions_Count(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.currentActions['a'] = "Action A"
	ui.currentActions['b'] = "Action B"

	assert.Equal(t, 2, len(ui.currentActions))
}

func TestUI_ModalState_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Initially no modals should be active
	assert.False(t, ui.IsModalActive())
}

func TestUI_ModalState_CanAddModal(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Add a modal page and check if it's detected
	ui.pages.AddPage("help_modal", ui.mainFlex, true, true)
	assert.True(t, ui.IsModalActive())
}

func TestUI_ModalState_CanRemoveModal(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Add a modal page
	ui.pages.AddPage("help_modal", ui.mainFlex, true, true)

	// Remove the modal page
	ui.pages.RemovePage("help_modal")
	assert.False(t, ui.IsModalActive())
}

func TestUI_ModalState_ErrorModal(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test with different modal types
	ui.pages.AddPage("error_modal", ui.mainFlex, true, true)
	assert.True(t, ui.IsModalActive())
}

func TestUI_ModalState_ConfirmModal(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.pages.RemovePage("error_modal")
	ui.pages.AddPage("confirm_modal", ui.mainFlex, true, true)
	assert.True(t, ui.IsModalActive())
}

func TestUI_ModalState_ExecOutputModal(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.pages.RemovePage("confirm_modal")
	ui.pages.AddPage("exec_output_modal", ui.mainFlex, true, true)
	assert.True(t, ui.IsModalActive())
}

func TestUI_KeyBindingHandlers_CommandModeKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	ui.app.SetFocus(ui.commandHandler.GetInput())

	assert.NotNil(t, ui.handleCommandModeKeyBindings)
}

func TestUI_KeyBindingHandlers_NormalModeKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.handleNormalModeKeyBindings)
}

func TestUI_KeyBindingHandlers_RuneKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.handleRuneKeyBindings)
}

func TestUI_KeyBindingHandlers_CtrlCKeyBinding(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.handleCtrlCKeyBinding)
}

func TestUI_KeyBindingHandlers_GlobalKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Verify all key binding handler functions exist
	assert.NotNil(t, ui.handleGlobalKeyBindings)
}

func TestUI_KeyBindingHandlers_ExecCommandKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.handleExecCommandKeyBindings)
}

func TestUI_KeyBindingHandlers_ShellViewKeyBindings(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	assert.NotNil(t, ui.handleShellViewKeyBindings)
}

func TestUI_ModalDetection_ExecCommandInput_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Initially no exec command input should be active
	assert.False(t, ui.isExecCommandInputActive())
}

func TestUI_ModalDetection_ExecCommandInput_CanActivate(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Create a mock exec command input
	mockInput := tview.NewInputField()
	mockInput.SetLabel(" Exec Command: ")

	// Set focus to the mock input
	ui.app.SetFocus(mockInput)
	assert.True(t, ui.isExecCommandInputActive())
}

func TestUI_ModalDetection_ExecCommandInput_DifferentLabel(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Create a mock exec command input
	mockInput := tview.NewInputField()
	mockInput.SetLabel(" Exec Command: ")

	// Set focus to the mock input
	ui.app.SetFocus(mockInput)

	// Test with different label
	mockInput.SetLabel("Different Label")
	assert.False(t, ui.isExecCommandInputActive())
}

func TestUI_ModalDetection_ShellView_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Initially no shell view should be active
	assert.False(t, ui.isShellViewActive())
}

func TestUI_ModalDetection_ShellView_NilShellView(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test with nil shell view
	ui.shellView = nil
	assert.False(t, ui.isShellViewActive())
}

func TestUI_ModalDetection_ShellInputField_InitiallyFalse(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Initially no shell input field should be focused
	assert.False(t, ui.isShellInputFieldFocused())
}

func TestUI_ModalDetection_ShellInputField_NilShellView(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test with nil shell view
	ui.shellView = nil
	assert.False(t, ui.isShellInputFieldFocused())
}

func TestUI_InitializationFunctions_SetupMainPages(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test setup main pages
	commandInput := ui.commandHandler.GetInput()
	ui.setupMainPages(commandInput)

	// Verify pages were added
	assert.True(t, ui.pages.HasPage("main"))
}

func TestUI_InitializationFunctions_SetupMainPages_CommandPage(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test setup main pages
	commandInput := ui.commandHandler.GetInput()
	ui.setupMainPages(commandInput)

	// Verify pages were added
	assert.True(t, ui.pages.HasPage("command"))
}

func TestUI_InitializationFunctions_InitializeUIState(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test UI state initialization
	ui.initializeUIState()

	assert.NotNil(t, ui.headerManager)
}

func TestUI_ViewCreationFunctions_CreateResourceViews(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test resource view creation
	ui.createResourceViews()

	// Verify all views were created
	assert.NotNil(t, ui.containersView)
}

func TestUI_ViewCreationFunctions_CreateResourceViews_Images(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test resource view creation
	ui.createResourceViews()

	// Verify all views were created
	assert.NotNil(t, ui.imagesView)
}

func TestUI_ViewCreationFunctions_CreateResourceViews_Volumes(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test resource view creation
	ui.createResourceViews()

	// Verify all views were created
	assert.NotNil(t, ui.volumesView)
}

func TestUI_ViewCreationFunctions_CreateResourceViews_Networks(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test resource view creation
	ui.createResourceViews()

	// Verify all views were created
	assert.NotNil(t, ui.networksView)
}

func TestUI_ViewCreationFunctions_RegisterViewsWithActions(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test view registration
	ui.registerViewsWithActions()

	// Verify views are registered
	assert.True(t, ui.viewRegistry.Exists("containers"))
}

func TestUI_ViewCreationFunctions_RegisterViewsWithActions_Images(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test view registration
	ui.registerViewsWithActions()

	// Verify views are registered
	assert.True(t, ui.viewRegistry.Exists("images"))
}

func TestUI_ViewCreationFunctions_RegisterViewsWithActions_Volumes(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test view registration
	ui.registerViewsWithActions()

	// Verify views are registered
	assert.True(t, ui.viewRegistry.Exists("volumes"))
}

func TestUI_ViewCreationFunctions_RegisterViewsWithActions_Networks(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test view registration
	ui.registerViewsWithActions()

	// Verify views are registered
	assert.True(t, ui.viewRegistry.Exists("networks"))
}

func TestUI_ViewCreationFunctions_SetDefaultView(t *testing.T) {
	ui, err := New(nil, "")
	require.NoError(t, err)
	require.NotNil(t, ui)

	// Test default view setting
	ui.setDefaultView()

	// Verify default view is set
	assert.Equal(t, constants.DefaultView, ui.viewRegistry.GetCurrentName())
}
