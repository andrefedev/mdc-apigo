package okgrpc

import (
	"context"
	"log"

	"apigo/internal/app"
	v1 "apigo/protobuf/gen/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"googlemaps.github.io/maps"
)

func (s Server) Code(ctx context.Context, req *v1.CodeReq) (*v1.CodeRes, error) {
	input := app.NewCodeInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ref, _, err := s.useservice.Code(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeRes{Ref: ref}, nil
}

func (s Server) CodeDetail(ctx context.Context, req *v1.CodeDetailReq) (*v1.CodeDetailRes, error) {
	input := app.NewCodeDetailInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	res, err := s.useservice.CodeDetail(ctx, input)
	if err != nil {
		return nil, err
	}

	code := res.ToProto()
	return &v1.CodeDetailRes{Code: code}, nil
}

func (s Server) CodeVerify(ctx context.Context, req *v1.CodeVerifyReq) (*v1.CodeVerifyRes, error) {
	input := app.NewCodeVerifyInput(req)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	uid, idk, err := s.useservice.CodeVerify(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.CodeVerifyRes{UserRef: uid, AccessToken: idk}, nil
}

// USER__

func (s Server) Userme(ctx context.Context, _ *v1.UsermeReq) (*v1.UsermeRes, error) {
	session, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.useservice.UserDetail(ctx, session.UserRef)
	if err != nil {
		return nil, err
	}

	return &v1.UsermeRes{User: user.ToProto()}, nil
}

func (s Server) UserDetail(ctx context.Context, req *v1.UserDetailReq) (*v1.UserDetailRes, error) {
	_, err := requireStaffUser(ctx)
	if err != nil {
		return nil, err
	}

	ref := req.GetRef()
	if err := uuid.Validate(ref); err != nil {
		return nil, err
	}

	user, err := s.useservice.UserDetail(ctx, ref)
	if err != nil {
		return nil, err
	}

	return &v1.UserDetailRes{User: user.ToProto()}, nil
}

func (s Server) UserListAll(ctx context.Context, req *v1.UserListAllReq) (*v1.UserListAllRes, error) {
	// Auth
	_, err := requireStaffUser(ctx)
	if err != nil {
		return nil, err
	}

	// filter
	f := req.GetFilter()
	filter := app.NewUserFilterInput(f)
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	// pagination
	p := req.GetPaging()
	paging := app.NewUserPagingInput(p)
	if err := paging.Validate(); err != nil {
		return nil, err
	}

	result, err := s.useservice.UserListAll(ctx, filter, paging)
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

func (s Server) UserCreate(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error) {
	_, err := requireRootUser(ctx)
	if err != nil {
		return nil, err
	}

	p := req.GetPayload()
	input := app.NewUserInsertInput(p)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := s.useservice.UserCreate(ctx, input)
	if err != nil {
		return nil, err
	}

	return &v1.UserCreateRes{User: result.ToProto()}, nil
}

func (s Server) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (*v1.UserUpdateRes, error) {
	_, err := requireRootUser(ctx)
	if err != nil {
		return nil, err
	}

	ref := req.GetRef()
	if err := uuid.Validate(ref); err != nil {
		return nil, err
	}

	payload := req.GetPayload()
	if payload == nil {
		return nil, ErrInvalidPayload
	}

	mask := req.GetUpdateMask()
	mask.Normalize()
	if !mask.IsValid(payload) {
		return nil, ErrInvalidUpdateMask
	}

	paths := mask.GetPaths()
	inputData := app.NewUserUpdateInput(payload)
	if err := inputData.Validate(paths); err != nil {
		return nil, err
	}

	result, err := s.useservice.UserUpdate(ctx, ref, paths, inputData)
	if err != nil {
		return nil, err
	}

	return &v1.UserUpdateRes{User: result.ToProto()}, nil
}

// GOOGLE_MAPS__

func (s Server) PlaceDetail(ctx context.Context, req *v1.PlaceDetailReq) (*v1.PlaceDetailRes, error) {
	ref := req.GetRef()
	if len(ref) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no place ref")
	}

	token := req.GetToken()
	if len(token) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no place session token")
	}

	// parse token
	u, _ := uuid.Parse(token)
	mtoken := maps.PlaceAutocompleteSessionToken(u)

	result, err := s.vendor.Mapsx.PlaceDetails(ctx, ref, mtoken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// formatear a protobuf...
	return &v1.PlaceDetailRes{Result: result.ToProto()}, nil
}

func (s Server) ReverseGeocode(ctx context.Context, req *v1.ReverseGeocodeReq) (*v1.ReverseGeocodeRes, error) {
	lat := req.GetLat()
	lng := req.GetLng()

	result, err := s.vendor.Mapsx.ReverseGeocode(ctx, lat, lng)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// formatear a protobuf...
	return &v1.ReverseGeocodeRes{Result: result.ToProto()}, nil
}

func (s Server) PlaceAutocomplete(ctx context.Context, req *v1.PlaceAutocompleteReq) (*v1.PlaceAutocompleteRes, error) {
	query := req.GetQuery()
	if len(query) < 0 {
		return nil, status.Error(codes.InvalidArgument, "no query")
	}

	token := maps.NewPlaceAutocompleteSessionToken()
	res, err := s.vendor.Mapsx.PlaceAutocomplete(ctx, query, token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	results := make([]*v1.Prediction, len(res))
	for i, p := range res {
		results[i] = p.ToProto()
	}

	log.Printf("token: %v", uuid.UUID(token).String())

	// formatear a protobuf...
	return &v1.PlaceAutocompleteRes{Token: uuid.UUID(token).String(), Results: results}, nil
}
