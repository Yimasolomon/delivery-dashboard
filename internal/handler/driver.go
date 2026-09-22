package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"delivery-dashboard/internal/model"
	"delivery-dashboard/internal/service"
)

type DriverHandler struct {
	service   *service.DriverService
	templates *template.Template
}

func NewDriverHandler(
	service *service.DriverService,
	templates *template.Template,
) *DriverHandler {
	return &DriverHandler{
		service:   service,
		templates: templates,
	}
}

func (h *DriverHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	drivers, err := h.service.ListDrivers(r.Context())
	if err != nil {
		log.Printf("failed to list drivers: %v", err)
		http.Error(
			w,
			"Failed to load drivers",
			http.StatusInternalServerError,
		)
		return
	}

	data := struct {
		Title   string
		Drivers []model.Driver
	}{
		Title:   "Drivers",
		Drivers: drivers,
	}

	if err := h.templates.ExecuteTemplate(
		w,
		"drivers.html",
		data,
	); err != nil {
		log.Printf("failed to render drivers: %v", err)
	}
}

func (h *DriverHandler) CreateRoute(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.CreateForm(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *DriverHandler) CreateForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	data := struct {
		Title string
	}{
		Title: "Create Driver",
	}

	if err := h.templates.ExecuteTemplate(
		w,
		"driver-form.html",
		data,
	); err != nil {
		log.Printf("failed to render driver form: %v", err)
	}
}

func (h *DriverHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	driver := model.Driver{
		Name:                r.FormValue("name"),
		Phone:               r.FormValue("phone"),
		Vehicle:             r.FormValue("vehicle"),
		VehicleRegistration: r.FormValue("vehicle_registration"),
		Status:              r.FormValue("status"),
	}

	_, err := h.service.CreateDriver(
		r.Context(),
		driver,
	)
	if err != nil {
		log.Printf("failed to create driver: %v", err)
		http.Error(
			w,
			"Failed to create driver",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/drivers",
		http.StatusSeeOther,
	)
}

func (h *DriverHandler) EditRoute(
	w http.ResponseWriter,
	r *http.Request,
) {
	if strings.HasSuffix(r.URL.Path, "/delete") {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		h.Delete(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.EditForm(w, r)
	case http.MethodPost:
		h.Update(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *DriverHandler) EditForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := driverIDFromPath(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	driver, err := h.service.GetDriver(
		r.Context(),
		id,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title  string
		Driver model.Driver
	}{
		Title:  "Edit Driver",
		Driver: driver,
	}

	if err := h.templates.ExecuteTemplate(
		w,
		"driver-edit.html",
		data,
	); err != nil {
		log.Printf("failed to render driver edit form: %v", err)
	}
}

func (h *DriverHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := driverIDFromPath(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	driver := model.Driver{
		ID:                  id,
		Name:                r.FormValue("name"),
		Phone:               r.FormValue("phone"),
		Vehicle:             r.FormValue("vehicle"),
		VehicleRegistration: r.FormValue("vehicle_registration"),
		Status:              r.FormValue("status"),
	}

	if err := h.service.UpdateDriver(
		r.Context(),
		driver,
	); err != nil {
		log.Printf("failed to update driver: %v", err)
		http.Error(
			w,
			"Failed to update driver",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/drivers",
		http.StatusSeeOther,
	)
}

func (h *DriverHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := driverIDFromPath(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	hasDeliveries, err := h.service.HasDeliveries(
		r.Context(),
		id,
	)

	if err != nil {
		log.Printf(
			"failed to check driver deliveries: %v",
			err,
		)

		http.Error(
			w,
			"Failed to check driver deliveries",
			http.StatusInternalServerError,
		)

		return
	}

	if hasDeliveries {
		http.Error(
			w,
			"Cannot delete driver with existing deliveries",
			http.StatusConflict,
		)

		return
	}

	if err := h.service.DeleteDriver(
		r.Context(),
		id,
	); err != nil {
		log.Printf(
			"failed to delete driver: %v",
			err,
		)

		http.Error(
			w,
			"Failed to delete driver",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/drivers",
		http.StatusSeeOther,
	)
}

func driverIDFromPath(path string) (int64, error) {
	idString := strings.TrimPrefix(
		path,
		"/drivers/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/edit",
	)

	idString = strings.TrimSuffix(
		idString,
		"/delete",
	)

	return strconv.ParseInt(
		idString,
		10,
		64,
	)
}
