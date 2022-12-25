package service

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/image/domain"
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

func New(config Config) (domain.Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &service{
		repository: config.Repository,
	}, nil
}

func (s service) Put(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (s service) Get(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
