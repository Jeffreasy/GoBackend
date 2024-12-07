package database

import (
	"database/sql"
	"fmt"

	"github.com/Jeffreasy/GoBackend/configs"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *configs.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
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
