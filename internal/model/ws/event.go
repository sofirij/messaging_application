package ws

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	// inbound
	EventTypingStart = "typing.start"
	EventTypingStop  = "typing.stop"
	EventMessageRead = "message.read"

	// outbound
	EventMessageNew          = "message.new"
	EventMessageDeleted      = "message.deleted"
	EventUserTypingStart     = "user.typing.start"
	EventUserTypingStop      = "user.typing.stop"
	EventMessageSeen         = "message.seen"
	EventUserOnline          = "user.online"
	EventUserOffline         = "user.offline"
	EventConversationNew     = "conversation.new"
	EventMessageEdited       = "message.edited"
	EventMemberAdded         = "member.added"
	EventMemberRemoved       = "member.removed"
	EventBroadcastUserStatus = "broadcast_user_status"

	EventError = "error"
)
