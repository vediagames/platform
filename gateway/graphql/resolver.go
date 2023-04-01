package graphql

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"

	authdomain "github.com/vediagames/platform/auth/domain"
	bucketdomain "github.com/vediagames/platform/bucket/domain"
	categorydomain "github.com/vediagames/platform/category/domain"
	fetcherdomain "github.com/vediagames/platform/fetcher/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	"github.com/vediagames/platform/gateway/graphql/generated"
	notificationdomain "github.com/vediagames/platform/notification/domain"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	gameService     gamedomain.Service
	categoryService categorydomain.Service
	sectionService  sectiondomain.Service
	tagService      tagdomain.Service
	searchService   searchdomain.Service
	emailClient     notificationdomain.EmailClient
	bucketClient    bucketdomain.Client
	fetcherClient   fetcherdomain.Client
	authService     authdomain.Service
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
	AuthService     authdomain.Service
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

	if c.AuthService == nil {
		return fmt.Errorf("auth client is required")
	}

	return nil
}

func NewResolver(cfg Config) *Resolver {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &Resolver{
		gameService:     cfg.GameService,
		categoryService: cfg.CategoryService,
		sectionService:  cfg.SectionService,
		tagService:      cfg.TagService,
		searchService:   cfg.SearchService,
		emailClient:     cfg.EmailClient,
		bucketClient:    cfg.BucketClient,
		fetcherClient:   cfg.FetcherClient,
		authService:     cfg.AuthService,
	}
}

func NewSchema(r *Resolver) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: r,
	})
}
