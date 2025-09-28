package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"surf_challenge/internal/api/user"
	"surf_challenge/internal/container"
)

func New(sugar *zap.SugaredLogger, dependencies *container.AppContainer) http.Handler {
	router := chi.NewRouter()

	usersHandler := user.NewHandler(sugar, dependencies.UserService)

	router.Route(
		"/api/v1", func(r chi.Router) {

			r.Route(
				"/users", func(r chi.Router) {
					r.Get("/", usersHandler.GetUsers())
				},
			)
		},
	)

	return router
}
