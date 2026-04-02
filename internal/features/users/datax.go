package users

import (
	"apigo/internal/platforms/validatex/validationx"
	v1 "apigo/protobuf/gen/v1"
)

type PagingInput struct {
	Limit  int32
	Offset int32
}

func NewPagingInput(req *v1.UserListAllReq_Paging) *PagingInput {
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

type FilterInput struct {
	IsSuper   *bool
	IsStaff   *bool
	IsActive  *bool
	FlatQuery *string
}

func NewFilterDataInput(req *v1.UserListAllReq_Filter) *FilterInput {
	return &FilterInput{
		IsSuper:   req.IsSuper,
		IsStaff:   req.IsStaff,
		IsActive:  req.IsActive,
		FlatQuery: req.FlatQuery,
	}
}

func (r *FilterInput) Validate() error {
	const oper = "User.FilterDataInput.Validate"

	// Normalize
	if r.FlatQuery == nil {
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
