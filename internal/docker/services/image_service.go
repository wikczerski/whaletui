package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	domaintypes "github.com/wikczerski/whaletui/internal/docker/types"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// ImageService handles image-related operations
type ImageService struct {
	cli *client.Client
	log *slog.Logger
}

// NewImageService creates a new ImageService
func NewImageService(cli *client.Client, log *slog.Logger) *ImageService {
	return &ImageService{
		cli: cli,
		log: log,
	}
}

// ListImages lists all images
func (s *ImageService) ListImages(ctx context.Context) ([]domaintypes.Image, error) {
	images, err := s.getImageList(ctx)
	if err != nil {
		return nil, err
	}

	result := s.convertToImages(images)
	utils.SortImagesByCreationTime(result)
	return result, nil
}

// InspectImage inspects an image
func (s *ImageService) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	imageInfo, err := s.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("image inspect failed %s: %w", id, err)
	}
	return utils.MarshalToMap(imageInfo)
}

// RemoveImage removes an image
func (s *ImageService) RemoveImage(ctx context.Context, id string, force bool) error {
	if err := utils.ValidateID(id, "image ID"); err != nil {
		return err
	}

	opts := image.RemoveOptions{
		Force:         force,
		PruneChildren: true, // Remove dependent images by default
	}

	_, err := s.cli.ImageRemove(ctx, id, opts)
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %w", id, err)
	}

	return nil
}

// Helper methods

func (s *ImageService) getImageList(ctx context.Context) ([]image.Summary, error) {
	opts := image.ListOptions{}
	return s.cli.ImageList(ctx, opts)
}

func (s *ImageService) convertToImages(images []image.Summary) []domaintypes.Image {
	result := make([]domaintypes.Image, 0, len(images))
	for i := range images {
		img := &images[i]
		result = append(result, s.convertImage(img))
	}
	return result
}

func (s *ImageService) convertImage(img *image.Summary) domaintypes.Image {
	repo, tag := utils.ParseImageRepository(img.RepoTags)
	size := utils.FormatSize(img.Size)
	return domaintypes.Image{
		ID:         img.ID[7:19], // Remove "sha256:" prefix and truncate
		Repository: repo,
		Tag:        tag,
		Size:       size,
		Created:    time.Unix(img.Created, 0),
		Containers: int(img.Containers),
	}
}
