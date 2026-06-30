package handlers

import (
	"net/http"
	"time"

	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
)

type LanguageHandler struct {
	languageService *services.LanguageService
}

type languageResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

func NewLanguageHandler(languageService *services.LanguageService) *LanguageHandler {
	return &LanguageHandler{
		languageService: languageService,
	}
}

func (h *LanguageHandler) Index(w http.ResponseWriter, r *http.Request) {
	languages, err := h.languageService.GetAllLanguages()
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"languages": mapLanguages(languages),
	})
}

func mapLanguages(languages []dbModels.Language) []languageResponse {
	response := make([]languageResponse, 0, len(languages))

	for _, language := range languages {
		response = append(response, languageResponse{
			ID:        language.ID,
			Name:      language.Name,
			Code:      language.Code,
			CreatedAt: language.CreatedAt,
		})
	}

	return response
}
