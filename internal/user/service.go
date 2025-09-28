package user

import (
	"context"
	"errors"
	"log"

	"surf_challenge/internal/user/domain"
	"surf_challenge/internal/user/mapper"
	"surf_challenge/internal/user/storage"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error)
}

type userService struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) Service {
	return &userService{
		repo: repo,
	}
}

func (s *userService) QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error) {
	log.Println("QueryUsers called with query:", query)

	users, totalResults, err := s.repo.QueryUsers(ctx)
	if err != nil {
		return nil, nil, err
	}

	usersDomain, err := mapper.MapUsersEntToDomain(users)
	if err != nil {
		return nil, nil, err
	}

	return usersDomain, &domain.Results{TotalItems: totalResults}, nil
}
