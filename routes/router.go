package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(controller *controllers.CustomerController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", controllers.HealthCheck)

	RegisterCustomerRoutes(r, controller)

	return r
}
