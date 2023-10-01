package gamedistribution

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gosimple/slug"

	"github.com/vediagames/platform/fetcher/domain"
)

const amountOfGames = 14999

type game struct {
	Title        string        `json:"Title"`
	Md5          string        `json:"Md5"`
	Description  string        `json:"Description"`
	Instructions string        `json:"Instructions"`
	Type         string        `json:"Type"`
	Subtype      string        `json:"SubType"`
	Mobile       string        `json:"Mobile"`
	Mobilemode   interface{}   `json:"MobileMode"`
	Height       int           `json:"Height"`
	Width        int           `json:"Width"`
	HTTPS        bool          `json:"Https"`
	Status       int           `json:"Status"`
	URL          string        `json:"Url"`
	Asset        []string      `json:"Asset"`
	Category     []string      `json:"Category"`
	Tag          []string      `json:"Tag"`
	Bundle       []interface{} `json:"Bundle"`
	Company      string        `json:"Company"`
	Tubiaurl     string        `json:"TubiaUrl"`
}

func (g game) domain() domain.FetchedGame {
	slug := slug.Make(g.Title)

	return domain.FetchedGame{
		Name:        g.Title,
		URL:         fmt.Sprintf("%s?gd_sdk_referrer_url=https://vedia.games/game/%s", g.URL, slug),
		Description: g.Description,
		Controls:    g.Instructions,
		Mobile:      g.Mobile == "true",
		Height:      g.Height,
		Width:       g.Width,
		Categories:  g.Category,
		Tags:        g.Tag,
		Images:      g.Asset,
		Slug:        slug,
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

	return fmt.Sprintf("https://catalog.api.gamedistribution.com/api/v2.0/rss/All/?collection=all&categories=All&tags=All&subType=all&type=all&mobile=all&rewarded=all&amount=%d&page=%d&format=json", s.Limit, page)
}
