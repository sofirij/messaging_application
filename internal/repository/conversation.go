package repository

import (
	"context"

	"app/internal/model/db"

	"github.com/jackc/pgx/v5/pgxpool"
	_"github.com/georgysavva/scany/v2/pgxscan"
)

type ConversationRepository interface {
	Create(ctx context.Context, createdBy int, conversationType string, name *string) (*db.Conversation, error)
	GetByUserID(ctx context.Context, userID int) ([]db.Conversation, error)
	GetByID(ctx context.Context, conversationID int) (*db.Conversation, error)
	GetMembers(ctx context.Context, conversationID int) ([]db.ConversationMember, error)
	GetMember(ctx context.Context, conversationID int) (*db.ConversationMember, error)
	Delete(ctx context.Context, conversationID int) error
	UpdateName(ctx context.Context, name string) (*db.Conversation, error)
	UpdateAvatarURL(ctx context.Context, avatarURL *string) (*db.Conversation, error)
	AddMembers(ctx context.Context, conversationID int, userIDs []int) error
	RemoveMember(ctx context.Context, conversationID, userID int) error
}

type conversationRepository struct {
	db *pgxpool.Pool
}

func (c *conversationRepository) AddMembers(ctx context.Context, conversationID int, userIDs []int) error {
	panic("unimplemented")
}

// Create implements [ConversationRepository].
func (c *conversationRepository) Create(ctx context.Context, createdBy int, convType string, name *string) (*db.Conversation, error) {
	panic("unimplemented")
}

// Delete implements [ConversationRepository].
func (c *conversationRepository) Delete(ctx context.Context, convID int) error {
	panic("unimplemented")
}

// GetByID implements [ConversationRepository].
func (c *conversationRepository) GetByID(ctx context.Context, convID int) (*db.Conversation, error) {
	panic("unimplemented")
}

// GetByUserID implements [ConversationRepository].
func (c *conversationRepository) GetByUserID(ctx context.Context, userID int) ([]db.Conversation, error) {
	panic("unimplemented")
}

// GetMember implements [ConversationRepository].
func (c *conversationRepository) GetMember(ctx context.Context, convID int) (*db.ConversationMember, error) {
	panic("unimplemented")
}

// GetMembers implements [ConversationRepository].
func (c *conversationRepository) GetMembers(ctx context.Context, convID int) ([]db.ConversationMember, error) {
	panic("unimplemented")
}

// RemoveMember implements [ConversationRepository].
func (c *conversationRepository) RemoveMember(ctx context.Context, convID int, userID int) error {
	panic("unimplemented")
}

// UpdateAvatarURL implements [ConversationRepository].
func (c *conversationRepository) UpdateAvatarURL(ctx context.Context, avatarURL *string) (*db.Conversation, error) {
	panic("unimplemented")
}

// UpdateName implements [ConversationRepository].
func (c *conversationRepository) UpdateName(ctx context.Context, name string) (*db.Conversation, error) {
	panic("unimplemented")
}

func NewConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{db: db}
}
