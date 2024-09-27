package quiz

import (
	"context"
	"time"
)

type Quiz struct {
	Artists   []quizArtist `json:"artists"`
	Album     quizAlbum    `json:"album"`
	Track     quizSong     `json:"track"`
	CreatedAt time.Time    `json:"created_at"`
}

type quizArtist struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type quizAlbum struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	ReleaseDate string `json:"release_date"`
}

type quizSong struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AudioPreview string `json:"audio_preview"`
}
type Service interface {
	GetTodaysQuiz(context.Context) (Quiz, error)
}
