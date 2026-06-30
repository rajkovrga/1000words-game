package handlers

import (
	"net/http"

	apiMiddleware "1000words-game/internal/api/middleware"
	"1000words-game/internal/api/requests"
	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	gameModels "1000words-game/models/game"
)

type ProgressHandler struct {
	progressService *services.ProgressService
}

type progressResponse struct {
	ID int `json:"id"`

	TargetLanguage languageShortResponse `json:"target_language"`
	NativeLanguage languageShortResponse `json:"native_language"`

	CurrentLevel levelShortResponse `json:"current_level"`
}

type languageShortResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type levelShortResponse struct {
	ID          int    `json:"id"`
	LevelNumber int    `json:"level_number"`
	Name        string `json:"name"`
}

func NewProgressHandler(
	progressService *services.ProgressService,
) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

func (h *ProgressHandler) Index(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	options, err := h.progressService.GetUserProgressOptions(user.ID)
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"progress": mapProgressOptions(options),
	})
}

func (h *ProgressHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	var request requests.CreateProgressRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	progress, err := h.progressService.CreateProgress(
		user.ID,
		request.TargetLanguageCode,
		request.NativeLanguageCode,
	)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.Created(w, mapProgressOption(*progress))
}

func mapProgressOptions(options []gameModels.ProgressOption) []progressResponse {
	response := make([]progressResponse, 0, len(options))

	for _, option := range options {
		response = append(response, mapProgressOption(option))
	}

	return response
}

func mapProgressOption(option gameModels.ProgressOption) progressResponse {
	return progressResponse{
		ID: option.ProgressID,

		TargetLanguage: languageShortResponse{
			ID:   option.TargetLanguageID,
			Code: option.TargetCode,
			Name: option.TargetName,
		},

		NativeLanguage: languageShortResponse{
			ID:   option.NativeLanguageID,
			Code: option.NativeCode,
			Name: option.NativeName,
		},

		CurrentLevel: levelShortResponse{
			ID:          option.CurrentLevelID,
			LevelNumber: option.LevelNumber,
			Name:        "Level",
		},
	}
}
