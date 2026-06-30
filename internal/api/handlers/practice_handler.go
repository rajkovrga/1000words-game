package handlers

import (
	"net/http"

	"1000words-game/internal/api/requests"
	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	gameModels "1000words-game/models/game"
)

type PracticeHandler struct {
	practiceService *services.PracticeService
}

type practiceStartResponse struct {
	LevelID          int                        `json:"level_id"`
	TargetLanguageID int                        `json:"target_language_id"`
	NativeLanguageID int                        `json:"native_language_id"`
	Questions        []practiceQuestionResponse `json:"questions"`
	TotalQuestions   int                        `json:"total_questions"`
	Mode             string                     `json:"mode"`
}

type practiceQuestionResponse struct {
	WordID int    `json:"word_id"`
	Word   string `json:"word"`
	Answer string `json:"answer"`
}

func NewPracticeHandler(
	practiceService *services.PracticeService,
) *PracticeHandler {
	return &PracticeHandler{
		practiceService: practiceService,
	}
}

func (h *PracticeHandler) Start(w http.ResponseWriter, r *http.Request) {
	var request requests.StartPracticeRequest

	if err := requests.DecodeJSON(r, &request); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	session, err := h.practiceService.StartPractice(
		request.LevelNumber,
		request.TargetLanguageCode,
		request.NativeLanguageCode,
	)
	if err != nil {
		responses.BadRequest(w, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, mapPracticeSession(session))
}

func mapPracticeSession(session *gameModels.Session) practiceStartResponse {
	questions := make([]practiceQuestionResponse, 0, len(session.Questions))

	for _, question := range session.Questions {
		questions = append(questions, practiceQuestionResponse{
			WordID: question.WordID,
			Word:   question.Word,
			Answer: question.Answer,
		})
	}

	return practiceStartResponse{
		LevelID:          session.LevelID,
		TargetLanguageID: session.TargetLanguageID,
		NativeLanguageID: session.NativeLanguageID,
		Questions:        questions,
		TotalQuestions:   len(questions),
		Mode:             "practice",
	}
}
