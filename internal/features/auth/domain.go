package auth

import (
	"time"

	v1 "apigo/protobuf/gen/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	TokenHash   string     `db:"token_hash"`
	LastUsedAt  *time.Time `db:"last_used_at"`
	DateExpires time.Time  `db:"date_expires"`
	DateCreated time.Time  `db:"date_created"`
	DateRevoked *time.Time `db:"date_revoked"`
}

type Identity struct {
	SessionRef  string     `db:"session_id"`
	UserRef     string     `db:"uid"`
	DateExpires time.Time  `db:"date_expires"`
	DateRevoked *time.Time `db:"date_revoked"`
	IsActive    bool       `db:"is_active"`
	IsStaff     bool       `db:"is_staff"`
	IsSuper     bool       `db:"is_super"`
}

func (i *Identity) IsAuthenticated() bool {
	return i != nil && i.SessionRef != "" && i.UserRef != ""
}

func (i *Identity) CanAccessBackoffice() bool {
	return i != nil && i.IsActive && (i.IsStaff || i.IsSuper)
}

func (i *Identity) CanManageUsers() bool {
	return i != nil && i.IsActive && i.IsSuper
}
