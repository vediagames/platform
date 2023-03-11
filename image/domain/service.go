package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	Get(context.Context, GetRequest) (GetResponse, error)
}

type GetRequest struct {
	Slug     string
	Image    Image
	Original OriginalThumbnail
	Resource Resource
}

func (r GetRequest) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.Slug == "", fmt.Errorf("empty slug"))

	if ve := r.Image.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid image: %w", ve))
	}

	if ve := r.Original.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid original: %w", ve))
	}

	if ve := r.Resource.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid resource: %w", ve))
	}

	switch r.Resource {
	case ResourceGame:
		if r.Original == OriginalThumbnail128x128 {
			err.Add(fmt.Errorf("thumbnail 128x128 not available for games"))
		}
	case ResourceTag:
		if r.Original == OriginalThumbnail512x512 {
			err.Add(fmt.Errorf("thumbnail 512x512 not available for tags"))
		}
	}

	return err.Err()
}

type GetResponse struct {
	URL string
}

func (r GetResponse) Validate() error {
	var err zeroerror.Error

	err.AddIf(r.URL == "", fmt.Errorf("empty URL"))

	return err.Err()
}
