package main

import (
	"context"
	"log"
	"net/http"

	"github.com/eneipereira/go-order-service/config"
	"github.com/eneipereira/go-order-service/database"
	"github.com/eneipereira/go-order-service/routes"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	dbpool, err := database.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v", err)
	}
	defer dbpool.Close()

	router := routes.SetupRouter()

	log.Printf("Servidor iniciando na porta %s", cfg.ServerPort)

	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Fatalf("Unable to start the server: %s\n", err)
	}
}