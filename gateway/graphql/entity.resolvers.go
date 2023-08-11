package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	categorydomain "github.com/vediagames/platform/category/domain"
	gamedomain "github.com/vediagames/platform/game/domain"
	"github.com/vediagames/platform/gateway/graphql/generated"
	"github.com/vediagames/platform/gateway/graphql/model"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

// Tags is the resolver for the tags field.
func (r *gameResolver) Tags(ctx context.Context, obj *model.Game) (*model.Tags, error) {
	if len(obj.TagIDRefs) == 0 {
		return nil, nil
	}

	svcRes, err := r.tagService.List(ctx, tagdomain.ListRequest{
		Language: tagdomain.Language(obj.Language),
		Page:     1,
		Limit:    20,
		Sort:     tagdomain.SortingMethodName,
		IDRefs:   obj.TagIDRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	return model.Tags{}.FromDomain(svcRes.Data), nil
}

// Categories is the resolver for the categories field.
func (r *gameResolver) Categories(ctx context.Context, obj *model.Game) (*model.Categories, error) {
	if len(obj.CategoryIDRefs) == 0 {
		return nil, nil
	}

	svcRes, err := r.categoryService.List(ctx, categorydomain.ListRequest{
		Language: categorydomain.Language(obj.Language),
		Page:     1,
		Limit:    20,
		IDRefs:   obj.CategoryIDRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	return model.Categories{}.FromDomain(svcRes.Data), nil
}

// Thumbnail is the resolver for the thumbnail field.
func (r *gameResolver) Thumbnail(ctx context.Context, obj *model.Game, request model.ThumbnailRequest) (string, error) {
	svcRes, err := r.imageService.Get(ctx, request.Domain(obj.Slug, false))
	if err != nil {
		return "", fmt.Errorf("failed to get: %w", err)
	}

	return svcRes.URL, nil
}

// Video is the resolver for the video field.
func (r *gameResolver) Video(ctx context.Context, obj *model.Game, original model.OriginalVideo) (string, error) {
	return fmt.Sprintf("%s/games/%s/%s", r.contentURL, obj.Slug, original.FileName()), nil
}

// Thumbnail is the resolver for the thumbnail field.
func (r *searchItemResolver) Thumbnail(ctx context.Context, obj *model.SearchItem, request model.ThumbnailRequest) (string, error) {
	svcRes, err := r.imageService.Get(ctx, request.Domain(obj.Slug, obj.Type == model.SearchItemTypeTag))
	if err != nil {
		return "", fmt.Errorf("failed to get: %w", err)
	}

	return svcRes.URL, nil
}

// Video is the resolver for the video field.
func (r *searchItemResolver) Video(ctx context.Context, obj *model.SearchItem, original model.OriginalVideo) (string, error) {
	if obj.Type != model.SearchItemTypeGame {
		return "", nil
	}

	return fmt.Sprintf("%s/games/%s/%s", r.contentURL, obj.Slug, original.FileName()), nil
}

// Tags is the resolver for the tags field.
func (r *sectionResolver) Tags(ctx context.Context, obj *model.Section) (*model.Tags, error) {
	if len(obj.TagIDRefs) == 0 {
		return nil, nil
	}

	svcRes, err := r.tagService.List(ctx, tagdomain.ListRequest{
		Language: tagdomain.Language(obj.Language),
		Page:     1,
		Limit:    20,
		Sort:     tagdomain.SortingMethodName,
		IDRefs:   obj.TagIDRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	return model.Tags{}.FromDomain(svcRes.Data), nil
}

// Categories is the resolver for the categories field.
func (r *sectionResolver) Categories(ctx context.Context, obj *model.Section) (*model.Categories, error) {
	if len(obj.CategoryIDRefs) == 0 {
		return nil, nil
	}

	svcRes, err := r.categoryService.List(ctx, categorydomain.ListRequest{
		Language: categorydomain.Language(obj.Language),
		Page:     1,
		Limit:    20,
		IDRefs:   obj.CategoryIDRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	return model.Categories{}.FromDomain(svcRes.Data), nil
}

// Games is the resolver for the games field.
func (r *sectionResolver) Games(ctx context.Context, obj *model.Section) (*model.Games, error) {
	if obj.Games != nil && len(obj.Games.Data) != 0 {
		return obj.Games, nil
	}

	if len(obj.GameIDRefs) == 0 {
		return nil, nil
	}

	svcRes, err := r.gameService.List(ctx, gamedomain.ListRequest{
		Language: gamedomain.Language(obj.Language),
		Page:     1,
		Limit:    30,
		Sort:     gamedomain.SortingMethodMostLiked,
		IDRefs:   obj.TagIDRefs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	return model.Games{}.FromDomain(svcRes.Data), nil
}

// Thumbnail is the resolver for the thumbnail field.
func (r *tagResolver) Thumbnail(ctx context.Context, obj *model.Tag, request model.ThumbnailRequest) (string, error) {
	svcRes, err := r.imageService.Get(ctx, request.Domain(obj.Slug, true))
	if err != nil {
		return "", fmt.Errorf("failed to get: %w", err)
	}

	return svcRes.URL, nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// SearchItem returns generated.SearchItemResolver implementation.
func (r *Resolver) SearchItem() generated.SearchItemResolver { return &searchItemResolver{r} }

// Section returns generated.SectionResolver implementation.
func (r *Resolver) Section() generated.SectionResolver { return &sectionResolver{r} }

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type gameResolver struct{ *Resolver }
type searchItemResolver struct{ *Resolver }
type sectionResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
