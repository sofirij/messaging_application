package db

import (
	"time"
)

type Conversation struct {
	ID              int       `db:"id"`
	Type            string    `db:"type"`
	Name            *string   `db:"name"`
	AvatarURL       *string   `db:"avatar_url"`
	CreatedBy       int       `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	LastMessageID   *int      `db:"last_message_id"`
	LastMessageRead *int      `db:"last_message_read"`
}

type ConversationMember struct {
	ConversationID int        `db:"conversation_id"`
	UserID         int        `db:"user_id"`
	JoinedAt       time.Time  `db:"joined_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
	AfterCursor    *int       `db:"after_cursor"`
}
