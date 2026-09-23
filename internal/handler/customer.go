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

type CustomerHandler struct {
	service  *service.CustomerService
	template *template.Template
}

func NewCustomerHandler(
	service *service.CustomerService,
	template *template.Template,
) *CustomerHandler {
	return &CustomerHandler{
		service:  service,
		template: template,
	}
}

func (h *CustomerHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customers, err := h.service.ListCustomers(r.Context())

	if err != nil {
		log.Printf("failed to load customers: %v", err)
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Customers []model.Customer
	}{
		Title:     "Customers",
		Customers: customers,
	}

	if err := h.template.ExecuteTemplate(
		w,
		"customers.html",
		data,
	); err != nil {
		log.Printf("failed to render customers: %v", err)
		return
	}
}

func (h *CustomerHandler) CreateRoute(
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CustomerHandler) CreateForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	data := struct {
		Title string
	}{
		Title: "Create Customer",
	}

	if err := h.template.ExecuteTemplate(
		w,
		"customer-form.html",
		data,
	); err != nil {
		log.Printf("failed to render customer form: %v", err)
		http.Error(
			w,
			"Failed to render customer form",
			http.StatusInternalServerError,
		)
	}
}

func (h *CustomerHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"Invalid form submission",
			http.StatusBadRequest,
		)
		return
	}

	customer := model.Customer{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		Address: r.FormValue("address"),
	}

	if customer.Name == "" ||
		customer.Phone == "" {

		http.Error(
			w,
			"Name and phone are required",
			http.StatusBadRequest,
		)
		return
	}

	_, err := h.service.CreateCustomer(
		r.Context(),
		customer,
	)

	if err != nil {
		log.Printf(
			"failed to create customer: %v",
			err,
		)

		http.Error(
			w,
			"Failed to create customer",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/customers",
		http.StatusSeeOther,
	)
}

func (h *CustomerHandler) EditRoute(
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

func (h *CustomerHandler) EditForm(
	w http.ResponseWriter,
	r *http.Request,
) {

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/customers/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/edit",
	)

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Printf("customer id parsed: %d", id)

	customer, err := h.service.GetCustomer(
		r.Context(),
		id,
	)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title    string
		Customer model.Customer
	}{
		Title:    "Edit Customer",
		Customer: customer,
	}

	if err := h.template.ExecuteTemplate(
		w,
		"customer-edit.html",
		data,
	); err != nil {

		log.Printf(
			"failed to render customer edit: %v",
			err,
		)

		http.Error(
			w,
			"Failed to render page",
			http.StatusInternalServerError,
		)
	}
}

func (h *CustomerHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Printf("updating customer")

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/customers/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/edit",
	)

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	customer := model.Customer{
		ID: id,
		Name: strings.TrimSpace(
			r.FormValue("name"),
		),
		Phone: strings.TrimSpace(
			r.FormValue("phone"),
		),
		Email: strings.TrimSpace(
			r.FormValue("email"),
		),
		Address: strings.TrimSpace(
			r.FormValue("address"),
		),
	}

	err = h.service.UpdateCustomer(
		r.Context(),
		customer,
	)

	if err != nil {
		log.Printf(
			"failed to update customer: %v",
			err,
		)

		http.Error(
			w,
			"Failed to update customer",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/customers",
		http.StatusSeeOther,
	)
}

func (h *CustomerHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	
	idString := strings.TrimPrefix(
		r.URL.Path,
		"/customers/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/delete",
	)

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)

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
			"failed to check customer deliveries: %v",
			err,
		)

		http.Error(
			w,
			"Failed to check customer deliveries",
			http.StatusInternalServerError,
		)

		return
	}

	if hasDeliveries {
		http.Error(
			w,
			"Cannot delete customer with existing deliveries",
			http.StatusConflict,
		)

		return
	}

	if err := h.service.DeleteCustomer(
		r.Context(),
		id,
	); err != nil {
		log.Printf(
			"failed to delete customer: %v",
			err,
		)

		http.Error(
			w,
			"Failed to delete customer",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/customers",
		http.StatusSeeOther,
	)
}
