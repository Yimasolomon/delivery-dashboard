package model

import "time"

type Delivery struct {
	ID                 int64
	TrackingID         string
	CustomerID         int64
	DriverID           int64
	PickupAddress      string
	DeliveryAddress    string
	PackageDescription string
	Quantity           int
	Weight             float64
	DeliveryDate       time.Time
	EstimatedDelivery  time.Time
	DeliveredAt        *time.Time
	Status             string
	Notes              string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
