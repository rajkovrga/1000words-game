package routes

import (
	"1000words-game/internal/api/handlers"
	apiMiddleware "1000words-game/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterMeRoutes(router chi.Router, deps Dependencies) {
	meHandler := handlers.NewMeHandler(deps.AuthorizationService)

	router.With(apiMiddleware.RequireAuth(deps.AuthService)).
		Get("/", meHandler.Show)

	router.With(apiMiddleware.RequireAuth(deps.AuthService)).
		Get("/roles", meHandler.Roles)

	router.With(apiMiddleware.RequireAuth(deps.AuthService)).
		Get("/permissions", meHandler.Permissions)
}
