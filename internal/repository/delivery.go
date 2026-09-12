package repository

import (
	"context"
	"database/sql"

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

func (r *DeliveryRepository) List(ctx context.Context) ([]model.Delivery, error) {
	rows, err := r.db.QueryContext(ctx, `
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
		ORDER BY created_at DESC
	`)
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
	result, err := r.db.ExecContext(ctx, `
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
		delivery.TrackingID,
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
			delivered_at = ?,
			status = ?,
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
		delivery.DeliveredAt,
		delivery.Status,
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

	result, err := tx.ExecContext(ctx, `
		UPDATE deliveries
		SET
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, deliveryID)
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
