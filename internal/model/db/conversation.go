package db

import (
	"time"
)

type Conversation struct {
	ID            int        `db:"id"`
	Type          string     `db:"type"`
	Name          *string    `db:"name"`
	AvatarURL     *string    `db:"avatar_url"`
	CreatedBy     int        `db:"created_by"`
	CreatedAt     time.Time  `db:"created_at"`
	LastMessageAt *time.Time `db:"last_message_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

type ConversationMember struct {
	ConversationID int        `db:"conversation_id"`
	UserID         int        `db:"user_id"`
	Role           string     `db:"role"`
	JoinedAt       time.Time  `db:"joined_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}
