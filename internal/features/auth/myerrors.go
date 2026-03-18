package auth

import (
	"apigo/internal/platforms/apperr"
	"errors"
)

var errInvalidPhone = errors.New("lookups validation failed")

func ErrInvalidPhone(cause error) error {
	if cause == nil {
		cause = errInvalidPhone
	}

	return apperr.ValidationPublic(
		"Auth.ErrInvalidPhone",
		"auth.invalid_phone",
		"El número de teléfono no es válido",
		cause,
	)
}

func ErrInvalidCode(cause error) error {
	return apperr.ValidationPublic(
		"Auth.ErrInvalidCode",
		"auth.invalid_code",
		"El código ingresado no es válido",
		cause,
	)
}

func ErrCodeExpired(cause error) error {
	return apperr.ConflictPublic(
		"Auth.ErrCodeExpired",
		"auth.code_expired",
		"El código ya expiró",
		cause,
	)
}

// IDENTITY

func ErrIdentityNotFound(cause error) error {
	return apperr.NotFoundPublic(
		"Auth.ErrIdentityNotFound",
		"auth.identity_not_found",
		"Identidad no encontrada",
		cause,
	)
}

func ErrAuthenticationRequired(cause error) error {
	return apperr.UnauthorizedPublic(
		"Auth.ErrAuthenticationRequired",
		"auth.authentication_required",
		"Debes iniciar sesión para continuar",
		cause,
	)
}
