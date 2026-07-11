package service

import (
	"errors"
	"fmt"
	"github.com/eneipereira/go-order-service/model"
	"github.com/eneipereira/go-order-service/repository"
)


type CreateOrderItemRequest struct {
	ProductID string
	Quantity  int
}


type CreateOrderRequest struct {
	Customer string
	Items    []CreateOrderItemRequest
}


type OrderService struct {
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	nextOrderID int
}


func NewOrderService(productRepo repository.ProductRepository, orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		productRepo: productRepo,
		orderRepo:   orderRepo,
		nextOrderID: 1, 
	}
}


func (s *OrderService) generateID() string {
	id := fmt.Sprintf("PED-%03d", s.nextOrderID)
	s.nextOrderID++
	return id
}


func (s *OrderService) CreateOrder(req CreateOrderRequest) (*model.Order, error) {
	var validationErrors []error

	if req.Customer == "" {
		validationErrors = append(validationErrors, model.ErrInvalidCustomer)
	}
	if len(req.Items) == 0 {
		validationErrors = append(validationErrors, model.ErrEmptyOrder)
	}


	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	var orderItems []*model.OrderItem

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", itemReq.ProductID, model.ErrInvalidQuantity))
			continue
		}

		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", itemReq.ProductID, model.ErrProductNotFound))
			continue
		}

		if product.Stock < itemReq.Quantity {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", product.ID, model.ErrInsufficientStock))
		}

		orderItems = append(orderItems, &model.OrderItem{
			Product:  product,
			Quantity: itemReq.Quantity,
			Price:    product.Price,
		})
	}

	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	order := &model.Order{
		ID:       s.generateID(),
		Customer: req.Customer,
		Items:    orderItems,
		Status:   model.StatusPending,
	}

	for _, item := range order.Items {
		item.Product.ReduceStock(item.Quantity)
		s.productRepo.Save(item.Product)
	}

	if err := s.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("Error saving order: %w", err)
	}

	return order, nil
}


func (s *OrderService) PayOrder(orderID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Pay(); err != nil {
		return nil, err
	}

	s.orderRepo.Save(order)
	return order, nil
}


func (s *OrderService) CancelOrder(orderID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	
	for _, item := range order.Items {
		item.Product.AddStock(item.Quantity)
		s.productRepo.Save(item.Product)
	}

	s.orderRepo.Save(order)
	return order, nil
}


func (s *OrderService) FindOrderByID(orderID string) (*model.Order, error) {
	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) ListOrders(filters ...repository.OrderFilter) ([]*model.Order, error) {
	return s.orderRepo.List(filters...)
}