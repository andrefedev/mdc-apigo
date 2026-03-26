package auth

import (
	"errors"
)

var (
	ErrInvalidCode  = errors.New("auth invalid code")
	ErrInvalidPhone = errors.New("invalid phone number")

	ErrCodeExpired  = errors.New("auth code expired")
	ErrCodeNotFound = errors.New("auth code not found")

	ErrAuthenticationRequired = errors.New("auth authentication required")

	ErrUserNotFound    = errors.New("auth user not found")
	ErrSessionNotFound = errors.New("auth session not found")
)

func WrapInvalidCode(cause error) error {
	if cause == nil {
		return ErrInvalidCode
	}
	return errors.Join(ErrInvalidCode, cause)
}

func WrapCodeExpired(cause error) error {
	if cause == nil {
		return ErrCodeExpired
	}
	return errors.Join(ErrCodeExpired, cause)
}

func WrapCodeNotFound(cause error) error {
	if cause == nil {
		return ErrCodeNotFound
	}
	return errors.Join(ErrCodeNotFound, cause)
}

func WrapAuthenticationRequired(cause error) error {
	if cause == nil {
		return ErrAuthenticationRequired
	}
	return errors.Join(ErrAuthenticationRequired, cause)
}

func WrapUserNotFound(cause error) error {
	if cause == nil {
		return ErrUserNotFound
	}
	return errors.Join(ErrUserNotFound, cause)
}

func WrapSessionNotFound(cause error) error {
	if cause == nil {
		return ErrSessionNotFound
	}
	return errors.Join(ErrSessionNotFound, cause)
}
