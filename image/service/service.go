package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vediagames/zeroerror"

	"github.com/vediagames/platform/image/domain"
)

type service struct {
	url       string
	processor domain.Processor
	client    *http.Client
}

type Config struct {
	URL       string
	Processor domain.Processor
	Client    *http.Client
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Client == nil, fmt.Errorf("empty client"))
	err.AddIf(c.Processor == nil, fmt.Errorf("empty processor"))

	return err.Err()
}

func New(c Config) domain.Service {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &service{
		url:       c.URL,
		client:    c.Client,
		processor: c.Processor,
	}
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	ogImg := originalThumbnailImage(req.Original)
	ogImgPath := imagePath(req.Resource, req.Slug, ogImg)
	ogImgURL := imageURL(s.url, ogImgPath)

	imgPath := imagePath(req.Resource, req.Slug, req.Image)
	imgURL := imageURL(s.url, imgPath)

	if imgPath == ogImgPath {
		return domain.GetResponse{
			URL: ogImgURL,
		}, nil
	}

	headRes, err := http.Head(imgURL)
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to head: %w", err)
	}

	if headRes.StatusCode == http.StatusOK {
		return domain.GetResponse{
			URL: imgURL,
		}, nil
	}

	_, err = s.processor.Process(ctx, domain.ProcessRequest{
		OriginalImageURL: ogImgURL,
		Path:             imgPath,
		Image:            req.Image,
	})
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to process: %w", err)
	}

	res := domain.GetResponse{
		URL: imgURL,
	}

	if err := res.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func imagePath(r domain.Resource, slug string, img domain.Image) string {
	return fmt.Sprintf("%s/%s/%s", resourceToPath(r), slug, img.File())
}

func imageURL(url, imagePath string) string {
	return fmt.Sprintf("%s/%s", url, imagePath)
}

func originalThumbnailImage(o domain.OriginalThumbnail) domain.Image {
	switch o {
	case domain.OriginalThumbnail512x512:
		return domain.Image{
			Format: domain.FormatJpg,
			Width:  512,
			Height: 512,
		}
	case domain.OriginalThumbnail128x128:
		return domain.Image{
			Format: domain.FormatJpg,
			Width:  512,
			Height: 512,
		}
	default:
		return domain.Image{
			Format: domain.FormatJpg,
			Width:  512,
			Height: 384,
		}
	}
}

func resourceToPath(r domain.Resource) string {
	switch r {
	case domain.ResourceTag:
		return "tags"
	case domain.ResourceGame:
		return "games"
	}

	return ""
}
