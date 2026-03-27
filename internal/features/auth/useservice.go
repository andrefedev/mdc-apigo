package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/cryptox"
)

type Service struct {
	deps ServiceDeps
}

type ServiceDeps struct {
	Repository     *Repository
	MessageService *messages.Service
}

func NewService(deps ServiceDeps) *Service {
	return &Service{deps: deps}
}

func (s *Service) Code(ctx context.Context, input *CodeInput) (string, string, error) {
	oper := "Auth.Service.Code"

	code, err := cryptox.GenerateRandomNumberString(6)
	if err != nil {
		return "", "", fmt.Errorf("%s: generate otp: %w", oper, err)
	}

	data := &CodeInsertData{
		Code:  code,
		Phone: input.Phone,
	}
	if err := data.Validation(); err != nil {
		return "", "", fmt.Errorf("%s: %w", oper, err)
	}

	ref, err := s.deps.Repository.CodeInsert(ctx, data)
	if err != nil {
		return "", "", fmt.Errorf("%s: insert code: %w", oper, err)
	}

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
		return "", "", fmt.Errorf("%s: send template: %w", oper, err)
	}

	return ref, code, nil
}

func (s *Service) CodeVerify(ctx context.Context, input *CodeVerifyInput) (string, string, error) {
	const op = "Auth.Service.CodeVerify"

	var uid string
	var idk string
	// REQUIERE SESSION >>>
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		code, err := s.deps.Repository.CodeSelect(ctx, input.Ref)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if code.Code != input.Code {
			return fmt.Errorf("%s: %w", op, WrapInvalidCode(err))
		}
		if time.Now().After(code.DateExpired) {
			return fmt.Errorf("%s: %w", op, WrapCodeExpired(err))
		}

		// Es probable que el usuario no exista
		// por lo que se debe crearlo.
		uid, err = s.deps.Repository.UserRefByPhone(ctx, code.Phone)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		idk, err = cryptox.GenerateRandomString(32)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// session insert
		if _, err := s.deps.Repository.SessionInsert(
			ctx,
			&SessionInsertData{
				UserRef:   uid,
				TokenHash: cryptox.HashIdToken(idk),
			},
		); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// SUCCESS...
		if _, err := s.deps.Repository.CodeDelete(ctx, input.Ref); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return uid, idk, nil
}

func (s *Service) CodeDetail(ctx context.Context, input *CodeDetailInput) (*Code, error) {
	const op = "Auth.Service.CodeDetail"

	res, err := s.deps.Repository.CodeSelect(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(res.DateExpired) {
		return nil, WrapCodeExpired(err)
	}

	return res, nil
}

func (s *Service) SessionByIdToken(ctx context.Context, idk string) (*Session, error) {
	const op = "Auth.Service.SessionByIdToken"

	if idk == "" {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionRequired(nil))
	}

	idk = cryptox.HashIdToken(idk)
	session, err := s.deps.Repository.SessionSelectByToken(ctx, idk)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			err = WrapSessionRequired(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// session expirada
	if time.Now().After(session.DateExpired) {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionExpired(nil))
	}

	// Session revocada
	if session.DateRevoked != nil {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionRevoked(nil))
	}

	return session, nil
}
