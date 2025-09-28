package storage

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	_ "embed"

	"surf_challenge/internal/user/storage/entity"
)

var ErrUserNotFound = errors.New("user not found")

var (
	usersOnce  sync.Once
	usersCache []*entity.User
	totalUsers int
	errUsers   error
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=storage
type Repository interface {
	QueryUsers(ctx context.Context, id *int64, page int, size int) ([]*entity.User, int, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}

type userRepository struct{}

func NewRepository() Repository {
	return &userRepository{}
}

func (ur *userRepository) QueryUsers(_ context.Context, id *int64, page int, size int) ([]*entity.User, int, error) {
	users, totalResults, err := loadFileWithUsers()
	if err != nil {
		return nil, 0, err
	}

	if id != nil {
		for _, user := range users {
			if user.ID == *id {
				return []*entity.User{user}, 1, nil
			}
		}
	}

	offset := (page - 1) * size
	if offset >= totalResults {
		return []*entity.User{}, totalResults, nil
	}

	end := offset + size
	if end > totalResults {
		end = totalResults
	}

	return users[offset:end], totalResults, nil
}

func (ur *userRepository) GetUserByID(_ context.Context, id int64) (*entity.User, error) {
	users, _, err := loadFileWithUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

//go:embed db/users.json
var usersFile []byte

func loadFileWithUsers() ([]*entity.User, int, error) {
	usersOnce.Do(
		func() {
			usersCache, totalUsers, errUsers = parseUsersFile()
		},
	)

	return usersCache, totalUsers, errUsers
}

func parseUsersFile() ([]*entity.User, int, error) {
	var users []*entity.User

	err := json.Unmarshal(usersFile, &users)
	if err != nil {
		return nil, 0, err
	}

	return users, len(users), nil
}
