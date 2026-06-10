package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "testuser"
	password := "password123"

	register(t, app, username, password)
	login(t, app, username, password)
}

func TestLogin_WrongCredentials(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "testuser"
	password := "password123"

	register(t, app, username, password)

	body := fmt.Sprintf(`{"username": "%v", "password": "%v"}`, username, "wrongpassword") // use wrong password

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
