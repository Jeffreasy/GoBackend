package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // voor file migraties
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Jeffreasy/GoBackend/configs"
	"github.com/Jeffreasy/GoBackend/internal/auth"
	"github.com/Jeffreasy/GoBackend/internal/contact"
	"github.com/Jeffreasy/GoBackend/internal/database"
	"github.com/Jeffreasy/GoBackend/internal/email"
	"github.com/Jeffreasy/GoBackend/internal/registration"
	"github.com/Jeffreasy/GoBackend/pkg/validator"
)

func main() {
	// Probeer .env te laden
	if err := godotenv.Load(); err != nil {
		log.Println("Geen .env bestand gevonden, gebruik standaard configuraties")
	}

	// Config laden
	cfg := configs.LoadConfig()

	// Verbinding maken met de database
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Kan niet verbinden met de database: %v", err)
	}
	defer db.Close()

	// Migraties uitvoeren
	if err := runMigrations(cfg, db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Validator initialiseren
	v := validator.NewValidator()

	// Email service
	emailService := email.NewService(cfg)

	// Auth service & handler
	authService := auth.NewService(db, cfg)
	authHandler := auth.NewHandler(authService, v, emailService)

	// Contact service & handler
	contactService := contact.NewService(db)
	contactHandler := contact.NewHandler(contactService, v, emailService)

	// Registration service & handler
	registrationService := registration.NewService(db)
	registrationHandler := registration.NewHandler(registrationService, v)

	// Router initialiseren
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Openbare endpoints
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/contact", contactHandler.SubmitContact)
	r.Post("/registrations", registrationHandler.CreateRegistration)
	r.Get("/registrations", registrationHandler.GetRegistrations)

	// Beschermde routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware(cfg.JWTSecret))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Dit is een beschermde route, alleen toegankelijk met een geldige token."))
		})
	})

	// Healthcheck endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server gestart op %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server fout: %v", err)
	}
}

func runMigrations(cfg *configs.Config, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Error creating postgres driver: %v", err)
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	// Debug info
	wd, _ := os.Getwd()
	log.Printf("Working directory: %s", wd)

	// List all migration files
	files, err := os.ReadDir("/app/migrations")
	if err != nil {
		log.Printf("Error reading migrations directory: %v", err)
	} else {
		log.Printf("Found migration files:")
		for _, f := range files {
			log.Printf("- %s", f.Name())
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		log.Printf("Error creating migrate instance: %v", err)
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration error: %v", err)
		return fmt.Errorf("migration failed: %w", err)
	}

	// Verify tables exist
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'contacts'
		);
	`).Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking if contacts table exists: %v", err)
	} else {
		log.Printf("Contacts table exists: %v", tableExists)
	}

	log.Println("Migrations completed successfully!")
	return nil
}
