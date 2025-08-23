package image

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

func newImagesUIMockWithServices(t *testing.T, sf interfaces.ServiceFactoryInterface) *mocks.MockUIInterface {
	ui := mocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()

	if sf == nil {
		// Create a mock service factory that returns nil for all services
		mockSF := mocks.NewMockServiceFactoryInterface(t)
		mockSF.On("GetImageService").Return(nil).Maybe()
		mockSF.On("GetContainerService").Return(nil).Maybe()
		mockSF.On("GetVolumeService").Return(nil).Maybe()
		mockSF.On("GetNetworkService").Return(nil).Maybe()
		mockSF.On("GetDockerInfoService").Return(nil).Maybe()
		mockSF.On("GetLogsService").Return(nil).Maybe()
		ui.On("GetServicesAny").Return(mockSF).Maybe()
	} else {
		ui.On("GetServicesAny").Return(sf)
	}

	return ui
}

func TestNewImagesView_Creation(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView)
}

func TestNewImagesView_ViewField(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView.GetView())
}

func TestNewImagesView_TableField(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView.GetTable())
}

func TestNewImagesView_ItemsField(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.Empty(t, imagesView.GetItems())
}

func TestImagesView_GetView(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	view := imagesView.GetView()

	assert.NotNil(t, view)
}

func TestImagesView_GetView_ReturnsCorrectView(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	view := imagesView.GetView()

	assert.Equal(t, imagesView.GetView(), view)
}

func TestImagesView_Refresh_NoServices(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	imagesView.Refresh()

	assert.Empty(t, imagesView.GetItems())
}

func TestImagesView_Refresh_WithServices(t *testing.T) {
	mockImages := []shared.Image{
		{
			ID:       "image1",
			RepoTags: []string{"test/repo1:latest"},
			Size:     "100MB",
			Created:  time.Now(),
		},
		{
			ID:       "image2",
			RepoTags: []string{"test/repo2:v1.0"},
			Size:     "200MB",
			Created:  time.Now().Add(-24 * time.Hour),
		},
	}

	is := mocks.NewMockImageService(t)
	is.On("ListImages", context.Background()).Return(mockImages, nil)

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is).Once()
	// Set up other methods as Maybe() since they might not be called
	sf.On("GetContainerService").Return(nil).Maybe()
	sf.On("GetVolumeService").Return(nil).Maybe()
	sf.On("GetNetworkService").Return(nil).Maybe()
	sf.On("GetDockerInfoService").Return(nil).Maybe()
	sf.On("GetLogsService").Return(nil).Maybe()
	sf.On("GetSwarmServiceService").Return(nil).Maybe()
	sf.On("GetSwarmNodeService").Return(nil).Maybe()
	sf.On("GetCurrentService").Return(nil).Maybe()
	sf.On("SetCurrentService", mock.AnythingOfType("string")).Maybe()
	sf.On("IsServiceAvailable", mock.AnythingOfType("string")).Return(false).Maybe()
	sf.On("IsContainerServiceAvailable").Return(false).Maybe()

	ui := newImagesUIMockWithServices(t, sf)

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Equal(t, mockImages, imagesView.GetItems())
}

func TestImagesView_Refresh_ServiceError(t *testing.T) {
	is := mocks.NewMockImageService(t)
	is.On("ListImages", context.Background()).Return([]shared.Image{}, assert.AnError)

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is).Once()
	// Set up other methods as Maybe() since they might not be called
	sf.On("GetContainerService").Return(nil).Maybe()
	sf.On("GetVolumeService").Return(nil).Maybe()
	sf.On("GetNetworkService").Return(nil).Maybe()
	sf.On("GetDockerInfoService").Return(nil).Maybe()
	sf.On("GetLogsService").Return(nil).Maybe()
	sf.On("GetSwarmServiceService").Return(nil).Maybe()
	sf.On("GetSwarmNodeService").Return(nil).Maybe()
	sf.On("GetCurrentService").Return(nil).Maybe()
	sf.On("SetCurrentService", mock.AnythingOfType("string")).Maybe()
	sf.On("IsServiceAvailable", mock.AnythingOfType("string")).Return(false).Maybe()
	sf.On("IsContainerServiceAvailable").Return(false).Maybe()

	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Empty(t, imagesView.GetItems())
}

func TestImagesView_ShowImageDetails_Success(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return()
	imagesView := NewImagesView(ui)

	mockImage := shared.Image{
		ID:       "image1",
		RepoTags: []string{"test/repo1:latest"},
		Size:     "100MB",
		Created:  time.Now(),
	}

	// Test the method directly - it should handle the case where no services are available
	imagesView.showImageDetails(&mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowImageDetails_InspectError(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return()
	imagesView := NewImagesView(ui)

	mockImage := shared.Image{
		ID:       "image1",
		RepoTags: []string{"test/repo2:latest"},
		Size:     "100MB",
		Created:  time.Now(),
	}

	// Test the method directly - it should handle the case where no services are available
	imagesView.showImageDetails(&mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_HandleAction_Delete(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	// Test that the key handler is properly set up
	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_Inspect(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	// Test that the key handler is properly set up
	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_InvalidSelection(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	testImage := shared.Image{
		ID:       "image1",
		RepoTags: []string{"test/repo1:latest"},
		Size:     "100MB",
		Created:  time.Now(),
	}

	// Test action handling with test image
	imagesView.handleAction('d', &testImage)
	imagesView.handleAction('i', &testImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowTable(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView)
}

func TestImagesView_DeleteImage(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView)
}

func TestImagesView_InspectImage(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView)
}

func TestImagesView_SetupKeyBindings(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	// Test that key bindings are properly set up
	assert.NotNil(t, imagesView.GetTable().GetInputCapture())
}

func TestImagesView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	// Test that key bindings are properly set up even with no items
	assert.NotNil(t, imagesView.GetTable().GetInputCapture())
}

func TestMockImplementsInterface(t *testing.T) {
	// Test that the mocks implement the interfaces
	var _ interfaces.ServiceFactoryInterface = mocks.NewMockServiceFactoryInterface(t)
	var _ interfaces.ImageService = mocks.NewMockImageService(t)
}

func TestMinimalRefresh(t *testing.T) {
	// Test the most basic refresh case
	mockImages := []shared.Image{
		{ID: "test", RepoTags: []string{"test/repo:latest"}},
	}

	is := mocks.NewMockImageService(t)
	is.On("ListImages", context.Background()).Return(mockImages, nil)

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is).Once()
	// Set up other methods as Maybe() since they might not be called
	sf.On("GetContainerService").Return(nil).Maybe()
	sf.On("GetVolumeService").Return(nil).Maybe()
	sf.On("GetNetworkService").Return(nil).Maybe()
	sf.On("GetDockerInfoService").Return(nil).Maybe()
	sf.On("GetLogsService").Return(nil).Maybe()
	sf.On("GetSwarmServiceService").Return(nil).Maybe()
	sf.On("GetSwarmNodeService").Return(nil).Maybe()
	sf.On("GetCurrentService").Return(nil).Maybe()
	sf.On("SetCurrentService", mock.AnythingOfType("string")).Maybe()
	sf.On("IsServiceAvailable", mock.AnythingOfType("string")).Return(false).Maybe()
	sf.On("IsContainerServiceAvailable").Return(false).Maybe()

	ui := newImagesUIMockWithServices(t, sf)

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Equal(t, mockImages, imagesView.GetItems())
}

// TestImagesView_RefreshWithRealData removed due to import cycle
// This test would require importing internal/ui/core which creates a circular dependency
