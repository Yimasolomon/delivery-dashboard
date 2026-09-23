package repository

import (
	"context"
	"database/sql"
	"fmt"

	"delivery-dashboard/internal/model"
)

type DeliveryRepository struct {
	db *sql.DB
}

func NewDeliveryRepository(db *sql.DB) *DeliveryRepository {
	return &DeliveryRepository{
		db: db,
	}
}

func (r *DeliveryRepository) List(
	ctx context.Context,
	filter model.DeliveryFilter,
) ([]model.Delivery, error) {

	sortColumn := "d.created_at"
	sortOrder := "DESC"

	switch filter.SortBy {
	case "delivery_date":
		sortColumn = "d.delivery_date"
	case "created_at":
		sortColumn = "d.created_at"
	case "status":
		sortColumn = "d.status"
	case "tracking_id":
		sortColumn = "d.tracking_id"
	}

	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.PageSize

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			d.id,
			d.tracking_id,
			d.customer_id,
			d.driver_id,
			d.pickup_address,
			d.delivery_address,
			d.package_description,
			d.quantity,
			d.weight,
			d.delivery_date,
			d.estimated_delivery,
			d.delivered_at,
			d.status,
			d.notes,
			d.created_at,
			d.updated_at
		FROM deliveries d
		LEFT JOIN customers c ON c.id = d.customer_id
		WHERE
			(
				? = ''
				OR d.tracking_id LIKE '%' || ? || '%'
				OR c.name LIKE '%' || ? || '%'
				OR c.phone LIKE '%' || ? || '%'
				OR d.delivery_address LIKE '%' || ? || '%'
			)
			AND (
					? = ''
					OR (
							? != 'delayed'
							AND d.status = ?
					)
					OR (
							? = 'delayed'
							AND d.estimated_delivery < CURRENT_TIMESTAMP
							AND d.status NOT IN ('delivered', 'cancelled', 'failed')
					)
			)
			AND (
				? = 0
				OR d.driver_id = ?
			)
			AND (
				? = ''
				OR d.delivery_date >= ?
			)
			AND (
					? = ''
					OR d.delivery_date <= ?
			)
			ORDER BY `+sortColumn+` `+sortOrder+`
			LIMIT ? OFFSET ?
	`,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Status,
		filter.Status,
		filter.Status,
		filter.Status,
		filter.DriverID,
		filter.DriverID,
		filter.DateFrom,
		filter.DateFrom,
		filter.DateTo,
		filter.DateTo,
		filter.PageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []model.Delivery

	for rows.Next() {
		var delivery model.Delivery

		err := rows.Scan(
			&delivery.ID,
			&delivery.TrackingID,
			&delivery.CustomerID,
			&delivery.DriverID,
			&delivery.PickupAddress,
			&delivery.DeliveryAddress,
			&delivery.PackageDescription,
			&delivery.Quantity,
			&delivery.Weight,
			&delivery.DeliveryDate,
			&delivery.EstimatedDelivery,
			&delivery.DeliveredAt,
			&delivery.Status,
			&delivery.Notes,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryRepository) Count(
	ctx context.Context,
	filter model.DeliveryFilter,
) (int, error) {
	var count int

	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM deliveries d
		LEFT JOIN customers c ON c.id = d.customer_id
		WHERE
			(
				? = ''
				OR d.tracking_id LIKE '%' || ? || '%'
				OR c.name LIKE '%' || ? || '%'
				OR c.phone LIKE '%' || ? || '%'
				OR d.delivery_address LIKE '%' || ? || '%'
			)
			AND (
					? = ''
					OR (
							? != 'delayed'
							AND d.status = ?
					)
					OR (
							? = 'delayed'
							AND d.estimated_delivery < CURRENT_TIMESTAMP
							AND d.status NOT IN ('delivered', 'cancelled', 'failed')
					)
			)
			AND (
				? = 0
				OR d.driver_id = ?
			)
			AND (
				? = ''
				OR d.delivery_date >= ?
			)
			AND (
				? = ''
				OR d.delivery_date <= ?
			)
	`,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Search,
		filter.Status,
		filter.Status,
		filter.Status,
		filter.Status,
		filter.DriverID,
		filter.DriverID,
		filter.DateFrom,
		filter.DateFrom,
		filter.DateTo,
		filter.DateTo,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DeliveryRepository) GetByID(ctx context.Context, id int64) (model.Delivery, error) {
	var delivery model.Delivery

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			tracking_id,
			customer_id,
			driver_id,
			pickup_address,
			delivery_address,
			package_description,
			quantity,
			weight,
			delivery_date,
			estimated_delivery,
			delivered_at,
			status,
			notes,
			created_at,
			updated_at
		FROM deliveries
		WHERE id = ?
	`, id).Scan(
		&delivery.ID,
		&delivery.TrackingID,
		&delivery.CustomerID,
		&delivery.DriverID,
		&delivery.PickupAddress,
		&delivery.DeliveryAddress,
		&delivery.PackageDescription,
		&delivery.Quantity,
		&delivery.Weight,
		&delivery.DeliveryDate,
		&delivery.EstimatedDelivery,
		&delivery.DeliveredAt,
		&delivery.Status,
		&delivery.Notes,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err != nil {
		return model.Delivery{}, err
	}

	return delivery, nil
}

func (r *DeliveryRepository) Create(ctx context.Context, delivery model.Delivery) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var trackingNumber int64

	err = tx.QueryRowContext(ctx, `
		UPDATE tracking_sequence
		SET next_value = next_value + 1
		WHERE id = 1
		RETURNING next_value - 1
	`).Scan(&trackingNumber)
	if err != nil {
		return 0, err
	}

	trackingID := fmt.Sprintf("TRK-%06d", trackingNumber)

	result, err := tx.ExecContext(ctx, `
		INSERT INTO deliveries (
			tracking_id,
			customer_id,
			driver_id,
			pickup_address,
			delivery_address,
			package_description,
			quantity,
			weight,
			delivery_date,
			estimated_delivery,
			delivered_at,
			status,
			notes
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		trackingID,
		delivery.CustomerID,
		delivery.DriverID,
		delivery.PickupAddress,
		delivery.DeliveryAddress,
		delivery.PackageDescription,
		delivery.Quantity,
		delivery.Weight,
		delivery.DeliveryDate,
		delivery.EstimatedDelivery,
		delivery.DeliveredAt,
		delivery.Status,
		delivery.Notes,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DeliveryRepository) Update(ctx context.Context, delivery model.Delivery) error {
	result, err := r.db.ExecContext(ctx, `
        UPDATE deliveries
        SET
            customer_id = ?,
            driver_id = ?,
            pickup_address = ?,
            delivery_address = ?,
            package_description = ?,
            quantity = ?,
            weight = ?,
            delivery_date = ?,
            estimated_delivery = ?,
            notes = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `,
		delivery.CustomerID,
		delivery.DriverID,
		delivery.PickupAddress,
		delivery.DeliveryAddress,
		delivery.PackageDescription,
		delivery.Quantity,
		delivery.Weight,
		delivery.DeliveryDate,
		delivery.EstimatedDelivery,
		delivery.Notes,
		delivery.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *DeliveryRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM deliveries
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *DeliveryRepository) UpdateStatus(
	ctx context.Context,
	deliveryID int64,
	status string,
	note string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var result sql.Result

	if status == "delivered" {
		result, err = tx.ExecContext(ctx, `
			UPDATE deliveries
			SET
				status = ?,
				delivered_at = CURRENT_TIMESTAMP,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, status, deliveryID)
	} else {
		result, err = tx.ExecContext(ctx, `
			UPDATE deliveries
			SET
				status = ?,
				delivered_at = NULL,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, status, deliveryID)
	}

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO delivery_status_history (
			delivery_id,
			status,
			note
		)
		VALUES (?, ?, ?)
	`, deliveryID, status, note)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryRepository) GetStatusHistory(
	ctx context.Context,
	deliveryID int64,
) ([]model.StatusHistory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			delivery_id,
			status,
			note,
			created_at
		FROM delivery_status_history
		WHERE delivery_id = ?
		ORDER BY created_at ASC, id ASC
	`, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.StatusHistory

	for rows.Next() {
		var entry model.StatusHistory

		err := rows.Scan(
			&entry.ID,
			&entry.DeliveryID,
			&entry.Status,
			&entry.Note,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
