package order

type OrderDetail struct {
	Products        int32
	Discounts       int32
	Shipping        int32
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
