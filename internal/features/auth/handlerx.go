package auth

import (
	"apigo/internal/platforms/okhttpx"
	"log"
	"net/http"
)

func (h *Handler) code(w http.ResponseWriter, r *http.Request) {
	oper := "Auth.Handlex.Code"

	ctx := r.Context()
	var req CodeRequest
	if err := okhttpx.DecodeJson(r, &req, oper); err != nil {
		okhttpx.Fail(w, r, err)
		return
	}

	req.Normalize()
	if err := req.Validate(); err != nil {
		okhttpx.Fail(w, r, err)
		return
	}

	id, _, err := h.deps.Service.Code(ctx, req.Phone)
	if err != nil {
		okhttpx.Fail(w, r, err)
		return
	}

	okhttpx.Json(w, http.StatusCreated, map[string]string{"ref": id})
}

func (h *Handler) verify(w http.ResponseWriter, r *http.Request) {
	// oper := "Auth.Handlex.Verify"

	log.Printf("POST /auth/verify headers=%v", r.Header)
}
