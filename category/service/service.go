package service

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/category/domain"
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

func (s service) List(ctx context.Context, req domain.ListRequest) (domain.ListResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.ListResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	repoRes, err := s.repository.Find(ctx, domain.FindQuery(req))
	if err != nil {
		return domain.ListResponse{}, fmt.Errorf("failed to find: %w", err)
	}

	res := domain.ListResponse(repoRes)
	if err := res.Validate(); err != nil {
		return domain.ListResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	repoRes, err := s.repository.FindOne(ctx, domain.FindOneQuery(req))
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to find one: %w", err)
	}

	res := domain.GetResponse(repoRes)
	if err := res.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) IncreaseClick(ctx context.Context, req domain.IncreaseClickRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	err := s.repository.IncreaseField(ctx, domain.IncreaseFieldQuery{
		ID:       req.ID,
		Field:    domain.IncreasableFieldClicks,
		ByAmount: req.ByAmount,
	})
	if err != nil {
		return fmt.Errorf("failed to increase field: %w", err)
	}

	return nil
}
