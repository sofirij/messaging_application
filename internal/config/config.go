package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	AppHost              string
	AppPort              string
	DbURL                string
	IsDevelopment        bool
	IsProduction         bool
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	BcryptCost           int
	UploadDir            string
	UploadBaseURLPath    string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found, using system env vars")
	}

	env := getEnv("APP_ENV", "development")

	return &Config{
		AppEnv:               env,
		AppHost:              getEnv("APP_HOST", "192.168.0.11"),
		AppPort:              getEnv("APP_PORT", ":3000"),
		DbURL:                getEnv("DB_URL", ""),
		IsDevelopment:        env != "production",
		IsProduction:         env == "production",
		JWTSecret:            getEnv("JWT_SECRET", "e4365e8cabd16c1df0f0ce8c0cb6d6095353f0a3d835c57ed1a97b42bc733d84"),
		AccessTokenDuration:  getDurationEnv("ACCESS_TOKEN_DURATION", time.Minute*15),
		RefreshTokenDuration: getDurationEnv("REFRESH_TOKEN_DURATION", time.Hour*24*7),
		BcryptCost:           getIntEnv("BCRYPT_COST", 12),
		UploadDir:            getEnv("UPLOAD_DIR", "./uploads"),
		UploadBaseURLPath:    getEnv("UPLOAD_BASE_URL_PATH", "/uploads"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		log.Println("Invalid duration fallback to default duration")
		return fallback
	}

	return d
}

func getIntEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		log.Println("Invalid integer fallback to default interger")
		return fallback
	}

	return n
}
