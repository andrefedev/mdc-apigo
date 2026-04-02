package users

import (
	"apigo/internal/platforms/validatex/normalizex"
	"apigo/internal/platforms/validatex/validationx"
	v1 "apigo/protobuf/gen/v1"
)

// USER_PAGING_INPUT

type PagingInput struct {
	Limit  int32
	Offset int32
}

func NewPagingInput(req *v1.UserListAllReq_Paging) *PagingInput {
	if req == nil {
		return &PagingInput{}
	}

	return &PagingInput{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func (r *PagingInput) Validate() error {
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

type FilterInput struct {
	IsSuper   *bool
	IsStaff   *bool
	IsActive  *bool
	FlatQuery *string
}

func NewFilterDataInput(req *v1.UserListAllReq_Filter) *FilterInput {
	if req == nil {
		return &FilterInput{}
	}

	return &FilterInput{
		IsSuper:   req.IsSuper,
		IsStaff:   req.IsStaff,
		IsActive:  req.IsActive,
		FlatQuery: req.FlatQuery,
	}
}

func (r *FilterInput) Validate() error {
	const oper = "User.FilterInput.Validate"

	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(validationx.ClearString(*r.FlatQuery))
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

type InsertInput struct {
	Name     string
	Phone    string
	IsSuper  bool
	IsStaff  bool
	IsActive bool
}

func NewInsertInput(payload *v1.UserCreateReq_Payload) *InsertInput {
	if req == nil {
		return &InsertInput{}
	}

	return &InsertInput{
		Name:     payload.GetName(),
		Phone:    payload.GetPhone(),
		IsSuper:  payload.GetIsSuper(),
		IsStaff:  payload.GetIsStaff(),
		IsActive: payload.GetIsActive(),
	}
}

func (r *InsertInput) Validate() error {
	const op = "User.InsertInput.Validate"

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
