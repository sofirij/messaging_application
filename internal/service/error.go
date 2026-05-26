package service

const (
	ErrCodeUnauthorized = "unauthorized"
	ErrCodeForbidden    = "forbidden"
	ErrCodeNotFound     = "not_found"
	ErrCodeConflict     = "conflict"
	ErrCodeBadRequest   = "bad_request"
)

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}