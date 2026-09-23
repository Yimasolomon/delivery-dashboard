package repository

import (
	"context"
	"database/sql"

	"delivery-dashboard/internal/model"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) GetSummary(
	ctx context.Context,
	filter model.ReportFilter,
) (model.ReportSummary, error) {
	var summary model.ReportSummary

	err := r.db.QueryRowContext(ctx, `
		SELECT
    COUNT(*),
		COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'picked_up' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'in_transit' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'out_for_delivery' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(
			CASE
				WHEN estimated_delivery < CURRENT_TIMESTAMP
					AND status NOT IN ('delivered', 'cancelled', 'failed')
				THEN 1
				ELSE 0
			END
			), 0)
		FROM deliveries
		WHERE (
			? = ''
			OR substr(delivery_date, 1, 10) <= ?
		)
		AND (
			? = ''
			OR substr(delivery_date, 1, 10) >= ?
		)
	`,
		filter.StartDate,
		filter.StartDate,
		filter.EndDate,
		filter.EndDate,
	).Scan(
		&summary.Total,
		&summary.Pending,
		&summary.PickedUp,
		&summary.InTransit,
		&summary.OutForDelivery,
		&summary.Delivered,
		&summary.Failed,
		&summary.Cancelled,
		&summary.Delayed,
	)

	if err != nil {
		return model.ReportSummary{}, err
	}

	return summary, nil
}

func (r *ReportRepository) GetDriverPerformance(
	ctx context.Context,
	filter model.ReportFilter,
) ([]model.DriverPerformance, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			d.id,
			d.name,
			COUNT(del.id),
			COALESCE(SUM(CASE WHEN del.status = 'delivered' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN del.status = 'failed' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(
				CASE
					WHEN del.status IN (
						'pending',
						'picked_up',
						'in_transit',
						'out_for_delivery'
					)
					THEN 1
					ELSE 0
				END
			), 0)
		FROM drivers d
		LEFT JOIN deliveries del
			ON del.driver_id = d.id
			AND (
				? = ''
				OR substr(del.delivery_date, 1, 10) >= ?
			)
			AND (
				? = ''
				OR substr(del.delivery_date, 1, 10) >= ?
			)
		GROUP BY d.id, d.name
		ORDER BY d.name
	`,
		filter.StartDate,
		filter.StartDate,
		filter.EndDate,
		filter.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.DriverPerformance

	for rows.Next() {
		var item model.DriverPerformance

		if err := rows.Scan(
			&item.DriverID,
			&item.DriverName,
			&item.Total,
			&item.Delivered,
			&item.Failed,
			&item.Active,
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

func (r *ReportRepository) GetCustomerActivity(
	ctx context.Context,
	filter model.ReportFilter,
) ([]model.CustomerActivity, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			c.id,
			c.name,
			COUNT(del.id)
		FROM customers c
		LEFT JOIN deliveries del
			ON del.customer_id = c.id
			AND (
				? = ''
				OR substr(del.delivery_date, 1, 10) >= ?
			)
			AND (
				? = ''
				OR substr(del.delivery_date, 1, 10) >= ?
			)
		GROUP BY c.id, c.name
		ORDER BY c.name
	`,
		filter.StartDate,
		filter.StartDate,
		filter.EndDate,
		filter.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.CustomerActivity

	for rows.Next() {
		var item model.CustomerActivity

		if err := rows.Scan(
			&item.CustomerID,
			&item.CustomerName,
			&item.Total,
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
