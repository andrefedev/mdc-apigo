package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"errors"
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

type UserUpdateData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func NewUserUpdateData(input *UserUpdateData) *UserUpdateData {
	if input == nil {
		return &UserUpdateData{}
	}
	return &UserInsertData{
		Name:     input.Name,
		Phone:    input.Phone,
		IsSuper:  input.IsSuper,
		IsStaff:  input.IsStaff,
		IsActive: input.IsActive,
	}

}

// # USER ADDR DATA #

type AddrData struct {
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

func (r *AddrData) Validate(paths []string) error {
	if paths == nil || len(paths) == 0 {
		paths = []string{
			"pid", "lat", "lng", "name", "cmna", "route", "street",
			"neighb", "locality", "sublocal", "address1", "address2",
		}
	}

	for _, path := range paths {
		switch strings.TrimSpace(path) {
		case "pid":
			if r.Pid == "" {
				return errors.New("el place id es un campo obligatorio")
			}
		case "lat":
			if r.Lat == 0 {
				return errors.New("la latitud es un campo obligatorio")
			}
		case "lng":
			if r.Lng == 0 {
				return errors.New("la longitud es un campo obligatorio")
			}
		case "route":
			if r.Route == "" {
				return errors.New("el route es un campo obligatorio")
			}
		case "street":
			if r.Street == "" {
				return errors.New("el street es un campo obligatorio")
			}
		}
	}

	return nil
}
