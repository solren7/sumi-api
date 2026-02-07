package errorx

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	cause   error
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Newf(code int, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s error: %s", e.Code, e.Message, e.cause)
}

func (e *Error) WithCause(cause error) *Error {
	e.cause = cause
	return e
}
