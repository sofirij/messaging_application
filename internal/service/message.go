package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/repository"

	"github.com/jackc/pgx/v5"
)

const (
	attachmentLimit   = 10
	messageLimit      = 20
	editMessageWindow = 15 * time.Minute
)

type MessageService interface {
	Create(ctx context.Context, senderID, conversationID int, req request.MessageCreateRequest) (*response.MessageResponse, error)
	UpdateBody(ctx context.Context, senderID, messageID int, req request.MessageEditRequest) error
	GetByConversationID(ctx context.Context, userID, conversationID int, before *int, limit int) (*response.PaginatedMessageResponse, error)
	SoftDelete(ctx context.Context, senderID, messageID int) (*response.MessageResponse, error)
	MarkAsRead(ctx context.Context, userID, conversationID, messageID int) error
	GetByID(ctx context.Context, userID, messageID int) (*response.MessageResponse, error)
}

type messageService struct {
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

func NewMessageService(messageRepo repository.MessageRepository, conversationRepo repository.ConversationRepository) MessageService {
	return &messageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
	}
}

func (m *messageService) Create(ctx context.Context, senderID, conversationID int, req request.MessageCreateRequest) (*response.MessageResponse, error) {
	// too many attachments
	if len(req.Attachments) > attachmentLimit {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: fmt.Sprintf("more than %d attachments", attachmentLimit),
		}
	}

	// missing body and attachments
	if req.Body == nil && req.Attachments == nil {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "missing body",
		}
	}

	err := userInConversation(ctx, m.conversationRepo, conversationID, senderID)

	if err != nil {
		return nil, err
	}

	message, _, err := m.messageRepo.CreateWithAttachments(ctx, conversationID, senderID, req.ReplyToID, req.Body, req.Attachments)

	if err != nil {
		return nil, err
	}

	resp, err := getMessageByID(ctx, m.messageRepo, message.ID)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *messageService) UpdateBody(ctx context.Context, senderID int, messageID int, req request.MessageEditRequest) error {
	message, err := userOwnsMessage(ctx, m.messageRepo, messageID, senderID)

	if err != nil {
		return err
	}

	// message past edit window
	if time.Since(message.CreatedAt) > editMessageWindow {
		return &Error{
			Code:    ErrCodeForbidden,
			Message: "edit window has expired",
		}
	}

	// cannot edit deleted message
	if message.DeletedAt != nil {
		return &Error{
			Code:    ErrCodeNotFound,
			Message: "message not found",
		}
	}

	_, err = m.messageRepo.UpdateBody(ctx, messageID, req.Body)

	if err != nil {
		return err
	}

	return nil
}

func (m *messageService) SoftDelete(ctx context.Context, senderID, messageID int) (*response.MessageResponse, error) {
	message, err := userOwnsMessage(ctx, m.messageRepo, messageID, senderID)

	if err != nil {
		return nil, err
	}

	// cannot delete deleted message
	if message.DeletedAt != nil {
		return nil, &Error{
			Code:    ErrCodeNotFound,
			Message: "message not found",
		}
	}

	message, err = m.messageRepo.SoftDelete(ctx, messageID)

	if err != nil {
		return nil, err
	}

	resp, err := getMessageByID(ctx, m.messageRepo, message.ID)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
ensure user is in conversation, get messages, get message attachments, generate paginated response
*/
func (m *messageService) GetByConversationID(ctx context.Context, userID, conversationID int, before *int, limit int) (*response.PaginatedMessageResponse, error) {
	err := userInConversation(ctx, m.conversationRepo, conversationID, userID)

	if err != nil {
		return nil, err
	}

	if limit > messageLimit {
		limit = messageLimit
	}

	member, err := m.conversationRepo.GetMember(ctx, conversationID, userID)

	if err != nil {
		return nil, err
	}

	messages, err := m.messageRepo.GetByConversationID(ctx, conversationID, member.AfterCursor, before, limit+1)

	if err != nil {
		return nil, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	messageIDs := make([]int, len(messages))
	for i, message := range messages {
		messageIDs[i] = message.ID
	}

	attachments, err := m.messageRepo.GetAttachmentsByMessageIDs(ctx, messageIDs)

	if err != nil {
		return nil, err
	}

	attachmentMap := make(map[int][]db.MessageAttachment)
	for _, a := range attachments {
		attachmentMap[a.MessageID] = append(attachmentMap[a.MessageID], a)
	}

	// get reply metadata for messages with replies
	replyIDs := make([]int, 0)
	for _, message := range messages {
		if (message.ReplyToID != nil) {
			replyIDs = append(replyIDs, *message.ReplyToID)
		}
	}

	replies, err := m.messageRepo.GetByIDs(ctx, replyIDs)
	if err != nil {
		return nil, err
	}

	replyMap := make(map[int]*response.ReplyMetadata)
	for _, reply := range replies {
		replyMap[reply.ID] = &response.ReplyMetadata{
			ID: reply.ID,
			SenderID: reply.SenderID,
			Body: reply.Body,
		}
	}

	messageResp := make([]response.MessageResponse, len(messages))

	for i, message := range messages {
		messageAttachments := make([]response.MessageAttachment, len(attachmentMap[message.ID]))

		for j, messageAttachment := range attachmentMap[message.ID] {
			messageAttachments[j] = response.MessageAttachment{
				ID:       messageAttachment.ID,
				Type:     messageAttachment.Type,
				URL:      messageAttachment.URL,
				Filename: messageAttachment.Filename,
				Size:     messageAttachment.Size,
			}
		}

		var replyMetadata *response.ReplyMetadata
		if message.ReplyToID != nil {
			replyMetadata = replyMap[*message.ReplyToID]
		}

		messageResp[i] = response.MessageResponse{
			ID:             message.ID,
			ConversationID: message.ConversationID,
			SenderID:       message.SenderID,
			Reply:      	replyMetadata,
			Body:           message.Body,
			Deleted:        message.DeletedAt != nil,
			CreatedAt:      message.CreatedAt,
			Attachments:    messageAttachments,
		}
	}

	var nextCursor *int
	if len(messages) > 0 {
		lastID := messages[len(messages)-1].ID
		nextCursor = &lastID
	}

	resp := response.PaginatedMessageResponse{
		Messages:   messageResp,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return &resp, nil
}

func (m *messageService) MarkAsRead(ctx context.Context, userID, conversationID, messageID int) error {
	err := userInConversation(ctx, m.conversationRepo, conversationID, userID)

	if err != nil {
		return err
	}

	conversation, err := m.conversationRepo.GetByID(ctx, conversationID)

	if err != nil {
		return err
	}

	message, err := getMessageByID(ctx, m.messageRepo, messageID)

	if err != nil {
		return err
	}

	if message.ConversationID != conversation.ID {
		return &Error{
			Code:    ErrCodeBadRequest,
			Message: "message not in conversation",
		}
	}

	if _, err := userOwnsMessage(ctx, m.messageRepo, messageID, userID); err == nil {
		return &Error{
			Code: ErrCodeBadRequest,
			Message: "Cannot read your own message",
		}
	}

	err = m.conversationRepo.UpdateLastMessageRead(ctx, userID, conversationID, messageID)

	if err != nil {
		return err
	}

	return nil
}

func (m *messageService) GetByID(ctx context.Context, userID, messageID int) (*response.MessageResponse, error) {
	message, err := getMessageByID(ctx, m.messageRepo, messageID)
	if err != nil {
		return nil, err
	}

	err = userInConversation(ctx, m.conversationRepo, message.ConversationID, userID)

	if err != nil {
		return nil, err
	}

	return message, nil
}


func userOwnsMessage(ctx context.Context, repo repository.MessageRepository, messageID, userID int) (*db.Message, error) {
	message, err := repo.GetByID(ctx, messageID)

	if err != nil {
		return nil, err
	}

	// sender doesn't own message
	if message.SenderID != userID {
		return nil, &Error{
			Code:    ErrCodeForbidden,
			Message: "user doesn't own message",
		}
	}

	return message, nil
}

func getMessageByID(ctx context.Context, messageRepo repository.MessageRepository, messageID int) (*response.MessageResponse, error) {
	message, err := messageRepo.GetByID(ctx, messageID)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &Error{
			Code:    ErrCodeNotFound,
			Message: "message not found",
		}
	}

	attachments, err := messageRepo.GetAttachmentsByMessageIDs(ctx, []int{message.ID})

	if err != nil {
		return nil, err
	}

	attachmentsResp := make([]response.MessageAttachment, len(attachments))

	for i, attachment := range attachments {
		attachmentsResp[i] = response.MessageAttachment{
			ID:       attachment.ID,
			Type:     attachment.Type,
			URL:      attachment.URL,
			Filename: attachment.Filename,
			Size:     attachment.Size,
		}
	}

	var replyMetadata *response.ReplyMetadata
	if message.ReplyToID != nil {
		reply, err := messageRepo.GetByID(ctx, *message.ReplyToID)
		if err != nil {
			return nil, err
		}

		replyMetadata = &response.ReplyMetadata{
			ID: reply.ID,
			SenderID: reply.SenderID,
			Body: reply.Body,
		}
	}

	return &response.MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Reply:      	replyMetadata,
		Body:           message.Body,
		Deleted:        message.DeletedAt != nil,
		CreatedAt:      message.CreatedAt,
		Attachments:    attachmentsResp,
	}, nil
}
