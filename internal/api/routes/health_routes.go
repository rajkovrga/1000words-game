package routes

import (
	"1000words-game/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(router chi.Router, deps Dependencies) {
	healthHandler := handlers.NewHealthHandler(deps.DB, deps.Config)

	router.Get("/health", healthHandler.Check)
}
