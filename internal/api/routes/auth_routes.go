package routes

import (
	"1000words-game/internal/api/handlers"
	apiMiddleware "1000words-game/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(router chi.Router, deps Dependencies) {
	authHandler := handlers.NewAuthHandler(deps.AuthService)

	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)

	router.With(apiMiddleware.RequireAuth(deps.AuthService)).
		Post("/logout", authHandler.Logout)
}
