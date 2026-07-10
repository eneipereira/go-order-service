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


func (o *Order) Total() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
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