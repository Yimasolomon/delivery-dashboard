package repository

import (
	"context"
	"database/sql"

	"delivery-dashboard/internal/model"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{
		db: db,
	}
}

func (r *DashboardRepository) GetStats(
	ctx context.Context,
) (model.DashboardStats, error) {
	var stats model.DashboardStats

	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'in_transit' THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'out_for_delivery' THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END),
			SUM(
				CASE
					WHEN estimated_delivery < CURRENT_TIMESTAMP
						AND status NOT IN ('delivered', 'cancelled', 'failed')
					THEN 1
					ELSE 0
				END
			)
		FROM deliveries
	`).Scan(
		&stats.Total,
		&stats.Pending,
		&stats.InTransit,
		&stats.OutForDelivery,
		&stats.Delivered,
		&stats.Failed,
		&stats.Delayed,
	)

	if err != nil {
		return model.DashboardStats{}, err
	}

	return stats, nil
}

func (r *DashboardRepository) GetStatusCounts(
	ctx context.Context,
) ([]model.StatusCount, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			status,
			COUNT(*)
		FROM deliveries
		GROUP BY status
		ORDER BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.StatusCount

	for rows.Next() {
		var item model.StatusCount

		if err := rows.Scan(
			&item.Status,
			&item.Count,
		); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *DashboardRepository) GetDeliveryTrends(
	ctx context.Context,
) ([]model.DeliveryTrend, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(DATE(delivery_date), 'Unknown'),
			COUNT(*)
		FROM deliveries
		GROUP BY COALESCE(DATE(delivery_date), 'Unknown')
		ORDER BY COALESCE(DATE(delivery_date), 'Unknown')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.DeliveryTrend

	for rows.Next() {
		var item model.DeliveryTrend

		if err := rows.Scan(
			&item.Date,
			&item.Count,
		); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
