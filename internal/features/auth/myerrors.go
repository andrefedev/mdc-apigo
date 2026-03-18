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

	return apperr.Validation("Auth.ErrInvalidPhone", cause).WithPublic(
		"auth.invalid_phone",
		"El número de teléfono no es válido",
	)
}

func ErrInvalidCode(cause error) error {
	return apperr.Validation("Auth.ErrInvalidCode", cause).WithPublic(
		"auth.invalid_code",
		"El código ingresado no es válido",
	)
}

func ErrCodeExpired(cause error) error {
	return apperr.Conflict("Auth.ErrCodeExpired", cause).WithPublic(
		"auth.code_expired",
		"El código ya expiró",
	)
}

// IDENTITY

func ErrIdentityNotFound(cause error) error {
	return apperr.NotFound("Auth.ErrIdentityNotFound", cause).WithPublic(
		"auth.identity_not_found",
		"Identidad no encontrada",
	)
}

func ErrAuthenticationRequired(cause error) error {
	return apperr.Unauthorized("Auth.ErrAuthenticationRequired", cause).WithPublic(
		"auth.authentication_required",
		"Debes iniciar sesión para continuar",
	)
}
