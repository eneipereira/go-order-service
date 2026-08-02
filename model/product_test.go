package model

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestValidateProductName(t *testing.T) {
	testCases := []struct {
		name        string
		productName string
		expectedErr error
	}{
		{"nome válido", "Produto Válido", nil},
		{"nome vazio", "", ErrProductNameRequired},
		{"nome muito curto", "ab", ErrProductNameTooShort},
		{"nome muito longo", string(make([]byte, 256)), ErrProductNameTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProductName(tc.productName)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateProductPrice(t *testing.T) {
	validPrice := 10.0
	invalidPrice := 0.0

	testCases := []struct {
		name        string
		price       *float64
		expectedErr error
	}{
		{"preço válido", &validPrice, nil},
		{"preço inválido (zero)", &invalidPrice, ErrProductPriceTooLow},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProductPrice(tc.price)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateProductStock(t *testing.T) {
	validStock := 10
	invalidStock := 0
	negativeStock := -1

	testCases := []struct {
		name        string
		stock       *int
		expectedErr error
	}{
		{"estoque válido", &validStock, nil},
		{"estoque inválido (zero)", &invalidStock, ErrProductStockTooLow},
		{"estoque inválido (negativo)", &negativeStock, ErrProductStockTooLow},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProductStock(tc.stock)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestProduct_AddStock(t *testing.T) {
	stock := 10
	p := &Product{Stock: &stock}

	t.Run("deve adicionar ao estoque com sucesso", func(t *testing.T) {
		p.AddStock(5)

		if *p.Stock != 15 {
			t.Errorf("estoque = %d; esperado = 15", *p.Stock)
		}
	})
}

func TestProduct_ReduceStock(t *testing.T) {
	stock := 10
	p := &Product{
		ID:    uuid.New(),
		Stock: &stock,
	}

	t.Run("deve reduzir o estoque com sucesso", func(t *testing.T) {
		err := p.ReduceStock(5)
		if err != nil {
			t.Errorf("erro inesperado: %v", err)
		}
		if *p.Stock != 5 {
			t.Errorf("estoque = %d; esperado = 5", *p.Stock)
		}
	})

	t.Run("deve retornar erro por estoque insuficiente", func(t *testing.T) {
		err := p.ReduceStock(10)
		if !errors.Is(err, ErrInsufficientStock) {
			t.Errorf("erro = %v; esperado = %v", err, ErrInsufficientStock)
		}
		if *p.Stock != 5 {
			t.Errorf("estoque não deveria ter sido alterado, mas está %d", *p.Stock)
		}
	})
}

func TestNewProduct(t *testing.T) {
	testCases := []struct {
		name        string
		productName string
		price       float64
		stock       int
		expectedErr error
	}{
		{"deve criar produto com sucesso", "Produto Válido", 99.99, 10, nil},
		{"deve retornar erro de nome inválido", "a", 99.99, 10, ErrProductNameTooShort},
		{"deve retornar erro de preço inválido", "Produto Válido", 0, 10, ErrProductPriceTooLow},
		{"deve retornar erro de estoque inválido", "Produto Válido", 99.99, 0, ErrProductStockTooLow},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product, err := NewProduct(tc.productName, tc.price, tc.stock)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}

			if tc.expectedErr == nil && product == nil {
				t.Error("produto não deveria ser nulo")
			}
		})
	}
}
