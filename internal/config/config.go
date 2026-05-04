package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppHost       string
	AppPort       string
	DbURL         string
	IsDevelopment bool
	IsProduction  bool
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found, using system env vars")
	}

	env := getEnv("APP_HOST", "development")

	return &Config{
		AppEnv:        env,
		AppHost:       getEnv("APP_HOST", "192.168.0.11"),
		AppPort:       getEnv("APP_PORT", ":3000"),
		DbURL:         getEnv("DB_URL", ""),
		IsDevelopment: env == "development",
		IsProduction:  env == "production",
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
