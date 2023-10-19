package cart

import (
	"context"
	"testing"

	"github.com/ibanezv/go_trafilea_cart/internal/product"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCartCreate(t *testing.T) {
	var tests = []struct {
		name     string
		expected repository.CartDB
	}{
		{
			name:     "Cart Created Success",
			expected: repository.CartDB{UserID: testUserID, CartID: testCartID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			repo := repositoryMock{}
			cartServ := NewCartService(&repo)

			// when
			cart, err := cartServ.Create(context.Background(), testUserID)

			// then
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.CartID, cart.CartID)
			assert.Equal(t, tt.expected.UserID, cart.UserID)
		})

	}
}

func TestCartAddProduct(t *testing.T) {
	var tests = []struct {
		name          string
		CartID        string
		productUpdate ProductUpdate
		cartExpected  Cart
		expectedError bool
	}{
		{
			name:          "Add product to cart Success",
			CartID:        testCartID,
			productUpdate: ProductUpdate{ProductID: testProductID, Quantity: 3},
			cartExpected:  Cart{CartID: testCartID, UserID: testUserID, Products: []product.Product{{ProductID: testProductID}}},
			expectedError: false,
		},
		{
			name:          "Add product failed cart not found",
			CartID:        testCartIDNotFound,
			productUpdate: ProductUpdate{ProductID: testProductID, Quantity: 3},
			cartExpected:  Cart{CartID: "", UserID: "", Products: []product.Product{}},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			repo := repositoryMock{}
			cartServ := NewCartService(&repo)

			// when
			cart, err := cartServ.AddProduct(context.Background(), tt.CartID, tt.productUpdate)

			// then
			assert.Equal(t, tt.expectedError, err != nil)
			assert.Equal(t, tt.cartExpected.CartID, cart.CartID)
			assert.Equal(t, len(tt.cartExpected.Products), len(cart.Products))
		})

	}
}

func TestCartModifyProduct(t *testing.T) {
	var tests = []struct {
		name             string
		CartID           string
		productUpdate    ProductUpdate
		cartExpected     Cart
		expectedError    bool
		quantityExpected int32
	}{
		{
			name:             "Modify count of product to cart Success",
			CartID:           testCartID,
			productUpdate:    ProductUpdate{ProductID: testProductID, Quantity: 3},
			cartExpected:     Cart{CartID: testCartID, UserID: testUserID, Products: []product.Product{{ProductID: testProductID}}},
			expectedError:    false,
			quantityExpected: 3,
		},
		{
			name:          "Modify count of product to cart Failed product not found in cart",
			CartID:        testCartID,
			productUpdate: ProductUpdate{ProductID: testProductIDNotFound, Quantity: 3},
			cartExpected:  Cart{},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			repo := repositoryMock{}
			cartServ := NewCartService(&repo)

			// when
			cart, err := cartServ.ModifyProduct(context.Background(), tt.CartID, tt.productUpdate)

			// then
			assert.Equal(t, tt.expectedError, err != nil)
			if !tt.expectedError {
				assert.Equal(t, tt.cartExpected.CartID, cart.CartID)
				assert.Equal(t, tt.quantityExpected, cart.Products[0].Quantity)
			}

		})

	}
}
