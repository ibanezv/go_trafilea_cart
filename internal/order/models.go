package order

import "github.com/ibanezv/go_trafilea_cart/internal/repository"

type OrderDetail struct {
	Products        int32
	Discounts       int32
	Shipping        float32
	TotalOrderPrice float32
}

type Order struct {
	CartID string
	Totals OrderDetail
}

type DiscountsConfig struct {
	categoryPromo bool
	shippingPromo bool
	DiscountPromo bool
}

const (
	COFFEE     = "Coffee"
	EQUIPMENT  = "Equipment"
	ACCESORIES = "Accessories"
)

func (o Order) ToDB() repository.OrderDB {
	return repository.OrderDB{CartID: o.CartID, Totals: repository.OrderDetailDB{Products: o.Totals.Products,
		Discounts: o.Totals.Discounts, Shipping: o.Totals.Shipping, TotalOrderPrice: o.Totals.TotalOrderPrice}}
}
