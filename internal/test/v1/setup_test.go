package v1

import (
	"context"
	"testing"
	"time"

	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/repository"
	"app/internal/router"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool         *pgxpool.Pool
	cfg            *config.Config
	messageService service.MessageService
	testCfg        = fiber.TestConfig{
		Timeout: 100 * time.Second,
	}
	listenCfg = fiber.ListenConfig{
		DisableStartupMessage: true,
	}
)

const (
	envFilePath  = "../../../.env.test"
	pingInterval = 2 * time.Second
	pongTimeout  = 2 * time.Second
)

func setupApp(t *testing.T) *fiber.App {
	t.Helper()

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})

	// setup middleware
	app.Use(middleware.CORS([]string{cfg.FrontendURL}))
	app.Use(middleware.Compress())
	app.Use(middleware.Logger())

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
	conversationService := service.NewConversationService(conversationRepo, userRepo, messageRepo, hubService)
	uploadService := service.NewUploadService(cfg)
	ticketService := service.NewTicketService(storage)

	go hubService.Run()

	// setup handlers
	userHandler := handler.NewUserHandler(userService, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	conversationHandler := handler.NewConversationHandler(conversationService, hubService)
	messageHandler := handler.NewMessageHandler(messageService, hubService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	wsHandler := handler.NewWSHandler(hubService, ticketService, pingInterval, pongTimeout)

	// setup routers
	v1 := app.Group("/api/v1")
	router.RegisterUserRoutes(v1, userHandler, cfg, userService)
	router.RegisterConversationRoutes(v1, conversationHandler, messageHandler, cfg, userService)
	router.RegisterMessageRoutes(v1, messageHandler, cfg, userService)
	router.RegisterWSRoutes(v1, wsHandler, ticketService, cfg, userService)
	router.RegisterUploadRoutes(v1, uploadHandler, cfg, userService)
	router.RegisterAuthRoutes(v1, userHandler)

	return app
}

func truncateTables(t *testing.T) {
	t.Helper()

	query := `
		TRUNCATE TABLE messages, message_attachments, conversation_members,
		conversations, refresh_tokens, users RESTART IDENTITY CASCADE
	`

	var err error
	for range 3 { // to avoid any deadlocks
		_, err = dbPool.Exec(context.Background(), query)
		if err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("failed to truncate tables: %v", err)
}

func TestMain(m *testing.M) {
	cfg = config.Load(envFilePath)

	var err error

	dbPool, err = pgxpool.New(context.Background(), cfg.DbURL)

	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	m.Run()
}
