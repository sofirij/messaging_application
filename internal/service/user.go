package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"app/internal/config"
	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/serviceerror"
	"app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req request.UserAuthRequest) error
	Login(ctx context.Context, req request.UserAuthRequest) (*response.UserResponse, string, string, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (string, string, error)
	SoftDelete(ctx context.Context, userID int) error
	UpdateAvatarURL(ctx context.Context, userID int, req request.UserAvatarRequest) (*response.UserResponse, error)
	UpdateUsername(ctx context.Context, userID int, req request.UserUsernameRequest) (*response.UserResponse, error)
	GetByID(ctx context.Context, userID int) (*response.UserResponse, error)
	SearchByUsername(ctx context.Context, query string) ([]response.UserResponse, error)
}

type userService struct {
	userRepo             repository.UserRepository
	jwtSecret            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	bcryptCost           int
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

const (
	minPasswordLength = 8
	maxPasswordLength = 72
	minUsernameLength = 3
	maxUsernameLength = 15
	queryLimit = 100
)

func (u *userService) Register(ctx context.Context, req request.UserAuthRequest) error {

	err := validateUsername(req.Username)

	if err != nil {
		return err
	}

	err = validatePassword(req.Password)

	if err != nil {
		return err
	}

	existingUser, err := u.userRepo.GetByUsername(ctx, req.Username)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// duplicate username
	if existingUser != nil {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), u.bcryptCost)
	if err != nil {
		return err
	}

	_, err = u.userRepo.Create(ctx, req.Username, string(passwordHash))
	if err != nil {
		return err
	}

	return nil
}

/*
updates the last seen at, generates the access and refresh token, stores the refresh token
*/
func (u *userService) Login(ctx context.Context, req request.UserAuthRequest) (*response.UserResponse, string, string, error) {
	user, err := u.userRepo.GetByUsername(ctx, req.Username)

	// user doesn't exist
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", "", &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	if err != nil {
		return nil, "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	// invalid password
	if err != nil {
		return nil, "", "", &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	err = u.userRepo.UpdateLastSeenAt(ctx, user.ID)

	if err != nil {
		return nil, "", "", err
	}

	accessToken, _, err := generateJWT(user.ID, u.jwtSecret, u.accessTokenDuration)

	if err != nil {
		return nil, "", "", err
	}

	refreshToken, expiresAt, err := generateJWT(user.ID, u.jwtSecret, u.refreshTokenDuration)

	if err != nil {
		return nil, "", "", err
	}

	tokenHash := hashToken(refreshToken)

	_, err = u.userRepo.CreateRefreshToken(ctx, user.ID, tokenHash, *expiresAt)

	if err != nil {
		return nil, "", "", err
	}

	resp := response.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		IsOnline:   true,
		LastSeenAt: user.LastSeenAt,
		CreatedAt:  user.CreatedAt,
	}

	return &resp, accessToken, refreshToken, nil
}

/*
delete refresh token for user on a specific device
*/
func (u *userService) Logout(ctx context.Context, token string) error {
	tokenHash := hashToken(token)
	return u.userRepo.DeleteRefreshToken(ctx, tokenHash)
}

/*
revoke expired refresh token, renew access and refresh token
*/
func (u *userService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	tokenHash := hashToken(token)

	refreshToken, err := u.userRepo.GetRefreshToken(ctx, tokenHash)

	// refresh token is expired
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", &serviceerror.Error{
			Code:    "",
			Message: "",
		}
	}

	if err != nil {
		return "", "", err
	}

	accessTokenString, _, err := generateJWT(refreshToken.UserID, u.jwtSecret, u.accessTokenDuration)

	if err != nil {
		return "", "", err
	}

	refreshTokenString, expiresAt, err := generateJWT(refreshToken.UserID, u.jwtSecret, u.refreshTokenDuration)

	if err != nil {
		return "", "", err
	}

	newTokenHash := hashToken(refreshTokenString)

	_, err = u.userRepo.CreateRefreshToken(ctx, refreshToken.UserID, newTokenHash, *expiresAt)

	if err != nil {
		return "", "", err
	}

	err = u.userRepo.DeleteRefreshToken(ctx, tokenHash)

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

/*
soft delete user and delete user's refresh tokens
*/
func (u *userService) SoftDelete(ctx context.Context, userID int) error {
	err := u.userRepo.DeleteAllRefreshTokens(ctx, userID)

	if err != nil {
		return err
	}

	return u.userRepo.SoftDelete(ctx, userID)
}

func (u *userService) UpdateAvatarURL(ctx context.Context, userID int, req request.UserAvatarRequest) (*response.UserResponse, error) {
	user, err := u.userRepo.UpdateAvatarURL(ctx, userID, req.AvatarURL)

	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		IsOnline:   true,
		LastSeenAt: user.LastSeenAt,
		CreatedAt:  user.CreatedAt,
	}

	return &resp, nil
}


func (u *userService) UpdateUsername(ctx context.Context, userID int, req request.UserUsernameRequest) (*response.UserResponse, error) {
	user, err := u.userRepo.UpdateUsername(ctx, userID, req.Username)

	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		IsOnline:   true,
		LastSeenAt: user.LastSeenAt,
		CreatedAt:  user.CreatedAt,
	}

	return &resp, nil
}

func (u *userService) GetByID(ctx context.Context, userID int) (*response.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		IsOnline:   true,
		LastSeenAt: user.LastSeenAt,
		CreatedAt:  user.CreatedAt,
	}

	return &resp, nil
}

func (u *userService) SearchByUsername(ctx context.Context, query string) ([]response.UserResponse, error) {
	users, err := u.userRepo.SearchByUsername(ctx, query, queryLimit)

	if err != nil {
		return make([]response.UserResponse, 0), err
	}

	resp := make([]response.UserResponse, 0)

	for _, user := range users {
		resp = append(resp, response.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			AvatarURL:  user.AvatarURL,
			IsOnline:   true,
			LastSeenAt: user.LastSeenAt,
			CreatedAt:  user.CreatedAt,
		})
	}

	return resp, nil
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo:             userRepo,
		jwtSecret:            cfg.JWTSecret,
		accessTokenDuration:  cfg.AccessTokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
		bcryptCost:           cfg.BcryptCost,
	}
}

func validateUsername(username string) error {
	if len(username) < minUsernameLength {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	if len(username) > maxUsernameLength {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	matched := usernameRegex.MatchString(username)

	// invalid username
	if !matched {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	if len(password) > maxPasswordLength {
		return &serviceerror.Error{
			Code:    "",
			Message: "",
			Fields:  map[string]string{},
		}
	}

	return nil
}

func generateJWT(userID int, secret string, duration time.Duration) (string, *time.Time, error) {
	iat := time.Now().Unix()
	exp := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp,
		"iat": iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", nil, err
	}

	expiresAt := time.Unix(exp, 0)

	return tokenString, &expiresAt, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
