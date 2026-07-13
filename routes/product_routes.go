package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterProductRoutes(router *chi.Mux, controller *controllers.ProductController) {
	router.Route("/products", func(r chi.Router) {
		r.Post("/", controller.Create)
		r.Get("/", controller.FindAll)
		r.Get("/{id}", controller.FindByID)
	})
}
