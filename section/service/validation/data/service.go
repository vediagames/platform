package data

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/section/domain"
)

type service struct {
	svc domain.Service
}

type Config struct {
	Service domain.Service
}

func (c Config) Validate() error {
	if c.Service == nil {
		return fmt.Errorf("service is required")
	}

	return nil
}

func New(cfg Config) (domain.Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &service{
		svc: cfg.Service,
	}, nil
}

func (s service) List(ctx context.Context, req domain.ListRequest) (domain.ListResponse, error) {
	return s.svc.List(ctx, req)
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	return s.svc.Get(ctx, req)
}

func (s service) GetPlaced(ctx context.Context, req domain.GetPlacedRequest) (domain.GetPlacedResponse, error) {
	return s.svc.GetPlaced(ctx, req)
}

func (s service) EditPlaced(ctx context.Context, req domain.EditPlacedRequest) error {
	for placement, sectionID := range req.Placements {
		_, err := s.svc.Get(ctx, domain.GetRequest{
			Field:    domain.GetByFieldID,
			Value:    sectionID,
			Language: domain.LanguageEnglish,
		})
		if err != nil {
			return fmt.Errorf("failed to get section with id %d on placement %d: %w", sectionID, placement, err)
		}
	}

	return s.svc.EditPlaced(ctx, req)
}
