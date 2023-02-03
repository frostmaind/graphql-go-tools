package model

type Product struct {
	Upc string `json:"upc"`
}

func (Product) IsEntity() {}

func (Product) IsUserOrProduct() {}

type Review struct {
	Body    string
	Author  *User
	Product *Product
}

type User struct {
	ID string `json:"id"`
}

func (User) IsEntity() {}

func (User) IsUserOrProduct() {}
