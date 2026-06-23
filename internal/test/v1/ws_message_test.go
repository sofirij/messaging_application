package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/ws"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWSMessage_Create(t *testing.T) {
	app := setupApp(t)
	go listen(t, app)

	defer truncateTables(t)
	defer app.Shutdown()

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

	ticket := getTicket(t, app, accessCookie2)
	conn := connect(t, ticket)
	defer conn.Close()

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	text := "this is the original text"

	// user1 sends a message
	createMessage(t, app, accessCookie1, conversationID, text)

	var event ws.Event
	for {
		err := conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMessageNew {
			break
		}
	}

	var payload response.MessageResponse
	err := json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.ID)
}

func TestWSMessage_Delete(t *testing.T) {
	app := setupApp(t)
	go listen(t, app)

	defer truncateTables(t)
	defer app.Shutdown()

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

	ticket := getTicket(t, app, accessCookie2)
	conn := connect(t, ticket)
	defer conn.Close()

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	originalText := "this is the original text"

	// user1 sends a message
	createMessage(t, app, accessCookie1, conversationID, originalText)

	// get the messageID
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/conversations/%d/messages", conversationID), nil)
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
	req3 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/messages/%d", messageID), nil)
	req3.AddCookie(accessCookie1)
	req3.Header.Set("Content-Type", "application/json")

	resp3, err := app.Test(req3)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp3.StatusCode)

	// ensure user 2 receives the delete message
	var event ws.Event
	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMessageDeleted {
			break
		}
	}

	var payload ws.MessageDeletedPayload
	err = json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.ConversationID)
	assert.NotZero(t, payload.MessageID)
}

func TestWSMessage_Edit(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	go listen(t, app)
	defer app.Shutdown()

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

	ticket := getTicket(t, app, accessCookie2)
	conn := connect(t, ticket)
	defer conn.Close()

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	originalText := "this is the original text"
	editedText := "this is the edited text"

	// user1 sends a message
	createMessage(t, app, accessCookie1, conversationID, originalText)

	// get the messageID
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/conversations/%d/messages", conversationID), nil)
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

	req3 := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/messages/%d", messageID), bytes.NewReader(messageEditBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)

	resp3, err := app.Test(req3)

	require.NoError(t, err)

	require.Equal(t, http.StatusNoContent, resp3.StatusCode)

	// ensure user2 receives the message edited message
	var event ws.Event
	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMessageEdited {
			break
		}
	}

	var payload ws.MessageEditedPayload
	err = json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.MessageID)
}

func TestWSMessage_Read(t *testing.T) {
	app := setupApp(t)
	go listen(t, app)

	defer truncateTables(t)
	defer app.Shutdown()

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

	ticket := getTicket(t, app, accessCookie2)
	conn := connect(t, ticket)
	defer conn.Close()

	// get the id for user 2
	user2ID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{user2ID}, // user2's id
	}

	result := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	text := "this is the original text"

	// user1 sends a message
	result2 := createMessage(t, app, accessCookie1, conversationID, text)
	messageID := result2.Data.ID

	// user2 reads message
	reqPayload := ws.MessageReadPayload{
		ConversationID: conversationID,
		MessageID:      messageID,
	}

	payloadBytes, err := json.Marshal(reqPayload)

	require.NoError(t, err)

	event := ws.Event{
		Type:    ws.EventMessageRead,
		Payload: payloadBytes,
	}

	err = conn.WriteJSON(event)
	require.NoError(t, err)

	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMessageSeen {
			break
		}
		if event.Type == "error" {
			t.Log(event.Type)
			var payload ws.ErrorPayload
			json.Unmarshal(event.Payload, &payload)
			t.Logf("%v\n", payload)
		}
	}

	var respPayload ws.MessageSeenPayload
	err = json.Unmarshal(event.Payload, &respPayload)
	require.NoError(t, err)

	assert.NotZero(t, respPayload.ConversationID)
	assert.NotZero(t, respPayload.MessageID)
	assert.NotZero(t, respPayload.UserID)

	// user1 sends another message
	result2 = createMessage(t, app, accessCookie1, conversationID, text)
	messageID = result2.Data.ID

	// user 2 reads the next message
	reqPayload = ws.MessageReadPayload{
		ConversationID: conversationID,
		MessageID:      messageID,
	}

	payloadBytes, err = json.Marshal(reqPayload)

	require.NoError(t, err)

	event = ws.Event{
		Type:    ws.EventMessageRead,
		Payload: payloadBytes,
	}

	err = conn.WriteJSON(event)
	require.NoError(t, err)

	// test to ensure that reads work on messages after the first
	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMessageSeen {
			break
		}
	}
}
