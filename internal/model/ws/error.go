package ws

import ()

type ErrorPayload struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Ref     *string `json:"ref,omitempty"`
}
