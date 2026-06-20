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
	"app/internal/model/ws"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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