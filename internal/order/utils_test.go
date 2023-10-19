package order

import (
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

var (
	testUserID                = "3c3ac4f2-9d08-4373-85ca-60f8798b8f9b"
	testCartID                = "e14d7379-b6f3-4c36-9bbd-ad95bccc7593"
	testCartIDNotFound        = "11111-bg566-4c38-9bbx-11"
	testCartIDWithoutDiscount = "00000-bg566-4c38-9bbx-000000000"
	testProductID             = "5488"
	testProductID2            = "5489"
	testProductID3            = "5490"
	testProductID4            = "5400"
	testProductID5            = "5401"
	testProductID6            = "5402"
	testProductCartDB         = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID, Name: "Test Product 1", Category: "Coffee", Price: 11.5}, Quantity: 12}
	testProductCartDB2        = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID2, Name: "Test Product 2", Category: "Equipment", Price: 50.5}, Quantity: 4}
	testProductCartDB3        = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID3, Name: "Test Product 3", Category: "Accessories", Price: 22.5}, Quantity: 14}

	testProductCartDB4 = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID4, Name: "Test Product 4", Category: "Coffee", Price: 11.5}, Quantity: 1}
	testProductCartDB5 = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID5, Name: "Test Product 5", Category: "Equipment", Price: 1.5}, Quantity: 2}
	testProductCartDB6 = repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: testProductID6, Name: "Test Product 6", Category: "Accessories", Price: 22.5}, Quantity: 1}
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
		c.Products = append(c.Products, testProductCartDB, testProductCartDB2, testProductCartDB3)
		return c, nil
	}
	if cartID == testCartIDWithoutDiscount {
		c := repository.CartDB{UserID: testUserID, CartID: cartID}
		c.Products = append(c.Products, testProductCartDB4, testProductCartDB5, testProductCartDB6)
		return c, nil
	}
	return repository.CartDB{}, repository.ErrRecordNotFound
}

func (r *repositoryMock) OrderCreate(order repository.OrderDB) repository.OrderDB {
	return order
}

func (r *repositoryMock) OrderUpdate(order repository.OrderDB) (repository.OrderDB, error) {
	return repository.OrderDB{}, nil
}

func (r *repositoryMock) OrderGet(cartID string) (repository.OrderDB, error) {
	return repository.OrderDB{}, nil
}

func (r *repositoryMock) ProductGet(ProductID string) (repository.ProductDB, error) {
	return repository.ProductDB{}, repository.ErrRecordNotFound
}

func (r *repositoryMock) FillProducts() {
	panic("unimplemented")
}
