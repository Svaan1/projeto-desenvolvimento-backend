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

func (h *Handler) GetTodaysQuizHandler(w http.ResponseWriter, r *http.Request) {
	quiz := h.Service.GetTodaysQuiz()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}
