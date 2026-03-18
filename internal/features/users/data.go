package users

import (
	"apigo/internal/platforms/validatex/validationx"
	"errors"
	"strings"
	"time"
)

// # USER #

type Data struct {
	Idk       *string
	Name      string
	Phone     string
	IsSuper   bool
	IsStaff   bool
	IsActive  bool
	LastLogin *time.Time
}

type FilterData struct {
	IsSuper   *bool
	IsStaff   *bool
	IsActive  *bool
	FlatQuery *string
}

func (r *FilterData) Normalize() {
	if r.FlatQuery != nil {
		r.FlatQuery = new(validationx.ClearString(*r.FlatQuery))
	}
}

func (r *FilterData) Validation() error {
	return nil
}

type PagingData struct {
	Limit  int32
	Offset int32
}

func (r *PagingData) Normalize() {
	if r.Limit == 0 {
		r.Limit = 20
	}

}

func (r *PagingData) Validation() error {
	return nil
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

func (r *AddrData) Normalize() {

}

func (r *AddrData) Validation(paths []string) error {
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
