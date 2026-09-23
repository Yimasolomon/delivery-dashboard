package handler

import (
	"html/template"
	"net/http"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/service"
)

type ReportHandler struct {
	service   *service.ReportService
	templates *template.Template
}

func NewReportHandler(
	service *service.ReportService,
	templates *template.Template,
) *ReportHandler {
	return &ReportHandler{
		service:   service,
		templates: templates,
	}
}

func (h *ReportHandler) Reports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := model.ReportFilter{
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	summary, err := h.service.GetSummary(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to load report summary", http.StatusInternalServerError)
		return
	}

	driverPerformance, err := h.service.GetDriverPerformance(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to load driver performance", http.StatusInternalServerError)
		return
	}

	customerActivity, err := h.service.GetCustomerActivity(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to load customer activity", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title             string
		Filter            model.ReportFilter
		Summary           model.ReportSummary
		DriverPerformance []model.DriverPerformance
		CustomerActivity  []model.CustomerActivity
	}{
		Title:             "Reports",
		Filter:            filter,
		Summary:           summary,
		DriverPerformance: driverPerformance,
		CustomerActivity:  customerActivity,
	}

	if err := h.templates.ExecuteTemplate(w, "reports.html", data); err != nil {
		http.Error(w, "Failed to render reports page", http.StatusInternalServerError)
		return
	}
}
