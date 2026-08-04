package response

import (
	"time"
)

type MessageResponse struct {
	ID             int                 `json:"id"`
	ConversationID int                 `json:"conversation_id"`
	SenderID       int                 `json:"sender_id"`
	Reply          *ReplyMetadata      `json:"reply"`
	Body           *string             `json:"body"`
	Deleted        bool                `json:"deleted"`
	CreatedAt      time.Time           `json:"created_at"`
	Attachments    []MessageAttachment `json:"attachments"`
}

type PaginatedMessageResponse struct {
	Messages   []MessageResponse `json:"messages"`
	NextCursor *int              `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
}

type ReplyMetadata struct {
	ID             int     `json:"id"`
	SenderID       int     `json:"sender_id"`
	Body           *string `json:"body"`
}

type MessageAttachment struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
