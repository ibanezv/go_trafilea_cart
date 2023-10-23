//  Company Api:
//   version: 0.0.1
//   title: Trafilea Cart Api
//  Schemes: http, https
//  Host: localhost:5000
//  BasePath: /
//  Produces:
//    - application/json
//
// securityDefinitions:
//  apiKey:
//    type: apiKey
//    in: header
//    name: authorization
// swagger:meta

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/ibanezv/go_trafilea_cart/internal/order"
)

// swagger:operation POST /api/v1/cart postCart
// ---
// summary: Create a new Cart .
// ---
// description: Create a new cart for request userId.
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:

// responses:
//
//	"200":
//	  "$ref": "#/responses/CartResponse"
//	"404":
//	  "$ref": "#/responses/CartResponse"
func PostCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartRequest := CartRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&cartRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !cartRequest.validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cart, err := cartService.Create(ctx, cartRequest.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var body []byte
		resp := convertToCartResponse(cart)
		body, err = json.Marshal(resp)
		w.Header().Add("content-type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(body)
	}
}

// swagger:operation GET /api/v1/cart/{cartId} getCart
// ---
// summary: Get a cart by id.
// description: Get a cart by id of existing cart.
// parameters:
//   - name: cartId
//     in: path
//     description: cart identity
//     type: string
//     required: true
//
// responses:
//
//	"200":
//	  "$ref": "#/responses/CartResponse"
//	"404":
//	  "$ref": "#/responses/CartResponse"
func GetCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		ctx := r.Context()
		c, err := cartService.Get(ctx, cartID)
		w.Header().Add("content-type", "application/json")
		if err != nil {
			if errors.Is(err, cart.ErrCartNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(convertToCartResponse(c))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

// swagger:operation PUT /cart/{cartId} ProductUpdate
// ---
// summary: Add ammount of product in cart
// description: ccc.
// parameters:
//   - name: cartId
//     in: path
//     description: cart identity
//     type: string
//     required: true
func PutCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		ctx := r.Context()

		productRequest := productsUpdateRequest{}
		err := json.NewDecoder(r.Body).Decode(&productRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		productUpdate := cart.ProductUpdate{ProductID: productRequest.ProductID, Quantity: productRequest.Quantity}
		if !productRequest.validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		c, err := cartService.AddProduct(ctx, cartID, productUpdate)
		if err != nil {
			if errors.Is(err, cart.ErrCartNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(convertToCartResponse(c))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

// swagger:operation PUT /cart/{cartId}/product/{productId} productsUpdateRequest
// ---
// summary: Update ammount of product in cart
// description: update products.
// parameters:
//   - name: cartId
//     in: path
//     description: cart identity
//     type: string
//     required: true
//   - name: productId
//     in: path
//     description: product identity
//     type: string
//     required: true
func PutProductCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		productID := vars["productId"]
		ctx := r.Context()
		productRequest := productsUpdateRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&productRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		productRequest.ProductID = productID
		if !productRequest.validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		productUpdate := cart.ProductUpdate{ProductID: productRequest.ProductID, Quantity: productRequest.Quantity}
		c, err := cartService.ModifyProduct(ctx, cartID, productUpdate)
		if err != nil {
			if errors.Is(err, cart.ErrCartNotFound) || errors.Is(err, cart.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := json.Marshal(convertToCartResponse(c))
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

// swagger:operation POST /order orderRequest
// ---
// summary: ccc
// description: ccc.
// parameters:
//   - name: cartId
//     in: path
//     description: cart identity
//     type: string
//     required: true
func PostOrder(orderService order.Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orderReq := orderRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !orderReq.validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		order, err := orderService.Create(ctx, orderReq.CartID)
		if err != nil {
			if errors.Is(err, cart.ErrCartNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(convertToOrderResponse(order))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(body)
	}
}

// swagger:operation GET /order/{orderId}
// ---
// summary: ccc
// description: ccc.
// parameters:
//   - name: cartId
//     in: path
//     description: cart identity
//     type: string
//     required: true
func GetOrder(orderService order.Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		if cartID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		o, err := orderService.Get(ctx, cartID)
		if err != nil {
			if errors.Is(err, order.ErrOrderNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(convertToOrderResponse(o))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

func convertToCartResponse(c cart.Cart) CartResponse {
	resp := CartResponse{CartID: c.CartID, UseID: c.UserID}
	for _, p := range c.Products {
		prod := ProductResponse{ProductID: p.ProductID, Name: p.Name, Category: p.Category, Price: p.Price, Quantity: p.Quantity}
		resp.Products = append(resp.Products, prod)
	}
	return resp
}

func convertToOrderResponse(o order.Order) orderResponse {
	return orderResponse{CartID: o.CartID, Totals: OrderDetailResponse{Products: o.Totals.Products, Discounts: o.Totals.Discounts,
		Shipping: o.Totals.Shipping, TotalOrderPrice: o.Totals.TotalOrderPrice}}
}
