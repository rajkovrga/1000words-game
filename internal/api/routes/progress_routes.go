package routes

import (
	"1000words-game/internal/api/handlers"
	apiMiddleware "1000words-game/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterProgressRoutes(router chi.Router, deps Dependencies) {
	progressHandler := handlers.NewProgressHandler(deps.ProgressService)

	router.Group(func(protected chi.Router) {
		protected.Use(apiMiddleware.RequireAuth(deps.AuthService))

		protected.With(
			apiMiddleware.RequirePermission(
				deps.AuthorizationService,
				"progress.read",
			),
		).Get("/", progressHandler.Index)

		protected.With(
			apiMiddleware.RequirePermission(
				deps.AuthorizationService,
				"progress.create",
			),
		).Post("/", progressHandler.Create)
	})
}
