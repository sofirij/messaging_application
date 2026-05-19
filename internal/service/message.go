package service

import (
	"context"
	"errors"

	"app/internal/model/db"
	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/service"
	"app/internal/repository"

	"github.com/jackc/pgx/v5"
)

const attachmentLimit = 10

type MessageService interface {
	Create(ctx context.Context, conversationID int, req request.MessageCreateRequest) (*response.MessageResponse, error)
	UpdateBody(ctx context.Context, senderID, messageID, req request.MessageEditRequest) error
	GetByConversationID(ctx context.Context, userID, conversationID int, before *int, limit int) (*response.PaginatedMessageResponse, error)
	SoftDelete(ctx context.Context, senderID, messageID int) error
}

type messageService struct {
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

func (m *messageService) Create(ctx context.Context, senderID int, conversationID int, req request.MessageCreateRequest) (*response.MessageResponse, error) {
	// too many attachments
	if len(req.Attachments) > attachmentLimit {
		return nil, &service.Error{}
	}

	err := senderInConversation(ctx, m.conversationRepo, conversationID, senderID)

	if err != nil {
		return nil, err
	}

	message, attachments, err := m.messageRepo.CreateWithAttachments(ctx, conversationID, senderID, req.ReplyToID, req.Body, req.Attachments)

	if err != nil {
		return nil, err
	}

	attachmentResp := make([]response.MessageAttachment, 0)

	for _, attachment := range attachments {
		attachmentResp = append(attachmentResp, response.MessageAttachment{
			ID:       attachment.ID,
			Type:     attachment.Type,
			URL:      attachment.URL,
			Filename: attachment.Filename,
			Size:     attachment.Size,
		})
	}

	resp := response.MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReplyToID:      message.ReplyToID,
		Body:           message.Body,
		Deleted:        message.DeletedAt != nil,
		CreatedAt:      message.CreatedAt,
		Attachments:    attachmentResp,
	}

	return &resp, nil
}

func (m *messageService) UpdateBody(ctx context.Context, senderID int, messageID int, req request.MessageEditRequest) error {
	message, err := senderOwnsMessage(ctx, m.messageRepo, messageID, senderID)

	if err != nil {
		return err
	}

	// cannot edit deleted message
	if message.DeletedAt != nil {
		return &service.Error{}
	}

	message, err = m.messageRepo.UpdateBody(ctx, messageID, req.Body)

	if err != nil {
		return err
	}

	return nil
}

func (m *messageService) SoftDelete(ctx context.Context, senderID, messageID int) (*response.MessageResponse, error) {
	message, err := senderOwnsMessage(ctx, m.messageRepo, messageID, senderID)

	if err != nil {
		return nil, err
	}

	// cannot delete deleted message
	if message.DeletedAt != nil {
		return nil, &service.Error{}
	}

	message, err = m.messageRepo.SoftDelete(ctx, messageID)

	if err != nil {
		return nil, err
	}

	resp := response.MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReplyToID:      message.ReplyToID,
		Body:           message.Body,
		Deleted:        message.DeletedAt != nil,
		CreatedAt:      message.CreatedAt,
		Attachments:    make([]response.MessageAttachment, 0),
	}

	return &resp, nil
}

/*
ensure user is in conversation, get messages, get message attachments, generate paginated response
*/
func (m *messageService) GetByConversationID(ctx context.Context, userID, conversationID int, before *int, limit int) (*response.PaginatedMessageResponse, error) {
	err := senderInConversation(ctx, m.conversationRepo, conversationID, userID)

	if err != nil {
		return nil, err
	}

	messages, err := m.messageRepo.GetByConversationID(ctx, conversationID, before, limit)

	if err != nil {
		return nil, err
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

		messageResp[i] = response.MessageResponse{
			ID:             message.ID,
			ConversationID: message.ConversationID,
			SenderID:       message.SenderID,
			ReplyToID:      message.ReplyToID,
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
		Data:       messageResp,
		NextCursor: nextCursor,
		HasMore:    false,
	}

	return &resp, nil
}

func senderInConversation(ctx context.Context, repo repository.ConversationRepository, conversationID, userID int) error {
	_, err := repo.GetMember(ctx, conversationID, userID)

	// sender doesn't belong in the conversation
	if errors.Is(err, pgx.ErrNoRows) {
		return &service.Error{}
	}

	return nil
}

func senderOwnsMessage(ctx context.Context, repo repository.MessageRepository, messageID, senderID int) (*db.Message, error) {
	message, err := repo.GetByID(ctx, messageID)

	if err != nil {
		return nil, err
	}

	// sender doesn't own message
	if message.SenderID != senderID {
		return nil, &service.Error{}
	}

	return message, nil
}
