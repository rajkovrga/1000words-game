package routes

import (
	"1000words-game/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterLevelRoutes(router chi.Router, deps Dependencies) {
	levelHandler := handlers.NewLevelHandler(deps.GameService)

	router.Get("/", levelHandler.Index)
}
