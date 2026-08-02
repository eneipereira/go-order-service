package model

import (
	"errors"
	"testing"
)

func TestOrder_CalculateSubtotal(t *testing.T) {
	order := Order{
		Items: []OrderItem{
			{Price: 10.50, Quantity: 2},
			{Price: 5.00, Quantity: 3},
		},
	}
	expected := 36.00
	if total := order.CalculateSubtotal(); total != expected {
		t.Errorf("subtotal = %f; esperado = %f", total, expected)
	}
}

func TestOrder_CalculateDiscount(t *testing.T) {
	testCases := []struct {
		name             string
		items            []OrderItem
		expectedDiscount float64
	}{
		{
			"deve aplicar desconto de 10% para compras acima de 5000",
			[]OrderItem{{Price: 6000, Quantity: 1}},
			600.0,
		},
		{
			"não deve aplicar desconto para compras abaixo de 5000",
			[]OrderItem{{Price: 4000, Quantity: 1}},
			0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := Order{Items: tc.items}
			if discount := order.CalculateDiscount(); discount != tc.expectedDiscount {
				t.Errorf("desconto = %f; esperado = %f", discount, tc.expectedDiscount)
			}
		})
	}
}

func TestOrder_CalculateTotal(t *testing.T) {
	order := Order{
		Items: []OrderItem{{Price: 6000, Quantity: 1}},
	}
	order.CalculateTotal()
	expectedTotal := 5400.0
	if order.Total != expectedTotal {
		t.Errorf("total = %f; esperado = %f", order.Total, expectedTotal)
	}
}

func TestOrder_Pay(t *testing.T) {
	t.Run("deve pagar pedido pendente com sucesso", func(t *testing.T) {
		order := &Order{Status: StatusPending}
		err := order.Pay()
		if err != nil {
			t.Errorf("erro inesperado: %v", err)
		}
		if order.Status != StatusPaid {
			t.Errorf("status = %s; esperado = %s", order.Status, StatusPaid)
		}
	})

	t.Run("deve retornar erro ao pagar pedido não pendente", func(t *testing.T) {
		order := &Order{Status: StatusCanceled}
		err := order.Pay()
		if !errors.Is(err, ErrInvalidStatusChange) {
			t.Errorf("erro = %v; esperado = %v", err, ErrInvalidStatusChange)
		}
	})
}

func TestOrder_Cancel(t *testing.T) {
	t.Run("deve cancelar pedido pendente com sucesso", func(t *testing.T) {
		order := &Order{Status: StatusPending}
		err := order.Cancel()
		if err != nil {
			t.Errorf("erro inesperado: %v", err)
		}
		if order.Status != StatusCanceled {
			t.Errorf("status = %s; esperado = %s", order.Status, StatusCanceled)
		}
	})

	t.Run("deve retornar erro ao cancelar pedido não pendente", func(t *testing.T) {
		order := &Order{Status: StatusPaid}
		err := order.Cancel()
		if !errors.Is(err, ErrInvalidStatusChange) {
			t.Errorf("erro = %v; esperado = %v", err, ErrInvalidStatusChange)
		}
	})
}
