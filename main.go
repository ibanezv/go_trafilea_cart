package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ibanezv/go_trafilea_cart/cmd/api"
)

//   Product Api:
//    version: 0.1
//    title: Product Api
//   Schemes: http, https
//   Host:
//   BasePath: /api/v1
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   SecurityDefinitions:
//    Bearer:
//     type: apiKey
//     name: Authorization
//     in: header
//   swagger:meta
func main() {
	var port = flag.Int("port", 8080, "port")
	flag.Parse()

	errs := make(chan error, 2)
	go func() {
		server := api.NewServer(*port)
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("exit", <-errs)
}
