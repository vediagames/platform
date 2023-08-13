package service

import (
	"context"

	"github.com/vediagames/platform/auth/domain"
)

func NewZero() domain.Service {
	return &service{}
}

type service struct{}

func (s service) Authenticate(_ context.Context, _ domain.AuthenticateRequest) (domain.AuthenticateResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s service) ToContext(_ context.Context, _ domain.User) context.Context {
	panic("not implemented") // TODO: Implement
}

func (s service) FromContext(_ context.Context) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}
