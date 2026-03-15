package aerrx

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

// Error representa errores tecnicos. No contiene mensaje publico ni codigo de UI.
type Error struct {
	Kind  Kind
	Oper  string
	Cause error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Cause == nil {
		return fmt.Sprintf("%s [%s]", e.Oper, e.Kind)
	}

	return fmt.Sprintf("%s [%s]: %v", e.Oper, e.Kind, e.Cause)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func New(kind Kind, oper string, cause error) error {
	return &Error{
		Kind:  kind,
		Oper:  oper,
		Cause: cause,
	}
}

func Wrap(oper string, cause error) error {
	if cause == nil {
		return nil
	}

	kd := KindOf(cause)
	return &Error{
		Kind:  kd,
		Oper:  oper,
		Cause: cause,
	}
}

func KindOf(err error) Kind {
	for err != nil {
		if target, ok := errors.AsType[*Error](err); ok && target != nil {
			return target.Kind
		}
		err = errors.Unwrap(err)
	}

	return KindInternal
}

func OperOf(err error) string {
	for err != nil {
		if target, ok := errors.AsType[*Error](err); ok && target != nil && target.Oper != "" {
			return target.Oper
		}
		err = errors.Unwrap(err)
	}

	return ""
}

func IsKind(err error, kind Kind) bool {
	return KindOf(err) == kind
}
