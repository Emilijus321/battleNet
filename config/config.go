package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	Environment string
	TMDBAPIKey  string
	TMDBBaseURL string
}

func Load() *Config {
	// Load .env file (ignore errors if file doesn't exist)
	godotenv.Load()

	port := getEnv("PORT", "3000")

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/movieapp?sslmode=disable"),
		Port:        port,
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
		TMDBAPIKey:  getEnv("TMDB_API_KEY", "454e2fb464bfab80451faca174310afc"),
		TMDBBaseURL: getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
