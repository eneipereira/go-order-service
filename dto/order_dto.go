package dto

import (
	"time"

	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
)

type CreateOrderItemDTO struct {
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
}

type CreateOrderDTO struct {
	CustomerID uuid.UUID            `json:"customerId"`
	Items      []CreateOrderItemDTO `json:"items"`
}

type OrderItemResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

type OrderResponseDTO struct {
	ID         uuid.UUID              `json:"id"`
	CustomerID uuid.UUID              `json:"customerId"`
	Status     model.OrderStatus      `json:"status"`
	Items      []OrderItemResponseDTO `json:"items"`
	Total      float64                `json:"total"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

func NewOrderResponseDTO(order model.Order) OrderResponseDTO {
	items := make([]OrderItemResponseDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponseDTO{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return OrderResponseDTO{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		Items:      items,
		Total:      order.Total,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}
