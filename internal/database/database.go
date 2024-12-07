package database

import (
	"database/sql"
	"fmt"

	"github.com/Jeffreasy/GoBackend/configs"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *configs.Config) (*sql.DB, error) {
	// In productie willen we SSL gebruiken, maar tijdens development niet
	sslmode := "require"
	if cfg.DBHost == "localhost" || cfg.DBHost == "127.0.0.1" || cfg.DBHost == "postgres" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
