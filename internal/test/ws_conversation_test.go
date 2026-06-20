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

func TestWSConversation_Create(t *testing.T) {
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

	// get the id for user 2
	userID := getUserID(t, app, accessCookie2)

	ticket := getTicket(t, app, accessCookie2)
	conn := connect(t, ticket)
	defer conn.Close()

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	createConversation(t, app, accessCookie1, bodyStruct)

	// ensure user 2 receives the conversation new message
	var event ws.Event
	for {
		err := conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventConversationNew {
			break
		}
	}

	var payload response.ConversationResponse
	err := json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.ID)
}

func TestWSConversation_AddMember(t *testing.T) {
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

	user3 := request.UserAuthRequest{
		Username: "testuser3",
		Password: "password123",
	}

	register(t, app, user1)
	register(t, app, user2)
	register(t, app, user3)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)
	_, accessCookie3, _ := login(t, app, user3)

	ticket := getTicket(t, app, accessCookie3)
	conn := connect(t, ticket)
	defer conn.Close()

	// get user2's id
	userID := getUserID(t, app, accessCookie2)

	// user1 creates group conversation with user2
	convName := "randomname"

	bodyStruct := request.ConversationCreateRequest{
		Type:    "group",
		Name:    &convName,
		UserIDs: []int{userID},
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	// user1 adds user3 to conversation
	addMemberStruct := request.ConversationAddMemberRequest{
		UserIDs: []int{getUserID(t, app, accessCookie3)},
	}

	body, err := json.Marshal(addMemberStruct)

	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/conversations/%d/members", conversationID), bytes.NewReader(body))
	req.AddCookie(accessCookie1)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// added member should receive a message
	var event ws.Event
	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMemberAdded {
			break
		}
	}

	var payload ws.MemberAddedPayload
	err = json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.ConversationID)
}

func TestWSConversation_RemoveMember(t *testing.T) {
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

	// get user2's id
	userID := getUserID(t, app, accessCookie2)

	// user1 creates group conversation with user2
	convName := "randomname"

	bodyStruct := request.ConversationCreateRequest{
		Type:    "group",
		Name:    &convName,
		UserIDs: []int{userID},
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	// user1 removes user2
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/conversations/%d/members/%d", conversationID, userID), nil)
	req.AddCookie(accessCookie1)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	
	// ensure user 2 receives the message
	var event ws.Event
	for {
		err = conn.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventMemberRemoved {
			break
		}
	}

	var payload ws.MemberRemovedPayload
	err = json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	assert.NotZero(t, payload.ConversationID)
}