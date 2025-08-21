package image

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ImagesView displays and manages Docker images
type ImagesView struct {
	*shared.BaseView[Image]
	handlers *handlers.ActionHandlers
}

// NewImagesView creates a new images view
func NewImagesView(ui interfaces.UIInterface) *ImagesView {
	headers := []string{"ID", "Repository", "Tag", "Size", "Created", "Containers"}
	baseView := shared.NewBaseView[Image](ui, "images", headers)

	iv := &ImagesView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	iv.ListItems = iv.listImages
	iv.FormatRow = func(i Image) []string { return iv.formatImageRow(&i) }
	iv.GetItemID = func(i Image) string { return i.ID }
	iv.GetItemName = func(i Image) string { return i.Repository }
	iv.HandleKeyPress = func(key rune, i Image) { iv.handleImageKey(key, &i) }
	iv.ShowDetails = func(i Image) { iv.showImageDetails(&i) }
	iv.GetActions = iv.getImageActions

	return iv
}

func (iv *ImagesView) listImages(ctx context.Context) ([]Image, error) {
	services := iv.GetUI().GetServices()
	if services == nil {
		return []Image{}, nil
	}

	if imageService := services.GetImageService(); imageService != nil {
		images, err := imageService.ListImages(ctx)
		if err != nil {
			return nil, err
		}

		result := make([]Image, len(images))
		for i, image := range images {
			if img, ok := image.(Image); ok {
				result[i] = img
			}
		}

		return result, nil
	}

	return []Image{}, nil
}

func (iv *ImagesView) formatImageRow(image *Image) []string {
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

func (iv *ImagesView) handleImageKey(key rune, image *Image) {
	iv.handleAction(key, image)
}

func (iv *ImagesView) showImageDetails(image *Image) {
	ctx := context.Background()
	services := iv.GetUI().GetServices()
	if services == nil || services.GetImageService() == nil {
		iv.ShowItemDetails(*image, nil, fmt.Errorf("image service not available"))
		return
	}
	inspectData, err := services.GetImageService().InspectImage(ctx, image.ID)
	iv.ShowItemDetails(*image, inspectData, err)
}

func (iv *ImagesView) handleAction(key rune, image *Image) {
	services := iv.GetUI().GetServices()
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
	services := iv.GetUI().GetServices()
	if services == nil || services.GetImageService() == nil {
		return
	}
	iv.handlers.HandleResourceAction('d', "image", id, "",
		services.GetImageService().InspectImage, nil, func() { iv.Refresh() })
}

func (iv *ImagesView) inspectImage(id string) {
	services := iv.GetUI().GetServices()
	if services == nil || services.GetImageService() == nil {
		return
	}
	iv.handlers.HandleResourceAction('i', "image", id, "",
		services.GetImageService().InspectImage, nil, func() { iv.Refresh() })
}
