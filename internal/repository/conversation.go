package repository

import (
	"context"

	"app/internal/model/db"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationRepository interface {
	Create(ctx context.Context, createdBy int, conversationType string, name *string, userIDs []int) (*db.Conversation, error)
	GetByUserID(ctx context.Context, userID int) ([]db.Conversation, error)
	GetByID(ctx context.Context, conversationID int) (*db.Conversation, error)
	GetMembers(ctx context.Context, conversationID int) ([]db.ConversationMember, error)
	GetMembersByConversationIDs(ctx context.Context, conversationIDs []int) ([]db.ConversationMember, error)
	GetMember(ctx context.Context, conversationID int, userID int) (*db.ConversationMember, error)
	UpdateName(ctx context.Context, conversationID int, name string) (*db.Conversation, error)
	UpdateAvatarURL(ctx context.Context, conversationID int, avatarURL *string) (*db.Conversation, error)
	AddMembers(ctx context.Context, conversationID int, userIDs []int) error
	RemoveMember(ctx context.Context, conversationID, userID int) error
	SoftDeleteMember(ctx context.Context, conversationID, userID int) error
	ClearMessages(ctx context.Context, conversationID, userID int) error
}

type conversationRepository struct {
	db *pgxpool.Pool
}

func (c *conversationRepository) SoftDeleteMember(ctx context.Context, conversationID, userID int) error {
	query := `
		UPDATE conversation_members
		SET deleted_at = NOW()
		WHERE conversation_id = $1
		AND user_id = $2
	`

	_, err := c.db.Exec(ctx, query, conversationID, userID)

	return err
}

// set after_cursor to the last message_id in the conversation
func (c *conversationRepository) ClearMessages(ctx context.Context, conversationID, userID int) error {
	query := `
		UPDATE conversation_members
		SET after_cursor = (
			SELECT id FROM messages
			WHERE conversation_id = $1
			ORDER BY id DESC LIMIT 1
		)
		WHERE user_id = $2
		AND conversation_id = $1
	`

	_, err := c.db.Exec(ctx, query, conversationID, userID)

	return err
}

func (c *conversationRepository) AddMembers(ctx context.Context, conversationID int, userIDs []int) error {
	_, err := c.db.CopyFrom(
		ctx,
		pgx.Identifier{"conversation_members"},
		[]string{"conversation_id", "user_id"},
		pgx.CopyFromSlice(len(userIDs), func(i int) ([]any, error) {
			userID := userIDs[i]
			return []any{conversationID, userID}, nil
		}),
	)

	return err
}

func (c *conversationRepository) Create(ctx context.Context, createdBy int, conversationType string, conversationName *string, userIDs []int) (*db.Conversation, error) {
	tx, err := c.db.Begin(ctx)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO conversations
		(type, name, created_by)
		VALUES ($1, $2, $3)
		RETURNING *
	`

	var conversation db.Conversation

	err = pgxscan.Get(
		ctx,
		tx,
		&conversation,
		query,
		conversationType,
		conversationName,
		createdBy,
	)

	if err != nil {
		return nil, err
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"conversation_members"},
		[]string{"conversation_id", "user_id"},
		pgx.CopyFromSlice(len(userIDs), func(i int) ([]any, error) {
			userID := userIDs[i]
			return []any{conversation.ID, userID}, nil
		}),
	)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (c *conversationRepository) GetByID(ctx context.Context, conversationID int) (*db.Conversation, error) {
	query := `
		SELECT * FROM conversations
		WHERE id = $1
	`

	var conversation db.Conversation

	err := pgxscan.Get(ctx, c.db, &conversation, query, conversationID)

	if err != nil {
		return nil, err
	}

	return &conversation, err
}

func (c *conversationRepository) GetByUserID(ctx context.Context, userID int) ([]db.Conversation, error) {
	query := `
		SELECT c.* FROM conversations AS c
		JOIN conversation_members AS cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1
		AND (cm.deleted_at IS NULL OR c.last_message_at > cm.deleted_at)
		ORDER BY c.last_message_at DESC NULLS LAST
	`

	var conversations []db.Conversation

	err := pgxscan.Select(ctx, c.db, &conversations, query, userID)

	return conversations, err
}

func (c *conversationRepository) GetMember(ctx context.Context, conversationID int, userID int) (*db.ConversationMember, error) {
	query := `
		SELECT * FROM conversation_members
		WHERE conversation_id = $1
		AND user_id = $2
	`

	var conversationMember db.ConversationMember

	err := pgxscan.Get(ctx, c.db, &conversationMember, query, conversationID, userID)

	if err != nil {
		return nil, err
	}

	return &conversationMember, nil
}

func (c *conversationRepository) GetMembers(ctx context.Context, conversationID int) ([]db.ConversationMember, error) {
	query := `
		SELECT * FROM conversation_members
		WHERE conversation_id = $1
	`

	var members []db.ConversationMember

	err := pgxscan.Select(ctx, c.db, &members, query, conversationID)

	return members, err
}

func (c *conversationRepository) GetMembersByConversationIDs(ctx context.Context, conversationIDs []int) ([]db.ConversationMember, error) {
	query := `
		SELECT * FROM conversation_members
		WHERE conversation_id = ANY($1)
	`

	var members []db.ConversationMember

	err := pgxscan.Select(ctx, c.db, &members, query, conversationIDs)

	if err != nil {
		return nil, err
	}

	return members, nil
}

func (c *conversationRepository) RemoveMember(ctx context.Context, conversationID int, userID int) error {
	query := `
		DELETE FROM conversation_members
		WHERE user_id = $1
		AND conversation_id = $2
	`

	_, err := c.db.Exec(ctx, query, userID, conversationID)

	return err
}

func (c *conversationRepository) UpdateAvatarURL(ctx context.Context, conversationID int, avatarURL *string) (*db.Conversation, error) {
	query := `
		UPDATE conversations
		SET avatar_url = $1
		WHERE id = $2
		RETURNING *
	`

	var conversation db.Conversation

	err := pgxscan.Get(ctx, c.db, &conversation, query, avatarURL, conversationID)

	if err != nil {
		return nil, err
	}

	return &conversation, err
}

func (c *conversationRepository) UpdateName(ctx context.Context, conversationID int, name string) (*db.Conversation, error) {
	query := `
		UPDATE conversations
		SET name = $1
		WHERE id = $2
		RETURNING *
	`

	var conversation db.Conversation

	err := pgxscan.Get(ctx, c.db, &conversation, query, name, conversationID)

	if err != nil {
		return nil, err
	}

	return &conversation, err
}

func NewConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{db: db}
}
