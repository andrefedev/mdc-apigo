package auth

import (
	"context"

	v1 "apigo/protobuf/gen/v1"
)

type Handler struct {
	v1.UnimplementedAuthServiceServer
	deps HandlerDeps
}

type HandlerDeps struct {
	Service *Service
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{deps: deps}
}

func (h *Handler) Code(ctx context.Context, req *v1.CodeReq) (*v1.CodeRes, error) {
	input := NewCodeInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ref, _, err := h.deps.Service.Code(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeRes{Ref: ref}, nil
}

func (h *Handler) CodeDetail(ctx context.Context, req *v1.CodeDetailReq) (*v1.CodeDetailRes, error) {
	input := NewCodeDetailInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	res, err := h.deps.Service.CodeDetail(ctx, input)
	if err != nil {
		return nil, err
	}

	code := res.ToProto()
	return &v1.CodeDetailRes{Code: code}, nil
}

func (h *Handler) CodeVerify(ctx context.Context, req *v1.CodeVerifyReq) (*v1.CodeVerifyRes, error) {
	input := NewCodeVerifyInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	uid, idk, err := h.deps.Service.CodeVerify(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeVerifyRes{UserRef: uid, AccessToken: idk}, nil
}
