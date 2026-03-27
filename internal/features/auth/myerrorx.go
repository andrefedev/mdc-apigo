package auth

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcStatusError(err error) error {

	if errors.Is(err, ErrInvalidCode) {
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	}
	if errors.Is(err, ErrCodeExpired) {
		return status.Error(codes.FailedPrecondition, "El código ingresado ya expiró")
	}
	if errors.Is(err, ErrCodeNotFound) {
		return status.Error(codes.NotFound, "El código solicitado no existe")
	}

	if errors.Is(err, ErrInvalidPhone) {
		return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	}

	if errors.Is(err, ErrUserNotFound) {
		return status.Error(codes.NotFound, "Usuario no encontrado")
	}

	if errors.Is(err, ErrSessionNotFound) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, ErrSessionRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, ErrSessionRevoked) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, ErrSessionExpired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}

	return nil
}
