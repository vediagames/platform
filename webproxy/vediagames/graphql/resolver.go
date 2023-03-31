package graphql

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"

	authdomain "github.com/vediagames/platform/auth/domain"
	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	gateway "github.com/vediagames/platform/gateway/graphql/generated"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
	"github.com/vediagames/platform/webproxy/vediagames/graphql/generated"
)

type Resolver struct {
	gameService     gamedomain.Service
	categoryService categorydomain.Service
	sectionService  sectiondomain.Service
	tagService      tagdomain.Service
	searchService   searchdomain.Service
	authService     authdomain.Service
	gatewayResolver gateway.QueryResolver
}

type Config struct {
	GameService     gamedomain.Service
	CategoryService categorydomain.Service
	SectionService  sectiondomain.Service
	TagService      tagdomain.Service
	SearchService   searchdomain.Service
	AuthService     authdomain.Service
	GatewayResolver gateway.QueryResolver
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

	if c.AuthService == nil {
		return fmt.Errorf("auth client is required")
	}

	if c.GatewayResolver == nil {
		return fmt.Errorf("gateway resolver is required")
	}

	return nil
}

func NewResolver(cfg Config) Resolver {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return Resolver{
		gameService:     cfg.GameService,
		categoryService: cfg.CategoryService,
		sectionService:  cfg.SectionService,
		tagService:      cfg.TagService,
		searchService:   cfg.SearchService,
		authService:     cfg.AuthService,
		gatewayResolver: cfg.GatewayResolver,
	}
}

func NewSchema(r *Resolver) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: r,
	})
}
