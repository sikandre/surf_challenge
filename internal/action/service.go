package action

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"surf_challenge/internal/action/domain"
	"surf_challenge/internal/action/mapper"
	"surf_challenge/internal/action/storage"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=action
type Service interface {
	GetActionByUserID(ctx context.Context, userID int64) ([]*domain.Action, error)
}

type service struct {
	logger *zap.SugaredLogger
	repo   storage.Repository
}

func NewService(logger *zap.SugaredLogger, repo storage.Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s service) GetActionByUserID(ctx context.Context, userID int64) ([]*domain.Action, error) {
	s.logger.Infow("GetActionByUserID called", "userID", userID)

	actions, err := s.repo.GetActionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions by userID %d: %w", userID, err)
	}

	actionsDomain := mapper.MapActionsEntToDomain(actions)

	return actionsDomain, nil
}
