package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrders(t *testing.T) {
	var tests = []struct {
		name                  string
		cartID                string
		expectedError         bool
		expectedDiscount      int32
		expectedShippingPrice float32
		expectedTotalProducts int32
		expectedTotalPrice    float32
	}{
		{
			name:                  "Order Created Success",
			cartID:                testCartID,
			expectedError:         false,
			expectedTotalProducts: 31,
			expectedShippingPrice: 0,
			expectedDiscount:      10,
			expectedTotalPrice:    578,
		},
		{
			name:                  "Order Created failed because Cart not found",
			cartID:                testCartIDNotFound,
			expectedError:         true,
			expectedTotalProducts: 0,
			expectedShippingPrice: 0,
			expectedDiscount:      0,
			expectedTotalPrice:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			repo := repositoryMock{}
			orderServ := NewOrdersService(&repo)

			// when
			order, err := orderServ.Create(context.Background(), tt.cartID)

			// then
			assert.Equal(t, tt.expectedError, err != nil)
			if !tt.expectedError {
				assert.Equal(t, tt.expectedTotalProducts, order.Totals.Products)
				assert.Equal(t, tt.expectedDiscount, order.Totals.Discounts)
				assert.Equal(t, tt.expectedShippingPrice, order.Totals.Shipping)
				assert.Equal(t, tt.expectedTotalPrice, order.Totals.TotalOrderPrice)
			}
		})

	}
}
