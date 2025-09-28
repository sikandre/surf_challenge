package container

import "surf_challenge/internal/user"

type AppContainer struct {
	UserService user.Service
}

func NewAppContainer() *AppContainer {
	return &AppContainer{
		UserService: user.NewService(),
	}
}
