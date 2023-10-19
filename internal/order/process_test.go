package order

import (
	"context"
	"testing"

	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/stretchr/testify/assert"
)

func TestProcessOrder(t *testing.T) {
	var tests = []struct {
		name                  string
		cartID                string
		expectedDiscount      int32
		expectedShippingPrice float32
		expectedTotalProducts int32
		expectedTotalPrice    float32
	}{
		{
			name:                  "Order process without promo/discount",
			cartID:                testCartID,
			expectedTotalProducts: 4,
			expectedShippingPrice: 100,
			expectedDiscount:      0,
			expectedTotalPrice:    137,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			repo := repositoryMock{}
			cartService := cart.NewCartService(&repo)
			c, _ := cartService.Get(context.Background(), testCartIDWithoutDiscount)
			orderProcess := OrderProcess{}

			// when
			order := orderProcess.New(c).SetConfig().ApplyDiscounts().Build()

			// then
			assert.Equal(t, tt.expectedTotalProducts, order.Totals.Products)
			assert.Equal(t, tt.expectedDiscount, order.Totals.Discounts)
			assert.Equal(t, tt.expectedShippingPrice, order.Totals.Shipping)
			assert.Equal(t, tt.expectedTotalPrice, order.Totals.TotalOrderPrice)
		})
	}
}
