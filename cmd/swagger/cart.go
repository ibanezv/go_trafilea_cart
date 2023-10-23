package swagger

import "github.com/ibanezv/go_trafilea_cart/cmd/api"

// swagger:parameters postCart getCart
type swaggerCartRequest struct {
	// in:body
	api.CartRequest
}

// swagger:response CartResponse
type swaggerCartResponse struct {
	// in:body
	Body struct {
		// HTTP status code 200 - Status OK
		Code int `json:"code"`
		// CartResponse models
		api.CartResponse `json:"data"`
	}
}
