package routes

import (
	"github.com/eneipereira/go-order-service/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterOrderRoutes(router *chi.Mux, controller *controllers.OrderController) {
	router.Route("/orders", func(r chi.Router) {
		r.Post("/", controllers.ErrorMiddleware(controller.Create))
		r.Get("/", controllers.ErrorMiddleware(controller.FindAll))
		r.Get("/{id}", controllers.ErrorMiddleware(controller.FindByID))
		r.Post("/{id}/pay", controllers.ErrorMiddleware(controller.Pay))
		r.Post("/{id}/cancel", controllers.ErrorMiddleware(controller.Cancel))
	})
}
