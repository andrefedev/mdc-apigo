package auth

import (
	"apigo/internal/platforms/validatex/validationx"
	v1 "apigo/protobuf/gen/v1"
	"fmt"
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
		return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
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
