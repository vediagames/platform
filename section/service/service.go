package service

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/section/domain"
)

type service struct {
	repository             domain.Repository
	webPlacementRepository domain.PlacedRepository
}

type Config struct {
	Repository       domain.Repository
	PlacedRepository domain.PlacedRepository
}

func (c Config) Validate() error {
	if c.Repository == nil {
		return fmt.Errorf("repository is required")
	}

	if c.PlacedRepository == nil {
		return fmt.Errorf("website placement repository is required")
	}

	return nil
}

func New(cfg Config) (domain.Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &service{
		repository:             cfg.Repository,
		webPlacementRepository: cfg.PlacedRepository,
	}, nil
}

func (s service) List(ctx context.Context, req domain.ListRequest) (domain.ListResponse, error) {
	repoRes, err := s.repository.Find(ctx, domain.FindQuery(req))
	if err != nil {
		return domain.ListResponse{}, fmt.Errorf("failed to find sections: %w", err)
	}

	return domain.ListResponse(repoRes), nil
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	repoRes, err := s.repository.FindOne(ctx, domain.FindOneQuery(req))
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to find section: %w", err)
	}

	return domain.GetResponse(repoRes), nil
}

func (s service) GetPlaced(ctx context.Context, req domain.GetPlacedRequest) (domain.GetPlacedResponse, error) {
	repoRes, err := s.webPlacementRepository.Find(ctx, domain.PlacedFindQuery(req))
	if err != nil {
		return domain.GetPlacedResponse{}, fmt.Errorf("failed to find website placements: %w", err)
	}

	return domain.GetPlacedResponse(repoRes), nil
}

func (s service) EditPlaced(ctx context.Context, req domain.EditPlacedRequest) error {
	if err := s.webPlacementRepository.Update(ctx, domain.PlacedUpdateQuery(req)); err != nil {
		return fmt.Errorf("failed to edit website placements: %w", err)
	}

	return nil
}
