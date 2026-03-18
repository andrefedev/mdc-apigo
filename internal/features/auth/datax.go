package auth

import (
	"apigo/internal/platforms/apperr"
	"apigo/internal/platforms/validatex/validationx"
)

// REQUEST

type CodeRequest struct {
	Phone string `json:"phone"`
}

func (r *CodeRequest) Normalize() {
	r.Phone = validationx.ClearString(r.Phone)
}

func (r *CodeRequest) Validate() error {
	const oper = "Auth.CodeRequest.Validate"

	if !validationx.IsPhoneNumber(r.Phone) {
		return apperr.Validation(oper, nil).WithPublic(
			"auth.invalid_phone",
			"El número de teléfono no es válido",
		)
	}

	return nil
}
