package repository

import (
	"context"
	"time"

	"app/internal/model/db"

	_ "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*db.User, error)
	GetByID(ctx context.Context, userID int) (*db.User, error)
	GetByUsername(ctx context.Context, username string) (*db.User, error)
	UpdateAvatarURL(ctx context.Context, id int, avatarURL *string) (*db.User, error)
	UpdateUsername(ctx context.Context, id int, username string) (*db.User, error)
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
	UpdateLastSeenAt(ctx context.Context, id int) error
	Deactivate(ctx context.Context, id int) error
	SearchByUsername(ctx context.Context, query string) ([]db.User, error)

	CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*db.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

// Create implements [UserRepository].
func (repo *userRepository) Create(ctx context.Context, username string, passwordHash string) (*db.User, error) {
	panic("unimplemented")
}

// CreateRefreshToken implements [UserRepository].
func (u *userRepository) CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error) {
	panic("unimplemented")
}

// Deactivate implements [UserRepository].
func (u *userRepository) Deactivate(ctx context.Context, id int) error {
	panic("unimplemented")
}

// DeleteRefreshToken implements [UserRepository].
func (u *userRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	panic("unimplemented")
}

// GetByID implements [UserRepository].
func (u *userRepository) GetByID(ctx context.Context, userID int) (*db.User, error) {
	panic("unimplemented")
}

// GetByUsername implements [UserRepository].
func (u *userRepository) GetByUsername(ctx context.Context, username string) (*db.User, error) {
	panic("unimplemented")
}

// GetRefreshToken implements [UserRepository].
func (u *userRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*db.RefreshToken, error) {
	panic("unimplemented")
}

// SearchByUsername implements [UserRepository].
func (u *userRepository) SearchByUsername(ctx context.Context, query string) ([]db.User, error) {
	panic("unimplemented")
}

// UpdateAvatarURL implements [UserRepository].
func (u *userRepository) UpdateAvatarURL(ctx context.Context, id int, avatarURL *string) (*db.User, error) {
	panic("unimplemented")
}

// UpdateLastSeenAt implements [UserRepository].
func (u *userRepository) UpdateLastSeenAt(ctx context.Context, id int) error {
	panic("unimplemented")
}

// UpdatePassword implements [UserRepository].
func (u *userRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	panic("unimplemented")
}

// UpdateUsername implements [UserRepository].
func (u *userRepository) UpdateUsername(ctx context.Context, id int, username string) (*db.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}
