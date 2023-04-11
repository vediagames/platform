package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"

	"github.com/vediagames/platform/gateway/graphql/model"
	"github.com/vediagames/platform/webproxy/graphql/generated"
)

// Tags is the resolver for the tags field.
func (r *gameResolver) Tags(ctx context.Context, obj *model.Game) (*model.Tags, error) {
	return r.gatewayResolver.Game().Tags(ctx, obj)
}

// Categories is the resolver for the categories field.
func (r *gameResolver) Categories(ctx context.Context, obj *model.Game) (*model.Categories, error) {
	return r.gatewayResolver.Game().Categories(ctx, obj)
}

// Thumbnail is the resolver for the thumbnail field.
func (r *gameResolver) Thumbnail(ctx context.Context, obj *model.Game, request model.ThumbnailRequest) (string, error) {
	return r.gatewayResolver.Game().Thumbnail(ctx, obj, request)
}

// Video is the resolver for the video field.
func (r *gameResolver) Video(ctx context.Context, obj *model.Game, original model.OriginalVideo) (string, error) {
	return r.gatewayResolver.Game().Video(ctx, obj, original)
}

// Thumbnail is the resolver for the thumbnail field.
func (r *searchItemResolver) Thumbnail(ctx context.Context, obj *model.SearchItem, request model.ThumbnailRequest) (string, error) {
	return r.gatewayResolver.SearchItem().Thumbnail(ctx, obj, request)
}

// Video is the resolver for the video field.
func (r *searchItemResolver) Video(ctx context.Context, obj *model.SearchItem, original model.OriginalVideo) (string, error) {
	return r.gatewayResolver.SearchItem().Video(ctx, obj, original)
}

// Tags is the resolver for the tags field.
func (r *sectionResolver) Tags(ctx context.Context, obj *model.Section) (*model.Tags, error) {
	return r.gatewayResolver.Section().Tags(ctx, obj)
}

// Categories is the resolver for the categories field.
func (r *sectionResolver) Categories(ctx context.Context, obj *model.Section) (*model.Categories, error) {
	return r.gatewayResolver.Section().Categories(ctx, obj)
}

// Games is the resolver for the games field.
func (r *sectionResolver) Games(ctx context.Context, obj *model.Section) (*model.Games, error) {
	return r.gatewayResolver.Section().Games(ctx, obj)
}

// Thumbnail is the resolver for the thumbnail field.
func (r *tagResolver) Thumbnail(ctx context.Context, obj *model.Tag, request model.ThumbnailRequest) (string, error) {
	return r.gatewayResolver.Tag().Thumbnail(ctx, obj, request)
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
