package service

import (
	"context"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req dto.CreateProductDTO) (*model.Product, error) {
	product, err := model.NewProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) FindByID(ctx context.Context, id string) (*model.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProductService) FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return s.repo.FindAll(ctx, limit, offset)
}
