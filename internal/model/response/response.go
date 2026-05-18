package response

import ()

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
