package response

import (
	"time"
)

type ConversationResponse struct {
	ID                            int              `json:"id"`
	Name                          string           `json:"name"`
	AvatarURL                     *string          `json:"avatar_url"`
	Type                          string           `json:"type"`
	CreatedAt                     time.Time        `json:"created_at"`
	CreatedBy                     int              `json:"created_by"`
	LastMessageReadByUser         *int             `json:"last_message_read_by_user"`
	LastMessageReadInConversation *int             `json:"last_message_read_in_conversation"`
	LastMessageSentInConversation *MessageResponse `json:"last_message_sent_in_conversation"`
}
