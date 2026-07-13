package model

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	productNameMinLength = 3
	productNameMaxLength = 255
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     *float64  `json:"price"`
	Stock     *int      `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewProduct(name string, price float64, stock int) (*Product, error) {
	name = strings.TrimSpace(name)

	if err := ValidateProductName(name); err != nil {
		return nil, err
	}
	if err := ValidateProductPrice(&price); err != nil {
		return nil, err
	}
	if err := ValidateProductStock(&stock); err != nil {
		return nil, err
	}

	return &Product{
		Name:  name,
		Price: &price,
		Stock: &stock,
	}, nil
}

func (p *Product) ReduceStock(quantity int) error {
	if *p.Stock < quantity {
		return fmt.Errorf("%w: requested %d, but only %d in stock for product %s", ErrInsufficientStock, quantity, *p.Stock, p.ID)
	}
	*p.Stock -= quantity
	return nil
}

func (p *Product) AddStock(quantity int) {
	*p.Stock += quantity
}

func ValidateProductName(name string) error {
	switch length := utf8.RuneCountInString(name); {
	case length == 0:
		return ErrProductNameRequired
	case length < productNameMinLength:
		return ErrProductNameTooShort
	case length > productNameMaxLength:
		return ErrProductNameTooLong
	default:
		return nil
	}
}

func ValidateProductPrice(price *float64) error {
	if *price <= 0 {
		return ErrProductPriceTooLow
	}
	return nil
}

func ValidateProductStock(stock *int) error {
	if *stock <= 0 {
		return ErrProductStockTooLow
	}
	return nil
}
