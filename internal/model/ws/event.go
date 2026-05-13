package ws

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventMessageSend     = "message.send"
	EventMessageDelete   = "message.delete"
	EventTypingStart     = "typing.start"
	EventTypingStop      = "typing.stop"
	EventMessageNew      = "message.new"
	EventMessageDeleted  = "message.deleted"
	EventUserOnline      = "user.online"
	EventUserOffline     = "user.offline"
	EventUserTypingStart = "user.typing.start" // outbound
	EventUserTypingStop  = "user.typing.stop"  // outbound
	EventError           = "error"
)
