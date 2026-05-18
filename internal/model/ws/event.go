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
	EventMessageSend   = "message.send"
	EventMessageDelete = "message.delete"
	EventTypingStart   = "typing.start"
	EventTypingStop    = "typing.stop"

	// outbound
	EventMessageNew      = "message.new"
	EventMessageDeleted  = "message.deleted"
	EventUserOnline      = "user.online"
	EventUserOffline     = "user.offline"
	EventUserTypingStart = "user.typing.start"
	EventUserTypingStop  = "user.typing.stop"

	EventError = "error"
)
