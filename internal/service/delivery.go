package service

import (
	"context"
	"fmt"

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

func (s *DeliveryService) ListDeliveries(
	ctx context.Context,
	filter model.DeliveryFilter,
) ([]model.Delivery, error) {
	return s.repository.List(ctx, filter)
}

func (s *DeliveryService) GetDelivery(ctx context.Context, id int64) (model.Delivery, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *DeliveryService) GetStatusHistory(
	ctx context.Context,
	deliveryID int64,
) ([]model.StatusHistory, error) {
	return s.repository.GetStatusHistory(ctx, deliveryID)
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, delivery model.Delivery) (int64, error) {
	return s.repository.Create(ctx, delivery)
}

func (s *DeliveryService) UpdateDelivery(
	ctx context.Context,
	delivery model.Delivery,
) error {
	return s.repository.Update(ctx, delivery)
}

func (s *DeliveryService) DeleteDelivery(
	ctx context.Context,
	deliveryID int64,
) error {
	return s.repository.Delete(ctx, deliveryID)
}

func (s *DeliveryService) UpdateStatus(
	ctx context.Context,
	deliveryID int64,
	status string,
	note string,
) error {
	delivery, err := s.repository.GetByID(ctx, deliveryID)
	if err != nil {
		return err
	}

	allowedTransitions := map[string]map[string]bool{
		"pending": {
			"picked_up": true,
			"cancelled": true,
		},
		"picked_up": {
			"in_transit": true,
			"failed":     true,
		},
		"in_transit": {
			"out_for_delivery": true,
			"failed":           true,
		},
		"out_for_delivery": {
			"delivered": true,
			"failed":    true,
		},
	}

	if !allowedTransitions[delivery.Status][status] {
		return fmt.Errorf(
			"cannot change status from %q to %q",
			delivery.Status,
			status,
		)
	}

	return s.repository.UpdateStatus(
		ctx,
		deliveryID,
		status,
		note,
	)
}
