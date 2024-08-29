package quiz

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// GetTodaysQuizHandler returns today's quiz.
//
// Returns:
//   - A JSON object containing the quiz data.
func (h *Handler) GetTodaysQuizHandler(w http.ResponseWriter, r *http.Request) {
	quiz, err := h.Service.GetTodaysQuiz()
	if err != nil {
		http.Error(w, "Error getting today's quiz", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}
