package users

import (
	"apigo/internal/platforms/apperr"
)

func ErrUserNotFound(cause error) error {
	return apperr.NotFound("Users.ErrUserNotFound", cause).WithPublic(
		"users.user_not_found",
		"Usuario no encontrado",
	)
}

func ErrAuthenticationRequired(cause error) error {
	return apperr.Unauthorized("Users.ErrAuthenticationRequired", cause).WithPublic(
		"users.authentication_required",
		"Debes iniciar sesión para continuar",
	)
}
