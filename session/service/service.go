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

func New(cfg Config) (domain.Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &service{
		repository: cfg.Repository,
	}, nil
}

func (s service) Create(ctx context.Context) (domain.CreateResponse, error) {

	repoRes, err := s.repository.Create(ctx)
	if err != nil {
		return domain.CreateResponse{}, fmt.Errorf("failed to create: %w", err)
	}

	return domain.CreateResponse(repoRes), nil
}
