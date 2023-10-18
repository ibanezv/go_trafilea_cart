package order

import "github.com/ibanezv/go_trafilea_cart/internal/cart"

const defaultShippingPrice = float32(100)

type OrderProcess struct {
	cart                *cart.Cart
	order               Order
	config              DiscountsConfig
	shippingPrice       float32
	discountPercentage  int32
	ammountDiscount     float32
	ammountProductsGift int32
}

func (o *OrderProcess) New(cart cart.Cart) OrderBuilder {
	o.cart = &cart
	return o
}

func (o *OrderProcess) SetConfig() OrderBuilder {
	o.config = DiscountsConfig{categoryPromo: validatePromoByCategory(COFFEE, *o.cart),
		shippingPromo: validatePromoShipping(EQUIPMENT, *o.cart), DiscountPromo: validateDiscounts(ACCESORIES, *o.cart)}
	return o
}

func (o *OrderProcess) ApplyDiscounts() OrderBuilder {
	o.shippingPrice = defaultShippingPrice

	if o.config.categoryPromo {
		for _, p := range o.cart.Products {
			if p.Category == COFFEE {
				o.ammountProductsGift = 1
				o.ammountDiscount = (-1) * p.Price
				break
			}
		}
	}

	if o.config.shippingPromo {
		o.shippingPrice = 0
	}

	if o.config.DiscountPromo {
		o.discountPercentage = 10
	}
	return o
}

func (o *OrderProcess) Build() *Order {
	countTotalProducts := countTotalProducts(o.cart)
	ammountTotal := amountTotalPrice(o.cart)

	o.order = Order{CartID: o.cart.CartID, Totals: OrderDetail{Products: countTotalProducts + o.ammountProductsGift,
		Discounts: o.discountPercentage, Shipping: o.shippingPrice,
		TotalOrderPrice: ((ammountTotal - (ammountTotal * float32(o.discountPercentage))) + o.ammountDiscount) + o.shippingPrice}}

	return &o.order
}

func validatePromoByCategory(category string, cart cart.Cart) bool {
	return countProductsByCategory(category, cart) >= 2
}

func validatePromoShipping(category string, cart cart.Cart) bool {
	return countProductsByCategory(category, cart) > 3
}

func validateDiscounts(category string, cart cart.Cart) bool {
	return amountByCategory(category, cart) > 70
}

func countTotalProducts(cart *cart.Cart) int32 {
	q := int32(0)
	for _, product := range cart.Products {
		q += product.Quantity
	}
	return q
}

func amountTotalPrice(cart *cart.Cart) float32 {
	amount := float32(0)
	for _, product := range cart.Products {
		amount += (float32(product.Quantity) * product.Price)
	}
	return amount
}

func countProductsByCategory(category string, cart cart.Cart) int32 {
	q := int32(0)
	for _, product := range cart.Products {
		if product.Category == category {
			q += product.Quantity
		}
	}
	return q
}

func amountByCategory(category string, cart cart.Cart) float32 {
	q := float32(0)
	for _, product := range cart.Products {
		if product.Category == category {
			q += product.Price
		}
	}
	return q
}
