package okhttpx

import (
	"errors"
	"log/slog"
	"net/http"

	"apigo/internal/platforms/aerr/aerrx"
	"apigo/internal/platforms/aerr/perrx"
	"apigo/internal/platforms/apperr"
)

func ParseError(err error) (int, perrx.PublicError) {
	if err == nil {
		return http.StatusOK, perrx.PublicError{}
	}

	if _, ok := errors.AsType[*apperr.Error](err); ok {
		status := statusFromAppErrKind(apperr.KindOf(err))
		payload := apperr.ResponseOf(err)
		return status, perrx.PublicError{Code: payload.Code, Body: payload.Body}
	}

	kind := aerrx.KindOf(err)
	return statusFromKind(kind), perrx.FromError(err)
}

// HELPERS

func statusFromKind(kind aerrx.Kind) int {
	switch kind {
	case aerrx.KindNotFound:
		return http.StatusNotFound
	case aerrx.KindValidation:
		return http.StatusBadRequest
	case aerrx.KindUnauthorized:
		return http.StatusUnauthorized
	case aerrx.KindForbidden:
		return http.StatusForbidden
	case aerrx.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
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

	if myErr, ok := errors.AsType[*aerrx.Error](err); ok && myErr != nil {
		slog.ErrorContext(
			ctx,
			"internal error",
			"oper", myErr.Oper,
			"kind", myErr.Kind,
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
