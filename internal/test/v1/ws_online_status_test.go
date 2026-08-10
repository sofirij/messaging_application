package v1

import (
	"encoding/json"
	"testing"

	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/ws"

	"github.com/stretchr/testify/require"
)

func TestWSOnlineStatus_Online(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	listening := listen(t, app)
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

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    db.DirectConversation,
		UserIDs: []int{userID}, // user2's id
	}

	createConversation(t, app, accessCookie1, bodyStruct)

	// both user1 and user2 go online
	ticket := getTicket(t, app, accessCookie1)
	<-listening
	conn1 := connect(t, ticket)
	defer conn1.Close()

	ticket = getTicket(t, app, accessCookie2)
	conn2 := connect(t, ticket)
	defer conn2.Close()

	// when user2 goes online user1 should receive a message that user2 is onlinie
	var event ws.Event
	for {
		err := conn1.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventUserOnline {
			break
		}
	}

	var payload ws.UserOnlinePayload
	err := json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	require.NotZero(t, payload.UserID)
}

func TestWSOnlineStatus_Offline(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)

	listening := listen(t, app)
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

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    db.DirectConversation,
		UserIDs: []int{userID}, // user2's id
	}

	createConversation(t, app, accessCookie1, bodyStruct)

	// both user1 and user2 go online
	ticket := getTicket(t, app, accessCookie1)
	<-listening
	conn1 := connect(t, ticket)
	defer conn1.Close()

	ticket = getTicket(t, app, accessCookie2)
	conn2 := connect(t, ticket)
	conn2.Close()

	// when user2 goes offline user1 should receive a message that user2 is offline
	var event ws.Event
	for {
		err := conn1.ReadJSON(&event)
		require.NoError(t, err)

		if event.Type == ws.EventUserOffline {
			break
		}
	}

	var payload ws.UserOfflinePayload
	err := json.Unmarshal(event.Payload, &payload)
	require.NoError(t, err)

	require.NotZero(t, payload.UserID)
	require.NotZero(t, payload.LastSeenAt)
}
