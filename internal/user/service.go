package user

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"surf_challenge/internal/action"
	"surf_challenge/internal/user/domain"
	"surf_challenge/internal/user/mapper"
	"surf_challenge/internal/user/storage"
)

var ErrNotFound = errors.New("not found")

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user
type Service interface {
	QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error)
	GetUserActionCount(ctx context.Context, userID int64) (int, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type userService struct {
	logger        *zap.SugaredLogger
	repo          storage.Repository
	actionService action.Service
}

func NewService(
	logger *zap.SugaredLogger,
	repo storage.Repository,
	actionService action.Service,
) Service {
	return &userService{
		logger:        logger,
		repo:          repo,
		actionService: actionService,
	}
}

func (s *userService) QueryUsers(ctx context.Context, query domain.Query) ([]*domain.User, *domain.Results, error) {
	s.logger.Infow("QueryUsers called", "query", query)

	users, totalResults, err := s.repo.QueryUsers(ctx, query.ID, query.Page, query.PageSize)
	if err != nil {
		return nil, nil, err
	}

	usersDomain, err := mapper.MapUsersEntToDomain(users)
	if err != nil {
		return nil, nil, err
	}

	return usersDomain, &domain.Results{TotalItems: totalResults}, nil
}

func (s *userService) GetUserActionCount(ctx context.Context, userID int64) (int, error) {
	s.logger.Infow("GetUserActionCount called", "userID", userID)

	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return 0, ErrNotFound
		}

		return 0, fmt.Errorf("failed to get user by ID: %w", err)
	}

	actions, err := s.actionService.GetActionByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get actions by user ID: %w", err)
	}

	return len(actions), nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	s.logger.Infow("GetUserByID called", "id", id)

	userEnt, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	userDomain, err := mapper.MapUserEntToDomain(userEnt)
	if err != nil {
		return nil, fmt.Errorf("failed to map user entity to domain: %w", err)
	}

	return userDomain, nil
}
