package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
)

type mockCustomerRepo struct {
	saveFunc     func(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Customer, error)
}

func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockCustomerRepo) Save(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, customer)
	}
	return customer, nil
}
func (m *mockCustomerRepo) FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	return nil, nil
}

type mockProductRepo struct {
	createFunc   func(ctx context.Context, product *model.Product) (*model.Product, error)
	findByIDFunc func(ctx context.Context, id string) (*model.Product, error)
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*model.Product, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockProductRepo) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	return product, nil
}
func (m *mockProductRepo) FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return nil, nil
}

type mockOrderRepo struct {
	createFunc                     func(ctx context.Context, order *model.Order) error
	findByIDFunc                   func(ctx context.Context, id string) (*model.Order, error)
	findAllFunc                    func(ctx context.Context, limit, offset int) ([]*model.Order, error)
	updateStatusFunc               func(ctx context.Context, id string, status model.OrderStatus) error
	cancelOrderAndRestockItemsFunc func(ctx context.Context, order *model.Order) error
}

func (m *mockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	return m.createFunc(ctx, order)
}
func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*model.Order, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	return m.updateStatusFunc(ctx, id, status)
}
func (m *mockOrderRepo) CancelOrderAndRestockItems(ctx context.Context, order *model.Order) error {
	return m.cancelOrderAndRestockItemsFunc(ctx, order)
}
func (m *mockOrderRepo) FindAll(ctx context.Context, limit, offset int) ([]*model.Order, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, limit, offset)
	}
	return nil, nil
}

func TestOrderService_Create(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price := 10.0

	testCases := []struct {
		name           string
		req            dto.CreateOrderDTO
		customerRepo   *mockCustomerRepo
		productRepo    *mockProductRepo
		orderRepo      *mockOrderRepo
		expectedErr    error
		expectOrder    bool
		expectedStatus model.OrderStatus
	}{
		{
			name: "deve criar um pedido com sucesso",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 2},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			productRepo: &mockProductRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return &model.Product{ID: productID, Price: &price}, nil
				},
			},
			orderRepo: &mockOrderRepo{
				createFunc: func(ctx context.Context, order *model.Order) error {
					return nil
				},
			},
			expectedErr:    nil,
			expectOrder:    true,
			expectedStatus: model.StatusPending,
		},
		{
			name: "deve retornar erro se o pedido estiver vazio",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items:      []dto.CreateOrderItemDTO{},
			},
			customerRepo: &mockCustomerRepo{},
			productRepo:  &mockProductRepo{},
			orderRepo:    &mockOrderRepo{},
			expectedErr:  model.ErrEmptyOrder,
			expectOrder:  false,
		},
		{
			name: "deve retornar erro se o cliente não existir",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 1},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return nil, model.ErrCustomerNotFound
				},
			},
			productRepo: &mockProductRepo{},
			orderRepo:   &mockOrderRepo{},
			expectedErr: model.ErrInvalidCustomerID,
			expectOrder: false,
		},
		{
			name: "deve retornar erro se o repositório de cliente falhar",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 1},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return nil, errors.New("falha no banco")
				},
			},
			productRepo: &mockProductRepo{},
			orderRepo:   &mockOrderRepo{},
			expectedErr: errors.New("falha no banco"),
			expectOrder: false,
		},
		{
			name: "deve retornar erro se o produto não existir",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 1},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			productRepo: &mockProductRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return nil, model.ErrProductNotFound
				},
			},
			orderRepo:   &mockOrderRepo{},
			expectedErr: model.ErrProductNotFound,
			expectOrder: false,
		},
		{
			name: "deve retornar erro se o repositório de produto falhar",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 1},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			productRepo: &mockProductRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return nil, errors.New("falha no banco")
				},
			},
			orderRepo:   &mockOrderRepo{},
			expectedErr: errors.New("produto " + productID.String() + ": falha no banco"),
			expectOrder: false,
		},
		{
			name: "deve retornar erro se a criação do pedido falhar",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 2},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			productRepo: &mockProductRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Product, error) {
					return &model.Product{ID: productID, Price: &price}, nil
				},
			},
			orderRepo: &mockOrderRepo{
				createFunc: func(ctx context.Context, order *model.Order) error {
					return errors.New("falha ao salvar")
				},
			},
			expectedErr: errors.New("erro ao salvar o pedido: falha ao salvar"),
			expectOrder: false,
		},
		{
			name: "deve retornar erro se a quantidade for inválida",
			req: dto.CreateOrderDTO{
				CustomerID: customerID,
				Items: []dto.CreateOrderItemDTO{
					{ProductID: productID, Quantity: 0},
				},
			},
			customerRepo: &mockCustomerRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Customer, error) {
					return &model.Customer{ID: customerID}, nil
				},
			},
			productRepo: &mockProductRepo{},
			orderRepo:   &mockOrderRepo{},
			expectedErr: model.ErrInvalidQuantity,
			expectOrder: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewOrderService(tc.orderRepo, tc.productRepo, tc.customerRepo)
			order, err := service.Create(context.Background(), tc.req)

			if tc.expectedErr != nil {
				if err == nil || (!errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error()) {
					t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
				}
			} else if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}

			if tc.expectOrder && order == nil {
				t.Error("pedido não deveria ser nulo")
			}

			if tc.expectOrder && order != nil && order.Status != tc.expectedStatus {
				t.Errorf("status = %s; esperado = %s", order.Status, tc.expectedStatus)
			}

			if !tc.expectOrder && order != nil {
				t.Error("pedido deveria ser nulo")
			}
		})
	}
}

func TestOrderService_Pay(t *testing.T) {
	orderID := uuid.New()

	testCases := []struct {
		name        string
		orderRepo   *mockOrderRepo
		expectedErr error
	}{
		{
			name: "deve pagar um pedido com sucesso",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPending}, nil
				},
				updateStatusFunc: func(ctx context.Context, id string, status model.OrderStatus) error {
					return nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "deve retornar erro se o pedido não for encontrado",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return nil, model.ErrOrderNotFound
				},
			},
			expectedErr: model.ErrOrderNotFound,
		},
		{
			name: "deve retornar erro se o status do pedido for inválido",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPaid}, nil
				},
			},
			expectedErr: model.ErrInvalidStatusChange,
		},
		{
			name: "deve retornar erro se a atualização do status falhar",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPending}, nil
				},
				updateStatusFunc: func(ctx context.Context, id string, status model.OrderStatus) error {
					return errors.New("falha no banco")
				},
			},
			expectedErr: errors.New("falha no banco"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewOrderService(tc.orderRepo, nil, nil)
			order, err := service.Pay(context.Background(), orderID.String())

			if tc.expectedErr != nil {
				if err == nil || (!errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error()) {
					t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("erro inesperado: %v", err)
				}
				if order.Status != model.StatusPaid {
					t.Errorf("status = %s; esperado = %s", order.Status, model.StatusPaid)
				}
			}
		})
	}
}

func TestOrderService_Cancel(t *testing.T) {
	orderID := uuid.New()

	testCases := []struct {
		name        string
		orderRepo   *mockOrderRepo
		expectedErr error
	}{
		{
			name: "deve cancelar um pedido com sucesso",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				},
				cancelOrderAndRestockItemsFunc: func(ctx context.Context, order *model.Order) error { return nil },
			},
			expectedErr: nil,
		},
		{
			name: "deve retornar erro se o pedido não for encontrado",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) { return nil, model.ErrOrderNotFound },
			},
			expectedErr: model.ErrOrderNotFound,
		},
		{
			name: "deve retornar erro se o status do pedido for inválido",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPaid}, nil
				},
			},
			expectedErr: model.ErrInvalidStatusChange,
		},
		{
			name: "deve retornar erro se o cancelamento no repositório falhar",
			orderRepo: &mockOrderRepo{
				findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
					return &model.Order{ID: orderID, Status: model.StatusPending}, nil
				},
				cancelOrderAndRestockItemsFunc: func(ctx context.Context, order *model.Order) error { return errors.New("falha no banco") },
			},
			expectedErr: errors.New("falha no banco"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewOrderService(tc.orderRepo, nil, nil)
			order, err := service.Cancel(context.Background(), orderID.String())

			if tc.expectedErr != nil {
				if err == nil || (!errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error()) {
					t.Errorf("erro = %v; esperado = %v", err, tc.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("erro inesperado: %v", err)
				}
				if order.Status != model.StatusCanceled {
					t.Errorf("status = %s; esperado = %s", order.Status, model.StatusCanceled)
				}
			}
		})
	}
}

func TestOrderService_FindByID(t *testing.T) {
	orderID := uuid.New()
	service := NewOrderService(&mockOrderRepo{
		findByIDFunc: func(ctx context.Context, id string) (*model.Order, error) {
			if id == orderID.String() {
				return &model.Order{ID: orderID}, nil
			}
			return nil, model.ErrOrderNotFound
		},
	}, nil, nil)

	_, err := service.FindByID(context.Background(), orderID.String())
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestOrderService_FindAll(t *testing.T) {
	service := NewOrderService(&mockOrderRepo{
		findAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Order, error) {
			return []*model.Order{}, nil
		},
	}, nil, nil)

	_, err := service.FindAll(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}
