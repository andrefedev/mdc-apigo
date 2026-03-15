package derrx

import (
	"errors"

	"apigo/internal/platforms/aerr/aerrx"
)

// Error is the public/client-facing error contract.
// It intentionally contains only code and safe message.
type Error struct {
	Code  string
	Body  string
	Cause error
}

// Error pass test
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Code == "" {
		return e.Body
	}

	return e.Code + ": " + e.Body
}

// Unwrap pass test
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

// New creates a new Error.
func New(code, body string, cause error) error {
	return &Error{
		Code:  code,
		Body:  body,
		Cause: cause,
	}
}

// NewKind creates a public error and guarantees that the chain contains
// an aerrx error. This is the preferred constructor for transport/domain
// errors that need a stable HTTP status and public payload.
func NewKind(kind aerrx.Kind, oper, code, body string, cause error) error {
	return New(code, body, aerrx.New(kind, oper, cause))
}

func Validation(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindValidation, oper, code, body, cause)
}

func Unauthorized(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindUnauthorized, oper, code, body, cause)
}

func Forbidden(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindForbidden, oper, code, body, cause)
}

func NotFound(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindNotFound, oper, code, body, cause)
}

func Conflict(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindConflict, oper, code, body, cause)
}

func Internal(oper, code, body string, cause error) error {
	return NewKind(aerrx.KindInternal, oper, code, body, cause)
}

// CodeOf returns the error code.
func CodeOf(err error) string {
	for err != nil {
		if target, ok := errors.AsType[*Error](err); ok && target != nil && target.Code != "" {
			return target.Code
		}
		err = errors.Unwrap(err)
	}
	return ""
}

// BodyOf returns the error message.
func BodyOf(err error) string {
	for err != nil {
		if target, ok := errors.AsType[*Error](err); ok && target != nil && target.Body != "" {
			return target.Body
		}
		err = errors.Unwrap(err)
	}

	return ""
}
