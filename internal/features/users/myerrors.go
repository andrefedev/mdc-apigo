package users

import (
	"errors"
)

var (
	ErrUserNotFound           = errors.New("users user not found")
	ErrAuthenticationRequired = errors.New("users authentication required")
)

func WrapUserNotFound(cause error) error {
	if cause == nil {
		return ErrUserNotFound
	}
	return errors.Join(ErrUserNotFound, cause)
}

func WrapAuthenticationRequired(cause error) error {
	if cause == nil {
		return ErrAuthenticationRequired
	}
	return errors.Join(ErrAuthenticationRequired, cause)
}
