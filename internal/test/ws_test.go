package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const readDeadline = 5 * time.Second

func listen(t *testing.T, app *fiber.App) {
	err := app.Listen(cfg.AppHost + cfg.AppPort, listenCfg)
	require.NoError(t, err)
}

func connect(t *testing.T, ticket string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+cfg.AppHost+cfg.AppPort+"/ws?ticket="+ticket, nil)

	require.NoError(t, err)
	conn.SetReadDeadline(time.Now().Add(readDeadline))
	return conn
}

func getTicket(t *testing.T, app *fiber.App, accessCookie *http.Cookie) string {
	req := httptest.NewRequest(http.MethodGet, "/api/ws/ticket", nil)
	req.AddCookie(accessCookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)

	var result response.Response[response.TicketResponse]

	err = json.NewDecoder(resp.Body).Decode(&result)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotZero(t, result.Data.Ticket)

	return result.Data.Ticket
}

func TestWS_Connect(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	ticket := getTicket(t, app, accessCookie1)

	go listen(t, app)
	defer app.Shutdown()

	conn := connect(t, ticket)
	conn.Close()
}
