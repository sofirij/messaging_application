package response

import (
	"time"
)

type ConversationResponse struct {
	ID              int              `json:"id"`
	Name            string          `json:"name"`
	AvatarURL       *string          `json:"avatar_url"`
	Type            string           `json:"type"`
	LastMessageID   *int             `json:"last_message_id"`
	CreatedAt       time.Time        `json:"created_at"`
	CreatedBy       int              `json:"created_by"`
	LastMessageRead *int             `json:"last_message_read"`
	Members         []MemberResponse `json:"members,omitempty"`
}

type MemberResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at"`
}
