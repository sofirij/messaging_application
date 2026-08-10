package v1

import (
	"app/internal/model/request"
	"net/http"
	"net/http/httptest"
	"testing"

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

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	for _, cookie := range resp.Cookies() {
		require.NotEqual(t, cookie.Value, accessCookie.Value)
		require.NotEqual(t, cookie.Value, refreshCookie.Value)
	}

	// refresh token should be invalid now and fail
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	req2.Header.Set("Contetn-Type", "application/json")
	req2.AddCookie(accessCookie)
	req2.AddCookie(refreshCookie)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
}
