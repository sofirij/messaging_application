package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/repository"
	"app/internal/router"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/memory/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	assetMaxAge        = 365 * 24 * time.Hour
	assetCacheDuration = 1 * time.Minute
	uploadURLPath      = "/uploads"
	envFilePath        = ".env"
	pingInterval       = 15 * time.Second
	pongTimeout        = pingInterval + 15*time.Second
)

func main() {
	cfg := config.Load(envFilePath)

	log.Printf("starting in %s mode\n", cfg.AppEnv)

	app := fiber.New(fiber.Config{
		AppName:      "MyApp",
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  100 * time.Second,
		BodyLimit:    100 * 1024 * 1024,
	})

	// setup middleware
	app.Use(middleware.CORS(cfg.AppHost))
	app.Use(middleware.Compress())
	app.Use(middleware.Logger())

	// serve static files
	app.Use(uploadURLPath, static.New(cfg.UploadDir, static.Config{
		Browse:        false,
		Compress:      true,
		MaxAge:        int(assetMaxAge.Seconds()),
		CacheDuration: assetCacheDuration,
	}))

	// initialize db connection
	// todo configure db connection
	dbPool, err := pgxpool.New(context.Background(), cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)

	}

	// initialize storage interface
	storage := memory.New(memory.Config{
		GCInterval: 1 * time.Second,
	})

	// setup repositories
	userRepo := repository.NewUserRepository(dbPool)
	messageRepo := repository.NewMessageRepository(dbPool)
	conversationRepo := repository.NewConversationRepository(dbPool)

	// setup services
	messageService := service.NewMessageService(messageRepo, conversationRepo)
	hubService := service.NewHubService(conversationRepo, userRepo, messageService)
	userService := service.NewUserService(userRepo, hubService, cfg)
	conversationService := service.NewConversationService(conversationRepo, userRepo, cfg)
	uploadService := service.NewUploadService(cfg)
	ticketService := service.NewTicketService(storage)

	// setup handlers
	userHandler := handler.NewUserHandler(userService, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	conversationHandler := handler.NewConversationHandler(conversationService, hubService)
	messageHandler := handler.NewMessageHandler(messageService, hubService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	wsHandler := handler.NewWSHandler(hubService, ticketService, pingInterval, pongTimeout)

	// setup routers
	v1 := app.Group("/api/v1")
	router.RegisterUserRoutes(v1, userHandler, cfg.JWTSecret)
	router.RegisterConversationRoutes(v1, conversationHandler, messageHandler, cfg.JWTSecret)
	router.RegisterMessageRoutes(v1, messageHandler, cfg.JWTSecret)
	router.RegisterWSRoutes(v1, wsHandler, ticketService, cfg.JWTSecret)
	router.RegisterUploadRoutes(v1, uploadHandler, cfg.JWTSecret)
	router.RegisterAuthRoutes(v1, userHandler, cfg.JWTSecret)

	// start hubservice
	go hubService.Run()

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

	// shut down hub
	hubService.Stop()
}
