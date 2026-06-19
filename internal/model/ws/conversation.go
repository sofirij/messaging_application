package ws

// outbound
type MemberAddedPayload struct {
	ConversationID int `json:"conversation_id"`
	UserID int `json:"user_id"`
}

type MemberRemovedPayload struct {
	ConversationID int `json:"conversation_id"`
	UserID int `json:"user_id"`
}