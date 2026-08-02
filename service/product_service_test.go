package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
)

type mockProductRepoForProductService struct {
	createFunc   func(ctx context.Context, product *model.Product) (*model.Product, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Product, error)
	findAllFunc  func(ctx context.Context, limit, offset int) ([]*model.Product, error)
}

func (m *mockProductRepoForProductService) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	return m.createFunc(ctx, product)
}

func (m *mockProductRepoForProductService) FindByID(ctx context.Context, id string) (*model.Product, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockProductRepoForProductService) FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return m.findAllFunc(ctx, limit, offset)
}

func TestProductService_Create(t *testing.T) {
	price := 19.99
	stock := 100

	testCases := []struct {
		name        string
		req         dto.CreateProductDTO
		repo        *mockProductRepoForProductService
		expectedErr error
	}{
		{
			name: "deve criar um produto com sucesso",
			req: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: &price,
				Stock: &stock,
			},
			repo: &mockProductRepoForProductService{
				createFunc: func(ctx context.Context, product *model.Product) (*model.Product, error) {
					return product, nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "deve retornar erro se o preço for nulo",
			req: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: nil,
				Stock: &stock,
			},
			repo:        &mockProductRepoForProductService{},
			expectedErr: model.ErrProductPriceRequired,
		},
		{
			name: "deve retornar erro se o estoque for nulo",
			req: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: &price,
				Stock: nil,
			},
			repo:        &mockProductRepoForProductService{},
			expectedErr: model.ErrProductStockRequired,
		},
		{
			name: "deve retornar erro de validação do nome",
			req: dto.CreateProductDTO{
				Name:  "a",
				Price: &price,
				Stock: &stock,
			},
			repo:        &mockProductRepoForProductService{},
			expectedErr: model.ErrProductNameTooShort,
		},
		{
			name: "deve retornar erro do repositório",
			req: dto.CreateProductDTO{
				Name:  "Test Product",
				Price: &price,
				Stock: &stock,
			},
			repo: &mockProductRepoForProductService{
				createFunc: func(ctx context.Context, product *model.Product) (*model.Product, error) {
					return nil, errors.New("falha no banco")
				},
			},
			expectedErr: errors.New("falha no banco"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewProductService(tc.repo)
			product, err := service.Create(context.Background(), tc.req)

			if tc.expectedErr != nil {
				if err == nil || (!errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error()) {
					t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
				}
			} else if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}

			if tc.expectedErr == nil && product == nil {
				t.Error("produto não deveria ser nulo")
			}
		})
	}
}

func TestProductService_FindByID(t *testing.T) {
	productID := uuid.New()
	repo := &mockProductRepoForProductService{
		findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
			if id == productID.String() {
				return &model.Product{ID: productID}, nil
			}
			return nil, model.ErrProductNotFound
		},
	}
	service := NewProductService(repo)

	t.Run("deve encontrar produto por ID", func(t *testing.T) {
		product, err := service.FindByID(context.Background(), productID.String())
		if err != nil {
			t.Errorf("erro inesperado: %v", err)
		}
		if product == nil {
			t.Error("produto não deveria ser nulo")
		}
	})

	t.Run("deve retornar erro se produto não for encontrado", func(t *testing.T) {
		_, err := service.FindByID(context.Background(), uuid.New().String())
		if !errors.Is(err, model.ErrProductNotFound) {
			t.Errorf("erro = %v; esperado = %v", err, model.ErrProductNotFound)
		}
	})
}

func TestProductService_FindAll(t *testing.T) {
	repo := &mockProductRepoForProductService{
		findAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Product, error) {
			return []*model.Product{}, nil
		},
	}
	service := NewProductService(repo)

	_, err := service.FindAll(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}
