package image

import (
	"context"
	"fmt"
	"strings"

	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ImagesView displays and manages Docker images
type ImagesView struct {
	*shared.BaseView[shared.Image]
	handlers *handlers.ActionHandlers
}

// NewImagesView creates a new images view
func NewImagesView(ui interfaces.UIInterface) *ImagesView {
	headers := []string{"ID", "Repository", "Tag", "Size", "Created", "Containers"}
	baseView := shared.NewBaseView[shared.Image](ui, "images", headers)

	iv := &ImagesView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	iv.ListItems = iv.listImages
	iv.FormatRow = func(i shared.Image) []string { return iv.formatImageRow(&i) }
	iv.GetItemID = func(i shared.Image) string { return i.ID }
	iv.GetItemName = func(i shared.Image) string { return i.RepoTags[0] }
	iv.HandleKeyPress = func(key rune, i shared.Image) { iv.handleImageKey(key, &i) }
	iv.ShowDetails = func(i shared.Image) { iv.showImageDetails(&i) }
	iv.GetActions = iv.getImageActions

	return iv
}

func (iv *ImagesView) listImages(ctx context.Context) ([]shared.Image, error) {
	services := iv.GetUI().GetServicesAny()
	if services == nil {
		return []shared.Image{}, nil
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if imageService := serviceFactory.GetImageService(); imageService != nil {
			// Type assertion to get the ListImages method
			if imageService != nil {
				images, err := imageService.ListImages(ctx)
				if err != nil {
					return nil, err
				}
				return images, nil
			}
		}
	}

	return []shared.Image{}, nil
}

func (iv *ImagesView) formatImageRow(image *shared.Image) []string {
	repoTag := ""
	if len(image.RepoTags) > 0 {
		repoTag = image.RepoTags[0]
	}
	parts := strings.Split(repoTag, ":")
	repo := ""
	tag := ""
	if len(parts) >= 2 {
		repo = parts[0]
		tag = parts[1]
	}

	return []string{
		image.ID,
		repo,
		tag,
		image.Size,
		builders.FormatTime(image.Created),
		"0", // shared.Image doesn't have Containers field
	}
}

func (iv *ImagesView) getImageActions() map[rune]string {
	return map[rune]string{
		'd': "Delete",
		'i': "Inspect",
	}
}

func (iv *ImagesView) handleImageKey(key rune, image *shared.Image) {
	iv.handleAction(key, image)
}

func (iv *ImagesView) showImageDetails(image *shared.Image) {
	ctx := context.Background()
	services := iv.GetUI().GetServicesAny()
	if services == nil {
		iv.ShowItemDetails(*image, nil, fmt.Errorf("image service not available"))
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetImageService() any }); ok {
		if imageService := serviceFactory.GetImageService(); imageService != nil {
			// Type assertion to get the InspectImage method
			if inspectService, ok := imageService.(interface {
				InspectImage(context.Context, string) (map[string]any, error)
			}); ok {
				inspectData, err := inspectService.InspectImage(ctx, image.ID)
				iv.ShowItemDetails(*image, inspectData, err)
				return
			}
		}
	}

	iv.ShowItemDetails(*image, nil, fmt.Errorf("image service not available"))
}

func (iv *ImagesView) handleAction(key rune, image *shared.Image) {
	services := iv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetImageService() any }); ok {
		if serviceFactory.GetImageService() != nil {
			switch key {
			case 'd':
				iv.deleteImage(image.ID)
			case 'i':
				iv.inspectImage(image.ID)
			}
		}
	}
}

func (iv *ImagesView) deleteImage(id string) {
	services := iv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetImageService() any }); ok {
		if imageService := serviceFactory.GetImageService(); imageService != nil {
			// Type assertion to get the InspectImage method
			if inspectService, ok := imageService.(interface {
				InspectImage(context.Context, string) (map[string]any, error)
			}); ok {
				iv.handlers.HandleResourceAction('d', "image", id, "",
					inspectService.InspectImage, nil, func() { iv.Refresh() })
			}
		}
	}
}

func (iv *ImagesView) inspectImage(id string) {
	services := iv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetImageService() any }); ok {
		if imageService := serviceFactory.GetImageService(); imageService != nil {
			// Type assertion to get the InspectImage method
			if inspectService, ok := imageService.(interface {
				InspectImage(context.Context, string) (map[string]any, error)
			}); ok {
				iv.handlers.HandleResourceAction('i', "image", id, "",
					inspectService.InspectImage, nil, func() { iv.Refresh() })
			}
		}
	}
}
