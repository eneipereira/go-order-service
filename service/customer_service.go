package service

import (
	"context"

	"github.com/eneipereira/go-order-service/dto"
	"github.com/eneipereira/go-order-service/model"
	"golang.org/x/crypto/bcrypt"
)

type CustomerRepository interface {
	Save(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	FindByID(ctx context.Context, id string) (*model.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error)
}

type CustomerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, req dto.CustomerDTO) (*model.Customer, error) {
	err := model.ValidateCustomerPassword(req.Password)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	customer, err := model.NewCustomerNoPassword(req.Name, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	customer.PasswordHash = string(hashedPassword)

	return s.repo.Save(ctx, customer)
}

func (s *CustomerService) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CustomerService) FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	return s.repo.FindAll(ctx, limit, offset)
}
