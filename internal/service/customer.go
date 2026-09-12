package service

import (
	"context"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/repository"
)

type CustomerService struct {
	repository *repository.CustomerRepository
}

func NewCustomerService(repository *repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		repository: repository,
	}
}

func (s *CustomerService) ListCustomers(ctx context.Context) ([]model.Customer, error) {
	return s.repository.List(ctx)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id int64) (model.Customer, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer model.Customer) (int64, error) {
	return s.repository.Create(ctx, customer)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer model.Customer) error {
	return s.repository.Update(ctx, customer)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
