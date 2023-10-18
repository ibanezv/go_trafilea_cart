package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/ibanezv/go_trafilea_cart/internal/order"
)

func PostCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartRequest := cartRequest{}
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

		cart, err := cartService.Create(cartRequest.UserID)
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
		w.Write(body)
	}
}

func GetCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		c, err := cartService.Get(cartID)
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
		w.Write(body)
	}
}

func PutCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]

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

		c, err := cartService.AddProduct(cartID, productUpdate)
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
		w.Write(body)
	}
}

func PutProductCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		productID := vars["productId"]
		productRequest := productsUpdateRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&productRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		productUpdate := cart.ProductUpdate{ProductID: productID, Quantity: productRequest.Quantity}
		if !productRequest.validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		c, err := cartService.ModifyProduct(cartID, productUpdate)
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
		w.Write(body)
	}
}

func PostOrder(orderService order.Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		order, err := orderService.Create(orderReq.CartID)
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
		w.Write(body)
	}
}

func GetOrder(orderService order.Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		if cartID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		o, err := orderService.Get(cartID)
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
		w.Write(body)
	}
}

func convertToCartResponse(c cart.Cart) cartResponse {
	resp := cartResponse{CartID: c.CartID, UseID: c.UserID}
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
