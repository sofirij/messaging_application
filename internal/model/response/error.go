package response

import ()

type ErrorDetail struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
