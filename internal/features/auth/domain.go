package auth

import (
	"time"
)

type Code struct {
	Ref         string    `db:"id"`
	Code        string    `db:"code"`
	Phone       string    `db:"lookups"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

type Identity struct {
	UserRef  string `db:"id" json:"id"`
	IdToken  string `db:"idk" json:"-"`
	IsSuper  bool   `db:"is_super" json:"is_super"`
	IsStaff  bool   `db:"is_staff" json:"is_staff"`
	IsActive bool   `db:"is_active" json:"is_active"`
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
