package users

import (
	"apigo/internal/platforms/aerr/derrx"
)

func ErrUserNotFound(cause error) error {
	return derrx.NotFound(
		"Users.ErrUserNotFound",
		"users.user_not_found",
		"Usuario no encontrado",
		cause,
	)
}

func ErrAuthenticationRequired(cause error) error {
	return derrx.Unauthorized(
		"Users.ErrAuthenticationRequired",
		"users.authentication_required",
		"Debes iniciar sesión para continuar",
		cause,
	)
}
