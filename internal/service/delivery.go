package service

import (
	"context"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/repository"
)

type DeliveryService struct {
	repository *repository.DeliveryRepository
}

func NewDeliveryService(repository *repository.DeliveryRepository) *DeliveryService {
	return &DeliveryService{
		repository: repository,
	}
}

func (s *DeliveryService) ListDeliveries(ctx context.Context) ([]model.Delivery, error) {
	return s.repository.List(ctx)
}

func (s *DeliveryService) GetDelivery(ctx context.Context, id int64) (model.Delivery, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, delivery model.Delivery) (int64, error) {
	return s.repository.Create(ctx, delivery)
}
