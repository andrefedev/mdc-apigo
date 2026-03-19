package auth

import (
	"context"

	muydelcampov1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/okgrpcx"
)

type GrpcSvc struct {
	muydelcampov1.UnimplementedAuthServiceServer
	deps GrpcSvcDeps
}

type GrpcSvcDeps struct {
	Service *Service
}

func NewGrpcSvc(deps GrpcSvcDeps) *GrpcSvc {
	return &GrpcSvc{deps: deps}
}

func (h *GrpcSvc) Code(ctx context.Context, req *muydelcampov1.CodeReq) (*muydelcampov1.CodeRes, error) {
	input := CodeRequest{
		Phone: req.GetPhone(),
	}
	input.Normalize()

	if err := input.Validate(); err != nil {
		return nil, okgrpcx.StatusError(err)
	}

	ref, _, err := h.deps.Service.Code(ctx, input.Phone)
	if err != nil {
		return nil, okgrpcx.StatusError(err)
	}

	return &muydelcampov1.CodeRes{Ref: ref}, nil
}
