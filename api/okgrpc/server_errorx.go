package okgrpc

import (
	"apigo/internal/app"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcStatusError(err error) error {
	if errors.Is(err, app.ErrInvalidCode) {
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	}
	if errors.Is(err, app.ErrCodeExpired) {
		return status.Error(codes.FailedPrecondition, "El código ingresado ya expiró")
	}
	if errors.Is(err, app.ErrCodeNotFound) {
		return status.Error(codes.NotFound, "El código solicitado no existe")
	}

	if errors.Is(err, app.ErrInvalidPhone) {
		return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	}

	if errors.Is(err, app.ErrUserNotFound) {
		return status.Error(codes.NotFound, "Usuario no encontrado")
	}

	if errors.Is(err, app.ErrSessionNotFound) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionRevoked) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionExpired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrForbidden) {
		return status.Error(codes.PermissionDenied, "No tienes permisos para realizar esta acción")
	}

	return nil
}
