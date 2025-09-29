package action

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"

	"surf_challenge/internal/action/domain"
	"surf_challenge/internal/action/mapper"
	"surf_challenge/internal/action/storage"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=action
type Service interface {
	GetActionByUserID(ctx context.Context, userID int64) ([]*domain.Action, error)
	GetNextActionProbability(ctx context.Context, action string) (map[string]string, error)
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

	actionsDomain, err := mapper.MapActionsEntToDomain(actions)
	if err != nil {
		return nil, fmt.Errorf("failed to map actions to domain: %w", err)
	}

	return actionsDomain, nil
}

func (s service) GetNextActionProbability(ctx context.Context, action string) (map[string]string, error) {
	s.logger.Infow("GetNextActionProbability called", "action", action)

	actions, err := s.repo.GetAllActions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all actions: %w", err)
	}

	domainActions, err := mapper.MapActionsEntToDomain(actions)
	if err != nil {
		return nil, fmt.Errorf("failed to map actions to domain: %w", err)
	}

	userActionsMap := make(map[int64][]*domain.Action)
	for _, act := range domainActions {
		id := int64(act.UserID)
		userActionsMap[id] = append(userActionsMap[id], act)
	}

	for _, acts := range userActionsMap {
		slices.SortFunc(
			acts,
			sortByCreatedAt(),
		)
	}

	nextActionCount := make(map[string]int)
	totalOccurrences := 0

	for _, acts := range userActionsMap {
		if len(acts) < 2 { // need curr and next action
			continue
		}

		for i := 0; i+1 < len(acts); i++ {
			if strings.EqualFold(acts[i].Type, action) {
				nextActionCount[acts[i+1].Type]++
				totalOccurrences++
			}
		}
	}

	if totalOccurrences == 0 {
		return map[string]string{}, nil
	}

	probabilityMap := make(map[string]string)

	for actType, count := range nextActionCount {
		probability := float64(count) / float64(totalOccurrences)
		probabilityMap[actType] = fmt.Sprintf("%.2f", probability)
	}

	return probabilityMap, nil
}

func sortByCreatedAt() func(i *domain.Action, j *domain.Action) int {
	return func(i, j *domain.Action) int {
		if i.CreatedAt.Equal(j.CreatedAt) { // tie breaker
			if i.ID < j.ID {
				return -1
			}

			return 1
		}

		if i.CreatedAt.Before(j.CreatedAt) {
			return -1
		}

		return 1
	}
}
