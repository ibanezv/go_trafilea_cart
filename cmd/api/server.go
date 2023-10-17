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
	/* idleChan := make(chan struct{})

	go func(app *CartApp, idleChan chan struct{}) {
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		log.Println("shutdown:", sig)
		app.Stop(5 * time.Second)
		// Actual shutdown trigger.
		close(idleChan)
	}(app, idleChan) */

	return http.ListenAndServe(app.port, app.serveMux)
	//return app.server.ListenAndServe()
}

/* func (app *CartApp) Stop(t time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	err := app.server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
} */

func (app *CartApp) routeMapper() {
	app.serveMux.HandleFunc("/cart/{cartId}", GetCart(app.cartService)).Methods(http.MethodGet)
	app.serveMux.HandleFunc("/cart", PostCart(app.cartService)).Methods(http.MethodPost)
	app.serveMux.HandleFunc("/cart/{cartId}/product/{productId}", PutProductCart(app.cartService)).Methods(http.MethodPut)
	app.serveMux.HandleFunc("/cart/{cartId}", PutCart(app.cartService)).Methods(http.MethodPut)
}
