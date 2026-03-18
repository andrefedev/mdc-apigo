package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	deps HandlerDeps
}

type HandlerDeps struct {
	Service *Service
}

func NewHandler(deps HandlerDeps) http.Handler {
	r := chi.NewRouter()
	h := &Handler{deps: deps}

	r.Post("/code", h.code)
	r.Post("/verify", h.verify)
	// r.Post("/revoke", h.revoke)

	return r
}
