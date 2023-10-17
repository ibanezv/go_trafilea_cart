package product

import (
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

type Products interface {
	Get(string) (Product, error)
}

type ProductsService struct {
	repo repository.Repository
}

func NewProductsService(repo repository.Repository) ProductsService {
	return ProductsService{repo: repo}
}

func (o *ProductsService) Get(cartID string) (Product, error) {
	return Product{}, nil
}
