package okgrpc

import (
	"apigo/internal/features/auth"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcStatusError(err error) error {
	if errors.Is(err, auth.ErrInvalidCode) {
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	}
	if errors.Is(err, auth.ErrCodeExpired) {
		return status.Error(codes.FailedPrecondition, "El código ingresado ya expiró")
	}
	if errors.Is(err, auth.ErrCodeNotFound) {
		return status.Error(codes.NotFound, "El código solicitado no existe")
	}

	if errors.Is(err, auth.ErrInvalidPhone) {
		return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	}

	if errors.Is(err, auth.ErrUserNotFound) {
		return status.Error(codes.NotFound, "Usuario no encontrado")
	}

	if errors.Is(err, auth.ErrSessionNotFound) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, auth.ErrSessionRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, auth.ErrSessionRevoked) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, auth.ErrSessionExpired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, auth.ErrForbidden) {
		return status.Error(codes.PermissionDenied, "No tienes permisos para realizar esta acción")
	}

	return nil
}
