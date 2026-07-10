package main

import (
	"errors"
	"fmt"
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
	productRepo ProductRepository
	orderRepo   OrderRepository
	nextOrderID int
}


func NewOrderService(productRepo ProductRepository, orderRepo OrderRepository) *OrderService {
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


func (s *OrderService) CreateOrder(req CreateOrderRequest) (*Order, error) {
	var validationErrors []error

	if req.Customer == "" {
		validationErrors = append(validationErrors, ErrInvalidCustomer)
	}
	if len(req.Items) == 0 {
		validationErrors = append(validationErrors, ErrEmptyOrder)
	}


	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	var orderItems []*OrderItem

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", itemReq.ProductID, ErrInvalidQuantity))
			continue
		}

		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", itemReq.ProductID, ErrProductNotFound))
			continue
		}

		if product.Stock < itemReq.Quantity {
			validationErrors = append(validationErrors, fmt.Errorf("item %s: %w", product.ID, ErrInsufficientStock))
		}

		orderItems = append(orderItems, &OrderItem{
			Product:  product,
			Quantity: itemReq.Quantity,
			Price:    product.Price,
		})
	}

	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	order := &Order{
		ID:       s.generateID(),
		Customer: req.Customer,
		Items:    orderItems,
		Status:   StatusPending,
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


func (s *OrderService) PayOrder(orderID string) (*Order, error) {
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


func (s *OrderService) CancelOrder(orderID string) (*Order, error) {
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


func (s *OrderService) FindOrderByID(orderID string) (*Order, error) {
	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) ListOrders(filters ...OrderFilter) ([]*Order, error) {
	return s.orderRepo.List(filters...)
}