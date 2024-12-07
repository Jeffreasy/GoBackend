package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Gebruik het absolute pad naar de migrations directory
	migrationPath := "file://migrations"
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		// Als de migrations directory niet in de root staat, probeer dan het pad relatief aan de binary
		migrationPath = "file:///app/migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		cfg.DBName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations completed successfully!")
	return nil
}
