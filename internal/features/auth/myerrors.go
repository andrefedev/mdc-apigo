package auth

import (
	"errors"
)

var (
	ErrInvalidCode  = errors.New("auth invalid code")
	ErrInvalidPhone = errors.New("invalid phone number")

	ErrCodeExpired  = errors.New("auth code expired")
	ErrCodeNotFound = errors.New("auth code not found")

	ErrUserNotFound = errors.New("auth user not found")

	ErrSessionNotFound = errors.New("auth session not found")

	ErrSessionExpired  = errors.New("auth session expired")
	ErrSessionRevoked  = errors.New("auth session revoked")
	ErrSessionRequired = errors.New("auth session required")
	ErrForbidden       = errors.New("auth forbidden")
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

func WrapSessionExpired(cause error) error {
	if cause == nil {
		return ErrSessionExpired
	}
	return errors.Join(ErrSessionExpired, cause)
}

func WrapSessionRevoked(cause error) error {
	if cause == nil {
		return ErrSessionRevoked
	}
	return errors.Join(ErrSessionRevoked, cause)
}

func WrapSessionRequired(cause error) error {
	if cause == nil {
		return ErrSessionRequired
	}
	return errors.Join(ErrSessionRequired, cause)
}

func WrapForbidden(cause error) error {
	if cause == nil {
		return ErrForbidden
	}
	return errors.Join(ErrForbidden, cause)
}
