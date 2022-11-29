package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/vediagames/vediagames.com/bff/graphql/generated"
	"github.com/vediagames/vediagames.com/bff/graphql/model"
	categorydomain "github.com/vediagames/vediagames.com/category/domain"
	gamedomain "github.com/vediagames/vediagames.com/game/domain"
	notificationdomain "github.com/vediagames/vediagames.com/notification/domain"
	searchdomain "github.com/vediagames/vediagames.com/search/domain"
	sectiondomain "github.com/vediagames/vediagames.com/section/domain"
	tagdomain "github.com/vediagames/vediagames.com/tag/domain"
)

// SendEmail is the resolver for the sendEmail field.
func (r *mutationResolver) SendEmail(ctx context.Context, request model.SendEmailRequest) (*bool, error) {
	err := r.emailClient.Email(ctx, notificationdomain.EmailRequest{
		To: notificationdomain.User{
			Email: "antonio.jelic@vediagames.com",
			Name:  "Antonio Jelic",
		},
		From: notificationdomain.User{
			Email: request.From,
			Name:  "vediagames.com Contact form",
		},
		Name:    request.Name,
		Subject: request.Subject,
		Body:    request.Body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to email: %w", err)
	}

	return pointerTrue(), nil
}

// RefreshMaterializedViews is the resolver for the refreshMaterializedViews field.
func (r *mutationResolver) RefreshMaterializedViews(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: RefreshMaterializedViews - refreshMaterializedViews"))
}

// GetMostPlayedGames is the resolver for the getMostPlayedGames field.
func (r *queryResolver) GetMostPlayedGames(ctx context.Context, request model.GetMostPlayedGamesRequest) (*model.ListGamesResponse, error) {
	gameRes, err := r.gameService.GetMostPlayedByDays(ctx, gamedomain.GetMostPlayedByDaysRequest{
		Page:     request.Page,
		Limit:    request.Limit,
		MaxDays:  request.MaxDays,
		Language: gamedomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get most played games: %w", err)
	}

	res := &model.ListGamesResponse{
		Data:  make([]model.Game, 0, len(gameRes.Data.Data)),
		Total: gameRes.Data.Total,
	}

	for _, game := range gameRes.Data.Data {
		gameModel, err := gameFromGame(game)
		if err != nil {
			return nil, fmt.Errorf("failed to convert game: %w", err)
		}

		res.Data = append(res.Data, gameModel)
	}

	return res, nil
}

// GetFreshGames is the resolver for the getFreshGames field.
func (r *queryResolver) GetFreshGames(ctx context.Context, request model.GetFreshGamesRequest) (*model.ListGamesResponse, error) {
	gameRes, err := r.gameService.GetFresh(ctx, gamedomain.GetFreshRequest{
		Language: gamedomain.Language(request.Language),
		Page:     request.Page,
		Limit:    request.Limit,
		MaxDays:  request.MaxDays,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get fresh games: %w", err)
	}

	res := &model.ListGamesResponse{
		Data:  make([]model.Game, 0, len(gameRes.Data.Data)),
		Total: gameRes.Data.Total,
	}

	for _, game := range gameRes.Data.Data {
		gameModel, err := gameFromGame(game)
		if err != nil {
			return nil, fmt.Errorf("failed to convert game: %w", err)
		}

		res.Data = append(res.Data, gameModel)
	}

	return res, nil
}

// AvailableLanguages is the resolver for the availableLanguages field.
func (r *queryResolver) AvailableLanguages(ctx context.Context) ([]*model.LanguageItem, error) {
	return []*model.LanguageItem{
		{
			Code: model.LanguageEn,
			Name: "English",
		},
		{
			Code: model.LanguageEs,
			Name: "Español",
		},
	}, nil
}

// ListGames is the resolver for the listGames field.
func (r *queryResolver) ListGames(ctx context.Context, request model.ListGamesRequest) (*model.ListGamesResponse, error) {
	gameRes, err := r.gameService.List(ctx, gamedomain.ListRequest{
		Language:        gamedomain.Language(request.Base.Language),
		Page:            request.Base.Page,
		Limit:           request.Base.Limit,
		AllowDeleted:    request.Base.AllowDeleted,
		AllowInvisible:  request.Base.AllowInvisible,
		Sort:            sortingMethodToDomain[gamedomain.SortingMethod](request.Sort),
		Categories:      request.Categories,
		Tags:            request.Tags,
		IDs:             request.Ids,
		ExcludedGameIDs: request.ExcludedGameIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	res := &model.ListGamesResponse{
		Data:  make([]model.Game, 0, len(gameRes.Data)),
		Total: gameRes.Total,
	}

	for _, game := range gameRes.Data {
		g, err := gameFromGame(game)
		if err != nil {
			return nil, fmt.Errorf("failed to convert game: %w", err)
		}

		res.Data = append(res.Data, g)
	}

	return res, nil
}

// GetGame is the resolver for the getGame field.
func (r *queryResolver) GetGame(ctx context.Context, request model.BaseGetRequest) (*model.GetGameResponse, error) {
	gameRes, err := r.gameService.Get(ctx, gamedomain.GetRequest{
		Field:    gamedomain.GetByField(request.Field),
		Value:    request.Value,
		Language: gamedomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	res, err := gameFromGame(gameRes.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert game: %w", err)
	}

	return &model.GetGameResponse{
		Data: &res,
	}, nil
}

// ListCategories is the resolver for the listCategories field.
func (r *queryResolver) ListCategories(ctx context.Context, request model.BaseListRequest) (*model.ListCategoriesResponse, error) {
	categoryRes, err := r.categoryService.List(ctx, categorydomain.ListRequest{
		Language:       categorydomain.Language(request.Language),
		Page:           request.Page,
		Limit:          request.Limit,
		AllowDeleted:   request.AllowDeleted,
		AllowInvisible: request.AllowInvisible,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	res := &model.ListCategoriesResponse{
		Data:  make([]*model.Category, 0, len(categoryRes.Data)),
		Total: categoryRes.Total,
	}

	for _, category := range categoryRes.Data {
		c, err := categoryFromCategory(category)
		if err != nil {
			return nil, fmt.Errorf("failed to convert category: %w", err)
		}

		res.Data = append(res.Data, &c)
	}

	return &model.ListCategoriesResponse{
		Data:  res.Data,
		Total: res.Total,
	}, nil
}

// GetCategory is the resolver for the getCategory field.
func (r *queryResolver) GetCategory(ctx context.Context, request model.BaseGetRequest) (*model.GetCategoryResponse, error) {
	categoryRes, err := r.categoryService.Get(ctx, categorydomain.GetRequest{
		Field:    categorydomain.GetByField(request.Field),
		Value:    request.Value,
		Language: categorydomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	res, err := categoryFromCategory(categoryRes.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert category: %w", err)
	}

	return &model.GetCategoryResponse{
		Data: &res,
	}, nil
}

// ListTags is the resolver for the listTags field.
func (r *queryResolver) ListTags(ctx context.Context, request model.ListTagsRequest) (*model.ListTagsResponse, error) {
	tagRes, err := r.tagService.List(ctx, tagdomain.ListRequest{
		Language:       tagdomain.Language(request.Base.Language),
		Page:           request.Base.Page,
		Limit:          request.Base.Limit,
		AllowDeleted:   request.Base.AllowDeleted,
		AllowInvisible: request.Base.AllowInvisible,
		Sort:           sortingMethodToTag(request.Sort),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	res := &model.ListTagsResponse{
		Data:  make([]*model.Tag, 0, len(tagRes.Data)),
		Total: tagRes.Total,
	}

	for _, tag := range tagRes.Data {
		t, err := tagFromTag(tag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag: %w", err)
		}

		res.Data = append(res.Data, &t)
	}

	return &model.ListTagsResponse{
		Data:  res.Data,
		Total: res.Total,
	}, nil
}

// GetTag is the resolver for the getTag field.
func (r *queryResolver) GetTag(ctx context.Context, request model.BaseGetRequest) (*model.GetTagResponse, error) {
	tagRes, err := r.tagService.Get(ctx, tagdomain.GetRequest{
		Field:    tagdomain.GetByField(request.Field),
		Value:    request.Value,
		Language: tagdomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	res, err := tagFromTag(tagRes.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tag: %w", err)
	}

	return &model.GetTagResponse{
		Data: &res,
	}, nil
}

// ListSections is the resolver for the listSections field.
func (r *queryResolver) ListSections(ctx context.Context, request model.BaseListRequest) (*model.ListSectionsResponse, error) {
	sectionRes, err := r.sectionService.List(ctx, sectiondomain.ListRequest{
		Language:       sectiondomain.Language(request.Language),
		Page:           request.Page,
		Limit:          request.Limit,
		AllowDeleted:   request.AllowDeleted,
		AllowInvisible: request.AllowInvisible,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list sections: %w", err)
	}

	res := &model.ListSectionsResponse{
		Total: sectionRes.Total,
		Data:  make([]*model.Section, 0, len(sectionRes.Data)),
	}

	for _, ss := range sectionRes.Data {
		section, err := r.sectionFromSection(ctx, ss, request.Language)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		res.Data = append(res.Data, &section)
	}

	return res, nil
}

// GetSection is the resolver for the getSection field.
func (r *queryResolver) GetSection(ctx context.Context, request model.BaseGetRequest) (*model.GetSectionResponse, error) {
	sectionRes, err := r.sectionService.Get(ctx, sectiondomain.GetRequest{
		Field:    sectiondomain.GetByField(request.Field),
		Value:    request.Value,
		Language: sectiondomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get section: %w", err)
	}

	res, err := r.sectionFromSection(ctx, sectionRes.Data, request.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to convert section: %w", err)
	}

	return &model.GetSectionResponse{
		Data: &res,
	}, nil
}

// GetWebsiteSectionsPlacement is the resolver for the getWebsiteSectionsPlacement field.
func (r *queryResolver) GetWebsiteSectionsPlacement(ctx context.Context, language model.Language) (*model.GetWebsiteSectionsPlacementResponse, error) {
	sectionRes, err := r.sectionService.GetWebsitePlacements(ctx, sectiondomain.GetWebsitePlacementsRequest{
		Language: sectiondomain.Language(language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get website sections placement: %w", err)
	}

	res := &model.GetWebsiteSectionsPlacementResponse{
		Data: make([]*model.GetWebsiteSectionPlacement, 0, len(sectionRes.Data)),
	}

	for _, ss := range sectionRes.Data {
		section, err := r.sectionFromSection(ctx, ss.Section, language)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		res.Data = append(res.Data, &model.GetWebsiteSectionPlacement{
			Section:         &section,
			PlacementNumber: ss.PlacementNumber,
		})
	}

	return res, nil
}

// FetchRandomGame is the resolver for the fetchRandomGame field.
func (r *queryResolver) FetchRandomGame(ctx context.Context) (*model.FetchedGame, error) {
	fetcherRes, err := r.fetcherClient.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game: %w", err)
	}

	res := model.FetchedGame(fetcherRes)

	return &res, nil
}

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, language model.Language, query string, maxGames int, maxTags int, allowDeleted bool, allowInvisible bool) (*model.SearchResponse, error) {
	searchRes, err := r.searchService.Search(ctx, searchdomain.SearchRequest{
		Query:          query,
		MaxGames:       maxGames,
		MaxTags:        maxTags,
		AllowDeleted:   allowDeleted,
		AllowInvisible: allowInvisible,
		Language:       searchdomain.Language(language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	res, err := searchFromSearch(searchRes.Games, searchRes.Tags, searchRes.Total, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert search: %w", err)
	}

	return &res, nil
}

// FullSearch is the resolver for the fullSearch field.
func (r *queryResolver) FullSearch(ctx context.Context, request model.FullSearchRequest) (*model.SearchResponse, error) {
	if request.Sort == nil {
		*request.Sort = model.SortingMethodMostRelevant
	}

	searchRes, err := r.searchService.FullSearch(ctx, searchdomain.FullSearchRequest{
		Query:          request.Query,
		Page:           request.Page,
		Limit:          request.Limit,
		AllowDeleted:   request.AllowDeleted,
		AllowInvisible: request.AllowInvisible,
		Sort:           sortingMethodToDomain[searchdomain.SortingMethod](request.Sort),
		Language:       searchdomain.Language(request.Language),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	res, err := searchFromSearch(searchRes.Games, searchRes.Tags, searchRes.Total, true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert search: %w", err)
	}

	return &res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
