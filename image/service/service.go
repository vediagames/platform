package service

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/config"
	"github.com/vediagames/vediagames.com/image/domain"
)

type service struct {
	processor domain.Processor
}

type Config struct {
	Processor domain.Processor
	Cfg       config.Imagor
}

func (c Config) Validate() error {
	if c.Cfg.Secret == "" {
		return fmt.Errorf("secret required")
	}
	if c.Cfg.URL == "" {
		return fmt.Errorf("url required")
	}

	return nil
}

func New(config Config) (domain.Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &service{
		processor: config.Processor,
	}, nil
}

func (s service) Put(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (s service) Get(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
