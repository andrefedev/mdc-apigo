package auth

import (
	"time"

	"apigo/internal/platforms/validatex/validationx"
)

// INTERNO

type CodeInsertData struct {
	Code  string `db:"code"`
	Phone string `db:"phone"`
}

func (d *CodeInsertData) Normalize() {
	d.Phone = validationx.ClearString(d.Phone)
}

func (d *CodeInsertData) Validation() error {
	// Normalize
	d.Phone = validationx.ClearString(d.Phone)

	// Validation
	if !validationx.IsPhoneNumber(d.Phone) {
		return ErrInvalidPhone
	}

	return nil
}

type SessionInsertData struct {
	UserRef     string    `db:"uid"`
	TokenHash   string    `db:"token_hash"`
	DateExpires time.Time `db:"date_expires"`
}
