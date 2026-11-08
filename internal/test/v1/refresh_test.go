package v1

import (
	"app/internal/model/request"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

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
	_, _, refreshCookie := login(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	var accessCookie *http.Cookie

	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			accessCookie = c
		}
	}

	require.NotNil(t, accessCookie)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie)

	resp2, err := app.Test(req2)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp2.StatusCode)
}

func TestRefresh_ConcurrentRefresh(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	_, _, refreshCookie := login(t, app, bodyStruct)

	var wg sync.WaitGroup
	const refreshes = 10
	results := make([]int, refreshes)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(refreshCookie)

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := app.Test(req)
			require.NoError(t, err)
			results[i] = resp.StatusCode
		}(i)
	}

	wg.Wait()

	successCount := 0
	for _, status := range results {
		if status == http.StatusNoContent {
			successCount++
		}
	}

	require.Equal(t, refreshes, successCount)
}
