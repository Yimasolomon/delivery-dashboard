package model

import "time"

type Customer struct {
	ID        int64
	Name      string
	Phone     string
	Email     string
	Address   string
	CreatedAt time.Time
}
