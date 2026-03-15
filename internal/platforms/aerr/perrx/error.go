package perrx

import (
	"apigo/internal/platforms/aerr/aerrx"
	"apigo/internal/platforms/aerr/derrx"
)

type PublicError struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

func FromError(err error) PublicError {
	if err == nil {
		return PublicError{}
	}

	code := derrx.CodeOf(err)
	body := derrx.BodyOf(err)
	kind := aerrx.KindOf(err)

	return PublicError{
		Code: fallbackCode(code, kind),
		Body: fallbackBody(body, kind),
	}
}

func fallbackCode(code string, kind aerrx.Kind) string {
	if code != "" {
		return code
	}

	return string(kind)
}

func fallbackBody(body string, kind aerrx.Kind) string {
	if body != "" {
		return body
	}

	switch kind {
	case aerrx.KindNotFound:
		return "No se encontró el recurso solicitado"
	case aerrx.KindValidation:
		return "Los datos enviados no son válidos"
	case aerrx.KindUnauthorized:
		return "Debes iniciar sesión"
	case aerrx.KindForbidden:
		return "No tienes permisos para realizar esta acción"
	case aerrx.KindConflict:
		return "La operación entra en conflicto con el estado actual"
	default:
		return "Ha ocurrido un error inesperado"
	}
}
