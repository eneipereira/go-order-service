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

type ProductRepository interface {
	Save(ctx context.Context, product *model.Product) (*model.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductController struct {
	Repo ProductRepository
}

func NewProductController(repo ProductRepository) *ProductController {
	return &ProductController{Repo: repo}
}

func (c *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product, err := model.NewProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	savedProduct, err := c.Repo.Save(r.Context(), product)
	if err != nil {
		log.Printf("Error saving product: %v", err)
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewProductResponseDTO(*savedProduct))
}

func (c *ProductController) FindAllProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := c.Repo.FindAll(r.Context(), limit, offset)
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

func (c *ProductController) FindProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := c.Repo.FindByID(r.Context(), id)
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
