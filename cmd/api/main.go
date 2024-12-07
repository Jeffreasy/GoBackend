package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // nodig voor "file://migrations"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Jeffreasy/GoBackend/configs"
	"github.com/Jeffreasy/GoBackend/internal/auth"
	"github.com/Jeffreasy/GoBackend/internal/contact"
	"github.com/Jeffreasy/GoBackend/internal/database"
	"github.com/Jeffreasy/GoBackend/internal/email"
	"github.com/Jeffreasy/GoBackend/pkg/validator"
)

func main() {
	// Load .env variabelen (indien aanwezig)
	if err := godotenv.Load(); err != nil {
		log.Println("Geen .env bestand gevonden, gebruik standaard configuraties")
	}

	// Laad de configuratie uit environment variabelen of defaults
	cfg := configs.LoadConfig()

	// Maak verbinding met de database
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Kan niet verbinden met de database: %v", err)
	}
	defer db.Close()

	// Voer migraties uit voordat de server start
	if err := runMigrations(cfg, db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Initialiseert een validator
	v := validator.NewValidator()

	// Initialiseert de email service
	emailService := email.NewService(cfg)

	// Auth service & handler
	authService := auth.NewService(db, cfg)
	authHandler := auth.NewHandler(authService, v, emailService)

	// Contact service & handler
	contactService := contact.NewService(db)
	contactHandler := contact.NewHandler(contactService, v, emailService)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Openbare endpoints
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/contact", contactHandler.SubmitContact)

	// Beschermde routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware(cfg.JWTSecret))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Dit is een beschermde route, alleen toegankelijk met geldige token."))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connectie
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
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	// Controleer verschillende mogelijke locaties voor de migraties
	migrationPaths := []string{
		"file://migrations",      // Voor lokale development
		"file:///app/migrations", // Voor Docker container
		"./migrations",           // Alternatieve lokale path
		"../migrations",          // Nog een alternatief
	}

	var lastErr error
	for _, path := range migrationPaths {
		m, err := migrate.NewWithDatabaseInstance(
			path,
			cfg.DBName,
			driver,
		)
		if err != nil {
			lastErr = err
			log.Printf("Kon migraties niet laden van %s: %v", path, err)
			continue
		}

		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply")
				return nil
			}
			lastErr = err
			log.Printf("Migratie mislukt voor %s: %v", path, err)
			continue
		}

		log.Printf("Migraties succesvol uitgevoerd vanaf %s", path)
		return nil
	}

	return fmt.Errorf("alle migratie pogingen mislukt, laatste fout: %w", lastErr)
}
