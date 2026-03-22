package auth

import (
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/apperr"
	"apigo/internal/platforms/cryptox"
	"context"
)

type Service struct {
	deps ServiceDeps
}

type ServiceDeps struct {
	AuthRepository *Repository
	MessageService *messages.Service
	// UserRepository *users.Repo
	// WhatsApp *WhatsApp.CloudAPIClient
	// JWT *authn.JWT
}

func NewService(deps ServiceDeps) *Service {
	return &Service{deps: deps}
}

func (s *Service) Code(ctx context.Context, input *codeInput) (string, string, error) {
	oper := "Auth.Service.Code"

	// OTP + challenge
	code, err := cryptox.GenerateRandomNumberString(6)
	if err != nil {
		return "", "", apperr.Internal(oper, err)
	}

	// InsertData
	data := &CodeInsertData{
		Code:  code,
		Phone: input.Phone,
	}

	data.Normalize()
	if err := data.Validation(); err != nil {
		return "", "", ErrInvalidPhone(err)
	}

	ref, err := s.deps.AuthRepository.CodeInsert(ctx, data)
	if err != nil {
		return "", "", apperr.Wrap(oper, err)
	}

	// Send Code Verification...
	templ := &messages.TemplateMessageRequest{
		To:   data.Phone,
		Type: messages.TypeTemplate,
		Template: &messages.TemplContent{
			Name: "verify_code",
			Language: messages.TemplLang{
				Code: "es_CO",
			},
			Components: []messages.TemplComp{
				{
					Type: "body",
					Parameters: []messages.TemplParam{
						{
							Type: "text",
							Text: new(code),
						},
					},
				},
				{
					Type:    "button",
					SubType: new("url"),
					Index:   new(0),
					Parameters: []messages.TemplParam{
						{
							Type: "text",
							Text: new(code),
						},
					},
				},
			},
		},
	}
	if err := s.deps.MessageService.SendTemplate(ctx, templ); err != nil {
		return "", "", apperr.Wrap(oper, err)
	}

	// ELIMINAR EL CODIGO DE LA BASE DE DATOS ??
	return ref, code, nil
}

func (s *Service) CodeVerify(ctx context.Context, code string) (bool, error) {
	// pasado un ref y un code debemos poder validar..
	return false, nil
}

func (s *Service) IdentityByIdToken(ctx context.Context, idToken string) (*Identity, error) {
	const op = "Auth.Service.IdentityByIdToken"

	identity, err := s.deps.AuthRepository.IdentitySelectByIdToken(ctx, idToken)
	if err != nil {
		if apperr.IsKind(err, apperr.KindNotFound) {
			return nil, ErrAuthenticationRequired(err)
		}
		return nil, apperr.Wrap(op, err)
	}

	return identity, nil
}
