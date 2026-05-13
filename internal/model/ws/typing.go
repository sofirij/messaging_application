package ws

import ()

// inbound
type TypingPayloadInbound struct {
	ConversationID int `json:"conversation_id"`
}

// outbound
type TypingPayloadOutbound struct {
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}
