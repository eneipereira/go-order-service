package controllers

import (
	"context"
	"encoding/json"
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

// Create a new order
// @Summary      Create order
// @Description  Creates a new order.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order    body      dto.CreateOrderDTO  true  "Order to create"
// @Success      201      {object}  dto.OrderResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      409  {object}  object{error=string} "Conflict"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /orders [post]
func (c *OrderController) Create(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var req dto.CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrInvalidJSON
	}

	order, err := c.Service.Create(r.Context(), req)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewOrderResponseDTO(*order))
	return nil
}

// Find an order by ID
// @Summary      Get order by ID
// @Description  Get a single order by its unique ID.
// @Tags         Orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  dto.OrderResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      404  {object}  object{error=string} "Not Found"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /orders/{id} [get]
func (c *OrderController) FindByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	order, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
	return nil
}

// Find all orders
// @Summary      List orders
// @Description  Get all orders with pagination.
// @Tags         Orders
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200     {array}   dto.OrderResponseDTO
// @Failure      500     {object}  object{error=string} "Internal Server Error"
// @Router       /orders [get]
func (c *OrderController) FindAll(w http.ResponseWriter, r *http.Request) error {
	limit, err := getQueryParamAsInt(r, "limit", 10)
	if err != nil {
		return err
	}
	offset, err := getQueryParamAsInt(r, "offset", 0)
	if err != nil {
		return err
	}

	orders, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		return err
	}

	responseDTOs := make([]dto.OrderResponseDTO, len(orders))
	for i, order := range orders {
		responseDTOs[i] = dto.NewOrderResponseDTO(*order)
	}

	writeJSONResponse(w, http.StatusOK, responseDTOs)
	return nil
}

// Pay for an order
// @Summary      Pay for an order
// @Description  Marks a PENDING order as PAID.
// @Tags         Orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  dto.OrderResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      404  {object}  object{error=string} "Not Found"
// @Failure      409  {object}  object{error=string} "Conflict"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /orders/{id}/pay [post]
func (c *OrderController) Pay(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	order, err := c.Service.Pay(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
	return nil
}

// Cancel an order
// @Summary      Cancel an order
// @Description  Marks a PENDING order as CANCELED and restocks items.
// @Tags         Orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  dto.OrderResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      404  {object}  object{error=string} "Not Found"
// @Failure      409  {object}  object{error=string} "Conflict"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /orders/{id}/cancel [post]
func (c *OrderController) Cancel(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	order, err := c.Service.Cancel(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, dto.NewOrderResponseDTO(*order))
	return nil
}
