package db

import (
	"time"
)

type Message struct {
	ID             int        `db:"id"`
	ConversationID int        `db:"conversation_id"`
	SenderID       int        `db:"sender_id"`
	ReplyToID      *int       `db:"reply_to_id"`
	Body           *string    `db:"body"`
	DeletedAt      *time.Time `db:"deleted_at"`
	EditedAt       *time.Time `db:"edited_at"`
	CreatedAt      time.Time  `db:"created_at"`
}

type MessageAttachment struct {
	ID        int    `db:"id"`
	MessageID int    `db:"message_id"`
	Type      string `db:"type"`
	URL       string `db:"url"`
	Filename  string `db:"filename"`
	Size      int64  `db:"size"`
}

type MessageHide struct {
	MessageID int       `db:"message_id"`
	UserID    int       `db:"user_id"`
	HiddenAt  time.Time `db:"hidden_at"`
}
