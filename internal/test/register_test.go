package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "testuser"
	password := "password123"

	register(t, app, username, password)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "testuser"
	password := "password123"

	register(t, app, username, password)

	body := fmt.Sprintf(`{"username": "%v", "password": "%v"}`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestRegister_InvalidInput(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "sj"                                    // smaller than 3 characters
	body := fmt.Sprintf(`{"username": "%v"}`, username) // missing password

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
