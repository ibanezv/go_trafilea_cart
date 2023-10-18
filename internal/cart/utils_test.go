package cart

import (
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

var (
	testUserID            = "3c3ac4f2-9d08-4373-85ca-60f8798b8f9b"
	testCartID            = "e14d7379-b6f3-4c36-9bbd-ad95bccc7593"
	testCartIDNotFound    = "00000-b6f3-4c36-9bbd-000000000"
	testProductID         = "5488"
	testProductIDNotFound = "7566"
	testProductDB         = repository.ProductDB{ProductID: testProductID, Name: "Test Product", Category: "Coffee", Price: 11.5}
	testProductCartDB     = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID, Name: "Test Product", Category: "Coffee", Price: 11.5}, Quantity: 2}
)

type repositoryMock struct{}

func (r *repositoryMock) CartCreate(userID string) (repository.CartDB, error) {
	return repository.CartDB{UserID: testUserID, CartID: testCartID}, nil
}

func (r *repositoryMock) CartUpdate(cartDB repository.CartDB) (repository.CartDB, error) {
	return cartDB, nil
}

func (r *repositoryMock) CartGet(cartID string) (repository.CartDB, error) {
	if cartID == testCartID {
		c := repository.CartDB{UserID: testUserID, CartID: cartID}
		c.Products = append(c.Products, testProductCartDB)
		return c, nil
	}
	return repository.CartDB{}, repository.ErrRecordNotFound
}

func (r *repositoryMock) OrderCreate(order repository.OrderDB) repository.OrderDB {
	return repository.OrderDB{}
}

func (r *repositoryMock) OrderUpdate(order repository.OrderDB) (repository.OrderDB, error) {
	return repository.OrderDB{}, nil
}

func (r *repositoryMock) OrderGet(cartID string) (repository.OrderDB, error) {
	return repository.OrderDB{}, nil
}

func (r *repositoryMock) ProductGet(ProductID string) (repository.ProductDB, error) {
	if ProductID == testProductID {
		return testProductDB, nil
	}
	return repository.ProductDB{}, repository.ErrRecordNotFound
}
