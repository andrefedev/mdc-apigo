package apperr

import (
	"errors"
	"fmt"
)

// Kind clasifica el tipo semántico del error para mapeo entre capas.
type Kind string

// Kinds soportados actualmente por el contrato de aplicación.
const (
	KindNotFound     Kind = "not_found"
	KindInternal     Kind = "internal"
	KindConflict     Kind = "conflict"
	KindForbidden    Kind = "forbidden"
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
)

// Error es el error canónico de la aplicación.
//
// Op identifica la operación que originó o envolvió el error.
// Kind clasifica el error para reglas de transporte.
// Code y Body definen el contrato público opcional expuesto a clientes.
// Cause enlaza el error subyacente para diagnóstico y unwrap.
type Error struct {
	Op    string
	Kind  Kind
	Code  string
	Body  string
	Cause error
}

// Error implementa la interfaz error.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Code != "" && e.Body != "":
		return fmt.Sprintf("%s [%s] %s: %s", e.Op, e.Kind, e.Code, e.Body)
	case e.Body != "":
		return fmt.Sprintf("%s [%s]: %s", e.Op, e.Kind, e.Body)
	case e.Cause != nil:
		return fmt.Sprintf("%s [%s]: %v", e.Op, e.Kind, e.Cause)
	default:
		return fmt.Sprintf("%s [%s]", e.Op, e.Kind)
	}
}

// Unwrap expone la causa subyacente para compatibilidad con errors.Unwrap/As/Is.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

// Public construye un Error completo, incluyendo contrato público si aplica.
func Public(op string, kind Kind, code, body string, cause error) *Error {
	return &Error{
		Op:    op,
		Kind:  kind,
		Code:  code,
		Body:  body,
		Cause: cause,
	}
}

// WithPublic clona el error y le asigna Code y Body públicos.
func (e *Error) WithPublic(code, message string) *Error {
	if e == nil {
		return nil
	}

	cp := *e
	cp.Code = code
	cp.Body = message
	return &cp
}

// Internal crea un error técnico de tipo internal.
func Internal(op string, cause error) *Error {
	return Public(op, KindInternal, "", "", cause)
}

// NotFound crea un error técnico de tipo not_found.
func NotFound(op string, cause error) *Error {
	return Public(op, KindNotFound, "", "", cause)
}

// Validation crea un error técnico de tipo validation.
func Validation(op string, cause error) *Error {
	return Public(op, KindValidation, "", "", cause)
}

// Unauthorized crea un error técnico de tipo unauthorized.
func Unauthorized(op string, cause error) *Error {
	return Public(op, KindUnauthorized, "", "", cause)
}

// Forbidden crea un error técnico de tipo forbidden.
func Forbidden(op string, cause error) *Error {
	return Public(op, KindForbidden, "", "", cause)
}

// Conflict crea un error técnico de tipo conflict.
func Conflict(op string, cause error) *Error {
	return Public(op, KindConflict, "", "", cause)
}

// Wrap reanota un error existente con una nueva operación.
//
// Si la causa ya es un Error, preserva Kind, Code y Body, y encadena la causa.
// Si no lo es, la promueve a internal.
func Wrap(op string, cause error) error {
	if cause == nil {
		return nil
	}

	if appErr, ok := errors.AsType[*Error](cause); ok && appErr != nil {
		return &Error{
			Op:    op,
			Kind:  appErr.Kind,
			Code:  appErr.Code,
			Body:  appErr.Body,
			Cause: cause, // error
		}
	}

	return Internal(op, cause)
}

// KindOf retorna el primer Kind encontrado en la cadena de errores.
// Si no encuentra un Error, retorna KindInternal por default.
func KindOf(err error) Kind {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil {
			return appErr.Kind
		}
		err = errors.Unwrap(err)
	}

	return KindInternal
}

// CodeOf retorna el primer Code público no vacío encontrado en la cadena.
func CodeOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Code != "" {
			return appErr.Code
		}
		err = errors.Unwrap(err)
	}

	return ""
}

// MessageOf retorna el primer Body público no vacío encontrado en la cadena.
func MessageOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Body != "" {
			return appErr.Body
		}
		err = errors.Unwrap(err)
	}

	return ""
}

// OpOf retorna el primer Op no vacío encontrado en la cadena.
func OpOf(err error) string {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil && appErr.Op != "" {
			return appErr.Op
		}
		err = errors.Unwrap(err)
	}

	return ""
}

// IsKind reporta si KindOf(err) coincide con kind.
func IsKind(err error, kind Kind) bool {
	return KindOf(err) == kind
}

// CauseOf retorna la primera Cause asociada a un Error en la cadena.
func CauseOf(err error) error {
	for err != nil {
		if appErr, ok := errors.AsType[*Error](err); ok && appErr != nil {
			return appErr.Cause
		}
		err = errors.Unwrap(err)
	}

	return nil
}
