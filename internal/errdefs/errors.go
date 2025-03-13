package errdefs

import (
	"errors"
	"fmt"
)

type ErrorType int

// Errors
var (
	// Permissions
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNotFound           = errors.New("not found")
	ErrGroupAlreadyExists = errors.New("group name already exists")
	ErrInvalidInput       = NewCustomError(ErrTypeInvalidInput, "")
)

// More dynamic errors
const (
	ErrTypeInvalidInput ErrorType = iota
)

var errorTypeStrings = map[ErrorType]string{
	ErrTypeInvalidInput: "ErrInvalidInput",
}

func (e ErrorType) String() string {
	if s, ok := errorTypeStrings[e]; ok {
		return s
	}
	return "ErrUnknown"
}

type CustomError struct {
	Type    ErrorType
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

func (e *CustomError) Is(target error) bool {
	t, ok := target.(*CustomError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

func NewCustomError(t ErrorType, message string) *CustomError {
	return &CustomError{
		Type:    t,
		Message: message,
	}
}
