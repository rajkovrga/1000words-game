package routes

import "github.com/go-chi/chi/v5"

func RegisterV1Routes(router chi.Router, deps Dependencies) {
	RegisterHealthRoutes(router, deps)

	router.Route("/auth", func(auth chi.Router) {
		RegisterAuthRoutes(auth, deps)
	})

	router.Route("/me", func(me chi.Router) {
		RegisterMeRoutes(me, deps)
	})

	router.Route("/languages", func(languages chi.Router) {
		RegisterLanguageRoutes(languages, deps)
	})

	router.Route("/levels", func(levels chi.Router) {
		RegisterLevelRoutes(levels, deps)
	})

	router.Route("/practice", func(practice chi.Router) {
		RegisterPracticeRoutes(practice, deps)
	})

	router.Route("/progress", func(progress chi.Router) {
		RegisterProgressRoutes(progress, deps)
	})

	router.Route("/game", func(game chi.Router) {
		RegisterGameRoutes(game, deps)
	})
}
