package routes

import (
	"1000words-game/internal/api/handlers"
	apiMiddleware "1000words-game/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterGameRoutes(router chi.Router, deps Dependencies) {
	gameHandler := handlers.NewGameHandler(
		deps.ProgressService,
		deps.GameService,
	)

	router.Group(func(protected chi.Router) {
		protected.Use(apiMiddleware.RequireAuth(deps.AuthService))

		protected.With(
			apiMiddleware.RequirePermission(
				deps.AuthorizationService,
				"game.start",
			),
		).Post("/start", gameHandler.Start)

		protected.With(
			apiMiddleware.RequirePermission(
				deps.AuthorizationService,
				"game.finish",
			),
		).Post("/finish", gameHandler.Finish)
	})
}
