package users

import (
	"apigo/internal/platforms/apperr"
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
		if apperr.IsKind(err, apperr.KindNotFound) {
			return nil, ErrUserNotFound(err)
		}

		return nil, apperr.Wrap(op, err)
	}

	return user, nil
}
