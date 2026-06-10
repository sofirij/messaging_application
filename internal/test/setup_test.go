package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/model/response"
	"app/internal/repository"
	"app/internal/router"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/memory/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	dbPool     *pgxpool.Pool
	cfg        *config.Config
	hubService service.HubService
)

const envFilePath = "../../.env.test"

func setupApp(t *testing.T) *fiber.App {
	t.Helper()

	app := fiber.New(fiber.Config{
		AppName:      "MyApp",
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Second,
		BodyLimit:    100 * 1024 * 1024, // 100MB
	})

	// setup middleware
	app.Use(middleware.CORS(cfg.AppHost))
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
	conversationService := service.NewConversationService(conversationRepo, userRepo, cfg)
	uploadService := service.NewUploadService(cfg)
	ticketService := service.NewTicketService(storage)

	// setup handlers
	userHandler := handler.NewUserHandler(userService, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	conversationHandler := handler.NewConversationHandler(conversationService)
	messageHandler := handler.NewMessageHandler(messageService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	wsHandler := handler.NewWSHandler(hubService, ticketService)

	// setup routers
	router.RegisterUserRoutes(app, userHandler, cfg.JWTSecret)
	router.RegisterConversationRoutes(app, conversationHandler, messageHandler, cfg.JWTSecret)
	router.RegisterMessageRoutes(app, messageHandler, cfg.JWTSecret)
	router.RegisterWSRoutes(app, wsHandler, ticketService, cfg.JWTSecret)
	router.RegisterUploadRoutes(app, uploadHandler, cfg.JWTSecret)

	return app
}

func truncateTables(t *testing.T) {
	t.Helper()

	query := `
		TRUNCATE TABLE messages, message_attachments, conversation_members,
		conversations, refresh_tokens, users RESTART IDENTITY CASCADE
	`

	_, err := dbPool.Exec(context.Background(), query)

	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func login(t *testing.T, app *fiber.App, username, password string) (*http.Response, *http.Cookie, *http.Cookie) {
	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)
	require.Equal(t, username, result.Data.Username)

	var access *http.Cookie
	var refresh *http.Cookie

	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "access_token":
			access = cookie
		case "refresh_token":
			refresh = cookie
		}
	}

	if access == nil || refresh == nil {
		t.Fatal("access_token or refresh_token cookie not found")
		return nil, nil, nil
	}

	return resp, access, refresh
}

func register(t *testing.T, app *fiber.App, username, password string) *http.Response {
	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return resp
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
