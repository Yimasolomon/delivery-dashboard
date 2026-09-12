package repository

import (
	"context"
	"database/sql"

	"delivery-dashboard/internal/model"
)

type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository(db *sql.DB) *DriverRepository {
	return &DriverRepository{
		db: db,
	}
}

func (r *DriverRepository) List(ctx context.Context) ([]model.Driver, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name,
			phone,
			vehicle,
			vehicle_registration,
			status,
			created_at
		FROM drivers
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []model.Driver

	for rows.Next() {
		var driver model.Driver

		err := rows.Scan(
			&driver.ID,
			&driver.Name,
			&driver.Phone,
			&driver.Vehicle,
			&driver.VehicleRegistration,
			&driver.Status,
			&driver.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		drivers = append(drivers, driver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id int64) (model.Driver, error) {
	var driver model.Driver

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			name,
			phone,
			vehicle,
			vehicle_registration,
			status,
			created_at
		FROM drivers
		WHERE id = ?
	`, id).Scan(
		&driver.ID,
		&driver.Name,
		&driver.Phone,
		&driver.Vehicle,
		&driver.VehicleRegistration,
		&driver.Status,
		&driver.CreatedAt,
	)

	if err != nil {
		return model.Driver{}, err
	}

	return driver, nil
}

func (r *DriverRepository) Create(ctx context.Context, driver model.Driver) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO drivers (
			name,
			phone,
			vehicle,
			vehicle_registration,
			status
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		driver.Name,
		driver.Phone,
		driver.Vehicle,
		driver.VehicleRegistration,
		driver.Status,
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

func (r *DriverRepository) Update(ctx context.Context, driver model.Driver) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE drivers
		SET
			name = ?,
			phone = ?,
			vehicle = ?,
			vehicle_registration = ?,
			status = ?
		WHERE id = ?
	`,
		driver.Name,
		driver.Phone,
		driver.Vehicle,
		driver.VehicleRegistration,
		driver.Status,
		driver.ID,
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

func (r *DriverRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM drivers
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
