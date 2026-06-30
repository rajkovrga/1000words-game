package api

import (
	"database/sql"
	"net/http"
	"time"

	"1000words-game/internal/api/responses"
	"1000words-game/internal/api/routes"
	"1000words-game/internal/config"
	"1000words-game/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDependencies struct {
	DB *sql.DB

	Config config.Config

	AuthService          *services.AuthService
	AuthorizationService *services.AuthorizationService
	UserService          *services.UserService
	LanguageService      *services.LanguageService
	PracticeService      *services.PracticeService
	ProgressService      *services.ProgressService
	GameService          *services.GameService
}

func NewRouter(deps RouterDependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		responses.NotFound(w, "Route not found")
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		responses.Error(
			w,
			http.StatusMethodNotAllowed,
			"method_not_allowed",
			"Method not allowed",
		)
	})

	routes.Register(router, routes.Dependencies{
		DB: deps.DB,

		Config: deps.Config,

		AuthService:          deps.AuthService,
		AuthorizationService: deps.AuthorizationService,
		UserService:          deps.UserService,
		LanguageService:      deps.LanguageService,
		PracticeService:      deps.PracticeService,
		ProgressService:      deps.ProgressService,
		GameService:          deps.GameService,
	})

	return router
}
