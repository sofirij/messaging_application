package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func createConversation(t *testing.T, app *fiber.App, accessCookie *http.Cookie, bodyStruct request.ConversationCreateRequest) (*response.Response[response.ConversationResponse], *http.Response) {
	body, err := json.Marshal(bodyStruct)

	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/conversations", bytes.NewReader(body))
	req.AddCookie(accessCookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var conversation response.Response[response.ConversationResponse]

	err = json.NewDecoder(resp.Body).Decode(&conversation)
	require.NoError(t, err)

	return &conversation, resp
}

func getUserID(t *testing.T, app *fiber.App, accessCookie *http.Cookie) int {
	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var user response.Response[response.UserResponse]
	err = json.NewDecoder(resp.Body).Decode(&user)

	require.NoError(t, err)
	return user.Data.ID
}

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

func register(t *testing.T, app *fiber.App, bodyStruct request.UserAuthRequest) *http.Response {
	body, err := json.Marshal(bodyStruct)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return resp
}

func createMessage(t *testing.T, userID, conversationID int, text string) *response.MessageResponse {
	bodyStruct := request.MessageCreateRequest{
		Body: &text,
	}

	messageResponse, err := messageService.Create(context.Background(), userID, conversationID, bodyStruct)
	require.NoError(t, err)

	return messageResponse
}

func clearMessages(t *testing.T, app *fiber.App, accessCookie *http.Cookie, conversationID int) *http.Response {
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/conversations/%d/messages", conversationID), nil)
	req.AddCookie(accessCookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	return resp
}
