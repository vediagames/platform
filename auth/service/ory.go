package service

import (
	"context"
	"fmt"

	ory "github.com/ory/kratos-client-go"
	"github.com/vediagames/zeroerror"

	"github.com/vediagames/vediagames.com/auth/domain"
)

type oryService struct {
	client *ory.APIClient
}

type OryConfig struct {
	Client *ory.APIClient
}

func (c OryConfig) Validate() error {
	var err zeroerror.Error

	if c.Client == nil {
		err.Add(fmt.Errorf("client is required"))
	}

	return err.Err()
}

func NewOry(cfg OryConfig) domain.Service {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &oryService{
		client: cfg.Client,
	}
}

func (s oryService) Authenticate(ctx context.Context, req domain.AuthenticateRequest) (domain.AuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.AuthenticateResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	session, _, err := s.client.FrontendApi.ToSession(ctx).Cookie(req.Cookies).Execute()
	if (err != nil && session == nil) || (err == nil && !*session.Active) {
		return domain.AuthenticateResponse{}, fmt.Errorf("failed to log in: %w", err)
	}

	identity := session.GetIdentity()
	traits := identity.GetTraits().(map[string]any)

	user := domain.User{
		ID:        identity.GetId(),
		SessionID: session.GetId(),
		Username:  traits["username"].(string),
		Email:     traits["email"].(string),
		CreatedAt: identity.GetCreatedAt(),
		UpdatedAt: identity.GetUpdatedAt(),
	}

	res := domain.AuthenticateResponse{
		User: user,
	}

	if err := res.Validate(); err != nil {
		return domain.AuthenticateResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

type contextKey string

const contextKeyUser contextKey = "user"

func (s oryService) ToContext(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

func (s oryService) FromContext(ctx context.Context) (domain.User, error) {
	anyUser := ctx.Value(contextKeyUser)
	if anyUser != nil {
		return domain.User{}, fmt.Errorf("not found")
	}

	user, ok := anyUser.(domain.User)
	if !ok {
		return domain.User{}, fmt.Errorf("failed to convert to domain")
	}

	return user, nil
}
