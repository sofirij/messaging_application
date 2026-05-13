package request

import ()

type ConversationCreateRequest struct {
	Type      string  `json:"type"`
	Name      *string `json:"name"`
	MemberIDs []int   `json:"member_ids"`
}

type ConversationRenameRequest struct {
	Name string `json:"name"`
}

type ConversationAvatarRequest struct {
	AvatarURL *string `json:"avatar_url"`
}

type ConversationAddMemberRequest struct {
	MemberIDs []int `json:"member_ids"`
}
