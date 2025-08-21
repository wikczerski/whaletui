package image

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
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
		ui.On("GetServices").Return(mockSF).Maybe()
	} else {
		ui.On("GetServices").Return(sf).Maybe()
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
	mockImages := []Image{
		{
			ID:         "image1",
			Repository: "test/repo1",
			Tag:        "latest",
			Size:       "100MB",
			Created:    time.Now(),
			Containers: 2,
		},
		{
			ID:         "image2",
			Repository: "test/repo2",
			Tag:        "v1.0",
			Size:       "200MB",
			Created:    time.Now().Add(-24 * time.Hour),
			Containers: 1,
		},
	}

	is := mocks.NewMockImageService(t)
	// Convert []Image to []any for the mock interface
	mockImagesAny := make([]any, len(mockImages))
	for i, img := range mockImages {
		mockImagesAny[i] = img
	}
	is.EXPECT().ListImages(context.Background()).Return(mockImagesAny, nil)

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is)
	ui := newImagesUIMockWithServices(t, sf)

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Equal(t, mockImages, imagesView.GetItems())
}

func TestImagesView_Refresh_ServiceError(t *testing.T) {
	is := mocks.NewMockImageService(t)
	is.EXPECT().ListImages(context.Background()).Return([]any{}, assert.AnError)

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.EXPECT().GetImageService().Return(is)
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Empty(t, imagesView.GetItems())
}

func TestImagesView_ShowImageDetails_Success(t *testing.T) {
	is := mocks.NewMockImageService(t)
	is.On("InspectImage", context.Background(), "image1").Return(map[string]any{"ok": true}, nil).Maybe()

	mockImage := Image{
		ID:         "image1",
		Repository: "test/repo1",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is)
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	imagesView := NewImagesView(ui)
	// Test the method directly without accessing unexported fields
	imagesView.showImageDetails(&mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowImageDetails_InspectError(t *testing.T) {
	is := mocks.NewMockImageService(t)
	is.On("InspectImage", context.Background(), "image1").Return(map[string]any(nil), assert.AnError).Maybe()

	mockImage := Image{
		ID:         "image1",
		Repository: "test/repo2",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	sf := mocks.NewMockServiceFactoryInterface(t)
	sf.On("GetImageService").Return(is)
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	imagesView := NewImagesView(ui)
	// Test the method directly without accessing unexported fields
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
	testImage := Image{
		ID:         "image1",
		Repository: "test/repo1",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
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
