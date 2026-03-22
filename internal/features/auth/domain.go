package auth

import (
	"time"
)

type Code struct {
	Ref         string    `db:"id"`
	Code        string    `db:"code"`
	Phone       string    `db:"phone"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

type Session struct {
	Ref         string    `db:"id"`
	UserRef     string    `db:"uid"`
	TokenHash   string    `db:"token_hash"`
	LastUsedAt  time.Time `db:"last_used_at"`
	DateExpires time.Time `db:"date_expires"`
	DateCreated time.Time `db:"date_created"`
	DateRevoked time.Time `db:"date_revoked"`
}

type Identity struct {
	UserRef  string `db:"id"`
	IdToken  string `db:"idk"`
	IsSuper  bool   `db:"is_super"`
	IsStaff  bool   `db:"is_staff"`
	IsActive bool   `db:"is_active"`
}

func (i *Identity) IsAuthenticated() bool {
	return i != nil && i.UserRef != ""
}

func (i *Identity) CanAccessBackoffice() bool {
	return i != nil && i.IsActive && (i.IsStaff || i.IsSuper)
}

func (i *Identity) CanManageUsers() bool {
	return i != nil && i.IsActive && i.IsSuper
}
