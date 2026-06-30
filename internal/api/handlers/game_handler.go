package handlers

import (
	"fmt"
	"net/http"

	apiMiddleware "1000words-game/internal/api/middleware"
	"1000words-game/internal/api/requests"
	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
)

type GameHandler struct {
	progressService *services.ProgressService
	gameService     *services.GameService
}

type gameStartResponse struct {
	AttemptID           int                    `json:"attempt_id"`
	ProgressID          int                    `json:"progress_id"`
	Level               gameLevelResponse      `json:"level"`
	AllowedWrongAnswers int                    `json:"allowed_wrong_answers"`
	Questions           []gameQuestionResponse `json:"questions"`
	TotalQuestions      int                    `json:"total_questions"`
	Mode                string                 `json:"mode"`
}

type gameLevelResponse struct {
	ID          int    `json:"id"`
	LevelNumber int    `json:"level_number"`
	Name        string `json:"name"`
}

type gameQuestionResponse struct {
	WordID int    `json:"word_id"`
	Word   string `json:"word"`
	Answer string `json:"answer"`
}

type gameFinishResponse struct {
	AttemptID           int                `json:"attempt_id"`
	ProgressID          int                `json:"progress_id"`
	CorrectCount        int                `json:"correct_count"`
	WrongCount          int                `json:"wrong_count"`
	TotalQuestions      int                `json:"total_questions"`
	AllowedWrongAnswers int                `json:"allowed_wrong_answers"`
	Passed              bool               `json:"passed"`
	NextLevel           *gameLevelResponse `json:"next_level,omitempty"`
}

func NewGameHandler(
	progressService *services.ProgressService,
	gameService *services.GameService,
) *GameHandler {
	return &GameHandler{
		progressService: progressService,
		gameService:     gameService,
	}
}

func (h *GameHandler) Start(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	var request requests.StartGameRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	progress, err := h.progressService.GetUserProgressOptionByID(
		user.ID,
		request.ProgressID,
	)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	session, err := h.gameService.StartLevel(user.ID, *progress)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, mapGameStartResponse(session, progress))
}

func (h *GameHandler) Finish(w http.ResponseWriter, r *http.Request) {
	user, ok := apiMiddleware.CurrentUser(r)
	if !ok {
		responses.Unauthorized(w, "User is not authenticated")
		return
	}

	var request requests.FinishGameRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	progress, err := h.progressService.GetUserProgressOptionByID(
		user.ID,
		request.ProgressID,
	)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	result, err := h.gameService.FinishGame(
		user.ID,
		*progress,
		request.AttemptID,
		request.CorrectCount,
		request.WrongCount,
		request.TotalQuestions,
	)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, mapGameFinishResponse(result))
}

func mapGameStartResponse(
	session *gameModels.Session,
	progress *gameModels.ProgressOption,
) gameStartResponse {
	questions := make([]gameQuestionResponse, 0, len(session.Questions))

	for _, question := range session.Questions {
		questions = append(questions, gameQuestionResponse{
			WordID: question.WordID,
			Word:   question.Word,
			Answer: question.Answer,
		})
	}

	return gameStartResponse{
		AttemptID:  session.AttemptID,
		ProgressID: session.ProgressID,
		Level: gameLevelResponse{
			ID:          session.LevelID,
			LevelNumber: progress.LevelNumber,
			Name:        fmt.Sprintf("Level %d", progress.LevelNumber),
		},
		AllowedWrongAnswers: session.MaxWrongAnswers,
		Questions:           questions,
		TotalQuestions:      len(questions),
		Mode:                "game",
	}
}

func mapGameFinishResponse(result *gameModels.FinishGameResult) gameFinishResponse {
	var nextLevel *gameLevelResponse

	if result.NextLevel != nil {
		nextLevel = mapLevelModel(result.NextLevel)
	}

	return gameFinishResponse{
		AttemptID:           result.AttemptID,
		ProgressID:          result.ProgressID,
		CorrectCount:        result.CorrectCount,
		WrongCount:          result.WrongCount,
		TotalQuestions:      result.TotalQuestions,
		AllowedWrongAnswers: result.MaxWrongAnswers,
		Passed:              result.Passed,
		NextLevel:           nextLevel,
	}
}

func mapLevelModel(level *dbModels.Level) *gameLevelResponse {
	if level == nil {
		return nil
	}

	return &gameLevelResponse{
		ID:          level.ID,
		LevelNumber: level.LevelNumber,
		Name:        level.Name,
	}
}
