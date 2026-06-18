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

// indirectly tests the happy path for GetByID
func TestConversation_CreateDirect(t *testing.T) {
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
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)

	// ensure that the only 2 members in the conversation are user1 and user2
	require.NotZero(t, result.Data.ID)

	conversationID := result.Data.ID

	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)

	resp3, err := app.Test(req3)

	require.NoError(t, err)

	err = json.NewDecoder(resp3.Body).Decode(&result)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)
	assert.Equal(t, 2, len(result.Data.Members)) // there should be exactly 2 members in the conversation

	for _, member := range result.Data.Members {
		switch member.Username {
		case user1.Username:
		case user2.Username:
			continue
		default:
			t.Fatalf("Invalid member %s\n", member.Username)
		}
	}
}

func TestConversation_CreateDuplicateDirect(t *testing.T) {
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

	body, err := json.Marshal(user1)
	require.NoError(t, err)

	register(t, app, user1)
	register(t, app, user2)

	_, accessCookie1, _ := login(t, app, user1)
	_, accessCookie2, _ := login(t, app, user2)

	// get the id for user 2
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	createConversation(t, app, accessCookie1, bodyStruct)

	// get the id of user1
	userID = getUserID(t, app, accessCookie1)

	// user2 creates a direct conversation with user1
	bodyStruct = request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user1's id
	}

	body, err = json.Marshal(bodyStruct)
	require.NoError(t, err)

	req4 := httptest.NewRequest(http.MethodPost, "/api/conversations", bytes.NewReader(body))
	req4.Header.Set("Content-Type", "application/json")
	req4.AddCookie(accessCookie2)

	resp4, err := app.Test(req4)

	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp4.StatusCode)
}

func TestConversation_CreateGroup(t *testing.T) {
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
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	groupName := "somegroup"

	bodyStruct := request.ConversationCreateRequest{
		Type:    "group",
		Name:    &groupName,
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	require.NotZero(t, result.Data.ID)

	conversationID := result.Data.ID

	// ensure that the only 2 members in the conversation are user1 and user2
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)

	resp3, err := app.Test(req3)

	require.NoError(t, err)

	err = json.NewDecoder(resp3.Body).Decode(&result)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)
	assert.Equal(t, 2, len(result.Data.Members)) // there should be 2 members in the conversation

	for _, member := range result.Data.Members {
		switch member.Username {
		case user1.Username:
		case user2.Username:
			continue
		default:
			t.Fatalf("Invalid member %s\n", member.Username)
		}
	}
}

func TestConversation_UserNotInConversation(t *testing.T) {
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

	// get the id for user 2
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	require.NotZero(t, result.Data.ID)

	conversationID := result.Data.ID

	// user 3 tries to access the conversation
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie3)

	resp3, err := app.Test(req3)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp3.StatusCode)
}

func TestConversation_AddMember(t *testing.T) {
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

	// ensure conversation has the right members
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req2.AddCookie(accessCookie1)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	err = json.NewDecoder(resp2.Body).Decode(&result)

	require.NoError(t, err)

	for _, member := range result.Data.Members {
		switch member.Username {
		case user1.Username:
		case user2.Username:
		case user3.Username:
			continue
		default:
			t.Fatalf("Invalid member %s", member.Username)
		}
	}
}

func TestConversation_AddMemberToDirect(t *testing.T) {
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

	// get user2's id
	userID := getUserID(t, app, accessCookie2)

	// user1 creates group conversation with user2
	convName := "randomname"

	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
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
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestConversation_RemoveMember(t *testing.T) {
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

	// ensure that user2 is removed
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)

	err = json.NewDecoder(resp2.Body).Decode(&result)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	for _, member := range result.Data.Members {
		switch member.Username {
		case user1.Username:
			continue
		default:
			t.Fatalf("Invalid member %s\n", member.Username)
		}
	}
}

func TestConversation_ClearMessages(t *testing.T) {
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
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := result.Data.ID

	// get the id for user 1
	userID = getUserID(t, app, accessCookie1)

	// user1 sends a message
	createMessage(t, userID, conversationID, "this is a text")

	// user2 clears the message
	resp := clearMessages(t, app, accessCookie2, conversationID)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// user1 sends a message after it was cleared for user2
	lastMessage := "this is a text after messages where cleared"
	createMessage(t, userID, conversationID, lastMessage)

	// ensure its only the last message user2 can see
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/conversations/%d/messages", conversationID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie2)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	
	var result2 response.PaginatedMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&result2)

	require.NoError(t, err)

	assert.Equal(t, 1, len(result2.Data))
	assert.Equal(t, lastMessage, *result2.Data[0].Body)
}

func TestConversation_Delete(t *testing.T) {
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
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)

	conversationID := result.Data.ID

	// user2 deletes conversation
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie2)
	
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// user1 should still see the conversation on their list
	req2 := httptest.NewRequest(http.MethodGet, "/api/conversations", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.Response[[]response.ConversationResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	require.Equal(t, 1, len(result2.Data))

	// user2's conversation list should be empty
	req3 := httptest.NewRequest(http.MethodGet, "/api/conversations", nil)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie2)

	resp3, err := app.Test(req3)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result3 response.Response[[]response.ConversationResponse]
	err = json.NewDecoder(resp3.Body).Decode(&result3)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	assert.Equal(t, 0, len(result3.Data))
}

func TestConversation_GetAfterDelete(t *testing.T) {
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
	userID := getUserID(t, app, accessCookie2)

	// user1 creates a direct conversation with user2
	bodyStruct := request.ConversationCreateRequest{
		Type:    "direct",
		UserIDs: []int{userID}, // user2's id
	}

	result, _ := createConversation(t, app, accessCookie1, bodyStruct)

	conversationID := result.Data.ID

	// user2 deletes conversation
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/conversations/%d", conversationID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie2)
	
	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// user1 sends a message
	text := "texting after conversation was deleted"
	createMessage(t, getUserID(t, app, accessCookie1), conversationID, text)

	// the conversation should show up in user2's list
	req2 := httptest.NewRequest(http.MethodGet, "/api/conversations", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie2)

	resp2, err := app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.Response[[]response.ConversationResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	assert.Equal(t, 1, len(result2.Data))
}