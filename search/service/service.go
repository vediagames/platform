package service

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"

	gamedomain "github.com/vediagames/platform/game/domain"
	"github.com/vediagames/platform/search/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

type Config struct {
	GameService gamedomain.Service
	TagService  tagdomain.Service
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.GameService == nil, fmt.Errorf("empty game service"))
	err.AddIf(c.TagService == nil, fmt.Errorf("empty tag service"))

	return err.Err()
}

func New(cfg Config) domain.Service {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &service{
		gameService: cfg.GameService,
		tagService:  cfg.TagService,
	}
}

type service struct {
	gameService gamedomain.Service
	tagService  tagdomain.Service
}

func (s service) Search(ctx context.Context, req domain.SearchRequest) (domain.SearchResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.SearchResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	gameRes, err := s.gameService.Search(ctx, gamedomain.SearchRequest{
		Query:          req.Query,
		Max:            req.MaxGames,
		AllowDeleted:   req.AllowDeleted,
		AllowInvisible: req.AllowInvisible,
		Language:       gamedomain.Language(req.Language),
	})
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("failed to search games: %w", err)
	}

	tagRes, err := s.tagService.Search(ctx, tagdomain.SearchRequest{
		Query:          req.Query,
		Max:            req.MaxTags,
		AllowDeleted:   req.AllowDeleted,
		AllowInvisible: req.AllowInvisible,
		Language:       tagdomain.Language(req.Language),
	})
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("failed to search tags: %w", err)
	}

	return populateSearchItemsFromImplementations(gameRes.Data.Data, tagRes.Data.Data, tagRes.Data.Total+gameRes.Data.Total), nil
}

func (s service) FullSearch(ctx context.Context, req domain.FullSearchRequest) (domain.SearchResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.SearchResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	gameRes, err := s.gameService.FullSearch(ctx, gamedomain.FullSearchRequest{
		Query:          req.Query,
		Page:           req.Page,
		Limit:          req.Limit,
		AllowDeleted:   req.AllowDeleted,
		AllowInvisible: req.AllowInvisible,
		Sort:           gamedomain.SortingMethod(req.Sort),
		Language:       gamedomain.Language(req.Language),
	})
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("failed to search games: %w", err)
	}

	tagRes, err := s.tagService.FullSearch(ctx, tagdomain.FullSearchRequest{
		Query:          req.Query,
		Page:           req.Page,
		Limit:          req.Limit,
		AllowDeleted:   req.AllowDeleted,
		AllowInvisible: req.AllowInvisible,
		Sort:           tagdomain.SortingMethod(req.Sort),
		Language:       tagdomain.Language(req.Language),
	})
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("failed to search tags: %w", err)
	}

	return populateSearchItemsFromImplementations(gameRes.Data.Data, tagRes.Data.Data, tagRes.Data.Total+gameRes.Data.Total), nil
}

func populateSearchItemsFromImplementations(games []gamedomain.Game, tag []tagdomain.Tag, total int) domain.SearchResponse {
	res := domain.SearchResponse{
		Games: make([]domain.SearchItem, 0, len(games)),
		Tags:  make([]domain.SearchItem, 0, len(tag)),
		Total: total,
	}

	for _, game := range games {
		res.Games = append(res.Games, domain.SearchItem{
			ID:               game.ID,
			Slug:             game.Slug,
			Name:             game.Name,
			ShortDescription: game.ShortDescription,
		})
	}

	for _, tag := range tag {
		res.Tags = append(res.Tags, domain.SearchItem{
			ID:               tag.ID,
			Slug:             tag.Slug,
			Name:             tag.Name,
			ShortDescription: tag.Description,
		})
	}

	return res
}
