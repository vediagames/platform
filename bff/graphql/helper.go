package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/vediagames/platform/bff/graphql/model"
	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

func searchFromSearch(games, tags []searchdomain.SearchItem, total int, fullSearch bool) (model.SearchResponse, error) {
	res := model.SearchResponse{
		Games: make([]*model.SearchItem, 0, len(games)),
		Tags:  make([]*model.SearchItem, 0, len(tags)),
		Total: total,
	}

	for _, game := range games {
		thumb := thumbnail512x512
		if fullSearch {
			thumb = thumbnail512x384
		}

		gameThumbnailPath, err := pathGame.Thumbnail(game.Slug, thumb)
		if err != nil {
			return model.SearchResponse{}, fmt.Errorf("failed to get game thumbnail path: %w", err)
		}

		res.Games = append(res.Games, &model.SearchItem{
			ID:               game.ID,
			ShortDescription: game.ShortDescription,
			Name:             game.Name,
			Slug:             game.Slug,
			Type:             model.SearchItemTypeGame,
			Thumbnail512x384: gameThumbnailPath,
		})
	}

	for _, tag := range tags {
		thumb := thumbnail128x128
		if fullSearch {
			thumb = thumbnail512x384
		}

		tagThumbnailPath, err := pathTag.Thumbnail(tag.Slug, thumb)
		if err != nil {
			return model.SearchResponse{}, fmt.Errorf("failed to get tag thumbnail path: %w for tag: %q", err, tag.Slug)
		}

		res.Tags = append(res.Tags, &model.SearchItem{
			ID:               tag.ID,
			ShortDescription: tag.ShortDescription,
			Name:             tag.Name,
			Slug:             tag.Slug,
			Type:             model.SearchItemTypeTag,
			Thumbnail512x384: tagThumbnailPath,
		})
	}

	return res, nil
}

func stringToPointer(s string) *string {
	return &s
}

func pointerTrue() *bool {
	v := true
	return &v
}

func (r *queryResolver) sectionFromSection(ctx context.Context, s sectiondomain.Section, l model.Language) (model.Section, error) {
	tags := make([]model.ComplimentaryTag, 0, len(s.TagIDRefs.Data))
	for _, tag := range s.TagIDRefs.Data {
		thumb, err := pathTag.Thumbnail(tag.Slug, thumbnail128x128)
		if err != nil {
			return model.Section{}, fmt.Errorf("failed to get thumbnail: %w", err)
		}

		tags = append(tags, model.ComplimentaryTag{
			ID:               tag.ID,
			Slug:             tag.Slug,
			Name:             tag.Name,
			Description:      &tag.Description,
			Thumbnail128x128: thumb,
		})
	}

	categories := make([]model.ComplimentaryCategory, 0, len(s.CategoryIDRefs.Data))
	for _, category := range s.CategoryIDRefs.Data {
		categories = append(categories, model.ComplimentaryCategory{
			ID:          category.ID,
			Slug:        category.Slug,
			Name:        category.Name,
			Description: &category.Description,
		})
	}

	pageURL := "/continue-playing"

	if s.Slug != "continue-playing" {
		var paramsFilter filterParams

		switch s.Slug {
		case "newest":
			paramsFilter = filterParams{
				Sort: model.SortingMethodNewest,
			}
		default:
			paramsFilter = filterParams{
				Tags:       s.TagIDRefs.IDs(),
				Categories: s.CategoryIDRefs.IDs(),
				Games:      s.Games,
			}
		}

		params, err := filterParamsInBase64(paramsFilter)
		if err != nil {
			return model.Section{}, fmt.Errorf("failed to encode filter params: %w", err)
		}

		pageURL = fmt.Sprintf("/filter?title=%s&params=%s", url.QueryEscape(s.Name), params)
	}

	gamesRes, err := r.ListGames(ctx, model.ListGamesRequest{
		Base: &model.BaseListRequest{
			Language: l,
			Page:     1,
			Limit:    10,
		},
		Ids: s.Games,
	})
	if err != nil {
		return model.Section{}, fmt.Errorf("failed to get games: %w", err)
	}

	return model.Section{
		ID:               s.ID,
		Language:         model.Language(s.Language),
		Slug:             s.Slug,
		Name:             s.Name,
		ShortDescription: &s.ShortDescription,
		Description:      &s.Description,
		Tags: &model.ComplimentaryTags{
			Data: tags,
		},
		Categories: &model.ComplimentaryCategories{
			Data: categories,
		},
		Games:       gamesRes,
		Status:      model.Status(s.Status),
		CreatedAt:   s.CreatedAt.String(),
		DeletedAt:   stringToPointer(s.DeletedAt.String()),
		PublishedAt: stringToPointer(s.PublishedAt.String()),
		Content:     &s.Content,
		PageUrl:     pageURL,
	}, nil
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

func sortingMethodToDomain[T gamedomain.SortingMethod | searchdomain.SortingMethod](s *model.SortingMethod) T {
	if s == nil {
		return T(model.SortingMethodID.String())
	}

	str := strings.Replace(s.String(), "_", "-", -1)

	return T(str)
}

func gameFromGame(game gamedomain.Game) (model.Game, error) {
	thumb512x384, err := pathGame.Thumbnail(game.Slug, thumbnail512x384)
	if err != nil {
		return model.Game{}, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb512x512, err := pathGame.Thumbnail(game.Slug, thumbnail512x512)
	if err != nil {
		return model.Game{}, fmt.Errorf("failed to get 512x512 thumbnail: %w", err)
	}

	tags := make([]model.ComplimentaryTag, 0, len(game.Tags))
	for _, tag := range game.Tags {
		tagThumb128x128, err := pathTag.Thumbnail(tag.Slug, thumbnail128x128)
		if err != nil {
			return model.Game{}, fmt.Errorf("failed to get tag 128x128 thumbnail: %w", err)
		}

		tags = append(tags, model.ComplimentaryTag{
			ID:               tag.ID,
			Slug:             tag.Slug,
			Name:             tag.Name,
			Description:      &tag.Description,
			Thumbnail128x128: tagThumb128x128,
		})
	}

	categories := make([]model.ComplimentaryCategory, 0, len(game.Categories))
	for _, category := range game.Categories {
		categories = append(categories, model.ComplimentaryCategory{
			ID:          category.ID,
			Slug:        category.Slug,
			Name:        category.Name,
			Description: &category.Description,
		})
	}

	return model.Game{
		ID:               game.ID,
		Language:         model.Language(game.Language),
		Slug:             game.Slug,
		Name:             game.Name,
		Status:           model.Status(game.Status),
		CreatedAt:        game.CreatedAt.String(),
		DeletedAt:        stringToPointer(game.DeletedAt.String()),
		PublishedAt:      stringToPointer(game.PublishedAt.String()),
		URL:              game.URL,
		Width:            game.Width,
		Height:           game.Height,
		ShortDescription: &game.ShortDescription,
		Description:      &game.Description,
		Content:          &game.Content,
		Likes:            game.Likes,
		Dislikes:         game.Dislikes,
		Plays:            game.Plays,
		Weight:           game.Weight,
		Player1Controls:  &game.Player1Controls,
		Player2Controls:  &game.Player2Controls,
		Tags: &model.ComplimentaryTags{
			Data: tags,
		},
		Categories: &model.ComplimentaryCategories{
			Data: categories,
		},
		Mobile:           game.Mobile,
		Thumbnail512x384: thumb512x384,
		Thumbnail512x512: thumb512x512,
	}, nil
}

func categoryFromCategory(c categorydomain.Category) (model.Category, error) {
	return model.Category{
		ID:               c.ID,
		Language:         model.Language(c.Language),
		Slug:             c.Slug,
		Name:             c.Name,
		ShortDescription: &c.ShortDescription,
		Description:      &c.Description,
		Content:          &c.Content,
		Status:           model.Status(c.Status),
		Clicks:           c.Clicks,
		CreatedAt:        c.CreatedAt.String(),
		DeletedAt:        stringToPointer(c.DeletedAt.String()),
		PublishedAt:      stringToPointer(c.PublishedAt.String()),
	}, nil
}

func sortingMethodToTag(s *model.TagSortingMethod) tagdomain.SortingMethod {
	if s == nil {
		*s = model.TagSortingMethodID
	}

	return tagdomain.SortingMethod(s.String())
}

func tagFromTag(t tagdomain.Tag) (model.Tag, error) {
	thumb512x384, err := pathTag.Thumbnail(t.Slug, thumbnail512x384)
	if err != nil {
		return model.Tag{}, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb128x128, err := pathTag.Thumbnail(t.Slug, thumbnail128x128)
	if err != nil {
		return model.Tag{}, fmt.Errorf("failed to get 128x128 thumbnail: %w", err)
	}

	return model.Tag{
		ID:               t.ID,
		Language:         model.Language(t.Language),
		Slug:             t.Slug,
		Name:             t.Name,
		ShortDescription: &t.ShortDescription,
		Description:      &t.Description,
		Content:          &t.Content,
		Status:           model.Status(t.Status),
		Clicks:           t.Clicks,
		CreatedAt:        t.CreatedAt.String(),
		DeletedAt:        stringToPointer(t.DeletedAt.String()),
		PublishedAt:      stringToPointer(t.PublishedAt.String()),
		Thumbnail512x384: thumb512x384,
		Thumbnail128x128: thumb128x128,
	}, nil
}

func sortingMethodToPointer[T model.SortingMethod | model.TagSortingMethod](m T) *T {
	return &m
}
