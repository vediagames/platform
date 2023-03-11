package request

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
	if err := req.Validate(); err != nil {
		return domain.ListResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	res, err := s.svc.List(ctx, req)
	if err != nil {
		return domain.ListResponse{}, fmt.Errorf("failed to list sections: %w", err)
	}

	if err := res.Validate(); err != nil {
		return domain.ListResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	res, err := s.svc.Get(ctx, req)
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to get section: %w", err)
	}

	if err := res.Validate(); err != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) GetWebsitePlacements(ctx context.Context, req domain.GetWebsitePlacementsRequest) (domain.GetWebsitePlacementsResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetWebsitePlacementsResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	res, err := s.svc.GetWebsitePlacements(ctx, req)
	if err != nil {
		return domain.GetWebsitePlacementsResponse{}, fmt.Errorf("failed to get website placements: %w", err)
	}

	if err := res.Validate(); err != nil {
		return domain.GetWebsitePlacementsResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) EditWebsitePlacements(ctx context.Context, req domain.EditWebsitePlacementsRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	if err := s.svc.EditWebsitePlacements(ctx, req); err != nil {
		return fmt.Errorf("failed to edit website placements: %w", err)
	}

	return nil
}
