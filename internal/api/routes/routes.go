package routes

import "github.com/go-chi/chi/v5"

func Register(router chi.Router, deps Dependencies) {
	router.Route("/api", func(api chi.Router) {
		RegisterPublicRoutes(api, deps)

		api.Route("/v1", func(v1 chi.Router) {
			RegisterV1Routes(v1, deps)
		})
	})
}
