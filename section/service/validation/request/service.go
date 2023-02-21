package request

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

func (s service) GetPlaced(ctx context.Context, req domain.GetPlacedRequest) (domain.GetPlacedResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.GetPlacedResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	res, err := s.svc.GetPlaced(ctx, req)
	if err != nil {
		return domain.GetPlacedResponse{}, fmt.Errorf("failed to get website placements: %w", err)
	}

	if err := res.Validate(); err != nil {
		return domain.GetPlacedResponse{}, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}

func (s service) EditPlaced(ctx context.Context, req domain.EditPlacedRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	if err := s.svc.EditPlaced(ctx, req); err != nil {
		return fmt.Errorf("failed to edit website placements: %w", err)
	}

	return nil
}
