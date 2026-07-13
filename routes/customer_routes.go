package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterCustomerRoutes(router *chi.Mux, controller *controllers.CustomerController) {
	router.Route("/customers", func(r chi.Router) {
		r.Post("/", controller.Create)
		r.Get("/", controller.FindAll)
		r.Get("/{id}", controller.FindByID)
	})
}
