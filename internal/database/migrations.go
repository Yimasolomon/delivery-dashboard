package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT,
			address TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS drivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			vehicle TEXT NOT NULL,
			vehicle_registration TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'offline',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			CHECK (status IN ('available', 'on_delivery', 'offline'))
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS deliveries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tracking_id TEXT NOT NULL UNIQUE,

			customer_id INTEGER NOT NULL,
			driver_id INTEGER NOT NULL,

			pickup_address TEXT NOT NULL,
			delivery_address TEXT NOT NULL,
			package_description TEXT NOT NULL,

			quantity INTEGER NOT NULL,
			weight REAL NOT NULL,

			delivery_date DATE NOT NULL,
			estimated_delivery DATETIME NOT NULL,
			delivered_at DATETIME,

			status TEXT NOT NULL DEFAULT 'pending',

			notes TEXT,

			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (customer_id)
				REFERENCES customers(id),

			FOREIGN KEY (driver_id)
				REFERENCES drivers(id),

			CHECK (quantity > 0),
			CHECK (weight >= 0),

			CHECK (
				status IN (
					'pending',
					'picked_up',
					'in_transit',
					'out_for_delivery',
					'delivered',
					'failed',
					'cancelled'
				)
			)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS delivery_status_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,

			delivery_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			note TEXT,

			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (delivery_id)
				REFERENCES deliveries(id)
				ON DELETE CASCADE,

			CHECK (
				status IN (
					'pending',
					'picked_up',
					'in_transit',
					'out_for_delivery',
					'delivered',
					'failed',
					'cancelled'
				)
			)
		);
		`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("executing migration: %w", err)
		}
	}

	return createIndexes(db)
}

func createIndexes(db *sql.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_deliveries_tracking_id
			ON deliveries(tracking_id);`,

		`CREATE INDEX IF NOT EXISTS idx_deliveries_status
			ON deliveries(status);`,

		`CREATE INDEX IF NOT EXISTS idx_deliveries_customer_id
			ON deliveries(customer_id);`,

		`CREATE INDEX IF NOT EXISTS idx_deliveries_driver_id
			ON deliveries(driver_id);`,

		`CREATE INDEX IF NOT EXISTS idx_deliveries_created_at
			ON deliveries(created_at);`,

		`CREATE INDEX IF NOT EXISTS idx_deliveries_estimated_delivery
			ON deliveries(estimated_delivery);`,

		`CREATE INDEX IF NOT EXISTS idx_status_history_delivery_id
			ON delivery_status_history(delivery_id);`,
	}

	for _, statement := range indexes {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("creating database index: %w", err)
		}
	}

	return nil
}
