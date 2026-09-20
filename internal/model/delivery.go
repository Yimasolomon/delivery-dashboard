package model

import "time"

const (
	StatusPending        = "pending"
	StatusPickedUp       = "picked_up"
	StatusInTransit      = "in_transit"
	StatusOutForDelivery = "out_for_delivery"
	StatusDelivered      = "delivered"
	StatusFailed         = "failed"
	StatusCancelled      = "cancelled"
)

var DeliveryStatuses = []string{
	StatusPending,
	StatusPickedUp,
	StatusInTransit,
	StatusOutForDelivery,
	StatusDelivered,
	StatusFailed,
	StatusCancelled,
}

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

type DeliveryFilter struct {
	Search    string
	Status    string
	DriverID  int64
	DateFrom  string
	DateTo    string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}
