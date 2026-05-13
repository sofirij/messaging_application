package request

import ()

type UserAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAvatarRequest struct {
	AvatarURL *string `json:"avatar_url"`
}

type UserUsernameRequest struct {
	Username string `json:"username"`
}

type UserPasswordRequest struct {
	Password string `json:"password"`
}
