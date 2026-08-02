package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mockProductService struct {
	createFunc   func(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Product, error)
	findAllFunc  func(ctx context.Context, limit, offset int) ([]*model.Product, error)
}

func (m *mockProductService) Create(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error) {
	return m.createFunc(ctx, req)
}

func (m *mockProductService) FindByID(ctx context.Context, id string) (*model.Product, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockProductService) FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return m.findAllFunc(ctx, limit, offset)
}

func TestProductController_Create(t *testing.T) {
	price := 19.99
	stock := 100

	testCases := []struct {
		name           string
		service        *mockProductService
		body           interface{}
		expectedStatus int
	}{
		{
			name: "deve criar produto com sucesso",
			service: &mockProductService{
				createFunc: func(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error) {
					return &model.Product{}, nil
				},
			},
			body: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: &price,
				Stock: &stock,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "deve retornar erro com body inválido",
			service:        &mockProductService{},
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "deve retornar erro do serviço",
			service: &mockProductService{
				createFunc: func(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error) {
					return nil, errors.New("service error")
				},
			},
			body: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: &price,
				Stock: &stock,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewProductController(tc.service)
			handler := ErrorMiddleware(controller.Create)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("status = %d; esperado = %d", rr.Code, tc.expectedStatus)
			}
		})
	}
}

func TestProductController_FindByID(t *testing.T) {
	productID := uuid.New()

	testCases := []struct {
		name           string
		service        *mockProductService
		productID      string
		expectedStatus int
	}{
		{
			name: "deve encontrar produto por ID",
			service: &mockProductService{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return &model.Product{ID: productID}, nil
				},
			},
			productID:      productID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "deve retornar erro com ID inválido",
			service:        &mockProductService{},
			productID:      "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "deve retornar not found do serviço",
			service: &mockProductService{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return nil, model.ErrProductNotFound
				},
			},
			productID:      uuid.NewString(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewProductController(tc.service)
			handler := ErrorMiddleware(controller.FindByID)

			req := httptest.NewRequest(http.MethodGet, "/products/"+tc.productID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
			chi.RouteContext(req.Context()).URLParams.Add("id", tc.productID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("status = %d; esperado = %d", rr.Code, tc.expectedStatus)
			}
		})
	}
}

func TestProductController_FindAll(t *testing.T) {
	controller := NewProductController(&mockProductService{
		findAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Product, error) {
			return []*model.Product{}, nil
		},
	})
	handler := ErrorMiddleware(controller.FindAll)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d; esperado = %d", rr.Code, http.StatusOK)
	}
}
