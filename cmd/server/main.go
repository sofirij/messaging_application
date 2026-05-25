package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/config"
	"app/internal/middleware"
	
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const (
	assetMaxAge        = 365 * 24 * time.Hour
	assetCacheDuration = 1 * time.Minute
	uploadURLPath      = "/uploads"
)

func main() {
	cfg := config.Load()

	log.Printf("starting in %s mode\n", cfg.AppEnv)

	app := fiber.New(fiber.Config{
		AppName:                 "MyApp",
		IdleTimeout:             30 * time.Second,
		WriteTimeout:            5 * time.Second,
	})

	// setup middleware
	app.Use(middleware.CORS(cfg.AppHost))
	app.Use(middleware.Compress())
	app.Use(middleware.Logger())
	app.Use(uploadURLPath, static.New(cfg.UploadDir, static.Config{
		Browse:        false,
		Compress:      true,
		MaxAge:        int(assetMaxAge.Seconds()),
		CacheDuration: assetCacheDuration,
	}))

	// app listen
	go func() {
		log.Printf("running on port %s\n", cfg.AppPort)

		if err := app.Listen(cfg.AppHost+cfg.AppPort, fiber.ListenConfig{
			ShutdownTimeout: 1 * time.Minute,
		}); err != nil {
			log.Fatalf("server error, %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("server shutdown failed, %v\n", err)
	}
	log.Println("server shutdown gracefully")
}
