package ws

import (
	"app/internal/model/request"
)

// outbound
type MessageDeletedPayload struct {
	ConversationID int `json:"conversation_id"`
	MessageID      int `json:"message_id"`
}

type MessageSeenPayload struct {
	UserID         int `json:"user_id"`
	ConversationID int `json:"conversation_id"`
	MessageID      int `json:"message_id"`
}

// inbound
type MessageSendPayload struct {
	ConversationID int                         `json:"conversation_id"`
	ReplyToID      *int                        `json:"reply_to_id"`
	Body           *string                     `json:"body"`
	Attachments    []request.MessageAttachment `json:"attachments"`
}

type MessageDeletePayload struct {
	MessageID int `json:"message_id"`
}

type MessageReadPayload struct {
	ConversationID int `json:"conversation_id"`
	MessageID      int `json:"message_id"`
}
