package main

import (
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
	
	if req.Customer == "" {
		return nil, ErrInvalidCustomer
	}
	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	var orderItems []*OrderItem
	var productsToUpdate []*Product

	
	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("Error trying to find product with ID %s: %w", itemReq.ProductID, err)
		}

		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("%w: product with ID %s", ErrInsufficientStock, product.ID)
		}

		orderItems = append(orderItems, &OrderItem{
			Product:  product,
			Quantity: itemReq.Quantity,
			Price:    product.Price, 
		})
		productsToUpdate = append(productsToUpdate, product)
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

// ListOrders retorna uma lista de pedidos, opcionalmente filtrada.
func (s *OrderService) ListOrders(filters ...OrderFilter) ([]*Order, error) {
	return s.orderRepo.List(filters...)
}