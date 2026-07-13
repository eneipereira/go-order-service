package main

import (
	"context"
	"log"
	"net/http"

	"github.com/eneipereira/go-order-service/config"
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/eneipereira/go-order-service/database"
	"github.com/eneipereira/go-order-service/repository"
	"github.com/eneipereira/go-order-service/routes"
	"github.com/eneipereira/go-order-service/service"
)

// @title           Order Service API
// @version         1.0
// @description     API for managing customers, products, and orders.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func main() {
	ctx := context.Background()

	cfg := config.Load()

	dbpool, err := database.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v", err)
	}
	defer dbpool.Close()

	customerRepo := repository.NewPostgresCustomerRepository(dbpool)
	productRepo := repository.NewPostgresProductRepository(dbpool)
	orderRepo := repository.NewPostgresOrderRepository(dbpool)

	customerService := service.NewCustomerService(customerRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, customerRepo)

	customerController := controllers.NewCustomerController(customerService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	router := routes.SetupRouter(customerController, productController, orderController)

	log.Printf("Server is running on port %s", cfg.ServerPort)

	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Fatalf("Unable to start the server: %s\n", err)
	}
}
