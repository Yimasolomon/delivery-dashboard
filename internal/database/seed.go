package database

import (
	"database/sql"
	"fmt"
	"time"
)

func Seed(db *sql.DB) error {
	var count int

	err := db.QueryRow(`SELECT COUNT(*) FROM customers`).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking existing customers: %w", err)
	}

	if count > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("starting seed transaction: %w", err)
	}
	defer tx.Rollback()

	customerIDs, err := seedCustomers(tx)
	if err != nil {
		return err
	}

	driverIDs, err := seedDrivers(tx)
	if err != nil {
		return err
	}

	if err := seedDeliveries(tx, customerIDs, driverIDs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing seed transaction: %w", err)
	}

	return nil
}

func seedCustomers(tx *sql.Tx) ([]int64, error) {
	customers := []struct {
		name    string
		phone   string
		email   string
		address string
	}{
		{
			name:    "Daniel Okafor",
			phone:   "0900 000 0001",
			email:   "daniel@example.com",
			address: "12 Allen Avenue, Ikeja, Lagos",
		},
		{
			name:    "Amaka Eze",
			phone:   "0900 000 0002",
			email:   "amaka@example.com",
			address: "8 Admiralty Way, Lekki, Lagos",
		},
		{
			name:    "Tunde Adebayo",
			phone:   "0900 000 0003",
			email:   "tunde@example.com",
			address: "21 Herbert Macaulay Way, Yaba, Lagos",
		},
		{
			name:    "Sarah Williams",
			phone:   "0900 000 0004",
			email:   "sarah@example.com",
			address: "15 Awolowo Road, Ikoyi, Lagos",
		},
		{
			name:    "Chinedu Obi",
			phone:   "0900 000 0005",
			email:   "chinedu@example.com",
			address: "4 Adeniran Ogunsanya Street, Surulere, Lagos",
		},
		{
			name:    "Fatima Bello",
			phone:   "0900 000 0006",
			email:   "fatima@example.com",
			address: "18 Mobolaji Bank Anthony Way, Maryland, Lagos",
		},
		{
			name:    "Michael Johnson",
			phone:   "0900 000 0007",
			email:   "michael@example.com",
			address: "7 Gbagada Expressway, Gbagada, Lagos",
		},
		{
			name:    "Ngozi Nwosu",
			phone:   "0900 000 0008",
			email:   "ngozi@example.com",
			address: "10 Admiralty Road, Lekki, Lagos",
		},
		{
			name:    "Ibrahim Musa",
			phone:   "0900 000 0009",
			email:   "ibrahim@example.com",
			address: "6 Marina Road, Lagos Island, Lagos",
		},
		{
			name:    "Grace Adeyemi",
			phone:   "0900 000 0010",
			email:   "grace@example.com",
			address: "14 Lekki-Epe Expressway, Ajah, Lagos",
		},
	}

	ids := make([]int64, 0, len(customers))

	query := `
		INSERT INTO customers (name, phone, email, address)
		VALUES (?, ?, ?, ?)
	`

	for _, customer := range customers {
		result, err := tx.Exec(
			query,
			customer.name,
			customer.phone,
			customer.email,
			customer.address,
		)
		if err != nil {
			return nil, fmt.Errorf("seeding customer %q: %w", customer.name, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("getting customer ID: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func seedDrivers(tx *sql.Tx) ([]int64, error) {
	drivers := []struct {
		name                string
		phone               string
		vehicle             string
		vehicleRegistration string
		status              string
	}{
		{
			name:                "Emeka Johnson",
			phone:               "0900 100 0001",
			vehicle:             "Toyota Hiace",
			vehicleRegistration: "LAG-DEV-001",
			status:              "available",
		},
		{
			name:                "Yusuf Ibrahim",
			phone:               "0900 100 0002",
			vehicle:             "Ford Transit",
			vehicleRegistration: "LAG-DEV-002",
			status:              "on_delivery",
		},
		{
			name:                "Peter Okoro",
			phone:               "0900 100 0003",
			vehicle:             "Toyota Hilux",
			vehicleRegistration: "LAG-DEV-003",
			status:              "available",
		},
		{
			name:                "David Adekunle",
			phone:               "0900 100 0004",
			vehicle:             "Nissan NV200",
			vehicleRegistration: "LAG-DEV-004",
			status:              "offline",
		},
		{
			name:                "Musa Abdullahi",
			phone:               "0900 100 0005",
			vehicle:             "Toyota Corolla",
			vehicleRegistration: "LAG-DEV-005",
			status:              "available",
		},
		{
			name:                "Samuel Eze",
			phone:               "0900 100 0006",
			vehicle:             "Suzuki Carry",
			vehicleRegistration: "LAG-DEV-006",
			status:              "on_delivery",
		},
	}

	ids := make([]int64, 0, len(drivers))

	query := `
		INSERT INTO drivers
			(name, phone, vehicle, vehicle_registration, status)
		VALUES (?, ?, ?, ?, ?)
	`

	for _, driver := range drivers {
		result, err := tx.Exec(
			query,
			driver.name,
			driver.phone,
			driver.vehicle,
			driver.vehicleRegistration,
			driver.status,
		)
		if err != nil {
			return nil, fmt.Errorf("seeding driver %q: %w", driver.name, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("getting driver ID: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func seedDeliveries(tx *sql.Tx, customerIDs, driverIDs []int64) error {
	now := time.Now()

	deliveries := []struct {
		customerID        int64
		driverID          int64
		pickupAddress     string
		deliveryAddress   string
		packageDesc       string
		quantity          int
		weight            float64
		deliveryDate      time.Time
		estimatedDelivery time.Time
		deliveredAt       *time.Time
		status            string
		notes             string
	}{
		// 1
		{
			customerID:        customerIDs[0],
			driverID:          driverIDs[0],
			pickupAddress:     "Ikeja City Mall, Ikeja, Lagos",
			deliveryAddress:   "12 Admiralty Way, Lekki, Lagos",
			packageDesc:       "Electronics accessories",
			quantity:          3,
			weight:            4.5,
			deliveryDate:      now.AddDate(0, 0, -2),
			estimatedDelivery: now.AddDate(0, 0, -2).Add(4 * time.Hour),
			status:            "delivered",
			notes:             "Delivered successfully",
		},

		// 2
		{
			customerID:        customerIDs[1],
			driverID:          driverIDs[1],
			pickupAddress:     "Yaba, Lagos",
			deliveryAddress:   "Victoria Island, Lagos",
			packageDesc:       "Office supplies",
			quantity:          5,
			weight:            8.2,
			deliveryDate:      now.AddDate(0, 0, -1),
			estimatedDelivery: now.AddDate(0, 0, -1).Add(5 * time.Hour),
			status:            "delivered",
			notes:             "Delivered to reception",
		},

		// 3
		{
			customerID:        customerIDs[2],
			driverID:          driverIDs[2],
			pickupAddress:     "Surulere, Lagos",
			deliveryAddress:   "Maryland, Lagos",
			packageDesc:       "Clothing items",
			quantity:          2,
			weight:            3.1,
			deliveryDate:      now,
			estimatedDelivery: now.Add(3 * time.Hour),
			status:            "out_for_delivery",
			notes:             "Driver is heading to customer",
		},

		// 4
		{
			customerID:        customerIDs[3],
			driverID:          driverIDs[3],
			pickupAddress:     "Ikoyi, Lagos",
			deliveryAddress:   "Lagos Island, Lagos",
			packageDesc:       "Documents",
			quantity:          1,
			weight:            0.8,
			deliveryDate:      now,
			estimatedDelivery: now.Add(2 * time.Hour),
			status:            "in_transit",
			notes:             "Important business documents",
		},

		// 5
		{
			customerID:        customerIDs[4],
			driverID:          driverIDs[4],
			pickupAddress:     "Gbagada, Lagos",
			deliveryAddress:   "Ajah, Lagos",
			packageDesc:       "Household items",
			quantity:          4,
			weight:            12.5,
			deliveryDate:      now,
			estimatedDelivery: now.Add(4 * time.Hour),
			status:            "picked_up",
			notes:             "Picked up from sender",
		},

		// 6
		{
			customerID:        customerIDs[5],
			driverID:          driverIDs[5],
			pickupAddress:     "Lekki Phase 1, Lagos",
			deliveryAddress:   "Yaba, Lagos",
			packageDesc:       "Computer equipment",
			quantity:          2,
			weight:            6.7,
			deliveryDate:      now,
			estimatedDelivery: now.Add(2 * time.Hour),
			status:            "out_for_delivery",
			notes:             "Fragile equipment",
		},

		// 7
		{
			customerID:        customerIDs[6],
			driverID:          driverIDs[0],
			pickupAddress:     "Victoria Island, Lagos",
			deliveryAddress:   "Surulere, Lagos",
			packageDesc:       "Books",
			quantity:          10,
			weight:            9.4,
			deliveryDate:      now.AddDate(0, 0, -3),
			estimatedDelivery: now.AddDate(0, 0, -3).Add(6 * time.Hour),
			status:            "delivered",
			notes:             "Delivered to customer",
		},

		// 8
		{
			customerID:        customerIDs[7],
			driverID:          driverIDs[2],
			pickupAddress:     "Maryland, Lagos",
			deliveryAddress:   "Ikoyi, Lagos",
			packageDesc:       "Kitchen equipment",
			quantity:          3,
			weight:            7.8,
			deliveryDate:      now,
			estimatedDelivery: now.Add(5 * time.Hour),
			status:            "pending",
			notes:             "Awaiting pickup",
		},

		// 9
		{
			customerID:        customerIDs[8],
			driverID:          driverIDs[4],
			pickupAddress:     "Lagos Island, Lagos",
			deliveryAddress:   "Gbagada, Lagos",
			packageDesc:       "Personal care products",
			quantity:          6,
			weight:            5.6,
			deliveryDate:      now,
			estimatedDelivery: now.Add(3 * time.Hour),
			status:            "in_transit",
			notes:             "Package is currently in transit",
		},

		// 10
		{
			customerID:        customerIDs[9],
			driverID:          driverIDs[5],
			pickupAddress:     "Ajah, Lagos",
			deliveryAddress:   "Lekki, Lagos",
			packageDesc:       "Food packaging materials",
			quantity:          8,
			weight:            14.2,
			deliveryDate:      now.AddDate(0, 0, -1),
			estimatedDelivery: now.AddDate(0, 0, -1).Add(3 * time.Hour),
			status:            "failed",
			notes:             "Delivery attempt unsuccessful",
		},

		// 11
		{
			customerID:        customerIDs[0],
			driverID:          driverIDs[1],
			pickupAddress:     "Ikeja, Lagos",
			deliveryAddress:   "Ikoyi, Lagos",
			packageDesc:       "Mobile phones",
			quantity:          2,
			weight:            1.5,
			deliveryDate:      now,
			estimatedDelivery: now.Add(6 * time.Hour),
			status:            "pending",
			notes:             "Waiting for pickup",
		},

		// 12
		{
			customerID:        customerIDs[1],
			driverID:          driverIDs[3],
			pickupAddress:     "Lekki, Lagos",
			deliveryAddress:   "Lagos Island, Lagos",
			packageDesc:       "Fashion accessories",
			quantity:          7,
			weight:            4.3,
			deliveryDate:      now.AddDate(0, 0, -1),
			estimatedDelivery: now.AddDate(0, 0, -1).Add(4 * time.Hour),
			status:            "delivered",
			notes:             "Delivered without issues",
		},

		// 13
		{
			customerID:        customerIDs[2],
			driverID:          driverIDs[4],
			pickupAddress:     "Yaba, Lagos",
			deliveryAddress:   "Gbagada, Lagos",
			packageDesc:       "Printer cartridges",
			quantity:          4,
			weight:            3.6,
			deliveryDate:      now,
			estimatedDelivery: now.Add(4 * time.Hour),
			status:            "picked_up",
			notes:             "Package collected",
		},

		// 14
		{
			customerID:        customerIDs[3],
			driverID:          driverIDs[0],
			pickupAddress:     "Victoria Island, Lagos",
			deliveryAddress:   "Ajah, Lagos",
			packageDesc:       "Office furniture parts",
			quantity:          6,
			weight:            18.4,
			deliveryDate:      now,
			estimatedDelivery: now.Add(5 * time.Hour),
			status:            "in_transit",
			notes:             "Large package",
		},

		// 15
		{
			customerID:        customerIDs[4],
			driverID:          driverIDs[2],
			pickupAddress:     "Surulere, Lagos",
			deliveryAddress:   "Ikeja, Lagos",
			packageDesc:       "Marketing materials",
			quantity:          20,
			weight:            11.2,
			deliveryDate:      now.AddDate(0, 0, -2),
			estimatedDelivery: now.AddDate(0, 0, -2).Add(5 * time.Hour),
			status:            "delivered",
			notes:             "Delivered to office",
		},

		// 16
		{
			customerID:        customerIDs[5],
			driverID:          driverIDs[5],
			pickupAddress:     "Maryland, Lagos",
			deliveryAddress:   "Lekki, Lagos",
			packageDesc:       "Computer accessories",
			quantity:          5,
			weight:            5.8,
			deliveryDate:      now,
			estimatedDelivery: now.Add(4 * time.Hour),
			status:            "out_for_delivery",
			notes:             "Driver near destination",
		},

		// 17
		{
			customerID:        customerIDs[6],
			driverID:          driverIDs[1],
			pickupAddress:     "Gbagada, Lagos",
			deliveryAddress:   "Victoria Island, Lagos",
			packageDesc:       "Medical supplies",
			quantity:          9,
			weight:            10.7,
			deliveryDate:      now,
			estimatedDelivery: now.Add(3 * time.Hour),
			status:            "in_transit",
			notes:             "Priority delivery",
		},

		// 18
		{
			customerID:        customerIDs[7],
			driverID:          driverIDs[3],
			pickupAddress:     "Ikoyi, Lagos",
			deliveryAddress:   "Yaba, Lagos",
			packageDesc:       "Books and stationery",
			quantity:          12,
			weight:            8.9,
			deliveryDate:      now.AddDate(0, 0, -1),
			estimatedDelivery: now.AddDate(0, 0, -1).Add(3 * time.Hour),
			status:            "cancelled",
			notes:             "Customer cancelled the order",
		},

		// 19
		{
			customerID:        customerIDs[8],
			driverID:          driverIDs[4],
			pickupAddress:     "Lagos Island, Lagos",
			deliveryAddress:   "Maryland, Lagos",
			packageDesc:       "Cosmetics",
			quantity:          5,
			weight:            3.2,
			deliveryDate:      now.AddDate(0, 0, -2),
			estimatedDelivery: now.AddDate(0, 0, -2).Add(3 * time.Hour),
			status:            "failed",
			notes:             "Recipient was unavailable",
		},

		// 20
		{
			customerID:        customerIDs[9],
			driverID:          driverIDs[0],
			pickupAddress:     "Ajah, Lagos",
			deliveryAddress:   "Victoria Island, Lagos",
			packageDesc:       "Home appliances",
			quantity:          2,
			weight:            15.6,
			deliveryDate:      now,
			estimatedDelivery: now.Add(6 * time.Hour),
			status:            "pending",
			notes:             "Pickup scheduled",
		},

		// 21
		{
			customerID:        customerIDs[0],
			driverID:          driverIDs[2],
			pickupAddress:     "Ikeja, Lagos",
			deliveryAddress:   "Yaba, Lagos",
			packageDesc:       "Electronic components",
			quantity:          15,
			weight:            6.4,
			deliveryDate:      now,
			estimatedDelivery: now.Add(4 * time.Hour),
			status:            "picked_up",
			notes:             "Package picked up",
		},

		// 22
		{
			customerID:        customerIDs[1],
			driverID:          driverIDs[5],
			pickupAddress:     "Lekki, Lagos",
			deliveryAddress:   "Ikoyi, Lagos",
			packageDesc:       "Beauty products",
			quantity:          8,
			weight:            4.1,
			deliveryDate:      now,
			estimatedDelivery: now.Add(2 * time.Hour),
			status:            "out_for_delivery",
			notes:             "Expected shortly",
		},

		// 23
		{
			customerID:        customerIDs[2],
			driverID:          driverIDs[1],
			pickupAddress:     "Surulere, Lagos",
			deliveryAddress:   "Lagos Island, Lagos",
			packageDesc:       "Business documents",
			quantity:          2,
			weight:            1.2,
			deliveryDate:      now.AddDate(0, 0, -3),
			estimatedDelivery: now.AddDate(0, 0, -3).Add(4 * time.Hour),
			status:            "delivered",
			notes:             "Signed for by recipient",
		},

		// 24
		{
			customerID:        customerIDs[3],
			driverID:          driverIDs[4],
			pickupAddress:     "Ikoyi, Lagos",
			deliveryAddress:   "Gbagada, Lagos",
			packageDesc:       "Household electronics",
			quantity:          3,
			weight:            9.8,
			deliveryDate:      now,
			estimatedDelivery: now.Add(5 * time.Hour),
			status:            "pending",
			notes:             "Awaiting driver assignment",
		},
	}

	query := `
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
	`

	for i, delivery := range deliveries {
		trackingID := fmt.Sprintf("TRK-%06d", 100001+i)

		result, err := tx.Exec(
			query,
			trackingID,
			delivery.customerID,
			delivery.driverID,
			delivery.pickupAddress,
			delivery.deliveryAddress,
			delivery.packageDesc,
			delivery.quantity,
			delivery.weight,
			delivery.deliveryDate,
			delivery.estimatedDelivery,
			delivery.deliveredAt,
			delivery.status,
			delivery.notes,
		)
		if err != nil {
			return fmt.Errorf("seeding delivery %s: %w", trackingID, err)
		}

		deliveryID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("getting delivery ID for %s: %w", trackingID, err)
		}

		if err := seedStatusHistory(
			tx,
			deliveryID,
			delivery.status,
			delivery.notes,
		); err != nil {
			return err
		}
	}

	return nil
}

func seedStatusHistory(
	tx *sql.Tx,
	deliveryID int64,
	status string,
	note string,
) error {
	statuses := map[string][]string{
		"pending": {
			"pending",
		},
		"picked_up": {
			"pending",
			"picked_up",
		},
		"in_transit": {
			"pending",
			"picked_up",
			"in_transit",
		},
		"out_for_delivery": {
			"pending",
			"picked_up",
			"in_transit",
			"out_for_delivery",
		},
		"delivered": {
			"pending",
			"picked_up",
			"in_transit",
			"out_for_delivery",
			"delivered",
		},
		"failed": {
			"pending",
			"picked_up",
			"in_transit",
			"failed",
		},
		"cancelled": {
			"pending",
			"cancelled",
		},
	}

	history, ok := statuses[status]
	if !ok {
		return fmt.Errorf("unknown delivery status %q", status)
	}

	query := `
		INSERT INTO delivery_status_history (
			delivery_id,
			status,
			note
		)
		VALUES (?, ?, ?)
	`

	for i, historyStatus := range history {
		historyNote := ""

		if i == len(history)-1 {
			historyNote = note
		}

		if _, err := tx.Exec(
			query,
			deliveryID,
			historyStatus,
			historyNote,
		); err != nil {
			return fmt.Errorf(
				"seeding status history for delivery %d: %w",
				deliveryID,
				err,
			)
		}
	}

	return nil
}
