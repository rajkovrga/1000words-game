package routes

import (
	"1000words-game/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterPracticeRoutes(router chi.Router, deps Dependencies) {
	practiceHandler := handlers.NewPracticeHandler(deps.PracticeService)

	router.Post("/start", practiceHandler.Start)
}
