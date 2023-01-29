package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	GetThumbnail(context.Context, GetThumbnailRequest) (GetThumbnailResponse, error)
}

type GetThumbnailRequest struct {
	Path      string
	Slug      string
	Thumbnail Thumbnail
}

func (r GetThumbnailRequest) Validate() error {
	var err zeroerror.Error

	if r.Path == "" {
		err.Add(fmt.Errorf("empty path"))
	}

	if r.Slug == "" {
		err.Add(fmt.Errorf("empty slug"))
	}

	if r.Thumbnail.Format.IsValid() {
		err.Add(fmt.Errorf("invalid format"))
	}

	if r.Thumbnail.Height == 0 {
		err.Add(fmt.Errorf("invalid Thumbnail.Height"))
	}

	if r.Thumbnail.Width == 0 {
		err.Add(fmt.Errorf("invalid Thumbnail.Width"))
	}

	return err.Err()
}

type GetThumbnailResponse struct {
	ImageURL string
}

type Thumbnail struct {
	Original  Original
	Width     int
	Height    int
	Format    Format
	IsDefault bool
}

type Format string

const (
	FormatWebp Format = "webp"
	FormatJpeg Format = "jpeg"
	FormatPng  Format = "png"
)

var AllFormat = []Format{
	FormatWebp,
	FormatJpeg,
	FormatPng,
}

func (e Format) IsValid() bool {
	switch e {
	case FormatWebp, FormatJpeg, FormatPng:
		return true
	}
	return false
}

func (e Format) String() string {
	return string(e)
}

type Original string

const (
	OriginalThumbnail512x512 Original = "thumbnail512x512"
	OriginalThumbnail512x384 Original = "thumbnail512x384"
	OriginalThumbnail128x128 Original = "thumbnail128x128"
)

var AllOriginal = []Original{
	OriginalThumbnail512x512,
	OriginalThumbnail512x384,
	OriginalThumbnail128x128,
}

func (e Original) IsValid() bool {
	switch e {
	case OriginalThumbnail512x512, OriginalThumbnail512x384, OriginalThumbnail128x128:
		return true
	}
	return false
}

func (e Original) String() string {
	return string(e)
}
