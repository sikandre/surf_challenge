package container

import (
	"go.uber.org/zap"

	"surf_challenge/internal/user"
	"surf_challenge/internal/user/storage"
)

type AppContainer struct {
	UserService user.Service
}

func NewAppContainer(logger *zap.SugaredLogger) *AppContainer {
	usersRepository := storage.NewRepository()

	return &AppContainer{
		UserService: user.NewService(logger, usersRepository),
	}
}
