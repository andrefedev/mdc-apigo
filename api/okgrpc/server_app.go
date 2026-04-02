package okgrpc

import (
	"apigo/internal/features/users"
	"context"

	"apigo/internal/features/auth"
	v1 "apigo/protobuf/gen/v1"
)

func (s Server) Code(ctx context.Context, req *v1.CodeReq) (*v1.CodeRes, error) {
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

func (s Server) CodeDetail(ctx context.Context, req *v1.CodeDetailReq) (*v1.CodeDetailRes, error) {
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

func (s Server) CodeVerify(ctx context.Context, req *v1.CodeVerifyReq) (*v1.CodeVerifyRes, error) {
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

// ""

func (s Server) Userme(ctx context.Context, _ *v1.UsermeReq) (*v1.UsermeRes, error) {
	session, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.UserService.Get(ctx, session.UserRef)
	if err != nil {
		return nil, err
	}

	return &v1.UsermeRes{User: user.ToProto()}, nil
}

func (s Server) UserListAll(ctx context.Context, req *v1.UserListAllReq) (*v1.UserListAllRes, error) {
	// filter
	f := req.GetFilter()
	filter := users.NewFilterDataInput(f)
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	// pagination
	p := req.GetPaging()
	paging := users.NewPagingInput(p)
	if err := paging.Validate(); err != nil {
		return nil, err
	}

	result, err := s.UserService.GetAll(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	// CONVERtIR USERS
	userspb := make([]*v1.User, 0, len(result))
	for i := range result {
		userspb = append(userspb, result[i].ToProto())
	}

	return &v1.UserListAllRes{Users: userspb}, nil
}
