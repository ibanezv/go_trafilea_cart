package order

type OrderBuilder interface {
	SetConfig(config DiscountsConfig) OrderBuilder
	ApplyDiscount() OrderBuilder
	Build() *Order
}
