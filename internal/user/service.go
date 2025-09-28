package user

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"surf_challenge/internal/user/domain"
	"surf_challenge/internal/user/mapper"
	"surf_challenge/internal/user/storage"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error)
}

type userService struct {
	logger *zap.SugaredLogger
	repo   storage.Repository
}

func NewService(logger *zap.SugaredLogger, repo storage.Repository) Service {
	return &userService{
		logger: logger,
		repo:   repo,
	}
}

func (s *userService) QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error) {
	s.logger.Infow("QueryUsers called", "query", query)

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
