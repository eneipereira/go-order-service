package model

import (
	"errors"
	"testing"
)

func TestNewCustomerNoPassword(t *testing.T) {
	testCases := []struct {
		name         string
		customerName string
		email        string
		phone        string
		expectedErr  error
	}{
		{"deve criar cliente com sucesso", "John Doe", "john.doe@example.com", "11999999999", nil},
		{"deve retornar erro de nome inválido", "J", "john.doe@example.com", "11999999999", ErrCustomerNameTooShort},
		{"deve retornar erro de email inválido", "John Doe", "invalid-email", "11999999999", ErrCustomerEmailInvalid},
		{"deve retornar erro de telefone inválido", "John Doe", "john.doe@example.com", "123", ErrCustomerPhoneInvalid},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customer, err := NewCustomerNoPassword(tc.customerName, tc.email, tc.phone)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}

			if tc.expectedErr == nil && customer == nil {
				t.Error("cliente não deveria ser nulo")
			}
		})
	}
}

func TestValidateCustomerName(t *testing.T) {
	testCases := []struct {
		name         string
		customerName string
		expectedErr  error
	}{
		{"nome válido", "Cliente Válido", nil},
		{"nome vazio", "", ErrCustomerNameRequired},
		{"nome muito curto", "a", ErrCustomerNameTooShort},
		{"nome muito longo", string(make([]byte, 256)), ErrCustomerNameTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomerName(tc.customerName)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateCustomerEmail(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		expectedErr error
	}{
		{"email válido", "test@example.com", nil},
		{"email vazio", "", ErrCustomerEmailRequired},
		{"email inválido", "invalid-email", ErrCustomerEmailInvalid},
		{"email muito longo", string(make([]byte, 256)) + "@example.com", ErrCustomerEmailTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomerEmail(tc.email)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateCustomerPhone(t *testing.T) {
	testCases := []struct {
		name        string
		phone       string
		expectedErr error
	}{
		{"telefone válido", "(11) 99999-9999", nil},
		{"telefone vazio", "", ErrCustomerPhoneRequired},
		{"telefone com poucos dígitos", "12345", ErrCustomerPhoneInvalid},
		{"telefone com muitos dígitos", "1234567890123456", ErrCustomerPhoneInvalid},
		{"telefone com formato inválido", "abc", ErrCustomerPhoneInvalid},
		{"telefone muito longo", string(make([]byte, 31)), ErrCustomerPhoneTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomerPhone(tc.phone)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestValidateCustomerPassword(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{"senha válida", "password123", nil},
		{"senha vazia", "", ErrCustomerPasswordRequired},
		{"senha muito curta", "123", ErrCustomerPasswordTooShort},
		{"senha muito longa", string(make([]byte, 65)), ErrCustomerPasswordTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomerPassword(tc.password)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
			}
		})
	}
}
