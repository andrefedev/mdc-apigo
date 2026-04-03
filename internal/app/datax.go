package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"fmt"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/validatex/validationx"

	"github.com/google/uuid"
)

// CODE__

type CodeInput struct {
	Phone string
}

func NewCodeInput(req *v1.CodeReq) *CodeInput {
	return &CodeInput{
		Phone: req.GetPhone(),
	}
}

func (r *CodeInput) Validate() error {
	const oper = "App.CodeInput.Validate"

	// Normalize
	r.Phone = validationx.ClearString(r.Phone)

	// Validation
	if !validationx.IsPhoneNumber(r.Phone) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
	}

	return nil
}

type CodeDetailInput struct {
	Ref string
}

func NewCodeDetailInput(req *v1.CodeDetailReq) *CodeDetailInput {
	return &CodeDetailInput{
		Ref: req.GetRef(),
	}
}

func (r *CodeDetailInput) Validate() error {
	const oper = "App.CodeDetailInput.Validate"

	// Normalize
	r.Ref = validationx.ClearString(r.Ref)

	// Validation
	if err := uuid.Validate(r.Ref); err != nil {
		return fmt.Errorf("%s: %w", oper, err)
	}

	return nil
}

// ····

type CodeVerifyInput struct {
	Ref  string
	Code string
}

func NewCodeVerifyInput(req *v1.CodeVerifyReq) *CodeVerifyInput {
	return &CodeVerifyInput{
		Ref:  req.GetRef(),
		Code: req.GetCode(),
	}
}

func (r *CodeVerifyInput) Validate() error {
	const oper = "App.VerifyCodeInput.Validate"
	// Normalize
	// Validation
	if !validationx.IsOneTimeCode(r.Code) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidCode)
	}

	return nil
}

// # USER__ #

type UserPagingInput struct {
	Limit  int32
	Offset int32
}

func NewUserPagingInput(req *v1.UserListAllReq_Paging) *UserPagingInput {
	if req == nil {
		return &UserPagingInput{}
	}

	return &UserPagingInput{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func (r *UserPagingInput) Validate() error {
	// Normalize
	limit := int32(20)
	if r.Limit == 0 {
		r.Limit = limit
	}
	if r.Limit > limit {
		r.Limit = limit
	}

	// Validation

	return nil
}

// USER_FILTER_INPUT

type UserFilterInput struct {
	IsSuper   *bool
	IsStaff   *bool
	IsActive  *bool
	FlatQuery *string
}

func NewUserFilterInput(req *v1.UserListAllReq_Filter) *UserFilterInput {
	if req == nil {
		return &UserFilterInput{}
	}

	return &UserFilterInput{
		IsSuper:   req.IsSuper,
		IsStaff:   req.IsStaff,
		IsActive:  req.IsActive,
		FlatQuery: req.FlatQuery,
	}
}

func (r *UserFilterInput) Validate() error {
	const op = "App.UserFilterInput.Validate"

	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(normalizex.NormalizeName(*r.FlatQuery))
	}

	// Validation
	if r.FlatQuery != nil {
		if len(*r.FlatQuery) < 2 {
			// return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
		}
	}

	return nil
}

// USER_INSERT_INPUT

type UserInsertInput struct {
	Name     string
	Phone    string
	IsSuper  bool
	IsStaff  bool
	IsActive bool
}

func NewUserInsertInput(payload *v1.UserCreateReq_Payload) *UserInsertInput {
	if payload == nil {
		return &UserInsertInput{}
	}

	return &UserInsertInput{
		Name:     payload.GetName(),
		Phone:    payload.GetPhone(),
		IsSuper:  payload.GetIsSuper(),
		IsStaff:  payload.GetIsStaff(),
		IsActive: payload.GetIsActive(),
	}
}

func (r *UserInsertInput) Validate() error {
	const op = "App.InsertInput.Validate"

	// Normalize
	r.Name = normalizex.NormalizeName(r.Name)
	r.Phone = validationx.ClearString(r.Phone)

	// Validation

	if r.Name == "" {

	}

	if !validationx.IsPhoneNumber(r.Phone) {

	}

	return nil
}

// USER_UPDATE_INPUT

type UserUpdateInput struct {
	Name     string `db:"name"`
	Phone    string `db:"phone"`
	IsSuper  bool   `db:"is_super"`
	IsStaff  bool   `db:"is_staff"`
	IsActive bool   `db:"is_active"`
}

func NewUserUpdateInput(payload *v1.UserUpdateReq_Payload) *UserUpdateInput {
	if payload == nil {
		return &UserUpdateInput{}
	}

	return &UserUpdateInput{
		Name:     payload.GetName(),
		Phone:    payload.GetPhone(),
		IsSuper:  payload.GetIsSuper(),
		IsStaff:  payload.GetIsStaff(),
		IsActive: payload.GetIsActive(),
	}
}

func (r *UserUpdateInput) Validate() error {
	const op = "App.UserUpdateInput.Validate"

	// Normalize
	r.Name = normalizex.NormalizeName(r.Name)
	r.Phone = validationx.ClearString(r.Phone)

	// Validation
	if r.Name != "" {
		if len(r.Name) < 2 {
			return fmt.Errorf("%s: %w", op, ErrInvalidName)
		}
	}

	if r.Phone != "" {
		if !validationx.IsPhoneNumber(r.Phone) {
			return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
		}
	}

	return nil
}
