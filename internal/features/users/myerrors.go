package users

import (
	"errors"
)

var (
	ErrUserNotFound           = errors.New("users user not found")
	ErrAuthenticationRequired = errors.New("users authentication required")

	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidPhone = errors.New("invalid phone number")
)

func WrapUserNotFound(cause error) error {
	if cause == nil {
		return ErrUserNotFound
	}
	return errors.Join(ErrUserNotFound, cause)
}

func WrapLoginRequired(cause error) error {
	if cause == nil {
		return ErrAuthenticationRequired
	}
	return errors.Join(ErrAuthenticationRequired, cause)
}

func WrapInvalidName(cause error) error {
	if cause == nil {
		return ErrInvalidName
	}
	return errors.Join(ErrInvalidName, cause)
}

func WrapInvalidPhone(cause error) error {
	if cause == nil {
		return ErrInvalidPhone
	}

	return errors.Join(ErrInvalidPhone, cause)
}
