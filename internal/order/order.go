package order

import (
	"context"
	"errors"
	"log"

	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

var ErrOrderNotFound = errors.New("order not found")

type Orders interface {
	Create(context.Context, string) (Order, error)
	Get(context.Context, string) (Order, error)
}

type OrdersService struct {
	repo repository.Repository
}

func NewOrdersService(repo repository.Repository) OrdersService {
	return OrdersService{repo: repo}
}

func (o *OrdersService) Get(ctx context.Context, cartID string) (Order, error) {
	orderDB, err := o.repo.OrderGet(cartID)
	if err != nil {
		log.Printf("cart %s not found:%v", cartID, err)
		return Order{}, ErrOrderNotFound
	}
	return convertToOrder(orderDB), nil
}

func (o *OrdersService) Create(ctx context.Context, cartID string) (Order, error) {
	cartDB, err := o.repo.CartGet(cartID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("cart %s not found:%v", cartID, err)
			return Order{}, cart.ErrCartNotFound
		}
		return Order{}, err
	}

	process := OrderProcess{}
	order := process.New(cart.ConvertToCart(cartDB)).SetConfig().ApplyDiscounts().Build()
	orderDB := o.repo.OrderCreate(order.ToDB())

	return convertToOrder(orderDB), nil
}

func convertToOrder(order repository.OrderDB) Order {
	return Order{CartID: order.CartID, Totals: OrderDetail{Products: order.Totals.Products, Shipping: order.Totals.Shipping,
		Discounts: order.Totals.Discounts, TotalOrderPrice: order.Totals.TotalOrderPrice}}
}
