package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id string) (*model.Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
	CancelOrderAndRestockItems(ctx context.Context, order *model.Order) error
}

type OrderService struct {
	orderRepo    OrderRepository
	productRepo  ProductRepository
	customerRepo CustomerRepository
}

func NewOrderService(orderRepo OrderRepository, productRepo ProductRepository, customerRepo CustomerRepository) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, req dto.CreateOrderDTO) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, model.ErrEmptyOrder
	}

	_, err := s.customerRepo.FindByID(ctx, req.CustomerID.String())
	if err != nil {
		if errors.Is(err, model.ErrCustomerNotFound) {
			return nil, model.ErrInvalidCustomerID
		}
		return nil, err
	}

	order := &model.Order{
		CustomerID: req.CustomerID,
		Status:     model.StatusPending,
		Items:      make([]model.OrderItem, 0, len(req.Items)),
	}

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("produto %s: %w", itemReq.ProductID, model.ErrInvalidQuantity)
		}

		product, err := s.productRepo.FindByID(ctx, itemReq.ProductID.String())
		if err != nil {
			return nil, fmt.Errorf("produto %s: %w", itemReq.ProductID, err)
		}

		order.Items = append(order.Items, model.OrderItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			Price:     *product.Price,
		})
	}

	order.CalculateTotal()

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("erro ao salvar o pedido: %w", err)
	}

	return order, nil
}

func (s *OrderService) Pay(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Pay(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Cancel(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.CancelOrderAndRestockItems(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

func (s *OrderService) FindAll(ctx context.Context, limit, offset int) ([]*model.Order, error) {
	return s.orderRepo.FindAll(ctx, limit, offset)
}
