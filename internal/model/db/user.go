package db

import (
	"time"
)

type User struct {
	ID         int        `db:"id"`
	Username   string     `db:"username"`
	Password   string     `db:"password"`
	AvatarURL  *string    `db:"avatar_url"`
	LastSeenAt *time.Time `db:"last_seen_at"`
	CreatedAt  time.Time  `db:"created_at"`
	DeletedAt  time.Time  `db:"deleted_at"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
