package data

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/section/domain"
)

func New(svc domain.Service) domain.Service {
	if svc == nil {
		panic("empty service")
	}

	return &service{
		svc: svc,
	}
}

type service struct {
	svc domain.Service
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
