// config.go

// Deze file laadt de configuratievariabelen uit de omgevingsvariabelen (of gebruikt defaults als deze niet gezet zijn).
// Zo kunnen we eenvoudig wisselen tussen development, staging, en productie omgevingen, zonder code te veranderen.

package configs

import (
	"os"
)

type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "mydb"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", "user@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "emailpassword"),
		FromEmail:    getEnv("FROM_EMAIL", "no-reply@example.com"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
