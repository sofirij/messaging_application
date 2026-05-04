package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/config"
	"app/internal/middleware"
	"app/internal/router"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	log.Printf("starting in %s mode\n", cfg.AppEnv)

	app := fiber.New(fiber.Config{
		AppName:            "MyApp",
		EnableIPValidation: true,
		IdleTimeout:        1 * time.Minute,
		ReadTimeout:        1 * time.Minute,
		WriteTimeout:       1 * time.Minute,
	})

	middleware.Setup(app, cfg.AppHost)

	router.Setup(app)

	// app listen
	go func() {
		log.Printf("running on port %s\n", cfg.AppPort)

		if err := app.Listen(cfg.AppHost+cfg.AppPort, fiber.ListenConfig{
			CertFile:        "./cert.pem",
			CertKeyFile:     "./key.pem",
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
