package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"delivery-dashboard/internal/model"
)

type createDeliveryForm struct {
	CustomerID         string
	DriverID           string
	PickupAddress      string
	DeliveryAddress    string
	PackageDescription string
	Quantity           string
	Weight             string
	DeliveryDate       string
	EstimatedDelivery  string
	Notes              string
}

func (f createDeliveryForm) toModel() (model.Delivery, error) {
	customerID, err := strconv.ParseInt(f.CustomerID, 10, 64)
	if err != nil || customerID <= 0 {
		return model.Delivery{}, fmt.Errorf("invalid customer ID")
	}

	driverID, err := strconv.ParseInt(f.DriverID, 10, 64)
	if err != nil || driverID <= 0 {
		return model.Delivery{}, fmt.Errorf("invalid driver ID")
	}

	quantity, err := strconv.Atoi(f.Quantity)
	if err != nil || quantity <= 0 {
		return model.Delivery{}, fmt.Errorf("quantity must be greater than zero")
	}

	weight, err := strconv.ParseFloat(f.Weight, 64)
	if err != nil || weight < 0 {
		return model.Delivery{}, fmt.Errorf("weight must be zero or greater")
	}

	deliveryDate, err := time.Parse("2006-01-02", f.DeliveryDate)
	if err != nil {
		return model.Delivery{}, fmt.Errorf("invalid delivery date")
	}

	estimatedDelivery, err := time.Parse("2006-01-02T15:04", f.EstimatedDelivery)
	if err != nil {
		return model.Delivery{}, fmt.Errorf("invalid estimated delivery time")
	}

	return model.Delivery{
		CustomerID:         customerID,
		DriverID:           driverID,
		PickupAddress:      strings.TrimSpace(f.PickupAddress),
		DeliveryAddress:    strings.TrimSpace(f.DeliveryAddress),
		PackageDescription: strings.TrimSpace(f.PackageDescription),
		Quantity:           quantity,
		Weight:             weight,
		DeliveryDate:       deliveryDate,
		EstimatedDelivery:  estimatedDelivery,
		Status:             "pending",
		Notes:              strings.TrimSpace(f.Notes),
	}, nil
}
