package model

type DashboardStats struct {
	Total          int
	Pending        int
	InTransit      int
	OutForDelivery int
	Delivered      int
	Failed         int
	Delayed        int
}

type StatusCount struct {
	Status string
	Count  int
}

type DeliveryTrend struct {
	Date  string
	Count int
}
