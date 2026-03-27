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
	Ref         string     `db:"id"`
	UserRef     string     `db:"uid"`
	TokenHash   string     `db:"token_hash"`
	LastUsedAt  *time.Time `db:"last_used_at"`
	DateExpired time.Time  `db:"date_expired"`
	DateCreated time.Time  `db:"date_created"`
	DateRevoked *time.Time `db:"date_revoked"`
}
