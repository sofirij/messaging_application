package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getUserID(t *testing.T, app *fiber.App, accessCookie *http.Cookie) int {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var user response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&user)

	require.NoError(t, err)
	return user.Data.ID
}

func TestUser_Search(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	user2 := request.UserAuthRequest{
		Username: "testuser2",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)

	// user1 searches for user2
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/search?q=%s", user2.Username), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie1)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result response.Response[[]response.UserResponse]

	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)
	assert.Equal(t, 1, len(result.Data))

	assert.Equal(t, "testuser2", result.Data[0].Username)
}
