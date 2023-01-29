package service

import (
	"context"
	"fmt"

	bucketdomain "github.com/vediagames/vediagames.com/bucket/domain"
	"github.com/vediagames/vediagames.com/config"
	"github.com/vediagames/vediagames.com/image/domain"
	"github.com/vediagames/zeroerror"
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
	ImagorCfg config.Imagor
	Storage   bucketdomain.Client
}

func (c Config) Validate() error {
	var err zeroerror.Error

	if c.ImagorCfg.Secret == "" {
		err.Add(fmt.Errorf("empty imagor secret"))
	}

	if c.ImagorCfg.URL == "" {
		err.Add(fmt.Errorf("empty imagor url"))
	}

	if c.Processor == nil {
		err.Add(fmt.Errorf("empty processor"))
	}

	if c.Storage == nil {
		err.Add(fmt.Errorf("empty storage"))
	}

	if c.CDNURLs.bunny == "" {
		err.Add(fmt.Errorf("empty bunny cdn url"))
	}

	if c.CDNURLs.s3 == "" {
		err.Add(fmt.Errorf("empty s3 cdn url"))
	}

	return err.Err()
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

func (s *service) GetThumbnail(ctx context.Context, req domain.GetThumbnailRequest) (domain.GetThumbnailResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetThumbnailResponse{}, fmt.Errorf("invalid config: %w", err)
	}

	thumb, err := domain.GetExistingThumbnail(req)
	if err != nil {
		return domain.GetThumbnailResponse{}, fmt.Errorf("failed to generate url for existing image on S3: %w", err)
	}

	if req.Thumbnail.IsDefault {
		return domain.GetThumbnailResponse{
			ImageURL: fmt.Sprintf("%s/%s", s.CDNURLs.s3, thumb),
		}, err
	}

	image, err := s.processor.Process(ctx, domain.ProcessQuery{
		Thumbnail: req.Thumbnail,
		ImageURL:  thumb,
	})

	path := domain.GetImagePath(req.Path, req.Slug, req.Thumbnail.Format.String())
	if err := s.storage.Upload(ctx, path, image.File); err != nil {
		return domain.GetThumbnailResponse{}, fmt.Errorf("failed to upload the processed image to storage: %w", err)
	}

	//TODO:	return url uploaded imageURL
	return domain.GetThumbnailResponse{
		ImageURL: fmt.Sprintf("%s/%s", s.CDNURLs.bunny, thumb),
	}, err
}
