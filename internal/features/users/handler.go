package users

import (
	"apigo/internal/features/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	deps HandlerDeps
}

type HandlerDeps struct {
	Service  *Service
	Identity *auth.Middleware
}

func NewHandler(deps HandlerDeps) http.Handler {
	r := chi.NewRouter()
	h := &Handler{deps: deps}

	r.With(
		deps.Identity.AttachIdentity,
		deps.Identity.IsAuthenticated,
	).Get("/me", h.me)

	return r
}
