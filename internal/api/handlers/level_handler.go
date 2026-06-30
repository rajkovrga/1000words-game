package handlers

import (
	"net/http"
	"time"

	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
)

type LevelHandler struct {
	gameService *services.GameService
}

type levelResponse struct {
	ID              int       `json:"id"`
	LevelNumber     int       `json:"level_number"`
	Name            string    `json:"name"`
	WordsRequired   int       `json:"words_required"`
	MaxWrongAnswers int       `json:"max_wrong_answers"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewLevelHandler(gameService *services.GameService) *LevelHandler {
	return &LevelHandler{
		gameService: gameService,
	}
}

func (h *LevelHandler) Index(w http.ResponseWriter, r *http.Request) {
	levels, err := h.gameService.GetAvailableLevels()
	if err != nil {
		responses.InternalServerError(w)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"levels": mapLevels(levels),
	})
}

func mapLevels(levels []dbModels.Level) []levelResponse {
	response := make([]levelResponse, 0, len(levels))

	for _, level := range levels {
		response = append(response, levelResponse{
			ID:              level.ID,
			LevelNumber:     level.LevelNumber,
			Name:            level.Name,
			WordsRequired:   level.WordsRequired,
			MaxWrongAnswers: level.MaxWrongAnswers,
			CreatedAt:       level.CreatedAt,
		})
	}

	return response
}
