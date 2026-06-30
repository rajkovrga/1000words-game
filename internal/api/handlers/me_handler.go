package handlers

import (
	"net/http"

	apiMiddleware "1000words-game/internal/api/middleware"
	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
)

type MeHandler struct {
	authorizationService *services.AuthorizationService
}

type meResponse struct {
	ID          int      `json:"id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func NewMeHandler(
	authorizationService *services.AuthorizationService,
) *MeHandler {
	return &MeHandler{
		authorizationService: authorizationService,
	}
}

func (h *MeHandler) Show(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	roles, err := h.authorizationService.GetUserRoleCodes(user.ID)
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	permissions, err := h.authorizationService.GetUserPermissionCodes(user.ID)
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, meResponse{
		ID:          user.ID,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
	})
}

func (h *MeHandler) Roles(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	roles, err := h.authorizationService.GetUserRoleCodes(user.ID)
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"roles": roles,
	})
}

func (h *MeHandler) Permissions(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	permissions, err := h.authorizationService.GetUserPermissionCodes(user.ID)
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"permissions": permissions,
	})
}
