package domain

import "fmt"

type path string
type thumbnail string

const (
	PathGame         path      = "games"
	PathTag          path      = "tags"
	Thumbnail128x128 thumbnail = "thumb128x128"
	Thumbnail512x384 thumbnail = "thumb512x384"
	Thumbnail512x512 thumbnail = "thumb512x512"
)

func (t thumbnail) JPG() string {
	return fmt.Sprintf("%s.jpg", t)
}

func (p path) Path(slug string, file string) string {
	return fmt.Sprintf("%s/%s/%s", p, slug, file)
}

func (p path) Thumbnail(slug string, t thumbnail) (string, error) {
	switch p {
	case PathGame:
		if t == Thumbnail128x128 {
			return "", fmt.Errorf("thumbnail 128x128 not available for games")
		}
	case PathTag:
		if t == Thumbnail512x512 {
			return "", fmt.Errorf("thumbnail 512x512 not available for tags")
		}
	default:
		return "", fmt.Errorf("thumbnails not available for %s", p)
	}

	return fmt.Sprintf("https://images.vediagames.com/file/vg-images/%s", p.Path(slug, t.JPG())), nil
}

func GetExistingThumbnail(req GetThumbnailRequest) (string, bool, error) {
	imageName := thumbnail(fmt.Sprintf("thumb%dx%d", req.Thumbnail.Width, req.Thumbnail.Height))
	var p path = PathGame
	if path(req.Path) == PathTag {
		p = PathTag
	}

	imageURL, err := p.Thumbnail(req.Slug, thumbnail(imageName))

	if req.Thumbnail.Format == FormatJpeg {
		switch imageName {
		case Thumbnail128x128, Thumbnail512x384, Thumbnail512x512:
			return imageURL, true, err
		}
	}
	return imageURL, false, err
}
