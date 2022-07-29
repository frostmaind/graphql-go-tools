package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jensneuse/graphql-go-tools/examples/federation/reviews/graph/generated"
	"github.com/jensneuse/graphql-go-tools/examples/federation/reviews/graph/model"
)

// Reviews is the resolver for the reviews field.
func (r *productResolver) Reviews(ctx context.Context, obj *model.Product) ([]*model.Review, error) {
	var res []*model.Review

	for _, review := range reviews {
		if review.Product.Upc == obj.Upc {
			res = append(res, review)
		}
	}

	return res, nil
}

// UserOrProduct is the resolver for the userOrProduct field.
func (r *queryResolver) UserOrProduct(ctx context.Context) ([]model.UserOrProduct, error) {
	return []model.UserOrProduct{
		model.User{ID: "1"},
		model.Product{Upc: "top-1"},
		model.User{ID: "2"},
		model.Product{Upc: "top-2"},
		//model.User{ID: "3"},
	}, nil
}

// Reviews is the resolver for the reviews field.
func (r *userResolver) Reviews(ctx context.Context, obj *model.User) ([]*model.Review, error) {
	var res []*model.Review

	for _, review := range reviews {
		if review.Author.ID == obj.ID {
			res = append(res, review)
		}
	}

	return res, nil
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *userResolver) Username(ctx context.Context, obj *model.User) (string, error) {
	return fmt.Sprintf("User %s", obj.ID), nil
}
