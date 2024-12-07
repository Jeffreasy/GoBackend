package database

import (
    "database/sql"
    "log"
    "os"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
    var err error
    DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }

    log.Println("Successfully connected to database")
}