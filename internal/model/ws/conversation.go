package ws

import "app/internal/model/response"

// outbound
type MemberAddedPayload struct {
	ConversationID int                   `json:"conversation_id"`
	User           response.UserResponse `json:"user"`
}

type MemberRemovedPayload struct {
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}
