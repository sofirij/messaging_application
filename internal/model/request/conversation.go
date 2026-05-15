package request

import ()

type ConversationCreateRequest struct {
	Type    string  `json:"type"`
	Name    *string `json:"name"`
	UserIDs []int   `json:"user_ids"`
}

type ConversationRenameRequest struct {
	Name string `json:"name"`
}

type ConversationAvatarRequest struct {
	AvatarURL *string `json:"avatar_url"`
}

type ConversationAddMemberRequest struct {
	UserIDs []int `json:"user_ids"`
}
