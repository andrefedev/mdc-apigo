package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"fmt"
	"strings"
	"time"

	"apigo/internal/platforms/validatex/validationx"
)

// # CODE__ #

type CodeInsertData struct {
	Code  string `db:"code"`
	Phone string `db:"phone"`
}

func (d *CodeInsertData) Validate() error {
	const op = "App.CodeInsertData.Validate"

	// Normalize
	d.Phone = validationx.ClearString(d.Phone)

	// Validation
	if !validationx.IsPhoneNumber(d.Phone) {
		return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
	}

	return nil
}

// # SESSION__ #

type SessionInsertData struct {
	UserRef     string    `db:"uid"`
	TokenHash   string    `db:"token_hash"`
	DateExpired time.Time `db:"date_expired"`
}

// # USER__ #

type UserFilterData struct {
	IsSuper   *bool   `db:"is_super"`
	IsStaff   *bool   `db:"is_staff"`
	IsActive  *bool   `db:"is_active"`
	FlatQuery *string `db:"flat_query"`
}

func NewUserFilterData(input *UserFilterInput) *UserFilterData {
	if input == nil {
		return &UserFilterData{}
	}

	return &UserFilterData{
		IsSuper:   input.IsSuper,
		IsStaff:   input.IsStaff,
		IsActive:  input.IsActive,
		FlatQuery: input.FlatQuery,
	}
}

func (r *UserFilterData) Validate() error {
	const op = "App.UserFilterData.Validate"

	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(normalizex.NormalizeName(*r.FlatQuery))
	}

	// Validation

	return nil
}

type UserPagingData struct {
	Limit  int32
	Offset int32
}

func NewUserPagingData(input *UserPagingInput) *UserPagingData {
	if input == nil {
		return &UserPagingData{}
	}
	return &UserPagingData{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func (r *UserPagingData) Validate() error {
	const op = "App.UserPagingData.Validate"

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

// USER_INSERT_DATA__

type UserInsertData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func NewUserInsertData(input *UserInsertInput) *UserInsertData {
	if input == nil {
		return &UserInsertData{}
	}
	return &UserInsertData{
		Name:     input.Name,
		Phone:    input.Phone,
		IsSuper:  input.IsSuper,
		IsStaff:  input.IsStaff,
		IsActive: input.IsActive,
	}

}

func (r *UserInsertData) Validate() error {
	const op = "App.UserInsertData.Validate"

	// Normalize
	r.Name = normalizex.NormalizeName(r.Name)
	r.Phone = validationx.ClearString(r.Phone)

	// Validation

	if r.Name == "" {
		return fmt.Errorf("%s: %w", op, ErrInvalidName)
	}

	if !validationx.IsPhoneNumber(r.Phone) {
		return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
	}

	return nil
}

// USER_UPDATE_DATA__

type UserUpdateData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func NewUserUpdateData(input *UserUpdateInput) *UserUpdateData {
	if input == nil {
		return &UserUpdateData{}
	}
	return &UserUpdateData{
		Name:     input.Name,
		Phone:    input.Phone,
		IsSuper:  input.IsSuper,
		IsStaff:  input.IsStaff,
		IsActive: input.IsActive,
	}

}

func (r UserUpdateData) Validate(paths []string) error {
	const op = "App.UserUpdateData.Validate"

	for _, path := range paths {
		switch strings.TrimSpace(path) {
		case "name":
			r.Name = normalizex.NormalizeName(r.Name)
			if r.Name == "" {
				return fmt.Errorf("%s: %w", op, ErrInvalidName)
			}
		case "phone":
			r.Phone = validationx.ClearString(r.Phone)
			if !validationx.IsPhoneNumber(r.Phone) {
				return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
			}
		case "last_login":
			if r.LastLogin == nil {
				return fmt.Errorf("%s: %w", op, ErrInvalidLastLogin)
			}
		}
	}

	return nil
}

// USER_ADDR__

// USER_ADDR_INSERT_DATA__

type UserAddrInsertData struct {
	Pid       string
	Lat       float64
	Lng       float64
	Name      string
	Cmna      string
	Route     string
	Street    string
	Neighb    string
	Locality  string
	Sublocal  string
	Address1  string // casa / apto complemento
	Address2  string // instrucciones de entrega
	IsDefault bool
}

func (r *UserAddrInsertData) Validate() error {
	return nil
}

// USER_ADDR_UPDATE_DATA__

// GOOGLE_MAPS__
