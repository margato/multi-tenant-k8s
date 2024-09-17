package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Product struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListProducts() []Product {
	n := 20
	products := make([]Product, n)

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("Product-%d", i)
		products[i] = Product{
			Id:   uuid.NewString(),
			Name: name,
		}
	}
	return products
}
