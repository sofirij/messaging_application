package v1

import (
	"app/internal/model/request"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	_, accessCookie, refreshCookie := login(t, app, bodyStruct)

	time.Sleep(1 * time.Second) // tokens generated at the same time will be identical

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	var newAccessCookie *http.Cookie
	var newRefreshCookie *http.Cookie

	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case accessCookie.Name:
			newAccessCookie = cookie
		case refreshCookie.Name:
			newRefreshCookie = cookie
		}
	}

	assert.NotNil(t, newAccessCookie)
	assert.NotNil(t, newRefreshCookie)

	// access authorization protected content with old tokens should fail
	resp2, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
}
