package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
)

type mockCustomerRepoForCustomerService struct {
	saveFunc     func(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Customer, error)
	findAllFunc  func(ctx context.Context, limit, offset int) ([]*model.Customer, error)
}

func (m *mockCustomerRepoForCustomerService) Save(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	return m.saveFunc(ctx, customer)
}

func (m *mockCustomerRepoForCustomerService) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockCustomerRepoForCustomerService) FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	return m.findAllFunc(ctx, limit, offset)
}

func TestCustomerService_Create(t *testing.T) {
	testCases := []struct {
		name        string
		req         dto.CustomerDTO
		repo        *mockCustomerRepoForCustomerService
		expectedErr error
	}{
		{
			name: "deve criar um cliente com sucesso",
			req: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			repo: &mockCustomerRepoForCustomerService{
				saveFunc: func(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
					return customer, nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "deve retornar erro de validação de senha",
			req: dto.CustomerDTO{
				Password: "123",
			},
			repo:        &mockCustomerRepoForCustomerService{},
			expectedErr: model.ErrCustomerPasswordTooShort,
		},
		{
			name: "deve retornar erro de validação de email",
			req: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "invalid-email",
				Phone:    "1234567890",
				Password: "password123",
			},
			repo:        &mockCustomerRepoForCustomerService{},
			expectedErr: model.ErrCustomerEmailInvalid,
		},
		{
			name: "deve retornar erro do repositório",
			req: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			repo: &mockCustomerRepoForCustomerService{
				saveFunc: func(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
					return nil, errors.New("falha no banco")
				},
			},
			expectedErr: errors.New("falha no banco"),
		},
		{
			name: "deve retornar erro se o hash da senha falhar",
			req: dto.CustomerDTO{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			repo:        &mockCustomerRepoForCustomerService{},
			expectedErr: errors.New("falha de hash"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewCustomerService(tc.repo)

			if tc.name == "deve retornar erro se o hash da senha falhar" {
				service.hasher = func(password []byte, cost int) ([]byte, error) {
					return nil, errors.New("falha de hash")
				}
			}

			customer, err := service.Create(context.Background(), tc.req)

			if tc.expectedErr != nil {
				if err == nil || (!errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error()) {
					t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
				}
			} else if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}

			if tc.expectedErr == nil && customer == nil {
				t.Error("cliente não deveria ser nulo")
			}
		})
	}
}

func TestCustomerService_FindByID(t *testing.T) {
	customerID := uuid.New()
	repo := &mockCustomerRepoForCustomerService{
		findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
			if id == customerID.String() {
				return &model.Customer{ID: customerID}, nil
			}
			return nil, model.ErrCustomerNotFound
		},
	}
	service := NewCustomerService(repo)

	t.Run("deve encontrar cliente por ID", func(t *testing.T) {
		customer, err := service.FindByID(context.Background(), customerID.String())
		if err != nil {
			t.Errorf("erro inesperado: %v", err)
		}
		if customer == nil {
			t.Error("cliente não deveria ser nulo")
		}
	})

	t.Run("deve retornar erro se cliente não for encontrado", func(t *testing.T) {
		_, err := service.FindByID(context.Background(), uuid.New().String())
		if !errors.Is(err, model.ErrCustomerNotFound) {
			t.Errorf("erro = %v; esperado = %v", err, model.ErrCustomerNotFound)
		}
	})
}

func TestCustomerService_FindAll(t *testing.T) {
	repo := &mockCustomerRepoForCustomerService{
		findAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
			return []*model.Customer{}, nil
		},
	}
	service := NewCustomerService(repo)

	_, err := service.FindAll(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}
