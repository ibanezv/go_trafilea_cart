package order

import (
	"errors"

	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

const shippingPrice = 100

type Orders interface {
	Create(string) (Order, error)
	Get(string) (Order, error)
}

type OrdersService struct {
	repo repository.Repository
}

func NewOrdersService(repo repository.Repository) OrdersService {
	return OrdersService{repo: repo}
}

func (o *OrdersService) Get(cartID string) (Order, error) {
	return Order{}, nil
}

func (o *OrdersService) Create(cartID string) (Order, error) {
	cartDB, err := o.repo.CartGet(cartID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return Order{}, cart.ErrCartNotFound
		}
		return Order{}, err
	}

	order := createOrder(cart.ConvertToCart(cartDB))
	orderDB := o.repo.OrderCreate(convertToOrderDB(order))

	return convertToOrder(orderDB), nil
}

func createOrder(cart cart.Cart) Order {
	promoConfig := DiscountsConfig{categoryPromo: validatePromoByCategory(COFFEE, cart),
		shippingPromo: validatePromoShipping(EQUIPMENT, cart), DiscountPromo: validateDiscounts(ACCESORIES, cart)}

	shippingPrice, discountPercentage, ammountDiscount := applyPromos(promoConfig, &cart)

	countTotalProducts := countTotalProducts(cart)
	ammountTotal := amountTotalPrice(cart)

	order := Order{CartID: cart.CartID, Totals: OrderDetail{Products: countTotalProducts, Discounts: discountPercentage, Shipping: shippingPrice,
		TotalOrderPrice: ((ammountTotal - (ammountTotal * float32(discountPercentage))) + ammountDiscount) + float32(shippingPrice)}}

	return order
}

func applyPromos(config DiscountsConfig, cart *cart.Cart) (int32, int32, float32) {
	shippingPrice := int32(shippingPrice)
	discountPerc := int32(0)
	ammountDiscount := float32(0)

	if config.categoryPromo {
		for _, p := range cart.Products {
			if p.Category == COFFEE {
				p.Quantity++
				ammountDiscount = (-1) * p.Price
				break
			}
		}
	}

	if config.shippingPromo {
		shippingPrice = 0
	}

	if config.DiscountPromo {
		discountPerc = 10
	}
	return shippingPrice, discountPerc, ammountDiscount
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

func countTotalProducts(cart cart.Cart) int32 {
	q := int32(0)
	for _, product := range cart.Products {
		q += product.Quantity
	}
	return q
}

func amountTotalPrice(cart cart.Cart) float32 {
	amount := float32(0)
	for _, product := range cart.Products {
		amount += (float32(product.Quantity) * product.Price)
	}
	return amount
}

func convertToOrderDB(order Order) repository.OrderDB {
	return repository.OrderDB{CartID: order.CartID, Totals: repository.OrderDetailDB{Products: order.Totals.Products,
		Discounts: order.Totals.Discounts, Shipping: order.Totals.Shipping, TotalOrderPrice: order.Totals.TotalOrderPrice}}
}

func convertToOrder(order repository.OrderDB) Order {
	return Order{CartID: order.CartID, Totals: OrderDetail{Products: order.Totals.Products, Shipping: order.Totals.Shipping,
		Discounts: order.Totals.Discounts, TotalOrderPrice: order.Totals.TotalOrderPrice}}
}
