package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh_Success(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	username := "testuser"
	password := "password123"

	register(t, app, username, password)
	_, accessCookie, refreshCookie := login(t, app, username, password)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	var newAccessCookie *http.Cookie
	var newRefreshCookie *http.Cookie

	for _, cookie := range resp.Cookies() {
		switch (cookie.Name) {
		case "access_token":
			newAccessCookie = cookie
		case "refresh_token":
			newRefreshCookie = cookie
		}
	}

	assert.NotNil(t, newAccessCookie)
	assert.NotNil(t, newRefreshCookie)
	
	resp2, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
}