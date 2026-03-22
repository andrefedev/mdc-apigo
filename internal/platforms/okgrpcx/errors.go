package okgrpcx

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"apigo/internal/platforms/apperr"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const errorDomain = "muydelcampo.apigo"

func StatusError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[*apperr.Error](err); !ok && apperr.CodeOf(err) == "" && apperr.MessageOf(err) == "" {
		err = apperr.Internal("okgrpcx.StatusError", err)
	}

	public := apperr.ResponseOf(err)
	st := status.New(codeFromKind(apperr.KindOf(err)), public.Code)

	info := &errdetails.ErrorInfo{
		Reason: "APP_ERROR",
		Domain: errorDomain,
		Metadata: map[string]string{
			"app_code": public.Code,
			"app_kind": string(apperr.KindOf(err)),
		},
	}

	if op := apperr.OpOf(err); op != "" {
		info.Metadata["app_op"] = op
	}

	localized := &errdetails.LocalizedMessage{
		Locale:  "es",
		Message: public.Body,
	}

	withDetails, detailsErr := st.WithDetails(info, localized)
	if detailsErr != nil {
		return st.Err()
	}

	return withDetails.Err()
}

func codeFromKind(kind apperr.Kind) codes.Code {
	switch kind {
	case apperr.KindNotFound:
		return codes.NotFound
	case apperr.KindValidation:
		return codes.InvalidArgument
	case apperr.KindUnauthorized:
		return codes.Unauthenticated
	case apperr.KindForbidden:
		return codes.PermissionDenied
	case apperr.KindConflict:
		return codes.FailedPrecondition
	default:
		return codes.Internal
	}
}
