package fetcher

import (
	"fmt"
	"math/rand"

	"github.com/vediagames/vediagames.com/fetcher/domain"
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

func New(cfg Config) (domain.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &client{
		clients: cfg.Clients,
	}, nil
}

func (c client) Fetch() (domain.FetchedGame, error) {
	client := c.clients[rand.Intn(len(c.clients))]

	game, err := client.Fetch()
	if err != nil {
		return domain.FetchedGame{}, fmt.Errorf("error fetching game: %v", err)
	}

	return game, nil
}
