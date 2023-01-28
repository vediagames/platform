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

	return GetImagePath(string(p), slug, t.JPG()), nil
}

func GetImagePath(path, slug, fileFormat string) string {
	return fmt.Sprintf("%s/%s/%s", path, slug, fileFormat)
}

func GetExistingThumbnail(req GetThumbnailRequest) (string, error) {
	var p path = PathGame
	var imageName = fmt.Sprintf("thumbnail%s", req.Thumbnail.Original.String())

	if path(req.Path) == PathTag {
		p = PathTag
	}

	imageURL, err := p.Thumbnail(req.Slug, thumbnail(imageName))
	if err != nil {
		return "", fmt.Errorf("failed to get image s3 default url: %w", err)
	}

	return imageURL, err
}
