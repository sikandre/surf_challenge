package container

import (
	"go.uber.org/zap"

	"surf_challenge/internal/action"
	actionstorage "surf_challenge/internal/action/storage"
	"surf_challenge/internal/user"
	"surf_challenge/internal/user/storage"
)

type AppContainer struct {
	UserService   user.Service
	ActionService action.Service
}

func NewAppContainer(logger *zap.SugaredLogger) *AppContainer {
	usersRepository := storage.NewRepository()
	actionsRepository := actionstorage.NewRepository()

	actionService := action.NewService(logger, actionsRepository)

	return &AppContainer{
		UserService:   user.NewService(logger, usersRepository, actionService),
		ActionService: actionService,
	}
}
