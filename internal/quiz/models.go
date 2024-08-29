package quiz

import (
	"time"
)

type Quiz struct {
	Artists   []quizArtist `json:"artists"`
	Album     quizAlbum    `json:"album"`
	Track     quizSong     `json:"track"`
	CreatedAt time.Time    `json:"created_at"`
}

type quizArtist struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type quizAlbum struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	ReleaseDate string `json:"release_date"`
}

type quizSong struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	AudioPreview string `json:"audio_preview"`
}
type Service interface {
	GetTodaysQuiz() (Quiz, error)
}
