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
	OrderCreate(order OrderDB) OrderDB
	OrderUpdate(order OrderDB) (OrderDB, error)
	OrderGet(cartID string) (OrderDB, error)
	ProductGet(ProductID string) (ProductDB, error)
	FillProducts()
}

type InMemoryRepository struct {
	cart    map[string]CartDB
	order   map[string]OrderDB
	product map[string]ProductDB
	lock    sync.RWMutex
}

var ErrRecordNotFound = errors.New("record not found")

func (r *InMemoryRepository) CartCreate(userID string) (CartDB, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	cartID := uuid.Must(uuid.NewRandom()).String()
	cartDB := CartDB{UserID: userID, CartID: cartID}
	r.cart[cartID] = cartDB
	return r.cart[cartID], nil
}

func (r *InMemoryRepository) CartUpdate(cartDbUpdate CartDB) (CartDB, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
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

func (r *InMemoryRepository) OrderCreate(o OrderDB) OrderDB {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.order[o.CartID] = o
	return r.order[o.CartID]
}

func (r *InMemoryRepository) OrderUpdate(o OrderDB) (OrderDB, error) {
	return OrderDB{}, nil
}

func (r *InMemoryRepository) OrderGet(cartID string) (OrderDB, error) {
	if orderDB, found := r.order[cartID]; found {
		return orderDB, nil
	}
	return OrderDB{}, ErrRecordNotFound
}

func (r *InMemoryRepository) ProductGet(ProductID string) (ProductDB, error) {
	if productDB, found := r.product[ProductID]; found {
		return productDB, nil
	}
	return ProductDB{}, ErrRecordNotFound
}

func (r *InMemoryRepository) FillProducts() {
	r.product["1"] = ProductDB{ProductID: "1", Name: "producto-1", Category: "Coffee", Price: 15.0}
	r.product["2"] = ProductDB{ProductID: "2", Name: "producto-2", Category: "Equipment", Price: 22.0}
	r.product["3"] = ProductDB{ProductID: "3", Name: "producto-3", Category: "Accessories", Price: 19.0}
}

func NewRepository() Repository {
	return &InMemoryRepository{cart: make(map[string]CartDB), order: make(map[string]OrderDB), product: make(map[string]ProductDB), lock: sync.RWMutex{}}
}
