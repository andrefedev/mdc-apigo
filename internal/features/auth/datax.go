package auth

import (
	"apigo/internal/platforms/apperr"
	"apigo/internal/platforms/validatex/validationx"
	v1 "apigo/protobuf/gen/v1"
)

// REQUEST

type codeInput struct {
	Phone string
}

func codeInputFromGrpc(req *v1.CodeReq) *codeInput {
	return &codeInput{
		Phone: req.GetPhone(),
	}
}

func (r *codeInput) Validate() error {
	const oper = "Auth.CodeRequest.Validate"

	// Normalization
	r.Phone = validationx.ClearString(r.Phone)

	// Validation
	if !validationx.IsPhoneNumber(r.Phone) {
		return apperr.Validation(oper, nil).WithPublic(
			"auth.invalid_phone",
			"El número de teléfono no es válido",
		)
	}

	return nil
}

// ····

type VerifyCodeInput struct {
	Ref  string
	Code string
}

//func NewCodeInput(req *v1.CodeReq) *CodeInput {
//	return &CodeInput{
//		Phone: req.GetPhone(),
//	}
//}
//
//func (r *CodeInput) Normalize() {
//	r.Phone = validationx.ClearString(r.Phone)
//}
