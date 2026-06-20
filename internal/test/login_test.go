package test

import (
	"app/internal/model/request"
	"app/internal/model/response"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func login(t *testing.T, app *fiber.App, bodyStruct request.UserAuthRequest) (*http.Response, *http.Cookie, *http.Cookie) {
	body, err := json.Marshal(bodyStruct)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)
	require.Equal(t, bodyStruct.Username, result.Data.Username)

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

func TestLogin(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	login(t, app, bodyStruct)
}

func TestLogin_WrongCredentials(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)

	bodyStruct = request.UserAuthRequest{
		Username: "testuser",
		Password: "wrongpassword", // use wrong password
	}
	
	body, err := json.Marshal(bodyStruct)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
