package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ibanezv/go_trafilea_cart/cmd/api"
)

func main() {
	var port = flag.Int("port", 8080, "port")
	flag.Parse()

	server := api.NewServer(*port)
	if err := server.Run(); err != http.ErrServerClosed {
		panic(err)
	}
	log.Println("shutdown: completed")
}
