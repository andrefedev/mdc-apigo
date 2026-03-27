package auth

import (
	"fmt"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/validatex/validationx"

	"github.com/google/uuid"
)

// REQUEST

type CodeInput struct {
	Phone string
}

func NewCodeInput(req *v1.CodeReq) *CodeInput {
	return &CodeInput{
		Phone: req.GetPhone(),
	}
}

func (r *CodeInput) Validate() error {
	const oper = "Auth.CodeInput.Validate"

	// Normalize
	r.Phone = validationx.ClearString(r.Phone)

	// Validation
	if !validationx.IsPhoneNumber(r.Phone) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
	}

	return nil
}

type CodeDetailInput struct {
	Ref string
}

func NewCodeDetailInput(req *v1.CodeDetailReq) *CodeDetailInput {
	return &CodeDetailInput{
		Ref: req.GetRef(),
	}
}

func (r *CodeDetailInput) Validate() error {
	const oper = "Auth.CodeDetailInput.Validate"

	// Normalize
	r.Ref = validationx.ClearString(r.Ref)

	// Validation
	if err := uuid.Validate(r.Ref); err != nil {
		return fmt.Errorf("%s: %w", oper, err)
	}

	return nil
}

// ····

type CodeVerifyInput struct {
	Ref  string
	Code string
}

func NewCodeVerifyInput(req *v1.CodeVerifyReq) *CodeVerifyInput {
	return &CodeVerifyInput{
		Ref:  req.GetRef(),
		Code: req.GetCode(),
	}
}

func (r *CodeVerifyInput) Validate() error {
	const oper = "Auth.VerifyCodeInput.Validate"
	// Normalize
	// Validation
	if !validationx.IsOneTimeCode(r.Code) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidCode)
	}

	return nil
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
