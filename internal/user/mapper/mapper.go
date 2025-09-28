package mapper

import (
	"time"

	"surf_challenge/internal/user/domain"
	"surf_challenge/internal/user/storage/entity"
)

func MapUsersEntToDomain(users []*entity.User) ([]*domain.User, error) {
	domainUsers := make([]*domain.User, len(users))
	for i, u := range users {
		model, err := MapUserEntToDomain(u)
		if err != nil {
			return nil, err
		}

		domainUsers[i] = model
	}

	return domainUsers, nil
}

func MapUserEntToDomain(u *entity.User) (*domain.User, error) {
	createdAt, err := time.Parse(time.RFC3339, u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: createdAt,
	}, nil
}
