package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ibanezv/go_trafilea_cart/internal/cart"
	"github.com/ibanezv/go_trafilea_cart/internal/order"
	"github.com/ibanezv/go_trafilea_cart/internal/product"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

type CartApp struct {
	serveMux       *mux.Router
	port           string
	cartService    cart.Carts
	orderService   order.Orders
	productService product.Products
}

func NewServer(port int) CartApp {
	router := *mux.NewRouter()
	db := repository.NewRepository()
	cartService := cart.NewCartService(db)
	orderService := order.NewOrdersService(db)
	producService := product.NewProductsService(db)
	app := CartApp{serveMux: &router, port: fmt.Sprintf(":%d", port), cartService: &cartService, orderService: &orderService, productService: &producService}
	app.routeMapper()
	return app
}

func (app *CartApp) Run() error {
	return http.ListenAndServe(app.port, app.serveMux)
}

func (app *CartApp) routeMapper() {
	app.serveMux.HandleFunc("/api/v1/cart/{cartId}", GetCart(app.cartService)).Methods(http.MethodGet)
	app.serveMux.HandleFunc("/api/v1/cart", PostCart(app.cartService)).Methods(http.MethodPost)
	app.serveMux.HandleFunc("/api/v1/cart/{cartId}/product/{productId}", PutProductCart(app.cartService)).Methods(http.MethodPut)
	app.serveMux.HandleFunc("/api/v1/cart/{cartId}", PutCart(app.cartService)).Methods(http.MethodPut)
	app.serveMux.HandleFunc("/api/v1/order", PostOrder(app.orderService)).Methods(http.MethodPost)
	app.serveMux.HandleFunc("/api/v1/order/{cartId}", GetOrder(app.orderService)).Methods(http.MethodGet)
}
