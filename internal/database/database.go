package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeffreasy/GoBackend/configs"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *configs.Config) (*sql.DB, error) {
	// Bepaal SSL mode op basis van environment
	sslmode := "require" // Default voor productie (Render)

	// Voor lokale development, gebruik disable
	if cfg.DBHost == "localhost" || cfg.DBHost == "127.0.0.1" || cfg.DBHost == "postgres" {
		sslmode = "disable"
	}

	// Bouw de connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslmode,
	)

	// Open de database connectie
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test de connectie
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Configureer connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
