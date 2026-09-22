package model

type ReportSummary struct {
	Total          int
	Pending        int
	PickedUp       int
	InTransit      int
	OutForDelivery int
	Delivered      int
	Failed         int
	Cancelled      int
	Delayed        int
}

type DriverPerformance struct {
	DriverID   int64
	DriverName string
	Total      int
	Delivered  int
	Failed     int
	Active     int
}

type CustomerActivity struct {
	CustomerID   int64
	CustomerName string
	Total        int
}

type ReportFilter struct {
	StartDate string
	EndDate   string
}
