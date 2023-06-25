package graphql

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vediagames/zeroerror"

	authdomain "github.com/vediagames/platform/auth/domain"
	bucketdomain "github.com/vediagames/platform/bucket/domain"
	categorydomain "github.com/vediagames/platform/category/domain"
	fetcherdomain "github.com/vediagames/platform/fetcher/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	"github.com/vediagames/platform/gateway/graphql/generated"
	imagedomain "github.com/vediagames/platform/image/domain"
	notificationdomain "github.com/vediagames/platform/notification/domain"
	"github.com/vediagames/platform/quote"
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
	imageService    imagedomain.Service
	contentURL      string
	quoteService    quote.Service
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
	ImageService    imagedomain.Service
	ContentURL      string
	QuoteService    quote.Service
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.GameService == nil, fmt.Errorf("game service is required"))
	err.AddIf(c.CategoryService == nil, fmt.Errorf("category service is required"))
	err.AddIf(c.SectionService == nil, fmt.Errorf("section service is required"))
	err.AddIf(c.TagService == nil, fmt.Errorf("tag service is required"))
	err.AddIf(c.SearchService == nil, fmt.Errorf("search service is required"))
	err.AddIf(c.EmailClient == nil, fmt.Errorf("email client is required"))
	err.AddIf(c.BucketClient == nil, fmt.Errorf("bucket client is required"))
	err.AddIf(c.FetcherClient == nil, fmt.Errorf("fetcher client is required"))
	err.AddIf(c.AuthService == nil, fmt.Errorf("auth service is required"))
	err.AddIf(c.ImageService == nil, fmt.Errorf("image service is required"))
	err.AddIf(c.ContentURL == "", fmt.Errorf("content URL is required"))
	err.AddIf(c.QuoteService == nil, fmt.Errorf("quote service is required"))

	return err.Err()
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
		imageService:    cfg.ImageService,
		contentURL:      cfg.ContentURL,
		quoteService:    cfg.QuoteService,
	}
}

func NewSchema(r *Resolver) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: r,
	})
}
