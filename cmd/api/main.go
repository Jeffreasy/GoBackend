// main.go

// Dit is het startpunt van je applicatie.
// In deze file worden alle benodigde configuraties, database connecties, routes en services ingeladen en wordt de HTTP-server gestart.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"dklbackendGolang/configs"
	"dklbackendGolang/internal/auth"
	"dklbackendGolang/internal/contact"
	"dklbackendGolang/internal/database"
	"dklbackendGolang/internal/email"
	"dklbackendGolang/pkg/validator"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// Migrate imports
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // nodig voor "file://migrations"
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

	// Initialiseert een validator voor input-validatie van struct velden (email, password, etc.)
	v := validator.NewValidator()

	// Initialiseert de email service. Dit object wordt gebruikt om emails te verzenden (zoals bevestigingsmails).
	emailService := email.NewService(cfg)

	// Auth service & handler: Voor inschrijven, inloggen, JWT generatie
	authService := auth.NewService(db, cfg)
	authHandler := auth.NewHandler(authService, v, emailService)

	// Contact service & handler: Voor het verwerken van contactformulieren
	contactService := contact.NewService(db)
	contactHandler := contact.NewHandler(contactService, v, emailService)

	// Router (hier: chi) om routes voor de HTTP-server in te stellen
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Logt requests in de console
	r.Use(middleware.Recoverer) // Zorgt dat de server niet crasht bij panics

	// Openbare endpoints: Registreren, inloggen en contactformulier versturen
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/contact", contactHandler.SubmitContact)

	// Beschermde routes (voorbeeld): deze routes vereisen een geldig JWT
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

	// Zorg dat je "migrations" map op dezelfde niveau staat als je main.go of in de Dockerfile wordt gekopieerd
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
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
