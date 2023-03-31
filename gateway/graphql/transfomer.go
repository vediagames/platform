package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	"github.com/vediagames/platform/gateway/graphql/model"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

func (r *Resolver) gamesFromDomain(ctx context.Context, domain gamedomain.Games) (*model.Games, error) {
	games := &model.Games{
		Data:  make([]*model.Game, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainGame := range domain.Data {
		game, err := r.gameFromDomain(ctx, domainGame)
		if err != nil {
			return nil, fmt.Errorf("failed to convert game: %w", err)
		}

		games.Data = append(games.Data, game)
	}

	return games, nil
}

func (r *Resolver) gameFromDomain(ctx context.Context, domain gamedomain.Game) (*model.Game, error) {
	thumb512x384, err := pathGame.Thumbnail(domain.Slug, thumbnail512x384)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb512x512, err := pathGame.Thumbnail(domain.Slug, thumbnail512x512)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x512 thumbnail: %w", err)
	}

	var tags *model.Tags
	if len(domain.TagIDRefs) == 0 {
		tagRes, err := r.tagService.List(ctx, tagdomain.ListRequest{
			Language: tagdomain.Language(domain.Language),
			Page:     1,
			Limit:    20,
			Sort:     tagdomain.SortingMethodName,
			IDRefs:   tagdomain.IDs(domain.TagIDRefs),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}

		tags, err = r.tagsFromDomain(tagRes.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tags: %w", err)
		}
	}

	var categories *model.Categories
	if len(domain.CategoryIDRefs) == 0 {
		categoryRes, err := r.categoryService.List(ctx, categorydomain.ListRequest{
			Language: categorydomain.Language(domain.Language),
			Page:     1,
			Limit:    20,
			IDRefs:   categorydomain.IDs(domain.TagIDRefs),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list categories: %w", err)
		}

		categories = r.categoriesFromDomain(categoryRes.Data)
	}

	return &model.Game{
		ID:                domain.ID,
		Language:          model.Language(domain.Language),
		Slug:              domain.Slug,
		Name:              domain.Name,
		Status:            model.Status(domain.Status),
		CreatedAt:         domain.CreatedAt.String(),
		DeletedAt:         stringToPointer(domain.DeletedAt.String()),
		PublishedAt:       stringToPointer(domain.PublishedAt.String()),
		URL:               domain.URL,
		Width:             domain.Width,
		Height:            domain.Height,
		ShortDescription:  stringToPointer(domain.ShortDescription),
		Description:       stringToPointer(domain.Description),
		Content:           stringToPointer(domain.Content),
		Likes:             domain.Likes,
		Dislikes:          domain.Dislikes,
		Plays:             domain.Plays,
		Weight:            domain.Weight,
		Player1Controls:   stringToPointer(domain.Player1Controls),
		Player2Controls:   stringToPointer(domain.Player2Controls),
		Tags:              tags,
		Categories:        categories,
		Mobile:            domain.Mobile,
		Thumbnail512x384:  thumb512x384,
		Thumbnail512x512:  thumb512x512,
		PageURL:           fmt.Sprintf("/game/%s", domain.Slug),
		FullScreenPageURL: fmt.Sprintf("/game/%s/fullscreen", domain.Slug),
	}, nil
}

func (r *Resolver) tagsFromDomain(domain tagdomain.Tags) (*model.Tags, error) {
	tags := &model.Tags{
		Data:  make([]*model.Tag, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainTag := range domain.Data {
		tag, err := r.tagFromDomain(domainTag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag: %w", err)
		}

		tags.Data = append(tags.Data, tag)
	}

	return tags, nil
}

func (r *Resolver) tagFromDomain(domain tagdomain.Tag) (*model.Tag, error) {
	thumb512x384, err := pathTag.Thumbnail(domain.Slug, thumbnail512x384)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb128x128, err := pathTag.Thumbnail(domain.Slug, thumbnail128x128)
	if err != nil {
		return nil, fmt.Errorf("failed to get 128x128 thumbnail: %w", err)
	}

	return &model.Tag{
		ID:               domain.ID,
		Language:         model.Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           model.Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		Thumbnail512x384: thumb512x384,
		Thumbnail128x128: thumb128x128,
		PageURL:          fmt.Sprintf("/tag/%s", domain.Slug),
	}, nil
}

func (r *Resolver) categoriesFromDomain(domain categorydomain.Categories) *model.Categories {
	categories := &model.Categories{
		Data:  make([]*model.Category, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainCategory := range domain.Data {
		category := r.categoryFromDomain(domainCategory)
		categories.Data = append(categories.Data, category)
	}

	return categories
}

func (r *Resolver) categoryFromDomain(domain categorydomain.Category) *model.Category {
	return &model.Category{
		ID:               domain.ID,
		Language:         model.Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           model.Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		PageURL:          fmt.Sprintf("/tag/%s", domain.Slug),
	}
}

func (r *Resolver) sectionsFromDomain(ctx context.Context, domain sectiondomain.Sections) (*model.Sections, error) {
	sections := &model.Sections{
		Data:  make([]*model.Section, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainSection := range domain.Data {
		section, err := r.sectionFromDomain(ctx, domainSection)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		sections.Data = append(sections.Data, section)
	}

	return sections, nil
}

func (r *Resolver) sectionFromDomain(ctx context.Context, domain sectiondomain.Section) (*model.Section, error) {
	var games *model.Games
	if len(domain.GameIDRefs) == 0 {
		gameRes, err := r.gameService.List(ctx, gamedomain.ListRequest{
			Language: gamedomain.Language(domain.Language),
			Page:     1,
			Limit:    30,
			Sort:     gamedomain.SortingMethodMostLiked,
			IDRefs:   gamedomain.IDs(domain.GameIDRefs),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list games: %w", err)
		}

		games, err = r.gamesFromDomain(ctx, gameRes.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert games: %w", err)
		}
	}

	var tags *model.Tags
	if len(domain.TagIDRefs) == 0 {
		tagRes, err := r.tagService.List(ctx, tagdomain.ListRequest{
			Language: tagdomain.Language(domain.Language),
			Page:     1,
			Limit:    20,
			Sort:     tagdomain.SortingMethodName,
			IDRefs:   tagdomain.IDs(domain.TagIDRefs),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}

		tags, err = r.tagsFromDomain(tagRes.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tags: %w", err)
		}
	}

	var categories *model.Categories
	if len(domain.CategoryIDRefs) == 0 {
		categoryRes, err := r.categoryService.List(ctx, categorydomain.ListRequest{
			Language: categorydomain.Language(domain.Language),
			Page:     1,
			Limit:    20,
			IDRefs:   categorydomain.IDs(domain.TagIDRefs),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list categories: %w", err)
		}

		categories = r.categoriesFromDomain(categoryRes.Data)
	}

	pageURL := "/continue-playing"

	if domain.Slug != "continue-playing" {
		var paramsFilter filterParams

		switch domain.Slug {
		case "newest":
			paramsFilter = filterParams{
				Sort: model.SortingMethodNewest,
			}
		default:
			paramsFilter = filterParams{
				Tags:       domain.TagIDRefs,
				Categories: domain.CategoryIDRefs,
				Games:      domain.GameIDRefs,
			}
		}

		params, err := filterParamsInBase64(paramsFilter)
		if err != nil {
			return &model.Section{}, fmt.Errorf("failed to encode filter params: %w", err)
		}

		pageURL = fmt.Sprintf("/filter?title=%s&params=%s", url.QueryEscape(domain.Name), params)
	}

	return &model.Section{
		ID:               domain.ID,
		Language:         model.Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		Status:           model.Status(domain.Status),
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Tags:             tags,
		Categories:       categories,
		Games:            games,
		PageURL:          pageURL,
	}, nil
}

func (r *Resolver) placedSectionsFromDomain(ctx context.Context, domain sectiondomain.GetPlacedResponse) (*model.PlacedSectionsResponse, error) {
	placedSections := &model.PlacedSectionsResponse{
		PlacedSections: make([]*model.PlacedSection, 0, len(domain.Data)),
	}

	for _, domainSection := range domain.Data {
		section, err := r.sectionFromDomain(ctx, domainSection.Section)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		placedSections.PlacedSections = append(placedSections.PlacedSections, &model.PlacedSection{
			Section:   section,
			Placement: domainSection.PlacementNumber,
		})
	}

	return placedSections, nil
}

func (r *Resolver) searchFromDomain(domain searchdomain.SearchResponse) (*model.SearchResponse, error) {
	searchResponse := &model.SearchResponse{
		SearchItems: make([]*model.SearchItem, 0, len(domain.Games)+len(domain.Tags)),
		Total:       domain.Total,
	}

	for _, domainItem := range domain.Games {
		thumb512x384, err := pathGame.Thumbnail(domainItem.Slug, thumbnail512x384)
		if err != nil {
			return nil, fmt.Errorf("failed to get 512x384 thumbnail for %q: %w", domainItem.Slug, err)
		}

		searchResponse.SearchItems = append(searchResponse.SearchItems, &model.SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             model.SearchItemTypeGame,
			PageURL:          fmt.Sprintf("/game/%s", domainItem.Slug),
			Thumbnail512x384: thumb512x384,
		})
	}

	for _, domainItem := range domain.Tags {
		thumb512x384, err := pathTag.Thumbnail(domainItem.Slug, thumbnail512x384)
		if err != nil {
			return nil, fmt.Errorf("failed to get 512x384 thumbnail for %q: %w", domainItem.Slug, err)
		}

		searchResponse.SearchItems = append(searchResponse.SearchItems, &model.SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             model.SearchItemTypeTag,
			PageURL:          fmt.Sprintf("/tag/%s", domainItem.Slug),
			Thumbnail512x384: thumb512x384,
		})
	}

	return searchResponse, nil
}

func pointerTrue() *bool {
	v := true
	return &v
}

func stringToPointer(s string) *string {
	return &s
}

type filterParams struct {
	Sort       model.SortingMethod `json:"sort,omitempty"`
	Tags       []int               `json:"tags,omitempty"`
	Categories []int               `json:"categories,omitempty"`
	Games      []int               `json:"games,omitempty"`
}

func filterParamsInBase64(p filterParams) (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to marshal filter params: %w", err)
	}

	return base64.StdEncoding.EncodeToString(body), nil
}

func sortingMethodToDomain[T gamedomain.SortingMethod | searchdomain.SortingMethod | tagdomain.SortingMethod](s *model.SortingMethod) T {
	if s == nil {
		return T(model.SortingMethodID.String())
	}

	str := strings.Replace(s.String(), "_", "-", -1)

	return T(str)
}
