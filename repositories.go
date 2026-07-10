package main


type ProductRepository interface {
	Save(product *Product) error
	FindByID(id string) (*Product, error)
	List() ([]*Product, error)
}


type inMemoryProductRepository struct {
	products map[string]*Product
}


func NewInMemoryProductRepository() ProductRepository {
	return &inMemoryProductRepository{
		products: make(map[string]*Product),
	}
}

func (r *inMemoryProductRepository) Save(product *Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepository) FindByID(id string) (*Product, error) {
	product, ok := r.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (r *inMemoryProductRepository) List() ([]*Product, error) {
	var productList []*Product
	for _, p := range r.products {
		productList = append(productList, p)
	}
	return productList, nil
}

// OrderFilter é uma função usada para filtrar pedidos.
type OrderFilter func(order *Order) bool

type OrderRepository interface {
	Save(order *Order) error
	FindByID(id string) (*Order, error)
	List(filters ...OrderFilter) ([]*Order, error)
}


type inMemoryOrderRepository struct {
	orders map[string]*Order
}


func NewInMemoryOrderRepository() OrderRepository {
	return &inMemoryOrderRepository{
		orders: make(map[string]*Order),
	}
}

func (r *inMemoryOrderRepository) Save(order *Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryOrderRepository) FindByID(id string) (*Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (r *inMemoryOrderRepository) List(filters ...OrderFilter) ([]*Order, error) {
	var orderList []*Order
outer:
	for _, order := range r.orders {
		// Aplica todos os filtros para o pedido atual
		for _, filter := range filters {
			if !filter(order) {
				continue outer // Se não passar em um filtro, vai para o próximo pedido
			}
		}
		orderList = append(orderList, order)
	}
	return orderList, nil
}