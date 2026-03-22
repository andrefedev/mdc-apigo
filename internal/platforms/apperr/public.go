package apperr

// PublicError representa la respuesta pública serializable derivada de un error.
type PublicError struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

// ResponseOf resuelve el contrato público final de un error.
//
// Si el error no define Code o Body explícitos, usa defaults derivados de Kind.
func ResponseOf(err error) PublicError {
	if err == nil {
		return PublicError{}
	}

	kind := KindOf(err)
	code := CodeOf(err)
	body := MessageOf(err)

	if code == "" {
		code = string(kind)
	}

	if body == "" {
		body = fallbackMessage(kind)
	}

	return PublicError{Code: code, Body: body}
}

// fallbackMessage retorna el mensaje público default para un Kind.
func fallbackMessage(kind Kind) string {
	switch kind {
	case KindNotFound:
		return "No se encontró el recurso solicitado"
	case KindInternal:
		return "Ha ocurrido un error interno del servidor"
	case KindConflict:
		return "La operación entra en conflicto con el estado actual"
	case KindForbidden:
		return "No tienes permisos para realizar esta acción"
	case KindValidation:
		return "Los datos enviados no son válidos"
	case KindUnauthorized:
		return "Debes iniciar sesión"
	default:
		return "Ha ocurrido un error desconocido"
	}
}
