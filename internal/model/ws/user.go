package ws

import (
	"time"
)

// outbound
type UserOnlinePayload struct {
	UserID int `json:"user_id"`
}

type UserOfflinePayload struct {
	UserID     int       `json:"user_id"`
	LastSeenAt time.Time `json:"last_seen_at"`
}
