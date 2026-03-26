package auth

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcStatusError(err error) error {
	// ############
	// VALIDATION #
	// ############

	if errors.Is(err, ErrCodeExpired) {
		return status.Error(codes.FailedPrecondition, "El código ingresado ya expiró")
	}
	if errors.Is(err, ErrInvalidCode) {
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	}

	if errors.Is(err, ErrInvalidPhone) {
		return status.Error(codes.InvalidArgument, "El número de télefono no es válido")
	}
	if errors.Is(err, ErrUserNotFound) {
		return status.Error(codes.NotFound, "Usuario no encontrado")
	}
	if errors.Is(err, ErrAuthenticationRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}

	return nil

	//switch {
	//case errors.Is(err, apperr.ErrInvalidPhone):
	//	return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	//case errors.Is(err, ErrInvalidCode):
	//	return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	//case errors.Is(err, ErrCodeExpired):
	//	return status.Error(codes.FailedPrecondition, "El código ya expiró")
	//case errors.Is(err, ErrAuthenticationRequired):
	//	return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	//default:
	//	return nil
	//}
}
