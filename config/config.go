package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser   string
	DBPass   string
	DBHost   string
	DBPort   string
	DBName   string
	JWTSecret string
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		DBUser:    getEnv("DB_USER", "nuser"),
		DBPass:    getEnv("DB_PASS", "npass"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBName:    getEnv("DB_NAME", "ndb"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}

	if cfg.DBUser == "" {
		log.Fatal("DB_USER not set")
	}

	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("WARNING: JWT_SECRET not set or using default value. Change in production!")
	}

	return cfg
}
