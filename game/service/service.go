package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vediagames/vediagames.com/game/domain"
)

type service struct {
	repository      domain.Repository
	statsRepository domain.StatsRepository
	eventRepository domain.EventRepository
}

type Config struct {
	Repository      domain.Repository
	StatsRepository domain.StatsRepository
	EventRepository domain.EventRepository
}

func (c Config) Validate() error {
	if c.Repository == nil {
		return fmt.Errorf("repository is required")
	}

	if c.StatsRepository == nil {
		return fmt.Errorf("stats repository is required")
	}

	if c.EventRepository == nil {
		return fmt.Errorf("event repository is required")
	}

	return nil
}

func New(config Config) (domain.Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &service{
		repository:      config.Repository,
		statsRepository: config.StatsRepository,
		eventRepository: config.EventRepository,
	}, nil
}

func (s service) List(ctx context.Context, req domain.ListRequest) (domain.ListResponse, error) {
	if ve := req.Validate(); ve != nil {
		return domain.ListResponse{}, fmt.Errorf("invalid request: %w", ve)
	}

	repoRes, err := s.repository.Find(ctx, domain.FindQuery{
		Language:        req.Language,
		Page:            req.Page,
		Limit:           req.Limit,
		AllowDeleted:    req.AllowDeleted,
		AllowInvisible:  req.AllowInvisible,
		Categories:      req.Categories,
		Tags:            req.Tags,
		Sort:            req.Sort,
		IDs:             req.IDs,
		ExcludedGameIDs: req.ExcludedGameIDs,
		MobileOnly:      req.MobileOnly,
	})
	if err != nil {
		return domain.ListResponse{}, fmt.Errorf("failed to list games: %w", err)
	}

	res := domain.ListResponse(repoRes)

	return res, res.Validate()
}

func (s service) Get(ctx context.Context, req domain.GetRequest) (domain.GetResponse, error) {
	if ve := req.Validate(); ve != nil {
		return domain.GetResponse{}, fmt.Errorf("invalid request: %w", ve)
	}

	repoRes, err := s.repository.FindOne(ctx, domain.FindOneQuery(req))
	if err != nil {
		return domain.GetResponse{}, fmt.Errorf("failed to get game: %w", err)
	}

	res := domain.GetResponse(repoRes)

	return res, res.Validate()
}

func (s service) GetMostPlayedByDays(ctx context.Context, req domain.GetMostPlayedByDaysRequest) (domain.GetMostPlayedByDaysResponse, error) {
	if ve := req.Validate(); ve != nil {
		return domain.GetMostPlayedByDaysResponse{}, fmt.Errorf("invalid request: %w", ve)
	}

	repoRes, err := s.statsRepository.FindMostPlayedIDsByDate(ctx, domain.FindMostPlayedIDsByDateQuery{
		Page:      req.Page,
		Limit:     req.Limit,
		DateLimit: time.Now().AddDate(0, 0, -req.MaxDays),
	})
	if err != nil {
		return domain.GetMostPlayedByDaysResponse{}, fmt.Errorf("failed to find most played ids: %w", err)
	}

	listRes, err := s.repository.Find(ctx, domain.FindQuery{
		Language: req.Language,
		Page:     req.Page,
		Limit:    req.Limit,
		Sort:     domain.SortingMethodMostPopular,
		IDs:      repoRes.Data,
	})
	if err != nil {
		return domain.GetMostPlayedByDaysResponse{}, fmt.Errorf("failed to find: %w", err)
	}

	res := domain.GetMostPlayedByDaysResponse{
		Data: domain.ListResponse(listRes),
	}

	if ve := res.Validate(); ve != nil {
		return domain.GetMostPlayedByDaysResponse{}, fmt.Errorf("invalid response: %w", ve)
	}

	return res, nil
}

func (s service) GetFresh(ctx context.Context, req domain.GetFreshRequest) (domain.GetFreshResponse, error) {
	if ve := req.Validate(); ve != nil {
		return domain.GetFreshResponse{}, fmt.Errorf("invalid request: %w", ve)
	}

	createDateLimit := time.Now().AddDate(0, 0, -req.MaxDays)

	repoRes, err := s.repository.Find(ctx, domain.FindQuery{
		Language:        req.Language,
		Page:            req.Page,
		Limit:           req.Limit,
		CreateDateLimit: createDateLimit,
		Sort:            domain.SortingMethodNewest,
	})
	if err != nil {
		return domain.GetFreshResponse{}, fmt.Errorf("failed to get fresh: %w", err)
	}

	res := domain.GetFreshResponse{
		Data: domain.ListResponse(repoRes),
	}

	if ve := res.Validate(); ve != nil {
		return domain.GetFreshResponse{}, fmt.Errorf("invalid response: %w", ve)
	}

	return res, nil
}

func (s service) LogEvent(ctx context.Context, req domain.LogEventRequest) error {
	if ve := req.Validate(); ve != nil {
		return fmt.Errorf("invalid request: %w", ve)
	}

	err := s.eventRepository.Log(ctx, domain.LogQuery(req))
	if err != nil {
		return fmt.Errorf("failed to log event: %w", err)
	}

	return nil
}

func (s service) Search(ctx context.Context, request domain.SearchRequest) (domain.SearchResponse, error) {
	if err := request.Validate(); err != nil {
		return domain.SearchResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	repoRes, err := s.repository.Search(ctx, domain.SearchQuery(request))
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("failed to search: %w", err)
	}

	res := domain.SearchResponse(repoRes)

	return res, request.Validate()
}

func (s service) FullSearch(ctx context.Context, request domain.FullSearchRequest) (domain.FullSearchResponse, error) {
	if err := request.Validate(); err != nil {
		return domain.FullSearchResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	repoRes, err := s.repository.FullSearch(ctx, domain.FullSearchQuery(request))
	if err != nil {
		return domain.FullSearchResponse{}, fmt.Errorf("failed to full search: %w", err)
	}

	res := domain.FullSearchResponse(repoRes)

	return res, res.Validate()
}
