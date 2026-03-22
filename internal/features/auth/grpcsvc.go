package auth

import (
	"context"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/okgrpcx"
)

type GrpcSvc struct {
	v1.UnimplementedAuthServiceServer
	deps GrpcSvcDeps
}

type GrpcSvcDeps struct {
	Service *Service
}

func NewGrpcSvc(deps GrpcSvcDeps) *GrpcSvc {
	return &GrpcSvc{deps: deps}
}

func (h *GrpcSvc) Code(ctx context.Context, req *v1.CodeReq) (*v1.CodeRes, error) {
	input := codeInputFromGrpc(req)
	if err := input.Validate(); err != nil {
		return nil, okgrpcx.StatusError(err)
	}

	ref, _, err := h.deps.Service.Code(ctx, input)
	if err != nil {
		return nil, okgrpcx.StatusError(err)
	}

	return &v1.CodeRes{Ref: ref}, nil
}
