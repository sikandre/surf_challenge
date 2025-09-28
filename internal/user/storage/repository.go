package storage

import (
	"context"

	"surf_challenge/internal/user/storage/entity"
)

type Repository interface {
	QueryUsers(ctx context.Context) ([]*entity.User, int, error)
}

type userRepository struct{}

func NewRepository() Repository {
	return &userRepository{}
}

func (u *userRepository) QueryUsers(ctx context.Context) ([]*entity.User, int, error) {
	// TODO implement me
	panic("implement me")
}
