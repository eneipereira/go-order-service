package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterProductRoutes(router *chi.Mux, controller *controllers.ProductController) {
	router.Route("/products", func(r chi.Router) {
		r.Post("/", controllers.ErrorMiddleware(controller.Create))
		r.Get("/", controllers.ErrorMiddleware(controller.FindAll))
		r.Get("/{id}", controllers.ErrorMiddleware(controller.FindByID))
	})
}
