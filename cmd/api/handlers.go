package api

import (
	"encoding/json"
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
			w.Header().Add("Statuscode", "400")
			return
		}

		cart, err := cartService.Create(cartRequest.UserID)
		if err != nil {
			w.Header().Add("Statuscode", "500")
			return
		}

		var body []byte
		resp := convertToCartResponse(cart)
		body, err = json.Marshal(resp)
		w.Header().Add("content-type", "application/json")
		if err != nil {
			w.Header().Add("Statuscode", "500")
			return
		}
		w.Header().Add("StatusCode", "201")
		w.Write(body)
	}
}

func GetCart(cartService cart.Carts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cartID := vars["cartId"]
		cart, err := cartService.Get(cartID)
		w.Header().Add("content-type", "application/json")
		if err != nil {
			w.Header().Add("Statuscode", "400")
			return
		}

		body, err := json.Marshal(convertToCartResponse(cart))
		if err != nil {
			w.Header().Add("Statuscode", "500")
			return
		}

		w.Header().Add("StatusCode", "200")
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
			w.Header().Add("StatusCode", "400")
			return
		}
		productUpdate := cart.ProductUpdate{ProductID: productRequest.ProductID, Quantity: productRequest.Quantity}

		cart, err := cartService.AddProduct(cartID, productUpdate)
		if err != nil {
			w.Header().Add("StatusCode", "500")
			return
		}

		body, err := json.Marshal(convertToCartResponse(cart))
		if err != nil {
			w.Header().Add("StatusCode", "500")
			return
		}

		w.Header().Add("content-type", "application/json")
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
			w.Header().Add("StatusCode", "400")
		}
		productUpdate := cart.ProductUpdate{ProductID: productID, Quantity: productRequest.Quantity}

		cart, err := cartService.ModifyProduct(cartID, productUpdate)
		if err != nil {
			w.Header().Add("StatusCode", "400")
			return
		}

		body, _ := json.Marshal(convertToCartResponse(cart))
		w.Header().Add("content-type", "application/json")
		w.Write(body)
	}
}

func PostOrder(orderService order.Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderReq := orderRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			w.Header().Add("StatusCode", "400")
			return
		}
		order, err := orderService.Create(orderReq.CartID)
		if err != nil {
			w.Header().Add("StatusCode", "400")
			return
		}
		body, err := json.Marshal(convertToOrderResponse(order))
		if err != nil {
			w.Header().Add("StatusCode", "500")
			return
		}
		w.Header().Add("content-type", "application/json")
		w.Header().Add("StatusCode", "201")
		w.Write(body)
	}
}

func convertToCartResponse(c cart.Cart) cartResponse {
	resp := cartResponse{CartID: c.CartID, UseID: c.UserID}
	for _, p := range c.Products {
		prod := ProductResponse{ProductID: p.ProductID, Name: p.Name, Category: p.Category, Price: p.Price}
		resp.Products = append(resp.Products, prod)
	}
	return resp
}

func convertToOrderResponse(o order.Order) orderResponse {
	return orderResponse{CartID: o.CartID, Totals: OrderDetailResponse{Products: o.Totals.Products, Discounts: o.Totals.Discounts,
		Shipping: o.Totals.Shipping, TotalOrderPrice: o.Totals.TotalOrderPrice}}
}
