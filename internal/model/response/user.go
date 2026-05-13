package response

import (
	"time"
)

type UserResponse struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	AvatarURL  *string    `json:"avatar_url"`
	IsOnline   bool       `json:"is_online"`
	LastSeenAt *time.Time `json:"last_seen_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
