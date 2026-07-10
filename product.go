package main

import "fmt"

type Product struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return fmt.Errorf("%w: Product with ID %s", ErrInsufficientStock, p.ID)
	}
	p.Stock -= quantity
	return nil
}

func (p *Product) AddStock(quantity int) {
	p.Stock += quantity
}