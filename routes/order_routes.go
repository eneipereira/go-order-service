package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterOrderRoutes(router *chi.Mux, controller *controllers.OrderController) {
	router.Route("/orders", func(r chi.Router) {
		r.Post("/", controller.Create)
		r.Get("/", controller.FindAll)
		r.Get("/{id}", controller.FindByID)
		r.Post("/{id}/pay", controller.Pay)
		r.Post("/{id}/cancel", controller.Cancel)
	})
}