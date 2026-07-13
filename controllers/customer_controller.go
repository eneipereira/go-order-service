package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/go-chi/chi/v5"
)

type CustomerService interface {
	Create(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error)
	FindByID(ctx context.Context, id string) (*model.Customer, error)
}

type CustomerController struct {
	Service CustomerService
}

func NewCustomerController(service CustomerService) *CustomerController {
	return &CustomerController{Service: service}
}

func (c *CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CustomerDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	savedCustomer, err := c.Service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrEmailAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if isValidationErrorCust(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Error saving customer: %v", err)
		http.Error(w, "Could not create customer", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewCustomerResponseDTO(*savedCustomer))
}

func (c *CustomerController) FindAll(w http.ResponseWriter, r *http.Request) {
	limit, err := getQueryParamAsInt(r, "limit", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := getQueryParamAsInt(r, "offset", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customers, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Error listing customers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	customerResponses := make([]dto.CustomerResponseDTO, len(customers))
	for i, customer := range customers {
		customerResponses[i] = dto.NewCustomerResponseDTO(*customer)
	}

	writeJSONResponse(w, http.StatusOK, customerResponses)
}

func (c *CustomerController) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrCustomerNotFound) {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching customer by ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, dto.NewCustomerResponseDTO(*customer))
}

func isValidationErrorCust(err error) bool {
	return errors.Is(err, model.ErrNotNullViolation) ||
		errors.Is(err, model.ErrCustomerNameRequired) ||
		errors.Is(err, model.ErrCustomerNameTooShort) ||
		errors.Is(err, model.ErrCustomerNameTooLong) ||
		errors.Is(err, model.ErrCustomerEmailRequired) ||
		errors.Is(err, model.ErrCustomerEmailInvalid) ||
		errors.Is(err, model.ErrCustomerEmailTooLong) ||
		errors.Is(err, model.ErrCustomerPhoneRequired) ||
		errors.Is(err, model.ErrCustomerPhoneInvalid) ||
		errors.Is(err, model.ErrCustomerPhoneTooLong) ||
		errors.Is(err, model.ErrCustomerPasswordRequired) ||
		errors.Is(err, model.ErrCustomerPasswordTooShort) ||
		errors.Is(err, model.ErrCustomerPasswordTooLong)
}

func getQueryParamAsInt(r *http.Request, paramName string, defaultValue int) (int, error) {
	paramStr := r.URL.Query().Get(paramName)
	if paramStr == "" {
		return defaultValue, nil
	}
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, fmt.Errorf("parâmetro inválido '%s': deve ser um número inteiro", paramName)
	}
	return paramInt, nil
}

func getIDFromRequest(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	return id, nil
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
