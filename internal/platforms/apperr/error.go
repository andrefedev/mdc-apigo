package apperr

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindNotFound     Kind = "not_found"
	KindInternal     Kind = "internal"
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindConflict     Kind = "conflict"
)

type Error struct {
	Op      string
	Kind    Kind
	Code    string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Code != "" && e.Message != "":
		return fmt.Sprintf("%s [%s] %s: %s", e.Op, e.Kind, e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("%s [%s]: %s", e.Op, e.Kind, e.Message)
	case e.Cause != nil:
		return fmt.Sprintf("%s [%s]: %v", e.Op, e.Kind, e.Cause)
	default:
		return fmt.Sprintf("%s [%s]", e.Op, e.Kind)
	}
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func Public(op string, kind Kind, code, message string, cause error) *Error {
	return &Error{
		Op:      op,
		Kind:    kind,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func (e *Error) WithPublic(code, message string) *Error {
	if e == nil {
		return nil
	}

	cp := *e
	cp.Code = code
	cp.Message = message
	return &cp
}

func Internal(op string, cause error) *Error {
	return Public(op, KindInternal, "", "", cause)
}

func NotFound(op string, cause error) *Error {
	return Public(op, KindNotFound, "", "", cause)
}

func Validation(op string, cause error) *Error {
	return Public(op, KindValidation, "", "", cause)
}

func Unauthorized(op string, cause error) *Error {
	return Public(op, KindUnauthorized, "", "", cause)
}

func Forbidden(op string, cause error) *Error {
	return Public(op, KindForbidden, "", "", cause)
}

func Conflict(op string, cause error) *Error {
	return Public(op, KindConflict, "", "", cause)
}

func Wrap(op string, cause error) error {
	if cause == nil {
		return nil
	}

	if appErr, ok := errors.AsType[*Error](cause); ok && appErr != nil {
		return &Error{
			Op:      op,
			Kind:    appErr.Kind,
			Code:    appErr.Code,
			Message: appErr.Message,
			Cause:   cause,
		}
	}

	return Internal(op, cause)
}

func KindOf(err error) Kind {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil {
			return appErr.Kind
		}
		err = errors.Unwrap(err)
	}

	return KindInternal
}

func CodeOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Code != "" {
			return appErr.Code
		}
		err = errors.Unwrap(err)
	}

	return ""
}

func MessageOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Message != "" {
			return appErr.Message
		}
		err = errors.Unwrap(err)
	}

	return ""
}

func OpOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Op != "" {
			return appErr.Op
		}
		err = errors.Unwrap(err)
	}

	return ""
}

func CauseOf(err error) error {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil {
			return appErr.Cause
		}
		err = errors.Unwrap(err)
	}

	return nil
}

func IsKind(err error, kind Kind) bool {
	return KindOf(err) == kind
}

type PublicError struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

func ResponseOf(err error) PublicError {
	if err == nil {
		return PublicError{}
	}

	code := CodeOf(err)
	body := MessageOf(err)
	kind := KindOf(err)

	if code == "" {
		code = string(kind)
	}

	if body == "" {
		body = fallbackMessage(kind)
	}

	return PublicError{
		Code: code,
		Body: body,
	}
}

func fallbackMessage(kind Kind) string {
	switch kind {
	case KindNotFound:
		return "No se encontró el recurso solicitado"
	case KindValidation:
		return "Los datos enviados no son válidos"
	case KindUnauthorized:
		return "Debes iniciar sesión"
	case KindForbidden:
		return "No tienes permisos para realizar esta acción"
	case KindConflict:
		return "La operación entra en conflicto con el estado actual"
	default:
		return "Ha ocurrido un error inesperado"
	}
}
