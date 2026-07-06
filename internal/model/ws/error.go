package ws

import ()

type ErrorPayload struct {
	Message string `json:"message"`
	Ref     string `json:"ref,omitempty"`
}
