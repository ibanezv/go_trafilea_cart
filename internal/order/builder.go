package order

import "github.com/ibanezv/go_trafilea_cart/internal/cart"

type OrderBuilder interface {
	New(cart.Cart) OrderBuilder
	SetConfig() OrderBuilder
	ApplyDiscounts() OrderBuilder
	Build() *Order
}
