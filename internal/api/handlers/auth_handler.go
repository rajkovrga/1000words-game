package handlers

import (
	"net/http"
	"time"

	apiMiddleware "1000words-game/internal/api/middleware"
	"1000words-game/internal/api/requests"
	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type authUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type authResponse struct {
	Token     string           `json:"token"`
	ExpiresAt time.Time        `json:"expires_at"`
	User      authUserResponse `json:"user"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request requests.RegisterRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	result, err := h.authService.Register(request.Email, request.Password)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.Created(w, mapAuthResult(result))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request requests.LoginRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	result, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		responses.Unauthorized(w, "Invalid email or password")
		return
	}

	responses.JSON(w, http.StatusOK, mapAuthResult(result))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := apiMiddleware.BearerTokenFromRequest(r)
	if err != nil {
		responses.Unauthorized(w, "Authentication token is required")
		return
	}

	if err := h.authService.Logout(token); err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

func mapAuthResult(result *services.AuthResult) authResponse {
	return authResponse{
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
		User: authUserResponse{
			ID:    result.User.ID,
			Email: result.User.Email,
		},
	}
}
