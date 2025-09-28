package mapper

import (
	"surf_challenge/internal/action/domain"
	"surf_challenge/internal/action/storage/entity"
)

func MapActionsEntToDomain(actionEnt []*entity.Action) []*domain.Action {
	result := make([]*domain.Action, 0, len(actionEnt))

	for _, a := range actionEnt {
		result = append(result, MapActionEntToDomain(a))
	}

	return result
}

func MapActionEntToDomain(a *entity.Action) *domain.Action {
	return &domain.Action{
		ID:         a.ID,
		Type:       a.Type,
		UserID:     a.UserID,
		TargetUser: a.TargetUser,
		CreatedAt:  a.CreatedAt,
	}
}
