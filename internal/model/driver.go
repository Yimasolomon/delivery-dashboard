package model

import "time"

type Driver struct {
	ID                  int64
	Name                string
	Phone               string
	Vehicle             string
	VehicleRegistration string
	Status              string
	CreatedAt           time.Time
}
