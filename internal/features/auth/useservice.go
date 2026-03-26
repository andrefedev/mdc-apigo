package auth

import (
	"apigo/internal/features/users"
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/cryptox"
	"context"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	deps ServiceDeps
}

type ServiceDeps struct {
	Repository     *Repository
	MessageService *messages.Service
	// UserRepository *users.Repo
	// WhatsApp *WhatsApp.CloudAPIClient
	// JWT *authn.JWT
}

func NewService(deps ServiceDeps) *Service {
	return &Service{deps: deps}
}

func (s *Service) Code(ctx context.Context, input *CodeInput) (string, string, error) {
	oper := "Auth.Service.Code"

	// OTP + challenge
	code, err := cryptox.GenerateRandomNumberString(6)
	if err != nil {
		return "", "", fmt.Errorf("%s: generate otp: %w", oper, err)
	}

	// InsertData
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
		return "", "", fmt.Errorf("%s: send template: %w", oper, err)
	}

	// ELIMINAR EL CODIGO DE LA BASE DE DATOS ??
	return ref, code, nil
}

func (s *Service) CodeVerify(ctx context.Context, input CodeVerifyInput) (*VerifyCodeResult, error) {
	op := "Auth.Service.CodeVerify"

	result := new(VerifyCodeResult)

	res, err := s.deps.Repository.CodeSelect(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("%s: code select: %w", op, err)
	}
	if res.Code != input.Code {
		return nil, ErrInvalidCode
	}
	if time.Now().After(res.DateExpired) {
		return nil, ErrCodeExpired
	}

	// Puedo importar User???
	// usuario existe?
	uid, err := s.deps.Repository.UserRefSelectByPhone(ctx, res.Phone)
	if err != nil {
		return nil, fmt.Errorf("%s: user ref by phone: %w", op, err)
	}

	idToken, err := cryptox.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("%s: generate id token: %w", op, err)
	}

	// crear la session...
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {

		// transaccion
		// create session
		//create table auth_sessions
		//(
		//	id           uuid default uuid_generate_v4() not null
		//constraint users_sessions_pk
		//primary key,
		//uid          uuid                            not null
		//constraint users_sessions_users_id_fk
		//references users,
		//token_hash   varchar(40),
		//last_used_at timestamp with time zone,
		//date_expires timestamp with time zone,
		//date_created timestamp with time zone,
		//date_revoked timestamp with time zone
		//);
		//
		//alter table auth_sessions
		//owner to dev;

		return nil
	}); err != nil {
		return nil, err
	}

	if err := s.deps.AuthRepository.db.WithTx(ctx, func(txCtx context.Context) error {
		code, err := s.deps.AuthRepository.CodeSelect(txCtx, input.Ref)
		if err != nil {
			return fmt.Errorf("%s: code select: %w", op, err)
		}

		if time.Now().After(code.DateExpired) {
			return ErrCodeExpired
		}
		if code.Code != input.Code {
			return ErrInvalidCode
		}

		userRef, err := s.deps.AuthRepository.UserRefSelectByPhone(txCtx, code.Phone)
		if err != nil {
			return fmt.Errorf("%s: user by phone: %w", op, err)
		}

		idToken, err := cryptox.GenerateRandomString(32)
		if err != nil {
			return fmt.Errorf("%s: generate id token: %w", op, err)
		}

		if err := s.deps.AuthRepository.UserIdTokenUpdate(txCtx, userRef, idToken); err != nil {
			return fmt.Errorf("%s: update id token: %w", op, err)
		}

		if _, err := s.deps.AuthRepository.CodeDelete(txCtx, input.Ref); err != nil {
			return fmt.Errorf("%s: delete code: %w", op, err)
		}

		result.UserRef = userRef
		result.IdToken = idToken
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) IdentityByIdToken(ctx context.Context, idToken string) (*Identity, error) {
	const op = "Auth.Service.IdentityByIdToken"

	identity, err := s.deps.AuthRepository.IdentitySelectByIdToken(ctx, idToken)
	if err != nil {
		if errors.Is(err, ErrIdentityNotFound) {
			return nil, fmt.Errorf("%s: %w", op, WrapAuthenticationRequired(err))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return identity, nil
}
