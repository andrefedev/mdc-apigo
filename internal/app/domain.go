package app

import (
	"time"

	v1 "apigo/protobuf/gen/v1"

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
