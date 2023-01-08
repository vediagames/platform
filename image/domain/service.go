package domain

import (
	"context"
	"io"
)

type Processor interface {
	Process(context.Context, GetThumbnailRequest, string) (io.Reader, error)
}

type Service interface {
	GetThumbnail(context.Context, GetThumbnailRequest) (string, error)
}

type GetThumbnailRequest struct {
	Path      string
	Slug      string
	Thumbnail Thumbnail
}
type Thumbnail struct {
	Original Original `json:"original"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Format   Format   `json:"format"`
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
