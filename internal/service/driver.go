package service

import (
	"context"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/repository"
)

type DriverService struct {
	repository *repository.DriverRepository
}

func NewDriverService(repository *repository.DriverRepository) *DriverService {
	return &DriverService{
		repository: repository,
	}
}

func (s *DriverService) ListDrivers(ctx context.Context) ([]model.Driver, error) {
	return s.repository.List(ctx)
}

func (s *DriverService) GetDriver(ctx context.Context, id int64) (model.Driver, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *DriverService) CreateDriver(ctx context.Context, driver model.Driver) (int64, error) {
	return s.repository.Create(ctx, driver)
}

func (s *DriverService) UpdateDriver(ctx context.Context, driver model.Driver) error {
	return s.repository.Update(ctx, driver)
}

func (s *DriverService) DeleteDriver(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}