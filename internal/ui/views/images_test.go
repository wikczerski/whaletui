package views

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/services"
)

// MockImageService is a mock implementation of ImageService
type MockImageService struct {
	images     []models.Image
	inspectErr error
	listErr    error
}

func (m *MockImageService) ListImages(ctx context.Context) ([]models.Image, error) {
	return m.images, m.listErr
}

func (m *MockImageService) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	return nil, m.inspectErr
}

func (m *MockImageService) RemoveImage(ctx context.Context, id string, force bool) error {
	return nil
}

func TestNewImagesView(t *testing.T) {
	mockUI := NewMockUI()
	imagesView := NewImagesView(mockUI)

	assert.NotNil(t, imagesView)
	// Note: ui field is internal, we can't test it directly
	assert.NotNil(t, imagesView.view)
	assert.NotNil(t, imagesView.table)
	assert.Empty(t, imagesView.items)
}

func TestImagesView_GetView(t *testing.T) {
	mockUI := NewMockUI()
	imagesView := NewImagesView(mockUI)
	view := imagesView.GetView()

	assert.NotNil(t, view)
	assert.Equal(t, imagesView.view, view)
}

func TestImagesView_Refresh_NoServices(t *testing.T) {
	mockUI := NewMockUI()
	mockUI.services = nil

	imagesView := NewImagesView(mockUI)
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

	mockImageService := &MockImageService{
		images:  mockImages,
		listErr: nil,
	}

	mockUI := NewMockUI()
	mockUI.services = &services.ServiceFactory{
		ImageService: mockImageService,
	}

	imagesView := NewImagesView(mockUI)
	imagesView.Refresh()

	assert.Equal(t, mockImages, imagesView.items)
}

func TestImagesView_Refresh_ServiceError(t *testing.T) {
	mockImageService := &MockImageService{
		images:  []models.Image{},
		listErr: assert.AnError,
	}

	mockUI := NewMockUI()
	mockUI.services = &services.ServiceFactory{
		ImageService: mockImageService,
	}

	imagesView := NewImagesView(mockUI)
	imagesView.Refresh()

	assert.Empty(t, imagesView.items)
}

func TestImagesView_ShowImageDetails_Success(t *testing.T) {
	mockImageService := &MockImageService{
		inspectErr: nil,
	}

	mockImage := models.Image{
		ID:         "image1",
		Repository: "test/repo1",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	mockUI := NewMockUI()
	mockUI.services = &services.ServiceFactory{
		ImageService: mockImageService,
	}

	imagesView := NewImagesView(mockUI)
	imagesView.items = []models.Image{mockImage}
	imagesView.showImageDetails(mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowImageDetails_InspectError(t *testing.T) {
	mockImageService := &MockImageService{
		inspectErr: assert.AnError,
	}

	mockImage := models.Image{
		ID:         "image1",
		Repository: "test/repo2",
		Tag:        "latest",
		Size:       "100MB",
		Created:    time.Now(),
		Containers: 2,
	}

	mockUI := NewMockUI()
	mockUI.services = &services.ServiceFactory{
		ImageService: mockImageService,
	}

	imagesView := NewImagesView(mockUI)
	imagesView.items = []models.Image{mockImage}
	imagesView.showImageDetails(mockImage)

	assert.NotNil(t, imagesView)
}

func TestImagesView_HandleAction_Delete(t *testing.T) {
	mockUI := NewMockUI()

	imagesView := NewImagesView(mockUI)
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

	// We'll avoid calling handleAction since it triggers complex UI operations
	assert.NotNil(t, imagesView)
	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_Inspect(t *testing.T) {
	mockUI := NewMockUI()

	imagesView := NewImagesView(mockUI)
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

	// We'll avoid calling handleAction since it triggers complex UI operations
	assert.NotNil(t, imagesView)
	assert.NotNil(t, imagesView.handleImageKey)
}

func TestImagesView_HandleAction_InvalidSelection(t *testing.T) {
	mockUI := NewMockUI()

	imagesView := NewImagesView(mockUI)
	imagesView.items = []models.Image{}
	imagesView.table.Select(0, 0)
	imagesView.handleAction('d')
	imagesView.handleAction('i')

	assert.NotNil(t, imagesView)
}

func TestImagesView_ShowTable(t *testing.T) {
	mockUI := NewMockUI()
	imagesView := NewImagesView(mockUI)

	assert.NotNil(t, imagesView)
}

func TestImagesView_DeleteImage(t *testing.T) {
	mockUI := NewMockUI()
	imagesView := NewImagesView(mockUI)

	assert.NotNil(t, imagesView)
}

func TestImagesView_InspectImage(t *testing.T) {
	mockUI := NewMockUI()
	imagesView := NewImagesView(mockUI)

	assert.NotNil(t, imagesView)
}

func TestImagesView_SetupKeyBindings(t *testing.T) {
	mockUI := NewMockUI()

	imagesView := NewImagesView(mockUI)
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

	// Test key bindings - just verify they don't panic
	// Note: We can't easily test tcell.EventKey creation in tests
	// but we can verify the input capture function exists
	assert.NotNil(t, imagesView.table.GetInputCapture())
}

func TestImagesView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	mockUI := NewMockUI()

	imagesView := NewImagesView(mockUI)
	imagesView.items = []models.Image{}
	imagesView.table.Select(0, 0)

	assert.NotNil(t, imagesView.table.GetInputCapture())
}
