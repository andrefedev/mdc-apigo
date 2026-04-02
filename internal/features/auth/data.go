package auth

import (
	"time"

	"apigo/internal/platforms/validatex/validationx"
)

// ########
// # CODE #
// ########

type CodeInsertData struct {
	Code  string `db:"code"`
	Phone string `db:"phone"`
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

// ###########
// # SESSION #
// ###########

type SessionInsertData struct {
	UserRef     string    `db:"uid"`
	TokenHash   string    `db:"token_hash"`
	DateExpired time.Time `db:"date_expired"`
}
