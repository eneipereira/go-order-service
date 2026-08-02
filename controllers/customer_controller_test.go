package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mockCustomerService struct {
	createFunc   func(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Customer, error)
	findAllFunc  func(ctx context.Context, limit, offset int) ([]*model.Customer, error)
}

func (m *mockCustomerService) Create(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error) {
	return m.createFunc(ctx, req)
}

func (m *mockCustomerService) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockCustomerService) FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	return m.findAllFunc(ctx, limit, offset)
}

func TestCustomerController_Create(t *testing.T) {
	testCases := []struct {
		name           string
		service        *mockCustomerService
		body           interface{}
		expectedStatus int
	}{
		{
			name: "deve criar cliente com sucesso",
			service: &mockCustomerService{
				createFunc: func(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error) {
					return &model.Customer{}, nil
				},
			},
			body: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "deve retornar erro com body inválido",
			service:        &mockCustomerService{},
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "deve retornar erro de conflito do serviço",
			service: &mockCustomerService{
				createFunc: func(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error) {
					return nil, model.ErrEmailAlreadyExists
				},
			},
			body: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCustomerController(tc.service)
			handler := ErrorMiddleware(controller.Create)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("status = %d; esperado = %d", rr.Code, tc.expectedStatus)
			}
		})
	}
}

func TestCustomerController_FindByID(t *testing.T) {
	customerID := uuid.New()

	testCases := []struct {
		name           string
		service        *mockCustomerService
		customerID     string
		expectedStatus int
	}{
		{
			name: "deve encontrar cliente por ID",
			service: &mockCustomerService{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			customerID:     customerID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "deve retornar erro com ID inválido",
			service:        &mockCustomerService{},
			customerID:     "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "deve retornar not found do serviço",
			service: &mockCustomerService{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return nil, model.ErrCustomerNotFound
				},
			},
			customerID:     uuid.NewString(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCustomerController(tc.service)
			handler := ErrorMiddleware(controller.FindByID)

			req := httptest.NewRequest(http.MethodGet, "/customers/"+tc.customerID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
			chi.RouteContext(req.Context()).URLParams.Add("id", tc.customerID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("status = %d; esperado = %d", rr.Code, tc.expectedStatus)
			}
		})
	}
}

func TestCustomerController_FindAll(t *testing.T) {
	controller := NewCustomerController(&mockCustomerService{
		findAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
			return []*model.Customer{}, nil
		},
	})
	handler := ErrorMiddleware(controller.FindAll)

	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d; esperado = %d", rr.Code, http.StatusOK)
	}
}
