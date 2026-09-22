package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"delivery-dashboard/internal/database"
	"delivery-dashboard/internal/handler"
	"delivery-dashboard/internal/repository"
	"delivery-dashboard/internal/service"
)

type Application struct {
	DB        *sql.DB
	Templates *template.Template
}

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./data/deliveries.db")

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	deliveryRepository := repository.NewDeliveryRepository(db)
	deliveryService := service.NewDeliveryService(deliveryRepository)

	customerRepository := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepository)

	driverRepository := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(driverRepository)

	dashboardRepository := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepository)

	deliveryHandler := handler.NewDeliveryHandler(
		deliveryService,
		customerService,
		driverService,
		templates,
	)

	driverHandler := handler.NewDriverHandler(
		driverService,
		templates,
	)

	dashboardHandler := handler.NewDashboardHandler(
		dashboardService,
		templates,
	)

	customerHandler := handler.NewCustomerHandler(
		customerService,
		templates,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", dashboardHandler.Dashboard)

	mux.HandleFunc("/customers", customerHandler.List)
	mux.HandleFunc("/customers/create", customerHandler.CreateRoute)
	mux.HandleFunc("/customers/", customerHandler.EditRoute)

	mux.HandleFunc("/drivers", driverHandler.List)
	mux.HandleFunc("/drivers/create", driverHandler.CreateRoute)
	mux.HandleFunc("/drivers/", driverHandler.EditRoute)

	mux.HandleFunc("/deliveries", deliveryHandler.List)
	mux.HandleFunc("/deliveries/create", deliveryHandler.CreateRoute)
	mux.HandleFunc("/deliveries/", deliveryHandler.Details)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("delivery dashboard starting on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Delivery Dashboard",
	}

	if err := app.Templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("templqte error: %v", err)
		http.Error(w, "Intenal server error", http.StatusInternalServerError)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
