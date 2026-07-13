package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
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

// Create a new customer
// @Summary      Create customer
// @Description  Creates a new customer.
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        customer  body      dto.CustomerDTO  true  "Customer to create"
// @Success      201       {object}  dto.CustomerResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      409  {object}  object{error=string} "Conflict"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /customers [post]
func (c *CustomerController) Create(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var req dto.CustomerDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	savedCustomer, err := c.Service.Create(r.Context(), req)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewCustomerResponseDTO(*savedCustomer))
	return nil
}

// Find all customers
// @Summary      List customers
// @Description  Get all customers with pagination.
// @Tags         Customers
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200     {array}   dto.CustomerResponseDTO
// @Failure      500     {object}  object{error=string} "Internal Server Error"
// @Router       /customers [get]
func (c *CustomerController) FindAll(w http.ResponseWriter, r *http.Request) error {
	limit, err := getQueryParamAsInt(r, "limit", 10)
	if err != nil {
		return err
	}

	offset, err := getQueryParamAsInt(r, "offset", 0)
	if err != nil {
		return err
	}

	customers, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		return err
	}

	customerResponses := make([]dto.CustomerResponseDTO, len(customers))
	for i, customer := range customers {
		customerResponses[i] = dto.NewCustomerResponseDTO(*customer)
	}

	writeJSONResponse(w, http.StatusOK, customerResponses)
	return nil
}

// Find a customer by ID
// @Summary      Get customer by ID
// @Description  Get a single customer by its unique ID.
// @Tags         Customers
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Success      200  {object}  dto.CustomerResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      404  {object}  object{error=string} "Not Found"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /customers/{id} [get]
func (c *CustomerController) FindByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	customer, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, dto.NewCustomerResponseDTO(*customer))
	return nil
}
