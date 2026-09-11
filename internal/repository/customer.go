package repository

import (
	"context"
	"database/sql"

	"delivery-dashboard/internal/model"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name,
			phone,
			email,
			address,
			created_at
		FROM customers
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer

	for rows.Next() {
		var customer model.Customer

		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Phone,
			&customer.Email,
			&customer.Address,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (model.Customer, error) {
	var customer model.Customer

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			name,
			phone,
			email,
			address,
			created_at
		FROM customers
		WHERE id = ?
	`, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Phone,
		&customer.Email,
		&customer.Address,
		&customer.CreatedAt,
	)

	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepository) Create(ctx context.Context, customer model.Customer) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO customers (
			name,
			phone,
			email,
			address
		)
		VALUES (?, ?, ?, ?)
	`,
		customer.Name,
		customer.Phone,
		customer.Email,
		customer.Address,
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

func (r *CustomerRepository) Update(ctx context.Context, customer model.Customer) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE customers
		SET
			name = ?,
			phone = ?,
			email = ?,
			address = ?
		WHERE id = ?
	`,
		customer.Name,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.ID,
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

func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM customers
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
