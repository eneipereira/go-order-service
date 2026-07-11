package repository

import "github.com/eneipereira/go-order-service/model"

type ProductRepository interface {
	Save(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	List() ([]*model.Product, error)
}


type inMemoryProductRepository struct {
	products map[string]*model.Product
}


func NewInMemoryProductRepository() ProductRepository {
	return &inMemoryProductRepository{
		products: make(map[string]*model.Product),
	}
}

func (r *inMemoryProductRepository) Save(product *model.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepository) FindByID(id string) (*model.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return nil, model.ErrProductNotFound
	}
	return product, nil
}

func (r *inMemoryProductRepository) List() ([]*model.Product, error) {
	var productList []*model.Product
	for _, p := range r.products {
		productList = append(productList, p)
	}
	return productList, nil
}

type OrderFilter func(order *model.Order) bool

type OrderRepository interface {
	Save(order *model.Order) error
	FindByID(id string) (*model.Order, error)
	List(filters ...OrderFilter) ([]*model.Order, error)
}


type inMemoryOrderRepository struct {
	orders map[string]*model.Order
}


func NewInMemoryOrderRepository() OrderRepository {
	return &inMemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (r *inMemoryOrderRepository) Save(order *model.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryOrderRepository) FindByID(id string) (*model.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, model.ErrOrderNotFound
	}
	return order, nil
}

func (r *inMemoryOrderRepository) List(filters ...OrderFilter) ([]*model.Order, error) {
	var orderList []*model.Order
outer:
	for _, order := range r.orders {
	
		for _, filter := range filters {
			if !filter(order) {
				continue outer
			}
		}
		orderList = append(orderList, order)
	}
	return orderList, nil
}