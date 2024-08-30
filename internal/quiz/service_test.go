package quiz

import (
	"backendProject/internal/spotify"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var quizService *service

func init() {
	godotenv.Load("../../.env")
	spotifyService := spotify.NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	quizService = NewService(spotifyService)
}

func TestGetTodaysQuiz(t *testing.T) {
	// Get today's quiz
	quiz, err := quizService.GetTodaysQuiz()
	if err != nil {
		t.Errorf("Error getting today's quiz")
	}

	// Check if the quiz has the correct fields
	if len(quiz.Artists) == 0 {
		t.Errorf("Expected quiz to have artists, got 0")
	}
	if quiz.Album.Name == "" {
		t.Errorf("Expected quiz to have an album name, got empty string")
	}
	if quiz.Track.Name == "" {
		t.Errorf("Expected quiz to have a track name, got empty string")
	}
	if quiz.CreatedAt.IsZero() {
		t.Errorf("Expected quiz to have a creation time, got zero time")
	}
	if time.Since(quiz.CreatedAt) > 24*time.Hour {
		t.Errorf("Expected quiz to be created within the last 24 hours")
	}
}

func TestGetTodaysQuizTwice(t *testing.T) {
	// Get today's quiz
	quiz, err := quizService.GetTodaysQuiz()
	if err != nil {
		t.Errorf("Error getting today's quiz")
	}

	// Get today's quiz again
	quiz2, err := quizService.GetTodaysQuiz()
	if err != nil {
		t.Errorf("Error getting today's quiz")
	}

	// Check if both quizzes are the same
	if quiz.Track.ID != quiz2.Track.ID {
		t.Errorf("Expected both quizzes to have the same track ID, got different IDs")
	}
	if quiz.CreatedAt != quiz2.CreatedAt {
		t.Errorf("Expected both quizzes to have the same creation time, got different times")
	}
}
