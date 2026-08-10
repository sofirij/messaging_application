package v1

import (
	"app/internal/model/request"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefresh_ServerSideRefresh(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	bodyStruct := request.UserAuthRequest{
		Username: "testuser",
		Password: "password123",
	}

	register(t, app, bodyStruct)
	_, _, refreshCookie := login(t, app, bodyStruct)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(refreshCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
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
	results := make([]int, 10)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
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
		if status == http.StatusOK {
			successCount++
		}
	}

	require.Equal(t, 1, successCount)
}