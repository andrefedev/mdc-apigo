package auth

import (
	"apigo/internal/platforms/httpx"
	"net/http"
)

type codeResponse struct {
	ID string `json:"id"`
}

func (h *Handler) code(w http.ResponseWriter, r *http.Request) {
	oper := "Auth.Handlex.Code"

	ctx := r.Context()
	var req CodeRequest
	if err := httpx.DecodeJson(r, &req, oper); err != nil {
		httpx.Fail(w, r, err)
		return
	}

	req.Normalize()
	if err := req.Validate(); err != nil {
		httpx.Fail(w, r, err)
		return
	}

	id, _, err := h.deps.Service.Code(ctx, req.Phone)
	if err != nil {
		httpx.Fail(w, r, err)
		return
	}

	httpx.Json(w, http.StatusCreated, codeResponse{ID: id})
}
