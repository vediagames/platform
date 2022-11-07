package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/vediagames/vediagames.com/bff/graphql"
	"github.com/vediagames/vediagames.com/bucket/s3"
	categorypostgresql "github.com/vediagames/vediagames.com/category/postgresql"
	categoryservice "github.com/vediagames/vediagames.com/category/service"
	"github.com/vediagames/vediagames.com/config"
	"github.com/vediagames/vediagames.com/fetcher"
	fetcherdomain "github.com/vediagames/vediagames.com/fetcher/domain"
	"github.com/vediagames/vediagames.com/fetcher/gamedistribution"
	"github.com/vediagames/vediagames.com/fetcher/gamemonetize"
	gamepostgresql "github.com/vediagames/vediagames.com/game/postgresql"
	gameservice "github.com/vediagames/vediagames.com/game/service"
	"github.com/vediagames/vediagames.com/notification/sendinblue"
	searchservice "github.com/vediagames/vediagames.com/search/service"
	sectionpostgresql "github.com/vediagames/vediagames.com/section/postgresql"
	sectionservice "github.com/vediagames/vediagames.com/section/service"
	sectionvalidationdata "github.com/vediagames/vediagames.com/section/service/validation/data"
	sectionvalidationrequest "github.com/vediagames/vediagames.com/section/service/validation/request"
	tagpostgresql "github.com/vediagames/vediagames.com/tag/postgresql"
	tagservice "github.com/vediagames/vediagames.com/tag/service"
)

func ServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer(cmd.Context())
		},
	}
}

func startServer(ctx context.Context) error {
	cfg := ctx.Value(config.Key).(config.Config)

	db, err := sqlx.Open("postgres", cfg.PostgreSQL.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}

	gameRepository, err := gamepostgresql.New(gamepostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create game repository: %w", err)
	}

	gameStatsRepository, err := gamepostgresql.NewStatsRepository(gamepostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create game stats repository: %w", err)
	}

	gameEventRepository, err := gamepostgresql.NewEventRepository(gamepostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create game event repository: %w", err)
	}

	gameService, err := gameservice.New(gameservice.Config{
		Repository:      gameRepository,
		StatsRepository: gameStatsRepository,
		EventRepository: gameEventRepository,
	})
	if err != nil {
		return fmt.Errorf("failed to create game service: %w", err)
	}

	categoryRepository, err := categorypostgresql.New(categorypostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create category repository: %w", err)
	}

	categoryService, err := categoryservice.New(categoryservice.Config{
		Repository: categoryRepository,
	})
	if err != nil {
		return fmt.Errorf("failed to create category service: %w", err)
	}

	sectionRepository, err := sectionpostgresql.New(sectionpostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create section repository: %w", err)
	}

	websitePlacementRepository, err := sectionpostgresql.NewWebsitePlacementRepository(sectionpostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create website placement repository: %w", err)
	}

	sectionService, err := sectionservice.New(sectionservice.Config{
		Repository:                 sectionRepository,
		WebsitePlacementRepository: websitePlacementRepository,
	})
	if err != nil {
		return fmt.Errorf("failed to create section service: %w", err)
	}

	sectionValidationData, err := sectionvalidationdata.New(sectionvalidationdata.Config{
		Service: sectionService,
	})
	if err != nil {
		return fmt.Errorf("failed to create section validation data: %w", err)
	}

	sectionValidationRequest, err := sectionvalidationrequest.New(sectionvalidationrequest.Config{
		Service: sectionValidationData,
	})
	if err != nil {
		return fmt.Errorf("failed to create section validation request: %w", err)
	}

	tagRepository, err := tagpostgresql.New(tagpostgresql.Config{
		DB: db,
	})
	if err != nil {
		return fmt.Errorf("failed to create tag repository: %w", err)
	}

	tagService, err := tagservice.New(tagservice.Config{
		Repository: tagRepository,
	})
	if err != nil {
		return fmt.Errorf("failed to create tag service: %w", err)
	}

	searchService, err := searchservice.New(searchservice.Config{
		TagService:  tagService,
		GameService: gameService,
	})
	if err != nil {
		return fmt.Errorf("failed to create search service: %w", err)
	}

	emailClient := sendinblue.New(http.Client{
		Timeout: 10 * time.Second,
	}, cfg.SendInBlue.Key)

	fetcherClient, err := fetcher.New(fetcher.Config{
		Clients: []fetcherdomain.Client{
			gamedistribution.New(10),
			gamemonetize.New(10),
		},
	})

	bucketClient, err := s3.New(s3.Config{
		Key:      cfg.Bucket.Key,
		Secret:   cfg.Bucket.Secret,
		Region:   cfg.Bucket.Region,
		EndPoint: cfg.Bucket.EndPoint,
		Bucket:   cfg.Bucket.Bucket,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket client: %w", err)
	}

	cache, err := NewCache(ctx, cfg.RedisAddress, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to create cache")
	}

	resolver, err := graphql.NewResolver(graphql.Config{
		GameService:     gameService,
		CategoryService: categoryService,
		SectionService:  sectionValidationRequest,
		TagService:      tagService,
		SearchService:   searchService,
		EmailClient:     emailClient,
		BucketClient:    bucketClient,
		FetcherClient:   fetcherClient,
	})
	if err != nil {
		return fmt.Errorf("failed to create resolver %w", err)
	}

	gqlHandler := handler.New(graphql.NewSchema(&resolver))
	gqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{Cache: cache})

	httpCors := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowCredentials: true,
		Debug:            cfg.LogLevel == "debug",
	})

	router := chi.NewRouter()

	logger := zerolog.Ctx(ctx).With().Str("transport", "http").Logger()

	router.Use(httpCors.Handler)
	router.Use(loggerMiddleware(&logger))

	router.Handle("/graph", gqlHandler)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(r.Context()).Log().Msg("HELLO")
		w.WriteHeader(http.StatusOK)
	})

	logger.Info().
		Int("port", cfg.Port).
		Msgf("starting server on port %d", cfg.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func loggerMiddleware(logger *zerolog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ip := r.Header.Get("CF-Connecting-IP")
			if ip == "" {
				ip = r.Header.Get("X-Forwarded-For")

				if ip == "" {
					ip = r.RemoteAddr
				}
			}

			r.Header.Set("Real-IP", ip)

			l := logger.With().
				Str("method", r.Method).
				Str("url", r.RequestURI).
				Interface("client_ip", ip).
				Logger()

			ctx := l.WithContext(r.Context())

			h.ServeHTTP(w, r.WithContext(ctx))

			if ctx.Err() != nil {
				l.Error().
					Err(ctx.Err()).
					Msgf("failed request: %w", ctx.Err())
			}

			l.Info().TimeDiff("latency", time.Now(), start).Msg("finished request")
		})
	}
}

type Cache struct {
	client    redis.UniversalClient
	ttl       time.Duration
	apqPrefix string
}

func NewCache(ctx context.Context, redisAddress string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Cache{client: client, ttl: ttl}, nil
}

func (c Cache) Add(ctx context.Context, key string, value interface{}) {
	key = fmt.Sprintf("%s:%s", c.apqPrefix, key)

	_, err := c.client.Set(ctx, key, value, c.ttl).Result()
	if err != nil {
		zerolog.Ctx(ctx).Err(fmt.Errorf("failed to set cache for key %q: %w", key, err))
	}
}

func (c Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	key = fmt.Sprintf("%s:%s", c.apqPrefix, key)

	s, err := c.client.Get(ctx, key).Result()
	if err != nil {
		zerolog.Ctx(ctx).Err(fmt.Errorf("failed to get cache for key %q: %w", key, err))
		return nil, false
	}

	return s, true
}
