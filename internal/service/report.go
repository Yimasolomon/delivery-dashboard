package service

import (
	"context"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/repository"
)

type ReportService struct {
	repository *repository.ReportRepository
}

func NewReportService(
	repository *repository.ReportRepository,
) *ReportService {
	return &ReportService{
		repository: repository,
	}
}

func (s *ReportService) GetSummary(
	ctx context.Context,
	filter model.ReportFilter,
) (model.ReportSummary, error) {
	return s.repository.GetSummary(ctx, filter)
}

func (s *ReportService) GetDriverPerformance(
	ctx context.Context,
	filter model.ReportFilter,
) ([]model.DriverPerformance, error) {
	return s.repository.GetDriverPerformance(ctx, filter)
}

func (s *ReportService) GetCustomerActivity(
	ctx context.Context,
	filter model.ReportFilter,
) ([]model.CustomerActivity, error) {
	return s.repository.GetCustomerActivity(ctx, filter)
}
