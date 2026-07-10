package main

import "fmt"

type OrderStatus string

const (
	StatusPending  OrderStatus = "PENDING"
	StatusPaid     OrderStatus = "PAID"
	StatusCanceled OrderStatus = "CANCELED"
)


type OrderItem struct {
	Product  *Product
	Quantity int
	Price    float64 
}


type Order struct {
	ID       string
	Customer string
	Items    []*OrderItem
	Status   OrderStatus
}

func (o *Order) Subtotal() float64 {
	var subtotal float64
	for _, item := range o.Items {
		subtotal += item.Price * float64(item.Quantity)
	}
	return subtotal
}

func (o *Order) Discount() float64 {
	subtotal := o.Subtotal()
	if subtotal > 5000.00 {
		return subtotal * 0.10
	}
	return 0
}

func (o *Order) Total() float64 {
	return o.Subtotal() - o.Discount()
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