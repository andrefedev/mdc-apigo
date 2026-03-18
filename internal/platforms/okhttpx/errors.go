package okhttpx

import (
	"errors"
	"log/slog"
	"net/http"

	"apigo/internal/platforms/apperr"
)

func ParseError(err error) (int, apperr.PublicError) {
	if err == nil {
		return http.StatusOK, apperr.PublicError{}
	}

	if _, ok := errors.AsType[*apperr.Error](err); ok {
		status := statusFromAppErrKind(apperr.KindOf(err))
		return status, apperr.ResponseOf(err)
	}

	appErr := apperr.Internal("okhttpx.ParseError", err)
	return http.StatusInternalServerError, apperr.ResponseOf(appErr)
}

func statusFromAppErrKind(kind apperr.Kind) int {
	switch kind {
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindValidation:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func slogInternalError(r *http.Request, err error) {
	ctx := r.Context()

	if myErr, ok := errors.AsType[*apperr.Error](err); ok && myErr != nil {
		slog.ErrorContext(
			ctx,
			"internal error",
			"op", myErr.Op,
			"kind", myErr.Kind,
			"code", myErr.Code,
			"cause", myErr.Cause,
			"pathURL", r.URL.Path,
		)
		return
	}

	slog.ErrorContext(
		ctx,
		"internal error",
		"cause", err,
		"pathURL", r.URL.Path,
	)
}
