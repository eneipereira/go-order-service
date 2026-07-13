package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending  OrderStatus = "PENDING"
	StatusPaid     OrderStatus = "PAID"
	StatusCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

type Order struct {
	ID         uuid.UUID   `json:"id"`
	CustomerID uuid.UUID   `json:"customerId"`
	Status     OrderStatus `json:"status"`
	Items      []OrderItem `json:"items"`
	Total      float64     `json:"total"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

func (o *Order) CalculateSubtotal() float64 {
	var subtotal float64
	for _, item := range o.Items {
		subtotal += item.Price * float64(item.Quantity)
	}
	return subtotal
}

func (o *Order) CalculateDiscount() float64 {
	subtotal := o.CalculateSubtotal()
	if subtotal > 5000.00 {
		return subtotal * 0.10
	}
	return 0
}

func (o *Order) CalculateTotal() {
	o.Total = o.CalculateSubtotal() - o.CalculateDiscount()
}

func (o *Order) Pay() error {
	if o.Status != StatusPending {
		return fmt.Errorf("%w: unable to pay order with status %s", ErrInvalidStatusChange, o.Status)
	}
	o.Status = StatusPaid
	return nil
}

func (o *Order) Cancel() error {
	if o.Status != StatusPending {
		return fmt.Errorf("%w: unable to cancel order with status %s", ErrInvalidStatusChange, o.Status)
	}
	o.Status = StatusCanceled
	return nil
}
