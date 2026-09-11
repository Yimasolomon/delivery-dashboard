package model

import "time"

type StatusHistory struct {
	ID         int64
	DeliveryID int64
	Status     string
	Note       string
	CreatedAt  time.Time
}
