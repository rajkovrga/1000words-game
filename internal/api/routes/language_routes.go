package routes

import (
	"1000words-game/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterLanguageRoutes(router chi.Router, deps Dependencies) {
	languageHandler := handlers.NewLanguageHandler(deps.LanguageService)

	router.Get("/", languageHandler.Index)
}
