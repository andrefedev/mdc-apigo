package users

import (
	"apigo/internal/platforms/aerr/aerrx"
	"context"
)

type Service struct {
	deps ServiceDeps
}

type ServiceDeps struct {
	UserRepository *Repository
	// UserRepository *users.Repo
	// WhatsApp *WhatsApp.CloudAPIClient
	// JWT *authn.JWT
}

func NewService(deps ServiceDeps) *Service {
	return &Service{deps: deps}
}

func (s *Service) GetByRef(ctx context.Context, ref string) (*User, error) {
	const op = "Users.Service.GetByRef"

	user, err := s.deps.UserRepository.Select(ctx, ref)
	if err != nil {
		if aerrx.IsKind(err, aerrx.KindNotFound) {
			return nil, ErrUserNotFound(err)
		}

		return nil, aerrx.Wrap(op, err)
	}

	return user, nil
}
