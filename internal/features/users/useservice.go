package users

import (
	"context"
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

func (s *Service) Get(ctx context.Context, ref string) (*User, error) {
	const op = "Users.Service.Get"

	user, err := s.deps.UserRepository.Select(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) GetAll(ctx context.Context, filter *FilterInput, paging *PagingInput) ([]*User, error) {
	const op = "Users.Service.GetAll"

	// aqui se convierten
	f := _NewFilterData(filter)
	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	p := _NewPagingData(paging)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	users, err := s.deps.UserRepository.SelectAll(ctx, f, p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
