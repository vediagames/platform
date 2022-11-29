package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"github.com/vediagames/vediagames.com/bff/graphql/model"
	categorydomain "github.com/vediagames/vediagames.com/category/domain"
	gamedomain "github.com/vediagames/vediagames.com/game/domain"
	searchdomain "github.com/vediagames/vediagames.com/search/domain"
	sectiondomain "github.com/vediagames/vediagames.com/section/domain"
	tagdomain "github.com/vediagames/vediagames.com/tag/domain"
)

type path string

const (
	pathGame path = "games"
	pathTag  path = "tags"
)

func (p path) Path(slug string, file string) string {
	return fmt.Sprintf("%s/%s/%s", p, slug, file)
}

func (p path) Thumbnail(slug string, t thumbnail) (string, error) {
	switch p {
	case pathGame:
		if t == thumbnail128x128 {
			return "", fmt.Errorf("thumbnail 128x128 not available for games")
		}
	case pathTag:
		if t == thumbnail512x512 {
			return "", fmt.Errorf("thumbnail 512x512 not available for tags")
		}
	default:
		return "", fmt.Errorf("thumbnails not available for %s", p)
	}

	return fmt.Sprintf("https://images.vediagames.com/file/vg-images/%s", p.Path(slug, t.JPG())), nil
}

type thumbnail string

func (t thumbnail) JPG() string {
	return fmt.Sprintf("%s.jpg", t)
}

const (
	thumbnail128x128 thumbnail = "thumb128x128"
	thumbnail512x384 thumbnail = "thumb512x384"
	thumbnail512x512 thumbnail = "thumb512x512"
)

func searchFromSearch(games []searchdomain.SearchItem, tags []searchdomain.SearchItem, total int, fullSearch bool) (model.SearchResponse, error) {
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
			ShortDescription: game.ShortDescription,
			Name:             game.Name,
			Slug:             game.Slug,
			Type:             model.SearchItemTypeGame,
			Link:             fmt.Sprintf("/games/%s", game.Slug),
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
			ShortDescription: tag.ShortDescription,
			Name:             tag.Name,
			Slug:             tag.Slug,
			Type:             model.SearchItemTypeTag,
			Link:             fmt.Sprintf("/tag/%d?slug=%s&name=%s", tag.ID, tag.Slug, tag.Name),
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
	tags := make([]model.ComplimentaryTag, 0, len(s.Tags.Data))
	for _, tag := range s.Tags.Data {
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

	categories := make([]model.ComplimentaryCategory, 0, len(s.Categories.Data))
	for _, category := range s.Categories.Data {
		categories = append(categories, model.ComplimentaryCategory{
			ID:          category.ID,
			Slug:        category.Slug,
			Name:        category.Name,
			Description: &category.Description,
		})
	}

	link := "/continue-playing"

	if s.Slug != "continue-playing" {
		var paramsFilter filterParams

		switch s.Slug {
		case "newest":
			paramsFilter = filterParams{
				Sort: model.SortingMethodNewest,
			}
		default:
			paramsFilter = filterParams{
				Tags:       s.Tags.IDs(),
				Categories: s.Categories.IDs(),
				Games:      s.Games,
			}
		}

		params, err := filterParamsInBase64(paramsFilter)
		if err != nil {
			return model.Section{}, fmt.Errorf("failed to encode filter params: %w", err)
		}

		link = fmt.Sprintf("/filter?title=%s&params=%s", url.QueryEscape(s.Name), params)
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
		PageURL:     link,
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
		PageURL:          fmt.Sprintf("/game/%s", game.Slug),
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
		PageURL:          fmt.Sprintf("/category/%s?id=%d", c.Slug, c.ID),
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
		PageURL:          fmt.Sprintf("/tag/%d?slug=%s&name=%s", t.ID, t.Slug, t.Name),
	}, nil
}

func sortingMethodToPointer[T model.SortingMethod | model.TagSortingMethod](m T) *T {
	return &m
}
