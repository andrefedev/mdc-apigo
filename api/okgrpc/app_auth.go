package okgrpc

import (
	"apigo/internal/features/auth"
	"context"

	v1 "apigo/protobuf/gen/v1"
)

func (s *Server) Code(ctx context.Context, req *v1.CodeReq) (*v1.CodeRes, error) {
	input := auth.NewCodeInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ref, _, err := s.AuthService.Code(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeRes{Ref: ref}, nil
}

func (s *Server) CodeDetail(ctx context.Context, req *v1.CodeDetailReq) (*v1.CodeDetailRes, error) {
	input := auth.NewCodeDetailInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	res, err := s.AuthService.CodeDetail(ctx, input)
	if err != nil {
		return nil, err
	}

	code := res.ToProto()
	return &v1.CodeDetailRes{Code: code}, nil
}

func (s *Server) CodeVerify(ctx context.Context, req *v1.CodeVerifyReq) (*v1.CodeVerifyRes, error) {
	input := auth.NewCodeVerifyInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	uid, idk, err := s.AuthService.CodeVerify(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeVerifyRes{UserRef: uid, AccessToken: idk}, nil
}
