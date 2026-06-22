package test

import (
	"app/internal/model/request"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	_, accessCookie, refreshCookie := login(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	for _, cookie := range resp.Cookies() {
		assert.NotEqual(t, cookie.Value, accessCookie.Value)
		assert.NotEqual(t, cookie.Value, refreshCookie.Value)
	}

	// refresh token should be invalid now and fail
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req2.Header.Set("Contetn-Type", "application/json")
	req2.AddCookie(accessCookie)
	req2.AddCookie(refreshCookie)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
}
