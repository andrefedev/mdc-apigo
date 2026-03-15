package users

import (
	"apigo/internal/features/auth"
	"apigo/internal/platforms/httpx"
	"net/http"
)

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	identity, ok := auth.IdentityFromContext(ctx)
	if !ok || identity == nil {
		httpx.Fail(w, r, ErrAuthenticationRequired(nil))
		return
	}

	user, err := h.deps.Service.GetByRef(ctx, identity.UserRef)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}

	httpx.Json(w, http.StatusOK, user)
}
