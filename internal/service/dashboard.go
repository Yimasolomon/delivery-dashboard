package service

import (
	"context"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/repository"
)

type DashboardService struct {
	repository *repository.DashboardRepository
}

func NewDashboardService(
	repository *repository.DashboardRepository,
) *DashboardService {
	return &DashboardService{
		repository: repository,
	}
}

func (s *DashboardService) GetStats(
	ctx context.Context,
) (model.DashboardStats, error) {
	return s.repository.GetStats(ctx)
}

func (s *DashboardService) GetStatusCounts(
	ctx context.Context,
) ([]model.StatusCount, error) {
	return s.repository.GetStatusCounts(ctx)
}

func (s *DashboardService) GetDeliveryTrends(
	ctx context.Context,
) ([]model.DeliveryTrend, error) {
	return s.repository.GetDeliveryTrends(ctx)
}
