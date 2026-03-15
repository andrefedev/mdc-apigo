package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	deps HandlerDeps
}

type HandlerDeps struct{}

func NewHandler(deps HandlerDeps) http.Handler {
	r := chi.NewRouter()
	h := &Handler{deps: deps}

	r.Get("/webhook", h.receive)
	r.Post("/webhook", h.verify)

	return r
}
