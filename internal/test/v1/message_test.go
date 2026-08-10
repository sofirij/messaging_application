package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func createMessage(t *testing.T, app *fiber.App, accessCookie *http.Cookie, conversationID int, text string) response.MessageResponse {
	bodyStruct := request.MessageCreateRequest{
		Body: &text,
	}
	bodyBytes, err := json.Marshal(bodyStruct)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/conversations/%d/messages", conversationID), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie)

	resp, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result response.Response[response.MessageResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)

	require.NotNil(t, result.Data.Body)
	require.Equal(t, text, *result.Data.Body)

	return result.Data
}

func TestMessage_CreateWithAttachmentsAndReply(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

	replyText := "this is the reply text"

	replyMessage := createMessage(t, app, accessCookie1, conversationID, replyText)

	// upload the file to be used as attachment
	content, err := os.ReadFile(validFilePath)

	require.NoError(t, err)

	contents := make([][]byte, 1)
	filenames := make([]string, 1)

	contents[0] = content
	filenames[0] = validFilename

	req2 := createMultipartRequestMany(t, accessCookie1, "file", filenames, contents)

	resp2, err := app.Test(req2, testCfg)

	require.NoError(t, err)

	// get the uploaded file url
	var result2 response.Response[[]response.UploadResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	fileURL := result2.Data[0].URL
	require.NotEmpty(t, fileURL)

	text := "this is a message"

	attachment := request.MessageAttachment{
		URL:      fileURL,
		Filename: result2.Data[0].Filename,
		Size:     result2.Data[0].Size,
		Type:     result2.Data[0].Type,
	}

	messageBody := request.MessageCreateRequest{
		ReplyToID:   &replyMessage.ID,
		Body:        &text,
		Attachments: []request.MessageAttachment{attachment},
	}

	bodyBytes, err := json.Marshal(messageBody)

	require.NoError(t, err)

	// send the message with the attachment
	req3 := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/conversations/%d/messages", conversationID), bytes.NewReader(bodyBytes))
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)

	resp3, err := app.Test(req3)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp3.StatusCode)

	var result3 response.Response[response.MessageResponse]
	err = json.NewDecoder(resp3.Body).Decode(&result3)

	require.NoError(t, err)

	// require the message response
	require.NotNil(t, result3.Data.Body)
	require.Equal(t, text, *result3.Data.Body)
	require.Equal(t, fileURL, result3.Data.Attachments[0].URL)
	require.Equal(t, replyMessage.ID, result3.Data.Reply.ID)
}

func TestMessage_Create(t *testing.T) {
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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

	text := "this is the original text"

	// user1 sends a message
	createMessage(t, app, accessCookie1, conversationID, text)
}

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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

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

	var result2 response.Response[*response.PaginatedMessageResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data.Messages[0].ID

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

	// require the the message was edited
	resp2, err = app.Test(req2)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	require.NotNil(t, result2.Data.Messages[0].Body)
	require.Equal(t, editedText, *result2.Data.Messages[0].Body)
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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

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

	var result2 response.Response[*response.PaginatedMessageResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data.Messages[0].ID

	// user2 fails to edit the message
	messageEditReq := request.MessageEditRequest{
		Body: editedText,
	}

	messageEditBytes, err := json.Marshal(messageEditReq)

	require.NoError(t, err)

	req3 := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/messages/%d", messageID), bytes.NewReader(messageEditBytes))
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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

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

	var result2 response.Response[*response.PaginatedMessageResponse]
	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	messageID := result2.Data.Messages[0].ID

	// user1 deletes message
	req3 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/messages/%d", messageID), nil)
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

	require.Nil(t, result2.Data.Messages[0].Body)
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
		Type:    db.DirectConversation,
		UserIDs: []int{user2ID}, // user2's id
	}

	conversation := createConversation(t, app, accessCookie1, bodyStruct)
	conversationID := conversation.ID

	otherText := "this is the other text"
	firstText := "this is the first text"

	// user1 sends 2 batches of messages each batch has a limit of 20 messages
	createMessage(t, app, accessCookie1, conversationID, firstText)
	limit := 20
	for range limit {
		createMessage(t, app, accessCookie1, conversationID, otherText)
	}

	// test the 'before' query parameter
	// get the last 20 messages
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/conversations/%d/messages?limit=20", conversationID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(accessCookie1)

	resp, err := app.Test(req)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result response.Response[response.PaginatedMessageResponse]

	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)

	require.Nil(t, result.Data.NextCursor)
	require.NotNil(t, result.Data.PreviousCursor)

	// get the first message
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/conversations/%d/messages?before=%d&limit=1", conversationID, *result.Data.PreviousCursor), nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(accessCookie1)

	resp2, err := app.Test(req2)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 response.Response[response.PaginatedMessageResponse]

	err = json.NewDecoder(resp2.Body).Decode(&result2)

	require.NoError(t, err)

	require.NotNil(t, result2.Data.Messages[0].Body)
	require.Equal(t, firstText, *result2.Data.Messages[0].Body)

	// test the 'at' query parameter
	// get the first 20 messages
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/conversations/%d/messages?limit=20&at=%d", conversationID, result2.Data.Messages[0].ID), nil)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(accessCookie1)

	resp3, err := app.Test(req3)

	require.NoError(t, err)

	var result3 response.Response[response.PaginatedMessageResponse]

	err = json.NewDecoder(resp3.Body).Decode(&result3)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp3.StatusCode)
	require.Equal(t, result2.Data.Messages[0].ID, result3.Data.Messages[0].ID)
	require.Nil(t, result3.Data.PreviousCursor)
	require.NotNil(t, result3.Data.NextCursor)
}
