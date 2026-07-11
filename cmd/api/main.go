package main

import (
	"log"
	"net/http"

	"github.com/eneipereira/go-order-service/routes"
)

func main() {
	router := routes.SetupRouter()

	port := ":8080"
	log.Printf("Server is running on port %s", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}