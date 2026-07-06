package ws

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

type MessageEditedPayload struct {
	ConversationID int    `json:"conversation_id"`
	MessageID      int    `json:"message_id"`
	Body           string `json:"body"`
}

// inbound
type MessageReadPayload struct {
	ConversationID int `json:"conversation_id"`
	MessageID      int `json:"message_id"`
}
