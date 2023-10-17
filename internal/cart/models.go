package cart

import "github.com/ibanezv/go_trafilea_cart/internal/product"

type Cart struct {
	CartID   string
	UserID   string
	Products []product.Product
}

type ProductUpdate struct {
	ProductID string
	Quantity  int32
}
