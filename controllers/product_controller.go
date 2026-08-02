package controllers

import (
	"context"
	"encoding/json"
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
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var req dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrInvalidJSON
	}

	savedProduct, err := c.Service.Create(r.Context(), req)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusCreated, dto.NewProductResponseDTO(*savedProduct))
	return nil
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
func (c *ProductController) FindAll(w http.ResponseWriter, r *http.Request) error {
	limit, err := getQueryParamAsInt(r, "limit", 10)
	if err != nil {
		return err
	}

	offset, err := getQueryParamAsInt(r, "offset", 0)
	if err != nil {
		return err
	}

	products, err := c.Service.FindAll(r.Context(), limit, offset)
	if err != nil {
		return err
	}

	productResponses := make([]dto.ProductResponseDTO, len(products))
	for i, p := range products {
		productResponses[i] = dto.NewProductResponseDTO(*p)
	}

	writeJSONResponse(w, http.StatusOK, productResponses)
	return nil
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
func (c *ProductController) FindByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	product, err := c.Service.FindByID(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, dto.NewProductResponseDTO(*product))
	return nil
}
