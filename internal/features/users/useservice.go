package users

import (
	"context"
	"errors"
	"fmt"
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
		if errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
