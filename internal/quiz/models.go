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

func (q Quiz) String() string {
	var result string
	result += "Quiz created at: " + q.CreatedAt.String() + "\n"
	result += "Artists:\n"
	for _, artist := range q.Artists {
		result += artist.String() + "\n"
	}
	result += "Album: " + q.Album.String() + "\n"
	result += "Track: " + q.Track.String() + "\n"
	return result
}

func (a quizArtist) String() string {
	var result string
	result += a.Name + " with genres "
	for i, genre := range a.Genres {
		if i > 0 {
			result += ", "
		}
		result += genre
	}
	result += "\n"
	return result
}

func (a quizAlbum) String() string {
	var result string
	result += a.Name + " released on " + a.ReleaseDate + "\n"
	return result
}

func (s quizSong) String() string {
	var result string
	result += s.Name + "\n"
	return result
}
