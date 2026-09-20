package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"delivery-dashboard/internal/model"
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

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page < 1 {
		page = 1
	}

	pageSize := 10

	driverID, _ := strconv.ParseInt(
		r.URL.Query().Get("driver"),
		10,
		64,
	)

	filter := model.DeliveryFilter{
		Search:    strings.TrimSpace(r.URL.Query().Get("search")),
		Status:    strings.TrimSpace(r.URL.Query().Get("status")),
		DriverID:  driverID,
		DateFrom:  strings.TrimSpace(r.URL.Query().Get("date_from")),
		DateTo:    strings.TrimSpace(r.URL.Query().Get("date_to")),
		SortBy:    strings.TrimSpace(r.URL.Query().Get("sort")),
		SortOrder: strings.TrimSpace(r.URL.Query().Get("order")),
		Page:      page,
		PageSize:  pageSize,
	}

	deliveries, err := h.service.ListDeliveries(
		r.Context(),
		filter,
	)

	if err != nil {
		log.Printf("failed to load deliveries: %v", err)
		http.Error(w, "Failed to load deliveries", http.StatusInternalServerError)
		return
	}

	totalDeliveries, err := h.service.CountDeliveries(
		r.Context(),
		filter,
	)

	if err != nil {
		log.Printf("failed to count deliveries: %v", err)
		http.Error(w, "Failed to count deliveries", http.StatusInternalServerError)
		return
	}

	totalPages := (totalDeliveries + pageSize - 1) / pageSize

	drivers, err := h.driverService.ListDrivers(r.Context())
	if err != nil {
		log.Printf("failed to load drivers: %v", err)
		http.Error(w, "Failed to load drivers", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title            string
		Deliveries       interface{}
		Search           string
		Status           string
		DriverID         int64
		DateFrom         string
		DateTo           string
		SortBy           string
		SortOrder        string
		Page             int
		PageSize         int
		TotalDeliveries  int
		TotalPages       int
		HasPrevious      bool
		HasNext          bool
		PreviousPage     int
		NextPage         int
		DeliveryStatuses []string
		Drivers          interface{}
	}{
		Title:            "Deliveries",
		Deliveries:       deliveries,
		Search:           filter.Search,
		Status:           filter.Status,
		DriverID:         filter.DriverID,
		DateFrom:         filter.DateFrom,
		DateTo:           filter.DateTo,
		SortBy:           filter.SortBy,
		SortOrder:        filter.SortOrder,
		Page:             page,
		PageSize:         pageSize,
		TotalDeliveries:  totalDeliveries,
		TotalPages:       totalPages,
		HasPrevious:      page > 1,
		HasNext:          page < totalPages,
		PreviousPage:     page - 1,
		NextPage:         page + 1,
		DeliveryStatuses: model.DeliveryStatuses,
		Drivers:          drivers,
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

func (h *DeliveryHandler) Details(w http.ResponseWriter, r *http.Request) {

	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/deliveries/"), "/")

	if strings.HasSuffix(path, "/delete") {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		h.Delete(w, r)
		return
	}

	if strings.HasSuffix(path, "/edit") {
		switch r.Method {
		case http.MethodGet:
			h.EditForm(w, r)
		case http.MethodPost:
			h.Update(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if strings.HasSuffix(path, "/status") {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		h.UpdateStatus(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idString := strings.TrimPrefix(r.URL.Path, "/deliveries/")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	delivery, err := h.service.GetDelivery(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	customer, err := h.customerService.GetCustomer(
		r.Context(),
		delivery.CustomerID,
	)
	if err != nil {
		http.Error(w, "Failed to load customer", http.StatusInternalServerError)
		return
	}

	driver, err := h.driverService.GetDriver(
		r.Context(),
		delivery.DriverID,
	)
	if err != nil {
		http.Error(w, "Failed to load driver", http.StatusInternalServerError)
		return
	}

	history, err := h.service.GetStatusHistory(
		r.Context(),
		delivery.ID,
	)
	if err != nil {
		http.Error(w, "Failed to load status history", http.StatusInternalServerError)
		return
	}

	allowedStatuses := map[string][]string{
		"pending": {
			"picked_up",
			"cancelled",
		},
		"picked_up": {
			"in_transit",
			"failed",
		},
		"in_transit": {
			"out_for_delivery",
			"failed",
		},
		"out_for_delivery": {
			"delivered",
			"failed",
		},
	}

	data := struct {
		Title           string
		Delivery        interface{}
		Customer        interface{}
		Driver          interface{}
		History         interface{}
		AllowedStatuses []string
	}{
		Title:           "Delivery Details",
		Delivery:        delivery,
		Customer:        customer,
		Driver:          driver,
		History:         history,
		AllowedStatuses: allowedStatuses[delivery.Status],
	}

	if err := h.template.ExecuteTemplate(
		w,
		"delivery-details.html",
		data,
	); err != nil {
		log.Printf("failed to render delivery details: %v", err)
		http.Error(w, "Failed to render delivery details", http.StatusInternalServerError)
		return
	}
}

func (h *DeliveryHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idString := strings.TrimPrefix(r.URL.Path, "/deliveries/")
	idString = strings.TrimSuffix(idString, "/edit")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	delivery, err := h.service.GetDelivery(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
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
		Delivery  model.Delivery
		Customers []model.Customer
		Drivers   []model.Driver
	}{
		Title:     "Edit Delivery",
		Delivery:  delivery,
		Customers: customers,
		Drivers:   drivers,
	}

	if err := h.template.ExecuteTemplate(
		w,
		"delivery-edit.html",
		data,
	); err != nil {
		log.Printf("failed to render delivery edit form: %v", err)
		http.Error(w, "Failed to render delivery edit form", http.StatusInternalServerError)
		return
	}
}

func (h *DeliveryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/deliveries/")
	idString = strings.TrimSuffix(idString, "/edit")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	existing, err := h.service.GetDelivery(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	form := updateDeliveryForm{
		CustomerID:         r.FormValue("customer_id"),
		DriverID:           r.FormValue("driver_id"),
		PickupAddress:      r.FormValue("pickup_address"),
		DeliveryAddress:    r.FormValue("delivery_address"),
		PackageDescription: r.FormValue("package_description"),
		Quantity:           r.FormValue("quantity"),
		Weight:             r.FormValue("weight"),
		DeliveryDate:       r.FormValue("delivery_date"),
		EstimatedDelivery:  r.FormValue("estimated_delivery"),
		Notes:              r.FormValue("notes"),
	}

	delivery, err := form.toModel(existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateDelivery(r.Context(), delivery); err != nil {
		log.Printf("failed to update delivery %d: %v", id, err)
		http.Error(w, "Failed to update delivery", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/deliveries/%d", id),
		http.StatusSeeOther,
	)
}

func (h *DeliveryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/deliveries/")
	idString = strings.TrimSuffix(idString, "/delete")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.service.DeleteDelivery(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		log.Printf("failed to delete delivery %d: %v", id, err)
		http.Error(w, "Failed to delete delivery", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/deliveries", http.StatusSeeOther)
}

func (h *DeliveryHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/deliveries/")
	idString = strings.TrimSuffix(idString, "/status")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	status := strings.TrimSpace(r.FormValue("status"))
	note := strings.TrimSpace(r.FormValue("note"))

	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(
		r.Context(),
		id,
		status,
		note,
	); err != nil {
		log.Printf("failed to update status for delivery %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/deliveries/%d", id),
		http.StatusSeeOther,
	)
}
