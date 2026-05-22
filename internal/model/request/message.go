package request

import ()

type MessageCreateRequest struct {
	ReplyToID   *int                `json:"reply_to_id"`
	Body        *string             `json:"body"`
	Attachments []MessageAttachment `json:"attachments,omitempty"`
}

type MessageAttachment struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
}

type MessageEditRequest struct {
	Body string `json:"body"`
}
