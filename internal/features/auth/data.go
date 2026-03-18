package auth

import (
	"apigo/internal/platforms/validatex/validationx"
	"errors"
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
	if !validationx.IsPhoneNumber(d.Phone) {
		return errors.New("el número de télefono no es válido")
	}

	return nil
}
