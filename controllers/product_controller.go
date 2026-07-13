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

// Create a new product
// @Summary      Create product
// @Description  Creates a new product.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.CreateProductDTO  true  "Product to create"
// @Success      201      {object}  dto.ProductResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /products [post]
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	savedProduct, err := c.Service.Create(r.Context(), req)
	if err != nil {
		if isValidationErrorProd(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error saving product: %v", err)
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewProductResponseDTO(*savedProduct))
}

// Find all products
// @Summary      List products
// @Description  Get all products with pagination.
// @Tags         Products
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200     {array}   dto.ProductResponseDTO
// @Failure      500     {object}  object{error=string} "Internal Server Error"
// @Router       /products [get]
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

// Find a product by ID
// @Summary      Get product by ID
// @Description  Get a single product by its unique ID.
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  dto.ProductResponseDTO
// @Failure      400  {object}  object{error=string} "Bad Request"
// @Failure      404  {object}  object{error=string} "Not Found"
// @Failure      500  {object}  object{error=string} "Internal Server Error"
// @Router       /products/{id} [get]
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

func isValidationErrorProd(err error) bool {
	return errors.Is(err, model.ErrNotNullViolation) ||
		errors.Is(err, model.ErrProductNameRequired) ||
		errors.Is(err, model.ErrProductNameTooShort) ||
		errors.Is(err, model.ErrProductNameTooLong) ||
		errors.Is(err, model.ErrProductPriceRequired) ||
		errors.Is(err, model.ErrProductPriceTooLow) ||
		errors.Is(err, model.ErrProductStockRequired) ||
		errors.Is(err, model.ErrProductStockTooLow)
}
