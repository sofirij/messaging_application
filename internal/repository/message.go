package repository

import (
	"context"

	"app/internal/model/db"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository interface {
	CreateWithAttachments(ctx context.Context, conversationID, senderID int, replyToID *int, body *string, attachments []db.MessageAttachment) (*db.Message, []db.MessageAttachment, error)
	GetByID(ctx context.Context, messageID int) (*db.Message, error)
	GetByConversationID(ctx context.Context, conversationID int, before *int, limit int) ([]db.Message, error)
	GetAttachmentsByMessageIDs(ctx context.Context, messageIDs []int) ([]db.MessageAttachment, error)
	UpdateBody(ctx context.Context, messageID int, body string) (*db.Message, error)
	SoftDelete(ctx context.Context, messageID int) error
}

type messageRepository struct {
	db *pgxpool.Pool
}

func (m *messageRepository) CreateWithAttachments(ctx context.Context, conversationID int, senderID int, replyToID *int, body *string, attachments []db.MessageAttachment) (*db.Message, []db.MessageAttachment, error) {
	tx, err := m.db.Begin(ctx)

	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	var message db.Message
	messageAttachments := make([]db.MessageAttachment, 0)

	query := `
		INSERT INTO messages
		(conversation_id, sender_id, reply_to_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	err = pgxscan.Get(ctx, tx, &message, query, conversationID, senderID, replyToID, body)

	if err != nil {
		return nil, nil, err
	}

	query = `
		INSERT INTO message_attachments
		(message_id, type, url, filename, size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	for _, a := range attachments {
		var attachment db.MessageAttachment

		err := pgxscan.Get(ctx, tx, &attachment, query, message.ID, a.Type, a.URL, a.Filename, a.Size)

		if err != nil {
			return nil, nil, err
		}

		messageAttachments = append(messageAttachments, attachment)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return &message, messageAttachments, nil
}

func (m *messageRepository) SoftDelete(ctx context.Context, messageID int) error {
	query := `
		UPDATE messages
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	_, err := m.db.Exec(ctx, query, messageID)

	return err
}

func (m *messageRepository) GetAttachmentsByMessageIDs(ctx context.Context, messageIDs []int) ([]db.MessageAttachment, error) {
	query := `
		SELECT * FROM message_attachments
		WHERE message_id = ANY($1)
	`

	attachments := make([]db.MessageAttachment, 0)

	err := pgxscan.Select(ctx, m.db, &attachments, query, messageIDs)

	return attachments, err
}

func (m *messageRepository) GetByConversationID(ctx context.Context, conversationID int, before *int, limit int) ([]db.Message, error) {
	query := `
		SELECT * FROM messages
		WHERE conversation_id = $1
		AND ($2::int is NULL OR id < $2)
		AND deleted_at IS NULL
		ORDER BY id DESC
		LIMIT $3
	`

	messages := make([]db.Message, 0)

	err := pgxscan.Select(ctx, m.db, &messages, query, conversationID, before, limit)

	return messages, err
}

func (m *messageRepository) GetByID(ctx context.Context, messageID int) (*db.Message, error) {
	query := `
		SELECT * FROM messages
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var message db.Message

	err := pgxscan.Get(ctx, m.db, &message, query, messageID)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (m *messageRepository) UpdateBody(ctx context.Context, messageID int, body string) (*db.Message, error) {
	query := `
		UPDATE messages
		SET body = $1, edited_at = NOW()
		WHERE deleted_at IS NULL
		AND id = $2
		RETURNING *
	`

	var message db.Message

	err := pgxscan.Get(ctx, m.db, &message, query, body, messageID)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func NewMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &messageRepository{db: db}
}
