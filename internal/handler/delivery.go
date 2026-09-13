package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"delivery-dashboard/internal/service"
)

type DeliveryHandler struct {
	service         *service.DeliveryService
	customerService *service.CustomerService
	driverService   *service.DriverService
	template        *template.Template
}

func NewDeliveryHandler(
	service *service.DeliveryService,
	customerService *service.CustomerService,
	driverService *service.DriverService,
	template *template.Template,
) *DeliveryHandler {
	return &DeliveryHandler{
		service:         service,
		customerService: customerService,
		driverService:   driverService,
		template:        template,
	}
}

func (h *DeliveryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deliveries, err := h.service.ListDeliveries(r.Context())
	if err != nil {
		http.Error(w, "Failed to load deliveries", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title      string
		Deliveries interface{}
	}{
		Title:      "Deliveries",
		Deliveries: deliveries,
	}

	if err := h.template.ExecuteTemplate(w, "deliveries.html", data); err != nil {
		log.Printf("failed to render deliveries: %v", err)
		return
	}
}

func (h *DeliveryHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customers, err := h.customerService.ListCustomers(r.Context())
	if err != nil {
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}

	drivers, err := h.driverService.ListDrivers(r.Context())
	if err != nil {
		http.Error(w, "Failed to load drivers", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Customers interface{}
		Drivers   interface{}
	}{
		Title:     "Create Delivery",
		Customers: customers,
		Drivers:   drivers,
	}

	if err := h.template.ExecuteTemplate(w, "delivery-form.html", data); err != nil {
		log.Printf("failed to render delivery form: %v", err)
		http.Error(w, "Failed to render delivery form", http.StatusInternalServerError)
		return
	}
}

func (h *DeliveryHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.CreateForm(w, r)

	case http.MethodPost:
		h.Create(w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DeliveryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	form := createDeliveryForm{
		CustomerID:         strings.TrimSpace(r.FormValue("customer_id")),
		DriverID:           strings.TrimSpace(r.FormValue("driver_id")),
		PickupAddress:      strings.TrimSpace(r.FormValue("pickup_address")),
		DeliveryAddress:    strings.TrimSpace(r.FormValue("delivery_address")),
		PackageDescription: strings.TrimSpace(r.FormValue("package_description")),
		Quantity:           strings.TrimSpace(r.FormValue("quantity")),
		Weight:             strings.TrimSpace(r.FormValue("weight")),
		DeliveryDate:       strings.TrimSpace(r.FormValue("delivery_date")),
		EstimatedDelivery:  strings.TrimSpace(r.FormValue("estimated_delivery")),
		Notes:              strings.TrimSpace(r.FormValue("notes")),
	}

	if form.CustomerID == "" ||
		form.DriverID == "" ||
		form.PickupAddress == "" ||
		form.DeliveryAddress == "" ||
		form.PackageDescription == "" ||
		form.Quantity == "" ||
		form.Weight == "" ||
		form.DeliveryDate == "" ||
		form.EstimatedDelivery == "" {
		http.Error(w, "Required fields are missing", http.StatusBadRequest)
		return
	}

	delivery, err := form.toModel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.CreateDelivery(r.Context(), delivery)
	if err != nil {
		log.Printf("failed to create delivery: %v", err)
		http.Error(w, "Failed to create delivery", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/deliveries", http.StatusSeeOther)
}
