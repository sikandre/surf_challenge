package mapper

import (
	"time"

	"surf_challenge/internal/action/domain"
	"surf_challenge/internal/action/storage/entity"
)

func MapActionsEntToDomain(actionEnt []*entity.Action) ([]*domain.Action, error) {
	result := make([]*domain.Action, 0, len(actionEnt))

	for _, a := range actionEnt {
		action, err := MapActionEntToDomain(a)
		if err != nil {
			return nil, err
		}

		result = append(result, action)
	}

	return result, nil
}

func MapActionEntToDomain(a *entity.Action) (*domain.Action, error) {
	createdAt, err := time.Parse(time.RFC3339, a.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Action{
		ID:         a.ID,
		Type:       a.Type,
		UserID:     a.UserID,
		TargetUser: a.TargetUser,
		CreatedAt:  createdAt,
	}, nil
}
