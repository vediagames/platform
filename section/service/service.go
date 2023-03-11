package service

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"

	"github.com/vediagames/vediagames.com/section/domain"
)

type Config struct {
	Repository                 domain.Repository
	WebsitePlacementRepository domain.WebsitePlacementRepository
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Repository == nil, fmt.Errorf("empty repository"))
	err.AddIf(c.WebsitePlacementRepository == nil, fmt.Errorf("empty web placement repository"))

	return err.Err()
}

type service struct {
	repository             domain.Repository
	webPlacementRepository domain.WebsitePlacementRepository
}

func New(cfg Config) domain.Service {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &service{
		repository:             cfg.Repository,
		webPlacementRepository: cfg.WebsitePlacementRepository,
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

func (s service) GetWebsitePlacements(ctx context.Context, req domain.GetWebsitePlacementsRequest) (domain.GetWebsitePlacementsResponse, error) {
	repoRes, err := s.webPlacementRepository.Find(ctx, domain.WebsitePlacementFindQuery(req))
	if err != nil {
		return domain.GetWebsitePlacementsResponse{}, fmt.Errorf("failed to find website placements: %w", err)
	}

	return domain.GetWebsitePlacementsResponse(repoRes), nil
}

func (s service) EditWebsitePlacements(ctx context.Context, req domain.EditWebsitePlacementsRequest) error {
	if err := s.webPlacementRepository.Update(ctx, domain.WebsitePlacementUpdateQuery(req)); err != nil {
		return fmt.Errorf("failed to edit website placements: %w", err)
	}

	return nil
}
