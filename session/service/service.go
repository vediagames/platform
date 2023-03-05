package service

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/session/domain"
)

type service struct {
	repository domain.Repository
}

type Config struct {
	Repository domain.Repository
}

func (c Config) Validate() error {
	if c.Repository == nil {
		return fmt.Errorf("repository is required")
	}

	return nil
}

func New(cfg Config) domain.Service {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &service{
		repository: cfg.Repository,
	}
}

func (s service) Create(ctx context.Context, req domain.CreateRequest) (domain.CreateResponse, error) {
	if ve := req.Validate(); ve != nil {
		return domain.CreateResponse{}, fmt.Errorf("invalid request: %w", ve)
	}

	repoRes, err := s.repository.Insert(ctx, req.ToInsertQuery())
	if err != nil {
		return domain.CreateResponse{}, fmt.Errorf("failed to create: %w", err)
	}

	res := domain.CreateResponse(repoRes)
	if err := res.Validate(); err != nil {
		return domain.CreateResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}
