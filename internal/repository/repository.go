package repository

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Repository interface {
	CartCreate(userID string) (CartDB, error)
	CartUpdate(cartDB CartDB) (CartDB, error)
	CartGet(cartID string) (CartDB, error)
	OrderCreate(order OrderDB) (OrderDB, error)
	OrderUpdate(order OrderDB) (OrderDB, error)
	OrderGet(cartID string) (OrderDB, error)
	ProductGet(ProductID string) (ProductDB, error)
}

type InMemoryRepository struct {
	cart    map[string]CartDB
	order   map[string]OrderDB
	product map[string]ProductDB
	lock    sync.RWMutex
}

var ErrRecordNotFound = errors.New("record not found")

func (r *InMemoryRepository) CartCreate(userID string) (CartDB, error) {
	cartID := uuid.Must(uuid.NewRandom()).String()
	cartDB := CartDB{UserID: userID, CartID: cartID}
	r.cart[cartID] = cartDB
	return r.cart[cartID], nil
}

func (r *InMemoryRepository) CartUpdate(cartDbUpdate CartDB) (CartDB, error) {
	if cartDB, found := r.cart[cartDbUpdate.CartID]; found {
		r.cart[cartDB.CartID] = cartDbUpdate
		return r.cart[cartDB.CartID], nil
	}
	return CartDB{}, ErrRecordNotFound
}

func (r *InMemoryRepository) CartGet(cartID string) (CartDB, error) {
	if cartDB, found := r.cart[cartID]; found {
		return cartDB, nil
	}
	return CartDB{}, ErrRecordNotFound
}

func (r *InMemoryRepository) OrderCreate(o OrderDB) (OrderDB, error) {
	return OrderDB{}, nil
}

func (r *InMemoryRepository) OrderUpdate(o OrderDB) (OrderDB, error) {
	return OrderDB{}, nil
}

func (r *InMemoryRepository) OrderGet(cartID string) (OrderDB, error) {
	return OrderDB{}, nil
}

func (r *InMemoryRepository) ProductGet(ProductID string) (ProductDB, error) {
	if productDB, found := r.product[ProductID]; found {
		return productDB, nil
	}
	return ProductDB{}, ErrRecordNotFound
}

func NewRepository() Repository {
	productTable := make(map[string]ProductDB)
	productTable["45"] = ProductDB{ProductID: "45", Name: "producto-1", Category: "CAFE", Price: 10.0}
	return &InMemoryRepository{cart: make(map[string]CartDB), order: make(map[string]OrderDB), product: productTable, lock: sync.RWMutex{}}
}
