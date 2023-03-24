package fetcher

import (
	"fmt"
	"math/rand"

	"github.com/vediagames/platform/fetcher/domain"
)

type client struct {
	clients []domain.Client
}

type Config struct {
	Clients []domain.Client
}

func (c Config) Validate() error {
	if len(c.Clients) == 0 {
		return fmt.Errorf("no fetchers registered")
	}

	return nil
}

func New(cfg Config) domain.Client {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &client{
		clients: cfg.Clients,
	}
}

func (c client) Fetch() (domain.FetchedGame, error) {
	client := c.clients[rand.Intn(len(c.clients))]

	game, err := client.Fetch()
	if err != nil {
		return domain.FetchedGame{}, fmt.Errorf("error fetching game: %v", err)
	}

	return game, nil
}
