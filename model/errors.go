package model

import "errors"

var (
	ErrProductNotFound     = errors.New("Product Not Found")
	ErrOrderNotFound       = errors.New("Order Not Found")
	ErrInvalidQuantity     = errors.New("Invalid Quantity")
	ErrInsufficientStock   = errors.New("Insufficient Stock")
	ErrInvalidCustomer     = errors.New("Invalid Customer")
	ErrEmptyOrder          = errors.New("Order must have at least one item")
	ErrInvalidStatusChange = errors.New("Invalid Status Change")
)