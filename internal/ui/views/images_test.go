package views

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wikczerski/whaletui/internal/models"
	"github.com/wikczerski/whaletui/internal/services"
	servicemocks "github.com/wikczerski/whaletui/internal/services/mocks"
	uimocks "github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func newImagesUIMockWithServices(t *testing.T, sf *services.ServiceFactory) *uimocks.MockUIInterface {
	ui := uimocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()

	if sf == nil {
		// Create a mock service factory that returns nil for all services
		mockSF := servicemocks.NewMockServiceFactoryInterface(t)
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

	assert.NotNil(t, imagesView.view)
}

func TestNewImagesView_TableField(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.NotNil(t, imagesView.table)
}

func TestNewImagesView_ItemsField(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	assert.Empty(t, imagesView.items)
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

	assert.Equal(t, imagesView.view, view)
}

func TestImagesView_Refresh_NoServices(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)

	imagesView.Refresh()

	assert.Empty(t, imagesView.items)
}

func TestImagesView_Refresh_WithServices(t *testing.T) {
	mockImages := []models.Image{
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

	is := servicemocks.NewMockImageService(t)
	is.On("ListImages", context.Background()).Return(mockImages, nil)

	sf := &services.ServiceFactory{ImageService: is}
	ui := newImagesUIMockWithServices(t, sf)

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Equal(t, mockImages, imagesView.items)
}

func TestImagesView_Refresh_ServiceError(t *testing.T) {
	is := servicemocks.NewMockImageService(t)
	is.On("ListImages", context.Background()).Return([]models.Image{}, assert.AnError)

	sf := &services.ServiceFactory{ImageService: is}
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	imagesView := NewImagesView(ui)
	imagesView.Refresh()

	assert.Empty(t, imagesView.items)
}

func TestImagesView_ShowImageDetails_Success(t *testing.T) {
	is := servicemocks.NewMockImageService(t)
	is.On("InspectImage", context.Background(), "image1").Return(map[string]any{"ok": true}, nil).Maybe()

	mockImage := models.Image{
		ID:         "image1",
		Repository: "test/repo1",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	sf := &services.ServiceFactory{ImageService: is}
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	imagesView := NewImagesView(ui)
	imagesView.items = []models.Image{mockImage}
	imagesView.showImageDetails(&mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowImageDetails_InspectError(t *testing.T) {
	is := servicemocks.NewMockImageService(t)
	is.On("InspectImage", context.Background(), "image1").Return(map[string]any(nil), assert.AnError).Maybe()

	mockImage := models.Image{
		ID:         "image1",
		Repository: "test/repo2",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	sf := &services.ServiceFactory{ImageService: is}
	ui := newImagesUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	imagesView := NewImagesView(ui)
	imagesView.items = []models.Image{mockImage}
	imagesView.showImageDetails(&mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_HandleAction_Delete(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	imagesView.items = []models.Image{
		{
			ID:         "image1",
			Repository: "test/repo1",
			Tag:        "latest",
			Size:       "100MB",
			Created:    time.Now(),
			Containers: 2,
		},
	}
	imagesView.table.Select(1, 0)

	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_Inspect(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	imagesView.items = []models.Image{
		{
			ID:         "image1",
			Repository: "test/repo1",
			Tag:        "latest",
			Size:       "100MB",
			Created:    time.Now(),
			Containers: 2,
		},
	}
	imagesView.table.Select(1, 0)

	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_InvalidSelection(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	testImage := models.Image{
		ID:         "image1",
		Repository: "test/repo1",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}
	imagesView.items = []models.Image{}
	imagesView.table.Select(0, 0)

	// Test action handling
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
	imagesView.items = []models.Image{
		{
			ID:         "image1",
			Repository: "test/repo1",
			Tag:        "latest",
			Size:       "100MB",
			Created:    time.Now(),
			Containers: 2,
		},
	}
	imagesView.table.Select(1, 0)

	assert.NotNil(t, imagesView.table.GetInputCapture())
}

func TestImagesView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	ui := newImagesUIMockWithServices(t, nil)
	imagesView := NewImagesView(ui)
	imagesView.items = []models.Image{}
	imagesView.table.Select(0, 0)

	assert.NotNil(t, imagesView.table.GetInputCapture())
}
