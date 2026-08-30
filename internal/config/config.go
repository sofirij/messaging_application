package config

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"time"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	AppHost              string
	AppPort              string
	DbURL                string
	IsDevelopment        bool
	IsProduction         bool
	JWTPublicKey         ed25519.PublicKey
	JWTPrivateKey        ed25519.PrivateKey
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	BcryptCost           int
	UploadDir            string
	AllowedOrigins       []string
}

func Load(envFilePath string) *Config {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println("no .env file found, using system env vars")
	}

	env := getEnv("APP_ENV", "development")

	jwtPublicKey := getEnv("JWT_PUBLIC_KEY", "")
	jwtPrivateKey := getEnv("JWT_PRIVATE_KEY", "")

	publicKeyBytes, err := base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	key, err := x509.ParsePKIXPublicKey(publicKeyBytes)

	if err != nil {
		log.Fatal(err)
	}

	publicKey := key.(ed25519.PublicKey)

	privateKeyBytes, err := base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	key, err = x509.ParsePKCS8PrivateKey(privateKeyBytes)

	if err != nil {
		log.Fatal(err)
	}

	privateKey := key.(ed25519.PrivateKey)

	return &Config{
		AppEnv:               env,
		AppHost:              getEnv("APP_HOST", "localhost"),
		AppPort:              getEnv("APP_PORT", ":3000"),
		DbURL:                getEnv("DB_URL", ""),
		IsDevelopment:        env != "production",
		IsProduction:         env == "production",
		JWTPrivateKey:        privateKey,
		JWTPublicKey:         publicKey,
		AccessTokenDuration:  getDurationEnv("ACCESS_TOKEN_DURATION", time.Minute*15),
		RefreshTokenDuration: getDurationEnv("REFRESH_TOKEN_DURATION", time.Hour*24*7),
		BcryptCost:           getIntEnv("BCRYPT_COST", 12),
		UploadDir:            getEnv("UPLOAD_DIR", "./uploads"),
		AllowedOrigins:       strings.Split(getEnv("ALLOWED_ORIGINS", ""), " "),
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
