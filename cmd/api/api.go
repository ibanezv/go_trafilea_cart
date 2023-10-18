package api

type requestObject interface {
	validate() bool
}

type cartRequest struct {
	UserID string
}

type cartResponse struct {
	CartID   string            `json:"cart_id"`
	UseID    string            `json:"user_id"`
	Products []ProductResponse `json:"products"`
}

type ProductResponse struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Price     float32 `json:"price"`
	Quantity  int32   `json:"quantity"`
}

type productsUpdateRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type orderRequest struct {
	CartID string `json:"cart_id"`
}

type orderResponse struct {
	CartID string              `json:"cart_id"`
	Totals OrderDetailResponse `json:"totals"`
}

type OrderDetailResponse struct {
	Products        int32   `json:"products"`
	Discounts       int32   `json:"discounts"`
	Shipping        float32 `json:"shipping"`
	TotalOrderPrice float32 `json:"total_order_price"`
}

func (r productsUpdateRequest) validate() bool {
	return (r.ProductID != "" && r.Quantity >= 0)
}

func (r orderRequest) validate() bool {
	return r.CartID != ""
}

func (r cartRequest) validate() bool {
	return r.UserID != ""
}
