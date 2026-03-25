package users

import (
	"apigo/internal/features/auth"
	"apigo/internal/platforms/okhttpx"
	"net/http"
)

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	identity, ok := auth.IdentityFromContext(ctx)
	if !ok || identity == nil {
		okhttpx.Fail(w, r, ErrAuthenticationRequired)
		return
	}

	user, err := h.deps.Service.GetByRef(ctx, identity.UserRef)
	if err != nil {
		okhttpx.Fail(w, r, err)
		return
	}

	okhttpx.Json(w, http.StatusOK, user)
}
