package auth

import (
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/aerr/aerrx"
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

func (s *Service) Code(ctx context.Context, phone string) (string, string, error) {
	oper := "Auth.Service.Code"

	// OTP + challenge
	code, err := cryptox.GenerateRandomNumberString(6)
	if err != nil {
		return "", "", aerrx.New(aerrx.KindInternal, oper, err)
	}

	// guardar codigo en la base de datos
	data := &CodeInsertData{Code: code, Phone: phone}
	data.Normalize()
	if err := data.Validation(); err != nil {
		return "", "", ErrInvalidPhone(err)
	}

	res, err := s.deps.AuthRepository.CodeInsert(ctx, data)
	if err != nil {
		return "", "", aerrx.Wrap(oper, err)
	}

	// Send Code Verification...
	templ := &messages.SendTemplateMessage{
		To:   data.Phone,
		Type: messages.TypeTemplate,
		Template: &messages.TemplContent{
			Name: "verify_code",
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
		//if _, cleanupErr := s.deps.AuthRepository.CodeDelete(ctx, res); cleanupErr != nil {
		//	slog.ErrorContext(ctx, "auth cleanup code after whatsapp failure", "code_ref", res, "err", cleanupErr)
		//}
		return "", "", aerrx.Wrap(oper, err)
	}

	// ELIMINAR EL CODIGO DE LA BASE DE DATOS ??

	return res, code, nil
}

func (s *Service) IdentityByIdToken(ctx context.Context, idToken string) (*Identity, error) {
	const op = "Auth.Service.IdentityByIdToken"

	identity, err := s.deps.AuthRepository.IdentitySelectByIdToken(ctx, idToken)
	if err != nil {
		if aerrx.IsKind(err, aerrx.KindNotFound) {
			return nil, ErrIdentityNotFound(err)
		}
		return nil, aerrx.Wrap(op, err)
	}

	return identity, nil
}
