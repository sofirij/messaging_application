package repository

import (
	"context"

	"app/internal/model/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository interface {
	Create(ctx context.Context, conversationID, sernderID int, replyToID *int, body *string) (*db.Message, error)
	CreateAttachments(ctx context.Context, messageID int, attachments []db.MessageAttachment) error
	GetByID(ctx context.Context, messageID int) (*db.Message, error)
	GetByConversationID(ctx context.Context, conversationID, before *int, limit int) ([]db.Message, error)
	GetAttachmentsByMessageIDs(ctx context.Context, messageIDs []int) ([]db.MessageAttachment, error)
	UpdateBody(ctx context.Context, messageID int, body string) (*db.Message, error)
	Delete(ctx context.Context, messageID int) error
}

type messageRepository struct {
	db *pgxpool.Pool
}

// Create implements [MessageRepository].
func (m *messageRepository) Create(ctx context.Context, conversationID int, sernderID int, replyToID *int, body *string) (*db.Message, error) {
	panic("unimplemented")
}

// CreateAttachments implements [MessageRepository].
func (m *messageRepository) CreateAttachments(ctx context.Context, messageID int, attachments []db.MessageAttachment) error {
	panic("unimplemented")
}

// Delete implements [MessageRepository].
func (m *messageRepository) Delete(ctx context.Context, messageID int) error {
	panic("unimplemented")
}

// GetAttachmentsByMessageIDs implements [MessageRepository].
func (m *messageRepository) GetAttachmentsByMessageIDs(ctx context.Context, messageIDs []int) ([]db.MessageAttachment, error) {
	panic("unimplemented")
}

// GetByConversationID implements [MessageRepository].
func (m *messageRepository) GetByConversationID(ctx context.Context, conversationID *int, before *int, limit int) ([]db.Message, error) {
	panic("unimplemented")
}

// GetByID implements [MessageRepository].
func (m *messageRepository) GetByID(ctx context.Context, messageID int) (*db.Message, error) {
	panic("unimplemented")
}

// UpdateBody implements [MessageRepository].
func (m *messageRepository) UpdateBody(ctx context.Context, messageID int, body string) (*db.Message, error) {
	panic("unimplemented")
}

func NewMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &messageRepository{db: db}
}
