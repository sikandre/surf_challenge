package user

import (
	"context"
	"errors"

	"surf_challenge/internal/user/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service interface {
	QueryUsers(context.Context, domain.Query) ([]domain.User, *domain.Results, error)
}

type userService struct {
}

func NewService() Service {
	return &userService{}
}

func (s *userService) QueryUsers(context.Context, domain.Query) ([]domain.User, *domain.Results, error) {
	// Implement the logic to query users from the data source
	return nil, nil, nil
}
