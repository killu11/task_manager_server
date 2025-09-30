package entites

import "errors"

var (
	NameAlreadyUsingErr = NewInvalidCredsErr("name_already_used")
	InvalidPasswordErr  = NewInvalidCredsErr("invalid_password")
	InvalidLoginErr     = NewInvalidCredsErr("invalid_login")
	InvalidCredErr      = NewInvalidCredsErr("")
	UserNotExistErr     = errors.New("user_not_exist")
)

type InvalidCredsError struct {
	message string
}

func (e *InvalidCredsError) Error() string {
	return e.message
}

func NewInvalidCredsErr(msg string) *InvalidCredsError {
	return &InvalidCredsError{message: msg}
}
