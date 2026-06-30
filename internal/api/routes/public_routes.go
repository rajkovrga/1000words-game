package routes

import "github.com/go-chi/chi/v5"

func RegisterPublicRoutes(router chi.Router, deps Dependencies) {
	RegisterHealthRoutes(router, deps)
}
