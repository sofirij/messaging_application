package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_Edit(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	user2 := request.UserAuthRequest{
		Username: "testuser2",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	originalText := "this is the original text"
	editedText := "this is the edited text"

	user1ID := getUserID(t, app, accessCookie1)

	// user1 sends a message
	createMessage(t, user1ID, conversationID, originalText)

	// get the messageID
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages", conversationID), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.PaginatedMessageResponse
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data[0].ID

	// user1 edits the message
	messageEditReq := request.MessageEditRequest{
		Body: editedText,
	}

	messageEditBytes, err := json.Marshal(messageEditReq)

	require.NoError(t, err)

	req3 := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/messages/%d", messageID), bytes.NewReader(messageEditBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)
	
	resp3, err := app.Test(req3)

	require.NoError(t, err)

	require.Equal(t, http.StatusNoContent, resp3.StatusCode)

	// assert the the message was edited
	resp2, err = app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	assert.Equal(t, editedText, *result2.Data[0].Body)
}

func TestMessage_EditForbidden(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	user2 := request.UserAuthRequest{
		Username: "testuser2",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	originalText := "this is the original text"
	editedText := "this is the edited text"

	user1ID := getUserID(t, app, accessCookie1)

	// user1 sends a message
	createMessage(t, user1ID, conversationID, originalText)

	// get the messageID
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages", conversationID), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.PaginatedMessageResponse
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data[0].ID

	// user2 fails to edit the message
	messageEditReq := request.MessageEditRequest{
		Body: editedText,
	}

	messageEditBytes, err := json.Marshal(messageEditReq)

	require.NoError(t, err)

	req3 := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/messages/%d", messageID), bytes.NewReader(messageEditBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie2)
	
	resp3, err := app.Test(req3)

	require.NoError(t, err)

	require.Equal(t, http.StatusForbidden, resp3.StatusCode)
}

func TestMessage_Delete(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	user2 := request.UserAuthRequest{
		Username: "testuser2",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	originalText := "this is the original text"

	user1ID := getUserID(t, app, accessCookie1)

	// user1 sends a message
	createMessage(t, user1ID, conversationID, originalText)

	// get the messageID
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages", conversationID), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.PaginatedMessageResponse
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data[0].ID

	// user1 deletes message
	req3 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/messages/%d", messageID), nil)
	req3.AddCookie(accessCookie1)
	req3.Header.Set("Content-Type", "application/json")

	resp3, err := app.Test(req3)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp3.StatusCode)

	// ensure that the body of the message is empty
	resp2, err = app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	assert.Nil(t, result2.Data[0].Body)
}

func TestMessage_GetPagination(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	user2 := request.UserAuthRequest{
		Username: "testuser2",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	text := "this is the original text"
	nextCycleText := "this it the text after the first cycle of get messages"

	user1ID := getUserID(t, app, accessCookie1)

	// send messages
	createMessage(t, user1ID, conversationID, nextCycleText)
	limit := 20
	for range limit {
		createMessage(t, user1ID, conversationID, text)
	}

	// get the first batch of messages
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages?limit=20", conversationID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie1)

	resp, err := app.Test(req)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result2 response.PaginatedMessageResponse

	err = json.NewDecoder(resp.Body).Decode(&result2)

	require.NoError(t, err)

	lastMessageID := *result2.NextCursor
	
	// get the last batch
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages?before=%d&limit=1", conversationID, lastMessageID), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp2.StatusCode)

	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	assert.Equal(t, nextCycleText, *result2.Data[0].Body)
}