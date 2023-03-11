package service

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"

	"github.com/vediagames/vediagames.com/section/domain"
)

type Config struct {
	Repository       domain.Repository
	PlacedRepository domain.PlacedRepository
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Repository == nil, fmt.Errorf("empty repository"))
	err.AddIf(c.PlacedRepository == nil, fmt.Errorf("empty placed repository"))

	return err.Err()
}

type service struct {
	repository             domain.Repository
	placedRepository domain.PlacedRepository
}

func New(cfg Config) domain.Service {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &service{
		repository:             cfg.Repository,
		placedRepository: cfg.PlacedRepository,
	}
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
	repoRes, err := s.placedRepository.Find(ctx, domain.PlacedFindQuery(req))
	if err != nil {
		return domain.GetPlacedResponse{}, fmt.Errorf("failed to find website placements: %w", err)
	}

	return domain.GetPlacedResponse(repoRes), nil
}

func (s service) EditPlaced(ctx context.Context, req domain.EditPlacedRequest) error {
	if err := s.placedRepository.Update(ctx, domain.PlacedUpdateQuery(req)); err != nil {
		return fmt.Errorf("failed to edit website placements: %w", err)
	}

	return nil
}
