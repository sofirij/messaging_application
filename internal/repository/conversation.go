package repository

import (
	"context"

	"app/internal/model/db"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationRepository interface {
	Create(ctx context.Context, createdBy int, conversationType string, name *string, rolesMap map[int]string) (*db.Conversation, error)
	GetByUserID(ctx context.Context, userID int) ([]db.Conversation, error)
	GetByID(ctx context.Context, conversationID int) (*db.Conversation, error)
	GetMembers(ctx context.Context, conversationID int) ([]db.ConversationMember, error)
	GetMember(ctx context.Context, conversationID int, userID int) (*db.ConversationMember, error)
	SoftDelete(ctx context.Context, conversationID int) error
	UpdateName(ctx context.Context, conversationID int, name string) (*db.Conversation, error)
	UpdateAvatarURL(ctx context.Context, conversationID int, avatarURL *string) (*db.Conversation, error)
	AddMembers(ctx context.Context, conversationID int, role string, userIDs []int) error
	RemoveMember(ctx context.Context, conversationID, userID int) error
}

type conversationRepository struct {
	db *pgxpool.Pool
}

func (c *conversationRepository) AddMembers(ctx context.Context, conversationID int, role string, userIDs []int) error {
	_, err := c.db.CopyFrom(
		ctx,
		pgx.Identifier{"conversation_members"},
		[]string{"conversation_id", "user_id", "role"},
		pgx.CopyFromSlice(len(userIDs), func(i int) ([]any, error) {
			userID := userIDs[i]
			return []any{conversationID, userID, role}, nil
		}),
	)

	return err
}

func (c *conversationRepository) Create(ctx context.Context, createdBy int, conversationType string, conversationName *string, rolesMap map[int]string) (*db.Conversation, error) {
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

	rows := make([][]any, 0)
	for userID, userRole := range rolesMap {
		rows = append(rows, []any{conversation.ID, userID, userRole})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"conversation_members"},
		[]string{"conversation_id", "user_id", "role"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (c *conversationRepository) SoftDelete(ctx context.Context, conversationID int) error {
	query := `
		UPDATE conversations
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	_, err := c.db.Exec(ctx, query, conversationID)

	return err
}

func (c *conversationRepository) GetByID(ctx context.Context, conversationID int) (*db.Conversation, error) {
	query := `
		SELECT * FROM conversations
		WHERE id = $1
		AND deleted_at IS NULL
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
		AND c.deleted_at IS NULL
		ORDER BY c.last_message_at DESC NULLS LAST
	`

	conversations := make([]db.Conversation, 0)

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

	members := make([]db.ConversationMember, 0)

	err := pgxscan.Select(ctx, c.db, &members, query, conversationID)

	return members, err
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
		AND deleted_at IS NULL
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
		AND deleted_at IS NULL
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
