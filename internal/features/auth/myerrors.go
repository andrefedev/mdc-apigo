package auth

import (
	"errors"
)

var (
	ErrInvalidCode = errors.New("auth invalid code")

	ErrCodeExpired      = errors.New("auth code expired")
	ErrIdentityNotFound = errors.New("auth identity not found")

	ErrAuthenticationRequired = errors.New("auth authentication required")
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

func WrapIdentityNotFound(cause error) error {
	if cause == nil {
		return ErrIdentityNotFound
	}
	return errors.Join(ErrIdentityNotFound, cause)
}

func WrapAuthenticationRequired(cause error) error {
	if cause == nil {
		return ErrAuthenticationRequired
	}
	return errors.Join(ErrAuthenticationRequired, cause)
}
