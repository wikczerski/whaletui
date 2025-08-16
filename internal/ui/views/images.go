package views

import (
	"context"
	"fmt"

	"github.com/user/d5r/internal/models"
	"github.com/user/d5r/internal/ui/builders"
	"github.com/user/d5r/internal/ui/handlers"
	"github.com/user/d5r/internal/ui/interfaces"
)

type ImagesView struct {
	*BaseView[models.Image]
	handlers *handlers.ActionHandlers
}

func NewImagesView(ui interfaces.UIInterface) *ImagesView {
	headers := []string{"ID", "Repository", "Tag", "Size", "Created", "Containers"}
	baseView := NewBaseView[models.Image](ui, "images", headers)

	iv := &ImagesView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	iv.ListItems = iv.listImages
	iv.FormatRow = iv.formatImageRow
	iv.GetItemID = func(i models.Image) string { return i.ID }
	iv.GetItemName = func(i models.Image) string { return i.Repository }
	iv.HandleKeyPress = iv.handleImageKey
	iv.ShowDetails = iv.showImageDetails
	iv.GetActions = iv.getImageActions

	return iv
}

func (iv *ImagesView) listImages(ctx context.Context) ([]models.Image, error) {
	services := iv.ui.GetServices()
	if services == nil || services.ImageService == nil {
		return []models.Image{}, nil
	}
	return services.ImageService.ListImages(ctx)
}

func (iv *ImagesView) formatImageRow(image models.Image) []string {
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

func (iv *ImagesView) handleImageKey(key rune, image models.Image) {
	switch key {
	case 'd':
		iv.deleteImage(image.ID, image.Repository)
	case 'i':
		iv.inspectImage(image.ID)
	}
}

func (iv *ImagesView) showImageDetails(image models.Image) {
	ctx := context.Background()
	services := iv.ui.GetServices()
	inspectData, err := services.ImageService.InspectImage(ctx, image.ID)
	iv.ShowItemDetails(image, inspectData, err)
}

func (iv *ImagesView) deleteImage(id, repository string) {
	services := iv.ui.GetServices()
	iv.handlers.HandleResourceAction('d', "image", id, repository,
		services.ImageService.InspectImage, nil, func() { iv.Refresh() })
}

func (iv *ImagesView) inspectImage(id string) {
	services := iv.ui.GetServices()
	iv.handlers.HandleResourceAction('i', "image", id, "",
		services.ImageService.InspectImage, nil, func() { iv.Refresh() })
}
