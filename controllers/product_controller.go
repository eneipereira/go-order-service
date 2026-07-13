package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
)

type ProductService interface {
	Create(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductController struct {
	Service ProductService
}

func NewProductController(service ProductService) *ProductController {
	return &ProductController{Service: service}
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	savedProduct, err := c.Service.Create(r.Context(), req)
	if err != nil {
		log.Printf("Error saving product: %v", err)
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewProductResponseDTO(*savedProduct))
}

func (c *ProductController) FindAll(w http.ResponseWriter, r *http.Request) {
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

	products, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Error listing products: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	productResponses := make([]dto.ProductResponseDTO, len(products))
	for i, p := range products {
		productResponses[i] = dto.NewProductResponseDTO(*p)
	}

	writeJSONResponse(w, http.StatusOK, productResponses)
}

func (c *ProductController) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching product by ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, dto.NewProductResponseDTO(*product))
}
