package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	ory "github.com/ory/kratos-client-go"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"

	authdomain "github.com/vediagames/platform/auth/domain"
	authservice "github.com/vediagames/platform/auth/service"
	"github.com/vediagames/platform/bff/graphql"
	"github.com/vediagames/platform/bucket/bunny"
	categorypostgresql "github.com/vediagames/platform/category/postgresql"
	categoryservice "github.com/vediagames/platform/category/service"
	"github.com/vediagames/platform/config"
	"github.com/vediagames/platform/fetcher"
	fetcherdomain "github.com/vediagames/platform/fetcher/domain"
	"github.com/vediagames/platform/fetcher/gamedistribution"
	"github.com/vediagames/platform/fetcher/gamemonetize"
	gamepostgresql "github.com/vediagames/platform/game/postgresql"
	gameservice "github.com/vediagames/platform/game/service"
	"github.com/vediagames/platform/notification/sendinblue"
	searchservice "github.com/vediagames/platform/search/service"
	sectionpostgresql "github.com/vediagames/platform/section/postgresql"
	sectionservice "github.com/vediagames/platform/section/service"
	sectionvalidationdata "github.com/vediagames/platform/section/service/validation/data"
	sectionvalidationrequest "github.com/vediagames/platform/section/service/validation/request"
	sessionbigquery "github.com/vediagames/platform/session/bigquery"
	sessiondomain "github.com/vediagames/platform/session/domain"
	sessionservice "github.com/vediagames/platform/session/service"
	tagpostgresql "github.com/vediagames/platform/tag/postgresql"
	tagservice "github.com/vediagames/platform/tag/service"
)

const (
	tableID   = "tableID"
	datasetID = "datasetID"
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
	cfg := ctx.Value(config.ContextKey).(config.Config)

	db, err := sqlx.Open("postgres", cfg.PostgreSQL.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}

	gameService := gameservice.New(gameservice.Config{
		Repository: gamepostgresql.New(gamepostgresql.Config{
			DB: db,
		}),
		StatsRepository: gamepostgresql.NewStatsRepository(gamepostgresql.Config{
			DB: db,
		}),
		EventRepository: gamepostgresql.NewEventRepository(gamepostgresql.Config{
			DB: db,
		}),
	})

	categoryService := categoryservice.New(categoryservice.Config{
		Repository: categorypostgresql.New(categorypostgresql.Config{
			DB: db,
		}),
	})

	sectionService := sectionservice.New(sectionservice.Config{
		Repository: sectionpostgresql.New(sectionpostgresql.Config{
			DB: db,
		}),
		WebsitePlacementRepository: sectionpostgresql.NewWebsitePlacementRepository(sectionpostgresql.Config{
			DB: db,
		}),
	})

	sectionService = sectionvalidationrequest.New(sectionvalidationdata.New(sectionService))

	tagService := tagservice.New(tagservice.Config{
		Repository: tagpostgresql.New(tagpostgresql.Config{
			DB: db,
		}),
	})

	searchService := searchservice.New(searchservice.Config{
		TagService:  tagService,
		GameService: gameService,
	})

	client, err := bigquery.NewClient(ctx, cfg.BigQuery.ProjectID,
		option.WithCredentialsFile(cfg.BigQuery.CredentialsPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create bigquery client: %w", err)
	}

	sessionService := sessionservice.New(sessionservice.Config{
		Repository: sessionbigquery.New(sessionbigquery.Config{
			Client:    client,
			TableID:   tableID,
			DatasetID: datasetID,
		}),
	})

	emailClient := sendinblue.New(sendinblue.Config{
		Token: cfg.SendInBlue.Key,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	})

	fetcherClient := fetcher.New(fetcher.Config{
		Clients: []fetcherdomain.Client{
			gamedistribution.New(10),
			gamemonetize.New(10),
		},
	})

	bucketClient := bunny.New(bunny.Config{
		URL:       cfg.BunnyStorage.URL,
		AccessKey: cfg.BunnyStorage.AccessKey,
		Zone:      cfg.BunnyStorage.Zone,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	})

	cache := graphql.NewCache(ctx, cfg.RedisAddress, 24*time.Hour)

	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{
		{
			URL: cfg.Auth.KratosURL,
		},
	}

	authService := authservice.NewOry(authservice.OryConfig{
		Client: ory.NewAPIClient(c),
	})

	resolver := graphql.NewResolver(graphql.Config{
		GameService:     gameService,
		CategoryService: categoryService,
		SectionService:  sectionService,
		TagService:      tagService,
		SearchService:   searchService,
		EmailClient:     emailClient,
		BucketClient:    bucketClient,
		FetcherClient:   fetcherClient,
		AuthService:     authService,
	})

	gqlHandler := handler.New(graphql.NewSchema(&resolver))
	gqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: cache,
	})

	httpCors := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowCredentials: true,
		Debug:            cfg.LogLevel == "debug",
	})

	router := chi.NewRouter()

	logger := zerolog.Ctx(ctx).With().Str("transport", "http").Logger()

	router.Use(httpCors.Handler)
	router.Use(loggerMiddleware(&logger))
	router.Use(authMiddleware(authService))

	router.Handle("/graph", gqlHandler)
	router.Handle("/session/new", createSessionHandler(sessionService))
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
					Msgf("failed request: %e", ctx.Err())
			}

			l.Info().TimeDiff("latency", time.Now(), start).Msg("finished request")
		})
	}
}

func authMiddleware(s authdomain.Service) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Header.Get("Cookie")

			res, err := s.Authenticate(r.Context(), authdomain.AuthenticateRequest{
				Cookies: cookies,
			})
			if err != nil {
				zerolog.Ctx(r.Context()).Error().Msgf("failed to authenticate: %s", err)
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				s.ToContext(r.Context(), res.User),
			))
		})
	}
}

type sessionNewResponse struct {
	ID         string    `json:"id"`
	IP         string    `json:"ip"`
	Device     string    `json:"device"`
	PageURL    string    `json:"page_url"`
	CreatedAt  time.Time `json:"createdAt"`
	InsertedAt time.Time `json:"insertedAt"`
}

type sessionNewRequest struct {
	IP        string `json:"ip"`
	Device    string `json:"device"`
	PageURL   string `json:"page_url"`
	CreatedAt string `json:"created_at"`
}

func createSessionHandler(s sessiondomain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sessionNewRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to decode: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to parse: %s", err)
			http.Error(w, sessiondomain.ErrInvalidCreatedAt.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.Create(r.Context(), sessiondomain.CreateRequest{
			IP:        sessiondomain.IP(req.IP),
			Device:    sessiondomain.Device(req.Device),
			PageURL:   req.PageURL,
			CreatedAt: createdAt,
		})
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to create: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonRes, err := json.Marshal(sessionNewResponse{
			ID:         res.Session.ID,
			IP:         res.Session.IP.String(),
			Device:     res.Session.Device.String(),
			PageURL:    res.Session.PageURL,
			CreatedAt:  res.Session.CreatedAt,
			InsertedAt: res.Session.InsertedAt,
		})
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to marshal: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write(jsonRes); err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to write: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
