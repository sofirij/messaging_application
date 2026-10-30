package config

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strconv"
	"strings"
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
	
	publicKeyPath := getEnv("JWT_PUBLIC_KEY_PATH", "")
	publicKey, err := loadPublicKey(publicKeyPath)
	
	if err != nil {
		log.Fatal(err)
	}

	privateKeyPath := getEnv("JWT_PRIVATE_KEY_PATH", "")
	privateKey, err := loadPrivateKey(privateKeyPath)

	if err != nil {
		log.Fatal(err)
	}

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

func loadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	
	if err != nil {
		return nil, err
	}

	return key.(ed25519.PrivateKey), nil
}

func loadPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(ed25519.PublicKey), nil
}
