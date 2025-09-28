package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"surf_challenge/internal/api/user"
	"surf_challenge/internal/container"
)

func New(dependencies *container.AppContainer) http.Handler {
	router := chi.NewRouter()

	usersHandler := user.NewHandler(dependencies.UserService)

	router.Route(
		"/users", func(r chi.Router) {
			r.Get("/", usersHandler.GetUsers())
		},
	)

	return router
}
