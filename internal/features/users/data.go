package users

import (
	"apigo/internal/platforms/validatex/validationx"
	"errors"
	"strings"
	"time"
)

// # USER #

type FilterData struct {
	IsSuper   *bool   `db:"is_super"`
	IsStaff   *bool   `db:"is_staff"`
	IsActive  *bool   `db:"is_active"`
	FlatQuery *string `db:"flat_query"`
}

func _NewFilterData(input *FilterInput) *FilterData {
	return &FilterData{
		IsSuper:   input.IsSuper,
		IsStaff:   input.IsStaff,
		IsActive:  input.IsActive,
		FlatQuery: input.FlatQuery,
	}
}

func (r *FilterData) Validate() error {
	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(validationx.ClearString(*r.FlatQuery))
	}

	// Validation

	return nil
}

type PagingData struct {
	Limit  int32
	Offset int32
}

func _NewPagingData(input *PagingInput) *PagingData {
	return &PagingData{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func (r *PagingData) Validate() error {
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

type InsertData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func _NewInsertData(input *InsertInput)

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
