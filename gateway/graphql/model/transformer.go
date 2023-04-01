package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	searchdomain "github.com/vediagames/platform/search/domain"
	sectiondomain "github.com/vediagames/platform/section/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

func (g Games) IDs() []int {
	ids := make([]int, 0, len(g.Data))
	for _, e := range g.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (g Games) FromDomain(domain gamedomain.Games) (*Games, error) {
	games := &Games{
		Data:  make([]*Game, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainGame := range domain.Data {
		game, err := Game{}.FromDomain(domainGame)
		if err != nil {
			return nil, fmt.Errorf("failed to convert game: %w", err)
		}

		games.Data = append(games.Data, game)
	}

	return games, nil
}

func (g Game) FromDomain(domain gamedomain.Game) (*Game, error) {
	thumb512x384, err := pathGame.Thumbnail(domain.Slug, thumbnail512x384)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb512x512, err := pathGame.Thumbnail(domain.Slug, thumbnail512x512)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x512 thumbnail: %w", err)
	}

	return &Game{
		ID:                domain.ID,
		Language:          Language(domain.Language),
		Slug:              domain.Slug,
		Name:              domain.Name,
		Status:            Status(domain.Status),
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
		Mobile:            domain.Mobile,
		Thumbnail512x384:  thumb512x384,
		Thumbnail512x512:  thumb512x512,
		PageURL:           fmt.Sprintf("/game/%s", domain.Slug),
		FullScreenPageURL: fmt.Sprintf("/game/%s/fullscreen", domain.Slug),
		TagIDRefs:         domain.TagIDRefs,
		CategoryIDRefs:    domain.CategoryIDRefs,
	}, nil
}

func (t Tags) IDs() []int {
	ids := make([]int, 0, len(t.Data))
	for _, e := range t.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (t Tags) FromDomain(domain tagdomain.Tags) (*Tags, error) {
	tags := &Tags{
		Data:  make([]*Tag, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainTag := range domain.Data {
		tag, err := Tag{}.FromDomain(domainTag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag: %w", err)
		}

		tags.Data = append(tags.Data, tag)
	}

	return tags, nil
}

func (t Tag) FromDomain(domain tagdomain.Tag) (*Tag, error) {
	thumb512x384, err := pathTag.Thumbnail(domain.Slug, thumbnail512x384)
	if err != nil {
		return nil, fmt.Errorf("failed to get 512x384 thumbnail: %w", err)
	}

	thumb128x128, err := pathTag.Thumbnail(domain.Slug, thumbnail128x128)
	if err != nil {
		return nil, fmt.Errorf("failed to get 128x128 thumbnail: %w", err)
	}

	return &Tag{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		Thumbnail512x384: thumb512x384,
		Thumbnail128x128: thumb128x128,
		PageURL:          fmt.Sprintf("/tag/%s", domain.Slug),
	}, nil
}

func (c Categories) IDs() []int {
	ids := make([]int, 0, len(c.Data))
	for _, e := range c.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (c Categories) FromDomain(domain categorydomain.Categories) *Categories {
	categories := &Categories{
		Data:  make([]*Category, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainCategory := range domain.Data {
		category := Category{}.FromDomain(domainCategory)
		categories.Data = append(categories.Data, category)
	}

	return categories
}

func (c Category) FromDomain(domain categorydomain.Category) *Category {
	return &Category{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		Status:           Status(domain.Status),
		Clicks:           domain.Clicks,
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		PageURL:          fmt.Sprintf("/category/%s", domain.Slug),
	}
}

func (s Sections) FromDomain(domain sectiondomain.Sections) (*Sections, error) {
	sections := &Sections{
		Data:  make([]*Section, 0, len(domain.Data)),
		Total: domain.Total,
	}

	for _, domainSection := range domain.Data {
		section, err := Section{}.FromDomain(domainSection)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		sections.Data = append(sections.Data, section)
	}

	return sections, nil
}

func (S Section) FromDomain(domain sectiondomain.Section) (*Section, error) {
	pageURL := "/continue-playing"

	if domain.Slug != "continue-playing" {
		var paramsFilter filterParams

		switch domain.Slug {
		case "newest":
			paramsFilter = filterParams{
				Sort: SortingMethodNewest,
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
			return &Section{}, fmt.Errorf("failed to encode filter params: %w", err)
		}

		pageURL = fmt.Sprintf("/filter?title=%s&params=%s", url.QueryEscape(domain.Name), params)
	}

	return &Section{
		ID:               domain.ID,
		Language:         Language(domain.Language),
		Slug:             domain.Slug,
		Name:             domain.Name,
		Status:           Status(domain.Status),
		CreatedAt:        domain.CreatedAt.String(),
		DeletedAt:        stringToPointer(domain.DeletedAt.String()),
		PublishedAt:      stringToPointer(domain.PublishedAt.String()),
		ShortDescription: stringToPointer(domain.ShortDescription),
		Description:      stringToPointer(domain.Description),
		Content:          stringToPointer(domain.Content),
		PageURL:          pageURL,
		TagIDRefs:        domain.TagIDRefs,
		CategoryIDRefs:   domain.CategoryIDRefs,
		GameIDRefs:       domain.GameIDRefs,
	}, nil
}

func (s PlacedSections) FromDomain(domain sectiondomain.GetPlacedResponse) (*PlacedSections, error) {
	placedSections := &PlacedSections{
		Data: make([]*PlacedSection, 0, len(domain.Data)),
	}

	for _, domainSection := range domain.Data {
		section, err := Section{}.FromDomain(domainSection.Section)
		if err != nil {
			return nil, fmt.Errorf("failed to convert section: %w", err)
		}

		placedSections.Data = append(placedSections.Data, &PlacedSection{
			Section:   section,
			Placement: domainSection.PlacementNumber,
		})
	}

	return placedSections, nil
}

func (s SearchItems) FromDomain(domain searchdomain.SearchResponse) (*SearchItems, error) {
	searchResponse := &SearchItems{
		Data:  make([]*SearchItem, 0, len(domain.Games)+len(domain.Tags)),
		Total: domain.Total,
	}

	for _, domainItem := range domain.Games {
		thumb512x384, err := pathGame.Thumbnail(domainItem.Slug, thumbnail512x384)
		if err != nil {
			return nil, fmt.Errorf("failed to get 512x384 thumbnail for %q: %w", domainItem.Slug, err)
		}

		searchResponse.Data = append(searchResponse.Data, &SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             SearchItemTypeGame,
			PageURL:          fmt.Sprintf("/game/%s", domainItem.Slug),
			Thumbnail512x384: thumb512x384,
		})
	}

	for _, domainItem := range domain.Tags {
		thumb512x384, err := pathTag.Thumbnail(domainItem.Slug, thumbnail512x384)
		if err != nil {
			return nil, fmt.Errorf("failed to get 512x384 thumbnail for %q: %w", domainItem.Slug, err)
		}

		searchResponse.Data = append(searchResponse.Data, &SearchItem{
			ShortDescription: domainItem.ShortDescription,
			Name:             domainItem.Slug,
			Slug:             domainItem.Slug,
			Type:             SearchItemTypeTag,
			PageURL:          fmt.Sprintf("/tag/%s", domainItem.Slug),
			Thumbnail512x384: thumb512x384,
		})
	}

	return searchResponse, nil
}

func stringToPointer(s string) *string {
	return &s
}

type filterParams struct {
	Sort       SortingMethod `json:"sort,omitempty"`
	Tags       []int         `json:"tags,omitempty"`
	Categories []int         `json:"categories,omitempty"`
	Games      []int         `json:"games,omitempty"`
}

func filterParamsInBase64(p filterParams) (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to marshal filter params: %w", err)
	}

	return base64.StdEncoding.EncodeToString(body), nil
}

func (m *SortingMethod) Domain() string {
	if m == nil {
		return SortingMethodID.String()
	}

	str := strings.Replace(m.String(), "_", "-", -1)

	return str
}
