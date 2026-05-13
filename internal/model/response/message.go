package response

import (
	"time"
)

type MessageResponse struct {
	ID             int                 `json:"id"`
	ConversationID int                 `json:"conversation_id"`
	SenderID       int                 `json:"sender_id"`
	ReplyToID      *int                `json:"reply_to_id"`
	Body           *string             `json:"body"`
	Deleted        bool                `json:"deleted"`
	CreatedAt      time.Time           `json:"created_at"`
	Attachments    []MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
