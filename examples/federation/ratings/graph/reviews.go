package graph

import "github.com/jensneuse/federation-example/ratings/graph/model"

var Ratings = []*model.Rating{
	{
		Rating:    2,
		Product: &model.Product{Upc: "top-1"},
		Author:  &model.User{ID: "1234"},
	},
	{
		Rating:    4,
		Product: &model.Product{Upc: "top-1"},
		Author:  &model.User{ID: "1234"},
	},
	{
		Rating:    2,
		Product: &model.Product{Upc: "top-1"},
		Author:  &model.User{ID: "7777"},
	},
	{
		Rating:    3,
		Product: &model.Product{Upc: "top-2"},
		Author:  &model.User{ID: "1234"},
	},
	{
		Rating:    1,
		Product: &model.Product{Upc: "top-2"},
		Author:  &model.User{ID: "7777"},
	},
	{
		Rating:   4,
		Product: &model.Product{Upc: "top-2"},
		Author:  &model.User{ID: "6666"},
	},
	//{
	//	Body:    "This is the last straw. Hat you will wear. 11/10",
	//	Product: &model.Product{Upc: "top-3"},
	//	Author:  &model.User{ID: "6666"},
	//},
}
