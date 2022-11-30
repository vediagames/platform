package graphql

import "fmt"

func cateogryPageUrl(ID int, slug string) string {
	return fmt.Sprintf("/category/%s?id=%d", slug, ID)
}
func gamePageUrl(slug string) string {
	return fmt.Sprintf("/game/%s", slug)
}
func tagUrl(ID int, slug, name string) string {
	return fmt.Sprintf("/tag/%s?id=%d&name=%s", slug, ID, name)
}
func fullScreenPageUrl(name, url string) string {
	return fmt.Sprintf("/game/fullscreen?name=%s&url=%s", name, url)
}

type path string
type thumbnail string

const (
	pathGame         path      = "games"
	pathTag          path      = "tags"
	thumbnail128x128 thumbnail = "thumb128x128"
	thumbnail512x384 thumbnail = "thumb512x384"
	thumbnail512x512 thumbnail = "thumb512x512"
)

func (t thumbnail) JPG() string {
	return fmt.Sprintf("%s.jpg", t)
}

func (p path) Path(slug string, file string) string {
	return fmt.Sprintf("%s/%s/%s", p, slug, file)
}

func (p path) Thumbnail(slug string, t thumbnail) (string, error) {
	switch p {
	case pathGame:
		if t == thumbnail128x128 {
			return "", fmt.Errorf("thumbnail 128x128 not available for games")
		}
	case pathTag:
		if t == thumbnail512x512 {
			return "", fmt.Errorf("thumbnail 512x512 not available for tags")
		}
	default:
		return "", fmt.Errorf("thumbnails not available for %s", p)
	}

	return fmt.Sprintf("https://images.vediagames.com/file/vg-images/%s", p.Path(slug, t.JPG())), nil
}
