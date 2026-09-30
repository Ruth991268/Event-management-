package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds environment variables
type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	HasuraAdminSecret string
	HasuraEndpoint    string
}

// LoadConfig loads environment variables from .env or system env
func LoadConfig() *Config {
	// Try loading .env, fallback to system env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/eventify")

	// Automatically append sslmode=disable if missing
	if !strings.Contains(dbURL, "sslmode") {
		if strings.Contains(dbURL, "?") {
			dbURL += "&sslmode=disable"
		} else {
			dbURL += "?sslmode=disable"
		}
	}

	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       dbURL,
		JWTSecret:         getEnv("JWT_SECRET", "b7f3c29e8a4d1f92e7b1a4c3f8d9e2c5"),
		HasuraAdminSecret: getEnv("HASURA_ADMIN_SECRET", "myadminsecretkey"),
		HasuraEndpoint:    getEnv("HASURA_ENDPOINT", "http://localhost:8081"),
	}

	// Critical env validation
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if cfg.HasuraAdminSecret == "" {
		log.Fatal("HASURA_ADMIN_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return fallback
}
