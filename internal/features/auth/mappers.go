package auth

import (
	"time"
)

type _CodeRaw struct {
	Ref         string    `db:"id"`
	Code        string    `db:"code"`
	Phone       string    `db:"phone"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

func (raw *_CodeRaw) ToModel() *Code {
	if raw == nil {
		return nil
	}

	return &Code{
		Ref:         raw.Ref,
		Code:        raw.Code,
		Phone:       raw.Phone,
		DateCreated: raw.DateCreated,
		DateExpired: raw.DateExpired,
	}
}
