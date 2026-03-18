package okhttpx

import (
	"apigo/internal/platforms/apperr"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func Json(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Fail(w http.ResponseWriter, r *http.Request, err error) {
	status, payload := ParseError(err)

	if status >= http.StatusInternalServerError {
		slogInternalError(r, err)
	}

	// Opcional: Podrías loggear como "Info" o "Warn" (ej. fallos de validación, no autorizados).
	// slog.WarnContext(r.Context(), "Error de cliente", "status", status, "path", r.URL.Path, "error", err)

	Json(w, status, payload)
}

// DECODE

func DecodeJson(r *http.Request, dst any, oper string) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return mapDecodeError(oper, err)
	}

	if err := decoder.Decode(new(struct{})); err != io.EOF {
		return apperr.ValidationPublic(
			oper,
			"http.invalid_json",
			"El cuerpo JSON debe contener un solo objeto",
			err,
		)
	}

	return nil
}

func mapDecodeError(op string, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, io.EOF) {
		return apperr.ValidationPublic(
			op,
			"http.empty_body",
			"Debes enviar un cuerpo JSON",
			err,
		)
	}

	if _, ok := errors.AsType[*json.SyntaxError](err); ok {
		return apperr.ValidationPublic(
			op,
			"http.invalid_json",
			"El cuerpo JSON no es válido",
			err,
		)
	}

	if _, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		return apperr.ValidationPublic(
			op,
			"http.invalid_json",
			"El cuerpo JSON contiene tipos inválidos",
			err,
		)
	}

	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		return apperr.ValidationPublic(
			op,
			"http.unknown_field",
			"El cuerpo JSON contiene campos no permitidos",
			err,
		)
	}

	return apperr.ValidationPublic(
		op,
		"http.invalid_json",
		"No se pudo interpretar el cuerpo JSON",
		err,
	)
}
