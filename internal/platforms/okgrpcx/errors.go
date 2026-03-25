package okgrpcx

import (
	"errors"

	"apigo/internal/features/auth"
	"apigo/internal/features/users"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorDomain = "muydelcampo.apigo"

const (
	kindConflict     = "conflict"
	kindInternal     = "internal"
	kindNotFound     = "not_found"
	kindUnauthorized = "unauthorized"
	kindValidation   = "validation"
)

type appErrorSpec struct {
	grpcCode codes.Code
	kind     string
	code     string
	body     string
}

func StatusError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	spec := publicErrorSpec(err)
	st := status.New(spec.grpcCode, spec.code)

	info := &errdetails.ErrorInfo{
		Reason: "APP_ERROR",
		Domain: errorDomain,
		Metadata: map[string]string{
			"app_code": spec.code,
			"app_kind": spec.kind,
		},
	}

	localized := &errdetails.LocalizedMessage{
		Locale:  "es",
		Message: spec.body,
	}

	withDetails, detailsErr := st.WithDetails(info, localized)
	if detailsErr != nil {
		return st.Err()
	}

	return withDetails.Err()
}

func publicErrorSpec(err error) appErrorSpec {
	if spec, ok := appErrorSpecOf(err); ok {
		return spec
	}

	return appErrorSpec{
		grpcCode: codes.Internal,
		kind:     kindInternal,
		code:     "internal",
		body:     "Ha ocurrido un error interno del servidor",
	}
}

func appErrorSpecOf(err error) (appErrorSpec, bool) {
	switch {
	case errors.Is(err, auth.ErrInvalidPhone):
		return appErrorSpec{
			grpcCode: codes.InvalidArgument,
			kind:     kindValidation,
			code:     "auth.invalid_phone",
			body:     "El número de teléfono no es válido",
		}, true
	case errors.Is(err, auth.ErrInvalidCode):
		return appErrorSpec{
			grpcCode: codes.InvalidArgument,
			kind:     kindValidation,
			code:     "auth.invalid_code",
			body:     "El código ingresado no es válido",
		}, true
	case errors.Is(err, auth.ErrCodeExpired):
		return appErrorSpec{
			grpcCode: codes.FailedPrecondition,
			kind:     kindConflict,
			code:     "auth.code_expired",
			body:     "El código ya expiró",
		}, true
	case errors.Is(err, auth.ErrIdentityNotFound):
		return appErrorSpec{
			grpcCode: codes.NotFound,
			kind:     kindNotFound,
			code:     "auth.identity_not_found",
			body:     "Identidad no encontrada",
		}, true
	case errors.Is(err, auth.ErrAuthenticationRequired):
		return appErrorSpec{
			grpcCode: codes.Unauthenticated,
			kind:     kindUnauthorized,
			code:     "auth.authentication_required",
			body:     "Debes iniciar sesión para continuar",
		}, true
	case errors.Is(err, users.ErrUserNotFound):
		return appErrorSpec{
			grpcCode: codes.NotFound,
			kind:     kindNotFound,
			code:     "users.user_not_found",
			body:     "Usuario no encontrado",
		}, true
	case errors.Is(err, users.ErrAuthenticationRequired):
		return appErrorSpec{
			grpcCode: codes.Unauthenticated,
			kind:     kindUnauthorized,
			code:     "users.authentication_required",
			body:     "Debes iniciar sesión para continuar",
		}, true
	default:
		return appErrorSpec{}, false
	}
}
