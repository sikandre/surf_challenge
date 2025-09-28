package container

import (
	"surf_challenge/internal/user"
	"surf_challenge/internal/user/storage"
)

type AppContainer struct {
	UserService user.Service
}

func NewAppContainer() *AppContainer {
	usersRepository := storage.NewRepository()

	return &AppContainer{
		UserService: user.NewService(usersRepository),
	}
}
