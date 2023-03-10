package domain

import (
	"fmt"

	"github.com/vediagames/zeroerror"
)

type Image struct {
	Format Format
	Width  int
	Height int
}

func (i Image) Validate() error {
	var err zeroerror.Error

	err.AddIf(i.Width <= 0, fmt.Errorf("invalid width"))
	err.AddIf(i.Height <= 0, fmt.Errorf("invalid height"))

	err.AddIf(i.Width > 2000, fmt.Errorf("too large width"))
	err.AddIf(i.Height > 2000, fmt.Errorf("too large height"))

	if ve := i.Format.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid format: %w", ve))
	}

	return err.Err()
}

func (i Image) File() string {
	return fmt.Sprintf("thumb%dx%d.%s", i.Width, i.Height, i.Format.String())
}

type Resource string

func (f Resource) Validate() error {
	switch f {
	case ResourceTag, ResourceGame:
		return nil
	default:
		return fmt.Errorf("invalid value: %q", f)
	}
}

const (
	ResourceGame = Resource("game")
	ResourceTag  = Resource("tag")
)

type Format string

func (f Format) Validate() error {
	switch f {
	case FormatWebp, FormatJpg, FormatPng:
		return nil
	default:
		return fmt.Errorf("invalid value: %q", f)
	}
}

func (f Format) String() string {
	return string(f)
}

const (
	FormatWebp = Format("webp")
	FormatJpg  = Format("jpg")
	FormatPng  = Format("png")
)

type OriginalThumbnail string

const (
	OriginalThumbnail512x512 = OriginalThumbnail("thumbnail512x512")
	OriginalThumbnail512x384 = OriginalThumbnail("thumbnail512x384")
	OriginalThumbnail128x128 = OriginalThumbnail("thumbnail128x128")
)

func (o OriginalThumbnail) Validate() error {
	switch o {
	case OriginalThumbnail512x512, OriginalThumbnail512x384, OriginalThumbnail128x128:
		return nil
	default:
		return fmt.Errorf("invalid value: %q", o)
	}
}

func (o OriginalThumbnail) String() string {
	return string(o)
}
