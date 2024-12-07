// database.go

// Deze file is verantwoordelijk voor het opzetten van een verbinding met de Postgres database.
// Door gebruik te maken van een aparte file maken we het beheer en onderhoud van de databaseverbinding eenvoudiger.

package database

import (
	"database/sql"
	"fmt"

	"dklbackendGolang/configs"

	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *configs.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
