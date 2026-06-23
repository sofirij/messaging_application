package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthToken_Valid(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	_, accessCookie, _ := login(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)

	var result response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, bodyStruct.Username, result.Data.Username)
}

func TestAuthToken_Invalid(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: "somethinginvalid",
	})

	resp, err := app.Test(req)

	require.NoError(t, err)

	var result response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthToken_Missing(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
