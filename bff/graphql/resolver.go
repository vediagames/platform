// go generate can be run by: go generate ./...
//go:generate go run github.com/99designs/gqlgen generate

package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vediagames/vediagames.com/bff/graphql/generated"
	bucketdomain "github.com/vediagames/vediagames.com/bucket/domain"
	categorydomain "github.com/vediagames/vediagames.com/category/domain"
	fetcherdomain "github.com/vediagames/vediagames.com/fetcher/domain"
	gamedomain "github.com/vediagames/vediagames.com/game/domain"
	notificationdomain "github.com/vediagames/vediagames.com/notification/domain"
	searchdomain "github.com/vediagames/vediagames.com/search/domain"
	sectiondomain "github.com/vediagames/vediagames.com/section/domain"
	tagdomain "github.com/vediagames/vediagames.com/tag/domain"
)

type Resolver struct {
	gameService     gamedomain.Service
	categoryService categorydomain.Service
	sectionService  sectiondomain.Service
	tagService      tagdomain.Service
	searchService   searchdomain.Service
	emailClient     notificationdomain.EmailClient
	bucketClient    bucketdomain.Client
	fetcherClient   fetcherdomain.Client
}

type Config struct {
	GameService     gamedomain.Service
	CategoryService categorydomain.Service
	SectionService  sectiondomain.Service
	TagService      tagdomain.Service
	SearchService   searchdomain.Service
	EmailClient     notificationdomain.EmailClient
	BucketClient    bucketdomain.Client
	FetcherClient   fetcherdomain.Client
}

func (c Config) Validate() error {
	if c.GameService == nil {
		return fmt.Errorf("game service is required")
	}

	if c.CategoryService == nil {
		return fmt.Errorf("category service is required")
	}

	if c.SectionService == nil {
		return fmt.Errorf("section service is required")
	}

	if c.TagService == nil {
		return fmt.Errorf("tag service is required")
	}

	if c.SearchService == nil {
		return fmt.Errorf("search service is required")
	}

	if c.EmailClient == nil {
		return fmt.Errorf("email client is required")
	}

	if c.BucketClient == nil {
		return fmt.Errorf("bucket client is required")
	}

	if c.FetcherClient == nil {
		return fmt.Errorf("fetcher client is required")
	}

	return nil
}

func NewResolver(cfg Config) (Resolver, error) {
	if err := cfg.Validate(); err != nil {
		return Resolver{}, fmt.Errorf("invalid config: %w", err)
	}

	return Resolver{
		gameService:     cfg.GameService,
		categoryService: cfg.CategoryService,
		sectionService:  cfg.SectionService,
		tagService:      cfg.TagService,
		searchService:   cfg.SearchService,
		emailClient:     cfg.EmailClient,
		bucketClient:    cfg.BucketClient,
		fetcherClient:   cfg.FetcherClient,
	}, nil
}

func NewSchema(r *Resolver) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: r,
	})
}
