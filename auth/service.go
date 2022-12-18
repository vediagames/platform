package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ory "github.com/ory/kratos-client-go"
	"github.com/rs/zerolog"
	"github.com/vediagames/zeroerror"
)

type Service struct {
	ory *ory.APIClient
}

func New(url string) Service {
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{
		{
			URL: url,
		},
	}

	return Service{
		ory: ory.NewAPIClient(c),
	}
}

// save the session to display it on the dashboard
func withUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

func GetUser(ctx context.Context) *ory.Session {
	return ctx.Value(contextKeyUser).(*ory.Session)
}

func (s *Service) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Header.Get("Cookie")

			session, _, err := s.ory.FrontendApi.ToSession(r.Context()).Cookie(cookies).Execute()
			if (err != nil && session == nil) || (err == nil && !*session.Active) {
				zerolog.Ctx(r.Context()).Error().Msgf("failed to log in: %s", err)
				next.ServeHTTP(w, r)
				return
			}

			user, err := sessionToUser(session)
			if err != nil {
				zerolog.Ctx(r.Context()).Error().Msgf("failed to convert session to user: %s", err)
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				withUser(r.Context(), user),
			))
		})
	}
}

type contextKey string

const contextKeyUser contextKey = "user"

type User struct {
	ID        string
	SessionID string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Validate() error {
	var err zeroerror.Error

	if u.ID == "" {
		err.Add(fmt.Errorf("empty ID"))
	}

	if u.SessionID == "" {
		err.Add(fmt.Errorf("empty session ID"))
	}

	if u.Username == "" {
		err.Add(fmt.Errorf("empty email"))
	}

	if u.CreatedAt.IsZero() {
		err.Add(fmt.Errorf("invalid created at"))
	}

	if u.UpdatedAt.IsZero() {
		err.Add(fmt.Errorf("invalid updated at"))
	}

	return err.Err()
}

func sessionToUser(v *ory.Session) (User, error) {
	identity := v.GetIdentity()

	traits := identity.GetTraits().(map[string]any)

	user := User{
		ID:        v.Identity.GetId(),
		SessionID: v.GetId(),
		Username:  traits["username"].(string),
		Email:     traits["email"].(string),
		CreatedAt: identity.GetCreatedAt(),
		UpdatedAt: identity.GetUpdatedAt(),
	}

	if err := user.Validate(); err != nil {
		return User{}, fmt.Errorf("failed to validate user: %w", err)
	}

	return user, nil
}
