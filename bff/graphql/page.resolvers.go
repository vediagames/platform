package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/vediagames/vediagames.com/bff/domain"
	"github.com/vediagames/vediagames.com/bff/graphql/model"
	gamedomain "github.com/vediagames/vediagames.com/game/domain"
)

// GetHomePage is the resolver for the getHomePage field.
func (r *queryResolver) GetHomePage(ctx context.Context, request model.GetHomePageRequest) (*model.GetHomePageResponse, error) {
	totalGamesRes, err := r.gameService.List(ctx, gamedomain.ListRequest{
		Sort:           gamedomain.SortingMethodID,
		Language:       gamedomain.LanguageEnglish,
		Page:           1,
		Limit:          1,
		AllowDeleted:   false,
		AllowInvisible: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get total games: %w", err)
	}

	gamesAddedInLast7Days, err := r.GetFreshGames(ctx, model.GetFreshGamesRequest{
		Language: request.Language,
		Page:     1,
		Limit:    4,
		MaxDays:  7,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get fresh games: %w", err)
	}

	mostPlayedGamesInLast7Days, err := r.GetMostPlayedGames(ctx, model.GetMostPlayedGamesRequest{
		Language: request.Language,
		Page:     1,
		Limit:    4,
		MaxDays:  7,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get most played games in last 7 days: %w", err)
	}

	mostPlayedGames, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    40,
		},
		Sort: sortingMethodToPointer[model.SortingMethod](model.SortingMethodMostPopular),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get most played games: %w", err)
	}

	websiteSectionPlacements, err := r.GetWebsiteSectionsPlacement(ctx, request.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to get website sections placement: %w", err)
	}

	for _, websitePlacement := range websiteSectionPlacements.Data {
		var listGamesReq model.ListGamesRequest

		switch websitePlacement.Section.Slug {
		case "continue-playing":
			if len(request.LastPlayedGameIDs) <= 0 {
				continue
			}

			listGamesReq = model.ListGamesRequest{
				Base: &model.BaseListRequest{
					Language: request.Language,
					Page:     1,
					Limit:    10,
				},
				Ids: request.LastPlayedGameIDs,
			}
		case "newest":
			listGamesReq = model.ListGamesRequest{
				Base: &model.BaseListRequest{
					Language: request.Language,
					Page:     1,
					Limit:    10,
				},
				Sort: sortingMethodToPointer[model.SortingMethod](model.SortingMethodNewest),
			}
		default:
			if len(websitePlacement.Section.Categories.Data) == 0 &&
				len(websitePlacement.Section.Tags.Data) == 0 &&
				len(websitePlacement.Section.Games.Data) == 0 {
				zerolog.Ctx(ctx).Err(fmt.Errorf("invalid section: %s, no elements", websitePlacement.Section.Slug))
			}

			listGamesReq = model.ListGamesRequest{
				Base: &model.BaseListRequest{
					Language: request.Language,
					Page:     1,
					Limit:    10,
				},
				Sort:       sortingMethodToPointer[model.SortingMethod](model.SortingMethodRandom),
				Categories: websitePlacement.Section.Categories.IDs(),
				Tags:       websitePlacement.Section.Tags.IDs(),
				Ids:        websitePlacement.Section.Games.IDs(),
			}
		}

		gamesRes, err := r.ListGames(ctx, listGamesReq)
		if err != nil {
			return nil, fmt.Errorf("failed to list games for section %q: %w", websitePlacement.Section.Slug, err)
		}

		websitePlacement.Section.Games = gamesRes
	}

	tagsRes, err := r.ListTags(ctx, model.ListTagsRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    11,
		},
		Sort: sortingMethodToPointer[model.TagSortingMethod](model.TagSortingMethodRandom),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tagSections := make([]*model.TagSection, 0, len(tagsRes.Data))

	for _, tag := range tagsRes.Data {
		gamesRes, err := r.ListGames(ctx, model.ListGamesRequest{
			Base: &model.BaseListRequest{
				Language: request.Language,
				Page:     1,
				Limit:    7,
			},
			Tags: []int{tag.ID},
			Sort: sortingMethodToPointer[model.SortingMethod](model.SortingMethodRandom),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list games for tag %q: %w", tag.Name, err)
		}

		tagSections = append(tagSections, &model.TagSection{
			Games: gamesRes,
			Tag:   tag,
		})
	}

	return &model.GetHomePageResponse{
		TotalGames:                 totalGamesRes.Total,
		TotalGamesAddedInLast7Days: gamesAddedInLast7Days.Total,
		MostPlayedGamesInLast7Days: mostPlayedGamesInLast7Days,
		GamesAddedInLast7Days:      gamesAddedInLast7Days,
		MostPlayedGames:            mostPlayedGames,
		Sections:                   websiteSectionPlacements,
		TagSection:                 tagSections,
	}, nil
}

// GetGamePage is the resolver for the getGamePage field.
func (r *queryResolver) GetGamePage(ctx context.Context, request model.GetGamePageRequest) (*model.GetGamePageResponse, error) {
	var tagsToIgnore = map[int]string{
		12: "boys",
		13: "girls",
		21: "challenge",
		54: "brain",
	}

	gameRes, err := r.GetGame(ctx, model.GetGameTagRequest{
		Field:    model.GetByFieldSlug,
		Value:    request.Slug,
		Language: request.Language,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	tagIDs := make([]int, 0, len(gameRes.Data.Tags.Data))
	for _, tag := range gameRes.Data.Tags.Data {
		if _, ok := tagsToIgnore[tag.ID]; !ok {
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	categoryIDs := make([]int, 0, len(gameRes.Data.Categories.Data))
	for _, category := range gameRes.Data.Categories.Data {
		categoryIDs = append(categoryIDs, category.ID)
	}

	otherGamesRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    8,
		},
		Sort:            sortingMethodToPointer[model.SortingMethod](model.SortingMethodRandom),
		Categories:      categoryIDs,
		Tags:            tagIDs,
		ExcludedGameIDs: []int{gameRes.Data.ID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return &model.GetGamePageResponse{
		Game:       gameRes,
		OtherGames: otherGamesRes,
		IsLiked:    false,
		IsDisliked: false,
	}, nil
}

// GetContinuePlayingPage is the resolver for the getContinuePlayingPage field.
func (r *queryResolver) GetContinuePlayingPage(ctx context.Context, request model.GetContinuePlayingPageRequest) (*model.GetContinuePlayingPageResponse, error) {
	gameRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     request.Page,
			Limit:    15,
		},
		Ids: request.LastPlayedGameIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	// order by last played game ids
	orderedGames := make([]model.Game, len(gameRes.Data))
	for i, game := range gameRes.Data {
		for _, id := range request.LastPlayedGameIDs {
			if game.ID == id {
				orderedGames[i] = game
			}
		}
	}

	gameRes.Data = orderedGames

	return &model.GetContinuePlayingPageResponse{
		Data: gameRes,
	}, nil
}

// GetFilterPage is the resolver for the getFilterPage field.
func (r *queryResolver) GetFilterPage(ctx context.Context, request model.GetFilterPageRequest) (*model.GetFilterPageResponse, error) {
	gameRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     request.Page,
			Limit:    15,
		},
		Sort:       request.Sort,
		Categories: request.CategoryIDs,
		Tags:       request.TagIDs,
		Ids:        request.GameIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return &model.GetFilterPageResponse{
		Data: gameRes,
	}, nil
}

// GetSearchPage is the resolver for the getSearchPage field.
func (r *queryResolver) GetSearchPage(ctx context.Context, request model.GetSearchPageRequest) (*model.GetSearchPageResponse, error) {
	req := model.FullSearchRequest{
		Language:       request.Language,
		Query:          request.Query,
		Page:           request.Page,
		Limit:          15,
		Sort:           request.Sort,
		AllowDeleted:   false,
		AllowInvisible: false,
	}

	searchRes, err := r.FullSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to full search: %w", err)
	}

	items := searchRes.Games
	items = append(items, searchRes.Tags...)

	return &model.GetSearchPageResponse{
		Items:        items,
		Total:        searchRes.Total,
		ShowingRange: fmt.Sprintf("%d-%d", req.Limit*(req.Page-1)+1, req.Limit*(req.Page-1)+len(items)),
	}, nil
}

// GetSiteMapPage is the resolver for the getSiteMapPage field.
func (r *queryResolver) GetSiteMapPage(ctx context.Context, request model.GetSiteMapPageRequest) (*model.GetSiteMapPageResponse, error) {
	categoryRes, err := r.ListCategories(ctx, model.BaseListRequest{
		Language: request.Language,
		Page:     1,
		Limit:    1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return &model.GetSiteMapPageResponse{
		Categories: categoryRes,
	}, nil
}

// GetWizardPage is the resolver for the getWizardPage field.
func (r *queryResolver) GetWizardPage(ctx context.Context, request model.GetWizardPageRequest) (*model.GetWizardPageResponse, error) {
	categoryRes, err := r.ListCategories(ctx, model.BaseListRequest{
		Language: request.Language,
		Page:     1,
		Limit:    1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	gameRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    4,
		},
		Sort:       sortingMethodToPointer[model.SortingMethod](model.SortingMethodMostPopular),
		Categories: request.CategoryIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return &model.GetWizardPageResponse{
		Categories: categoryRes,
		Games:      gameRes,
	}, nil
}

// GetTagPage is the resolver for the getTagPage field.
func (r *queryResolver) GetTagPage(ctx context.Context, request model.GetTagPageRequest) (*model.GetTagPageResponse, error) {
	tagRes, err := r.GetTag(ctx, model.GetGameTagRequest{
		Field:    model.GetByFieldID,
		Value:    strconv.Itoa(request.TagID),
		Language: request.Language,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	gameRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     request.Page,
			Limit:    15,
		},
		Tags: []int{request.TagID},
		Sort: sortingMethodToPointer[model.SortingMethod](model.SortingMethodMostPopular),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return &model.GetTagPageResponse{
		Tag:   tagRes,
		Games: gameRes,
	}, nil
}

// GetTagsPage is the resolver for the getTagsPage field.
func (r *queryResolver) GetTagsPage(ctx context.Context, request model.GetTagsPageRequest) (*model.GetTagsPageResponse, error) {
	tagRes, err := r.ListTags(ctx, model.ListTagsRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     request.Page,
			Limit:    15,
		},
		Sort: sortingMethodToPointer[model.TagSortingMethod](model.TagSortingMethodName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return &model.GetTagsPageResponse{
		Data: tagRes,
	}, nil
}

// GetCategoryPage is the resolver for the getCategoryPage field.
func (r *queryResolver) GetCategoryPage(ctx context.Context, request model.GetCategoryPageRequest) (*model.GetCategoryPageResponse, error) {
	tagLayout, ok := domain.CategoryPageLayouts[request.Slug]
	if !ok {
		return nil, fmt.Errorf("invalid layout slug: %q", request.Slug)
	}

	categoryRes, err := r.GetCategory(ctx, model.BaseGetRequest{
		Field:    model.GetByFieldSlug,
		Value:    request.Slug,
		Language: request.Language,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	gamesRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    22,
		},
		Sort:       sortingMethodToPointer[model.SortingMethod](model.SortingMethodMostPopular),
		Categories: gamedomain.IDs{categoryRes.Data.ID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get first section games: %w", err)
	}

	tagSections := make([]*model.TagSection, 0, len(tagLayout))
	for _, tag := range tagLayout {
		tagRes, err := r.GetTag(ctx, model.GetGameTagRequest{
			Field:    model.GetByFieldSlug,
			Value:    tag,
			Language: request.Language,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get tag %q: %w", tag, err)
		}

		gamesRes, err := r.ListGames(ctx, model.ListGamesRequest{
			Base: &model.BaseListRequest{
				Language: request.Language,
				Page:     1,
				Limit:    7,
			},
			Sort: sortingMethodToPointer[model.SortingMethod](model.SortingMethodMostPopular),
			Tags: gamedomain.IDs{tagRes.Data.ID},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get games for tag %q: %w", tag, err)
		}

		tagSection := &model.TagSection{
			Games: gamesRes,
			Tag:   tagRes.Data,
		}

		tagSections = append(tagSections, tagSection)
	}

	randomTagsRes, err := r.ListTags(ctx, model.ListTagsRequest{
		Base: &model.BaseListRequest{
			Language: request.Language,
			Page:     1,
			Limit:    11,
		},
		Sort: sortingMethodToPointer[model.TagSortingMethod](model.TagSortingMethodRandom),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get random tags: %w", err)
	}

	firstSectionGames := gamesRes
	otherGames := &model.ListGamesResponse{}

	if len(gamesRes.Data) > 7 {
		firstSectionGames = &model.ListGamesResponse{
			Data:  gamesRes.Data[:7],
			Total: gamesRes.Total,
		}

		otherGames = &model.ListGamesResponse{
			Data:  gamesRes.Data[7:],
			Total: gamesRes.Total,
		}
	}

	return &model.GetCategoryPageResponse{
		Category:          categoryRes,
		FirstSectionGames: firstSectionGames,
		TagSections:       tagSections,
		Tags:              randomTagsRes,
		OtherGames:        otherGames,
	}, nil
}

// GetCategoriesPage is the resolver for the getCategoriesPage field.
func (r *queryResolver) GetCategoriesPage(ctx context.Context, request model.GetCategoriesPageRequest) (*model.GetCategoriesPageResponse, error) {
	categoryRes, err := r.ListCategories(ctx, model.BaseListRequest{
		Language: request.Language,
		Page:     1,
		Limit:    1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return &model.GetCategoriesPageResponse{
		Data: categoryRes,
	}, nil
}
