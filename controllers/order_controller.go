package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
)

type OrderService interface {
	Create(ctx context.Context, req dto.CreateOrderDTO) (*model.Order, error)
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Order, error)
	Pay(ctx context.Context, orderID string) (*model.Order, error)
	Cancel(ctx context.Context, orderID string) (*model.Order, error)
}

type OrderController struct {
	Service OrderService
}

func NewOrderController(service OrderService) *OrderController {
	return &OrderController{Service: service}
}

func (c *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	order, err := c.Service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInsufficientStock), errors.Is(err, model.ErrProductNotFound), errors.Is(err, model.ErrInvalidCustomerID):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Printf("Error creating order: %v", err)
			http.Error(w, "Could not create order", http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewOrderResponseDTO(*order))
}

func (c *OrderController) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching order by ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
}

func (c *OrderController) FindAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := getQueryParamAsInt(r, "limit", 10)
	offset, _ := getQueryParamAsInt(r, "offset", 0)

	orders, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responseDTOs := make([]dto.OrderResponseDTO, len(orders))
	for i, order := range orders {
		responseDTOs[i] = dto.NewOrderResponseDTO(*order)
	}

	writeJSONResponse(w, http.StatusOK, responseDTOs)
}

func (c *OrderController) Pay(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.Service.Pay(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrInvalidStatusChange) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("Error paying order: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
}

func (c *OrderController) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.Service.Cancel(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrInvalidStatusChange) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("Error canceling order: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
}