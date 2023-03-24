package gamemonetize

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/vediagames/platform/fetcher/domain"
)

const amountOfGames = 19009
const primaryResolution = "512x384"

var imgResolutions = []string{
	"512x512",
	"512x384",
	"512x340",
}

type game struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
	URL          string `json:"url"`
	Category     string `json:"category"`
	Tags         string `json:"tags"`
	Thumb        string `json:"thumb"`
	Width        string `json:"width"`
	Height       string `json:"height"`
}

func (g game) domain() domain.FetchedGame {
	width, err := strconv.Atoi(g.Width)
	if err != nil {
		width = 0
	}

	height, err := strconv.Atoi(g.Height)
	if err != nil {
		height = 0
	}

	return domain.FetchedGame{
		Name:        g.Title,
		URL:         g.URL,
		Description: g.Description,
		Controls:    g.Instructions,
		Mobile:      false,
		Height:      width,
		Width:       height,
		Categories:  strings.Split(g.Category, ", "),
		Tags:        strings.Split(g.Tags, ", "),
		Images:      getImages(g.Thumb),
	}
}

type client struct {
	Limit int
}

func New(limit int) domain.Client {
	return client{
		Limit: limit,
	}
}

func (s client) Fetch() (domain.FetchedGame, error) {
	var games []game

	r, err := http.Get(s.getUrl())
	if err != nil {
		return domain.FetchedGame{}, fmt.Errorf("failed to fetch games: %v", err)
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&games); err != nil {
		return domain.FetchedGame{}, fmt.Errorf("failed to decode games: %v", err)
	}

	if len(games) == 0 {
		return domain.FetchedGame{}, domain.ErrNoData
	}

	randGame := games[rand.Intn(len(games))]
	return randGame.domain(), nil

}

func (s client) getUrl() string {
	maxPages := amountOfGames / s.Limit

	page := rand.Intn(maxPages + 1)

	return fmt.Sprintf("https://gamemonetize.com/feed.php?format=0&num=%d&page=%d", s.Limit, page)
}

func getImages(thumb string) []string {
	var images []string

	images = append(images, thumb)

	for _, resolution := range imgResolutions {
		if resolution != primaryResolution {
			url := strings.ReplaceAll(thumb, primaryResolution, resolution)
			r, err := http.Get(url)
			if err == nil && r.Header.Get("content-type") == "image/jpeg" {
				images = append(images, url)
			}
		}
	}

	return images
}
