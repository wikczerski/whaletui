package views

import (
	"context"
	"fmt"

	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// ImagesView displays and manages Docker images
type ImagesView struct {
	*BaseView[models.Image]
	handlers *handlers.ActionHandlers
}

// NewImagesView creates a new images view
func NewImagesView(ui interfaces.UIInterface) *ImagesView {
	headers := []string{"ID", "Repository", "Tag", "Size", "Created", "Containers"}
	baseView := NewBaseView[models.Image](ui, "images", headers)

	iv := &ImagesView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	iv.ListItems = iv.listImages
	iv.FormatRow = func(i models.Image) []string { return iv.formatImageRow(&i) }
	iv.GetItemID = func(i models.Image) string { return i.ID }
	iv.GetItemName = func(i models.Image) string { return i.Repository }
	iv.HandleKeyPress = func(key rune, i models.Image) { iv.handleImageKey(key, &i) }
	iv.ShowDetails = func(i models.Image) { iv.showImageDetails(&i) }
	iv.GetActions = iv.getImageActions

	return iv
}

func (iv *ImagesView) listImages(ctx context.Context) ([]models.Image, error) {
	services := iv.ui.GetServices()
	if services == nil {
		return []models.Image{}, nil
	}

	if imageService := services.GetImageService(); imageService != nil {
		return imageService.ListImages(ctx)
	}

	return []models.Image{}, nil
}

func (iv *ImagesView) formatImageRow(image *models.Image) []string {
	return []string{
		image.ID,
		image.Repository,
		image.Tag,
		image.Size,
		builders.FormatTime(image.Created),
		fmt.Sprintf("%d", image.Containers),
	}
}

func (iv *ImagesView) getImageActions() map[rune]string {
	return map[rune]string{
		'd': "Delete",
		'i': "Inspect",
	}
}

func (iv *ImagesView) handleImageKey(key rune, image *models.Image) {
	iv.handleAction(key, image)
}

func (iv *ImagesView) showImageDetails(image *models.Image) {
	ctx := context.Background()
	services := iv.ui.GetServices()
	if services == nil || services.GetImageService() == nil {
		iv.ShowItemDetails(*image, nil, fmt.Errorf("image service not available"))
		return
	}
	inspectData, err := services.GetImageService().InspectImage(ctx, image.ID)
	iv.ShowItemDetails(*image, inspectData, err)
}

func (iv *ImagesView) handleAction(key rune, image *models.Image) {
	services := iv.ui.GetServices()
	if services == nil {
		return
	}

	if services.GetImageService() == nil {
		return
	}

	switch key {
	case 'd':
		iv.deleteImage(image.ID)
	case 'i':
		iv.inspectImage(image.ID)
	}
}

func (iv *ImagesView) deleteImage(id string) {
	services := iv.ui.GetServices()
	if services == nil || services.GetImageService() == nil {
		return
	}
	iv.handlers.HandleResourceAction('d', "image", id, "",
		services.GetImageService().InspectImage, nil, func() { iv.Refresh() })
}

func (iv *ImagesView) inspectImage(id string) {
	services := iv.ui.GetServices()
	if services == nil || services.GetImageService() == nil {
		return
	}
	iv.handlers.HandleResourceAction('i', "image", id, "",
		services.GetImageService().InspectImage, nil, func() { iv.Refresh() })
}
