package service

import (
	"context"
	"fmt"

	bucketdomain "github.com/vediagames/vediagames.com/bucket/domain"
	"github.com/vediagames/vediagames.com/config"
	"github.com/vediagames/vediagames.com/image/domain"
)

type service struct {
	processor domain.Processor
	storage   bucketdomain.Client
	CDNURLs   CDNURLs
}

type CDNURLs struct {
	s3    string
	bunny string
}
type Config struct {
	CDNURLs   CDNURLs
	Processor domain.Processor
	Cfg       config.Imagor
	Storage   bucketdomain.Client
}

func (c Config) Validate() error {
	if c.Cfg.Secret == "" {
		return fmt.Errorf("secret required")
	}

	if c.Cfg.URL == "" {
		return fmt.Errorf("url required")
	}

	if c.Processor == nil {
		return fmt.Errorf("image processor service required")
	}

	if c.Storage == nil {
		return fmt.Errorf("image stoage service required")
	}

	if c.CDNURLs.bunny == "" {
		return fmt.Errorf("bunny cdn url required")
	}

	if c.CDNURLs.s3 == "" {
		return fmt.Errorf("s3 cdn url required")
	}

	return nil
}

func New(config Config) (domain.Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &service{
		processor: config.Processor,
		storage:   config.Storage,
	}, nil
}

func (s *service) GetThumbnail(ctx context.Context, req domain.GetThumbnailRequest) (string, error) {
	// the resolution format already exists on S3
	thumb, err := domain.GetExistingThumbnail(req)
	if err != nil {
		return "", fmt.Errorf("failed to generate url for existing image on S3: %w", err)
	}
	if req.Thumbnail.IsDefault {
		// original requested
		return fmt.Sprintf("%s%s", s.CDNURLs.s3, thumb), err
	}
	// process new image as requested resolution & format
	image, err := s.processor.Process(ctx, req, thumb)

	// upload the processed image to bunny storage with the same path
	path := domain.GetImagePath(req.Path, req.Slug, req.Thumbnail.Format.String())
	if err := s.storage.Upload(ctx, path, image); err != nil {
		return "", fmt.Errorf("failed to upload the processed image to storage: %w", err)
	}
	//TODO:	return url uploaded imageURL
	return fmt.Sprintf("%s%s", s.CDNURLs.bunny, thumb), err
}
