package errors

import (
	"fmt"
)

type ErrorType string

const (
	ErrInternal     ErrorType = "INTERNAL"
	ErrNotFound     ErrorType = "NOT_FOUND"
	ErrUnauthorized ErrorType = "UNAUTHORIZED"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func New(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}
