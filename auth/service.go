package auth

import (
	"context"
	"log"
	"net/http"

	ory "github.com/ory/kratos-client-go"
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
func withUser(ctx context.Context, v *ory.Session) context.Context {
	return context.WithValue(ctx, userCtxKey, v)
}

func GetUser(ctx context.Context) *ory.Session {
	return ctx.Value(userCtxKey).(*ory.Session)
}

func (s *Service) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("handling middleware request\n")

			cookies := r.Header.Get("Cookie")

			session, _, err := s.ory.FrontendApi.ToSession(r.Context()).Cookie(cookies).Execute()
			if (err != nil && session == nil) || (err == nil && !*session.Active) {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				withUser(r.Context(), session),
			))
		})
	}
}

const userCtxKey = "user"

type User struct {
	Username string
	Email    string
}
