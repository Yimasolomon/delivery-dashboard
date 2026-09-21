package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/service"
)

type DashboardHandler struct {
	service  *service.DashboardService
	template *template.Template
}

func NewDashboardHandler(
	service *service.DashboardService,
	template *template.Template,
) *DashboardHandler {
	return &DashboardHandler{
		service:  service,
		template: template,
	}
}

func (h *DashboardHandler) Dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		log.Printf("failed to load dashboard stats: %v", err)
		http.Error(
			w,
			"Failed to load dashboard",
			http.StatusInternalServerError,
		)
		return
	}

	statusCounts, err := h.service.GetStatusCounts(r.Context())
	if err != nil {
		log.Printf("failed to load dashboard status counts: %v", err)
		http.Error(
			w,
			"Failed to load dashboard",
			http.StatusInternalServerError,
		)
		return
	}

	deliveryTrends, err := h.service.GetDeliveryTrends(r.Context())
	if err != nil {
		log.Printf("failed to load dashboard delivery trends: %v", err)
		http.Error(
			w,
			"Failed to load dashboard",
			http.StatusInternalServerError,
		)
		return
	}

	statusCountsJSON, err := json.Marshal(statusCounts)
	if err != nil {
		log.Printf("failed to encode status counts: %v", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	deliveryTrendsJSON, err := json.Marshal(deliveryTrends)
	if err != nil {
		log.Printf("failed to encode delivery trends: %v", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title              string
		Stats              model.DashboardStats
		StatusCountsJSON   template.JS
		DeliveryTrendsJSON template.JS
		StatusCounts       []model.StatusCount
		DeliveryTrends     []model.DeliveryTrend
	}{
		Title:              "Delivery Dashboard",
		Stats:              stats,
		StatusCounts:       statusCounts,
		DeliveryTrends:     deliveryTrends,
		StatusCountsJSON:   template.JS(statusCountsJSON),
		DeliveryTrendsJSON: template.JS(deliveryTrendsJSON),
	}

	if err := h.template.ExecuteTemplate(
		w,
		"dashboard.html",
		data,
	); err != nil {
		log.Printf("failed to render dashboard: %v", err)
		http.Error(
			w,
			"Failed to render dashboard",
			http.StatusInternalServerError,
		)
		return
	}
}
