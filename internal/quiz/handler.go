package quiz

import (
	"encoding/json"
	"net/http"
)

type QuizService interface {
	GetTodaysQuiz() Quiz
}

type Handler struct {
	Service QuizService
}

func NewHandler(s QuizService) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) GetTodaysQuizHandler(w http.ResponseWriter, r *http.Request) {
	quiz := h.Service.GetTodaysQuiz()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}
