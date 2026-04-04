package app

import (
	"time"

	v1 "apigo/protobuf/gen/v1"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AUTH__

type Code struct {
	Ref         string    `db:"id"`
	Code        string    `db:"code"`
	Phone       string    `db:"phone"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

func (r *Code) ToProto() *v1.Code {
	return &v1.Code{
		Ref:         r.Ref,
		Phone:       r.Phone,
		DateCreated: timestamppb.New(r.DateCreated),
		DateExpired: timestamppb.New(r.DateExpired),
	}
}

type Session struct {
	Ref         string     `db:"id"`
	UserRef     string     `db:"uid"`
	IsSuper     bool       `db:"is_super"`
	IsStaff     bool       `db:"is_staff"`
	IsActive    bool       `db:"is_active"`
	TokenHash   string     `db:"token_hash"`
	DateExpired time.Time  `db:"date_expired"`
	DateCreated time.Time  `db:"date_created"`
	DateRevoked *time.Time `db:"date_revoked"`
}

func (i *Session) IsRoot() bool {
	return i != nil && i.IsActive && i.IsSuper
}

func (i *Session) IsEmployee() bool {
	return i != nil && i.IsActive && (i.IsStaff || i.IsSuper)
}

// USER__

type User struct {
	Ref        string     `db:"id"`
	Name       string     `db:"name"`
	Phone      string     `db:"phone"`
	IsStaff    bool       `db:"is_staff"`
	IsSuper    bool       `db:"is_super"`
	IsActive   bool       `db:"is_active"`
	LastLogin  *time.Time `db:"last_login"`
	DateJoined time.Time  `db:"date_joined"`
}

func (u User) ToProto() *v1.User {
	var dateJoined *timestamppb.Timestamp
	if !u.DateJoined.IsZero() {
		dateJoined = timestamppb.New(u.DateJoined)
	}

	var lastLogin *timestamppb.Timestamp
	if u.LastLogin != nil && !u.LastLogin.IsZero() {
		lastLogin = timestamppb.New(*u.LastLogin)
	}

	return &v1.User{
		Ref:        u.Ref,
		Name:       u.Name,
		Phone:      u.Phone,
		IsSuper:    u.IsSuper,
		IsStaff:    u.IsStaff,
		IsActive:   u.IsActive,
		LastLogin:  lastLogin,
		DateJoined: dateJoined,
	}
}

// USER_ADDR__

type UserAddr struct {
	Ref         string     `db:"id"`
	Pid         string     `db:"pid"`
	Lat         float64    `db:"lat"`
	Lng         float64    `db:"lng"`
	Name        string     `db:"name"`
	Cmna        string     `db:"cmna"`
	Route       string     `db:"route"`
	Street      string     `db:"street"`
	Neighb      string     `db:"neighb"`
	Locality    string     `db:"locality"`
	Sublocal    string     `db:"sublocal"`
	Address1    string     `db:"address1"` // casa / apto complemento
	Address2    string     `db:"address2"` // instrucciones de entrega
	IsDefault   bool       `db:"is_default"`
	DateCreated time.Time  `db:"date_created"`
	DateUpdated *time.Time `db:"date_updated"`
}

func (u *UserAddr) ToProto() *v1.UserAddr {
	var dateCreated *timestamp.Timestamp
	if !u.DateCreated.IsZero() {
		dateCreated = timestamppb.New(u.DateCreated)
	}

	var dateUpdated *timestamp.Timestamp
	if u.DateUpdated != nil && !u.DateUpdated.IsZero() {
		dateUpdated = timestamppb.New(*u.DateUpdated)
	}

	return &v1.UserAddr{
		Ref:         u.Ref,
		Pid:         u.Pid,
		Lat:         u.Lat,
		Lng:         u.Lng,
		Name:        u.Name,
		Cmna:        u.Cmna,
		Route:       u.Route,
		Street:      u.Street,
		Neighb:      u.Neighb,
		Locality:    u.Locality,
		Sublocal:    u.Sublocal,
		Address1:    u.Address1,
		Address2:    u.Address2,
		IsDefault:   u.IsDefault,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}
