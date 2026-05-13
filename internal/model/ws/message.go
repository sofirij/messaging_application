package ws

import (
	"app/internal/model/request"
	"app/internal/model/response"
)

// outbound
type MessageNewPayload struct {
	Message response.MessageResponse `json:"message"`
}

type MessageDeletedPayload struct {
	MessageID int `json:"message_id"`
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
