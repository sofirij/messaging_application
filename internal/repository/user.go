package repository

import (
	"context"
	"time"

	"app/internal/model/db"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*db.User, error)
	GetByID(ctx context.Context, userID int) (*db.User, error)
	GetByIDs(ctx context.Context, userIDs []int) ([]db.User, error)
	GetByUsername(ctx context.Context, username string) (*db.User, error)
	UpdateAvatarURL(ctx context.Context, userID int, avatarURL *string) (*db.User, error)
	UpdateUsername(ctx context.Context, userID int, username string) (*db.User, error)
	UpdatePassword(ctx context.Context, userID int, passwordHash string) (*db.User, error)
	UpdateLastSeenAt(ctx context.Context, userID int) (*db.User, error)
	SoftDelete(ctx context.Context, userID int) error
	SearchByUsername(ctx context.Context, query string, limit int) ([]db.User, error)
	CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*db.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllRefreshTokens(ctx context.Context, userID int) error
	UpdateRefreshToken(ctx context.Context, userID int, oldTokenHash string, newTokenHash string, expiresAt time.Time) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func (u *userRepository) Create(ctx context.Context, username string, passwordHash string) (*db.User, error) {
	query := `
		INSERT INTO users
		(username, password)
		VALUES($1, $2) RETURNING *
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, username, passwordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) CreateRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens
		(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) RETURNING *
	`

	var refreshToken db.RefreshToken

	err := pgxscan.Get(ctx, u.db, &refreshToken, query, userID, tokenHash, expiresAt)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (u *userRepository) SoftDelete(ctx context.Context, userID int) error {
	query := `
		UPDATE users 
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	_, err := u.db.Exec(ctx, query, userID)

	return err
}

func (u *userRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
	`

	_, err := u.db.Exec(ctx, query, tokenHash)

	return err
}

func (u *userRepository) GetByID(ctx context.Context, userID int) (*db.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = $1
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetByIDs(ctx context.Context, userIDs []int) ([]db.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = ANY($1)
	`

	var users []db.User

	err := pgxscan.Select(ctx, u.db, &users, query, userIDs)

	if err != nil {
		return nil, err
	}

	return users, err
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*db.User, error) {
	query := `
		SELECT * FROM users
		WHERE LOWER(username) = LOWER($1)
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*db.RefreshToken, error) {
	query := `
		SELECT * FROM refresh_tokens
		WHERE token_hash = $1
		AND expires_at > NOW()
	`

	var refreshToken db.RefreshToken

	err := pgxscan.Get(ctx, u.db, &refreshToken, query, tokenHash)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (u *userRepository) SearchByUsername(ctx context.Context, q string, limit int) ([]db.User, error) {
	query := `
		SELECT * FROM users
		WHERE LOWER(username) LIKE '%' || LOWER($1) || '%'
		LIMIT $2
	`

	var users []db.User

	err := pgxscan.Select(ctx, u.db, &users, query, q, limit)

	return users, err
}

func (u *userRepository) UpdateAvatarURL(ctx context.Context, userID int, avatarURL *string) (*db.User, error) {
	query := `
		UPDATE users
		SET avatar_url = $1
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING *
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, avatarURL, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) UpdateLastSeenAt(ctx context.Context, userID int) (*db.User, error) {
	query := `
		UPDATE users
		SET last_seen_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING *
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, userID int, passwordHash string) (*db.User, error) {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING *
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, passwordHash, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) UpdateUsername(ctx context.Context, userID int, username string) (*db.User, error) {
	query := `
		UPDATE users
		SET username = $1
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING *
	`

	var user db.User

	err := pgxscan.Get(ctx, u.db, &user, query, username, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) DeleteAllRefreshTokens(ctx context.Context, userID int) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := u.db.Exec(ctx, query, userID)

	return err
}

func (u *userRepository) UpdateRefreshToken(ctx context.Context, userID int, oldTokenHash string, newTokenHash string, expiresAt time.Time) error {
	tx, err := u.db.Begin(ctx)

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT * FROM refresh_tokens
		WHERE token_hash = $1
		AND user_id = $2
		AND expires_at > NOW()
		FOR UPDATE
	`

	var token db.RefreshToken

	err = pgxscan.Get(ctx, tx, &token, query, oldTokenHash, userID)
	if err != nil {
		return err
	}

	query = `
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
	`

	_, err = tx.Exec(ctx, query, oldTokenHash)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO refresh_tokens
		(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(ctx, query, userID, newTokenHash, expiresAt)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}
