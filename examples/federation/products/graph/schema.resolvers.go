package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/wundergraph/graphql-go-tools/examples/federation/products/graph/generated"
	"github.com/wundergraph/graphql-go-tools/examples/federation/products/graph/model"
)

// TopProducts is the resolver for the topProducts field.
func (r *queryResolver) TopProducts(ctx context.Context, first *int) ([]*model.Product, error) {
	return hats, nil
}

// UpdatedPrice is the resolver for the updatedPrice field.
func (r *subscriptionResolver) UpdatedPrice(ctx context.Context) (<-chan *model.Product, error) {
	updatedPrice := make(chan *model.Product)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(updateInterval):
				rand.Seed(time.Now().UnixNano())
				product := hats[0]

				if randomnessEnabled {
					product = hats[rand.Intn(len(hats)-1)]
					product.Price = rand.Intn(maxPrice-minPrice+1) + minPrice
					updatedPrice <- product
					continue
				}

				product.Price = currentPrice
				currentPrice += 1
				updatedPrice <- product
			}
		}
	}()
	return updatedPrice, nil
}

// UpdateProductPrice is the resolver for the updateProductPrice field.
func (r *subscriptionResolver) UpdateProductPrice(ctx context.Context, upc string) (<-chan *model.Product, error) {
	updatedPrice := make(chan *model.Product)
	var product *model.Product

	for _, hat := range hats {
		if hat.Upc == upc {
			product = hat
			break
		}
	}

	if product == nil {
		return nil, fmt.Errorf("unknown product upc: %s", upc)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				rand.Seed(time.Now().UnixNano())
				min := 10
				max := 1499
				product.Price = rand.Intn(max-min+1) + min
				updatedPrice <- product
			}
		}
	}()

	return updatedPrice, nil
}

// Stock is the resolver for the stock field.
func (r *subscriptionResolver) Stock(ctx context.Context) (<-chan []*model.Product, error) {
	stock := make(chan []*model.Product)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				rand.Seed(time.Now().UnixNano())
				randIndex := rand.Intn(len(hats))

				if hats[randIndex].InStock > 0 {
					hats[randIndex].InStock--
				}

				stock <- hats
			}
		}
	}()

	return stock, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
