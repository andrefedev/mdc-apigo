package okgrpcx

import (
	"errors"

	"apigo/internal/apperr"
	"apigo/internal/features/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StatusError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, apperr.ErrInvalidPhone):
		return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	case errors.Is(err, auth.ErrInvalidCode):
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	case errors.Is(err, auth.ErrCodeExpired):
		return status.Error(codes.FailedPrecondition, "El código ya expiró")
	case errors.Is(err, auth.ErrIdentityNotFound):
		return status.Error(codes.NotFound, "Identidad no encontrada")
	case errors.Is(err, auth.ErrAuthenticationRequired):
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	default:
		return status.Error(codes.Internal, "Ha ocurrido un error interno del servidor")
	}
}
