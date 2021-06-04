package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jensneuse/federation-example/ratings/graph/generated"
	"github.com/jensneuse/federation-example/ratings/graph/model"
)

func (r *productResolver) Ratings(ctx context.Context, obj *model.Product) ([]*model.Rating, error) {
	var res []*model.Rating

	for _, Rating := range Ratings {
		if Rating.Product.Upc == obj.Upc {
			res = append(res, Rating)
		}
	}

	return res, nil
}

func (r *userResolver) Ratings(ctx context.Context, obj *model.User) ([]*model.Rating, error) {
	var res []*model.Rating

	for _, Rating := range Ratings {
		if Rating.Author.ID == obj.ID {
			res = append(res, Rating)
		}
	}

	return res, nil
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type productResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
