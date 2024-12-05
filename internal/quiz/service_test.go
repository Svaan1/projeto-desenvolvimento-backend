package quiz

import (
	"backendProject/internal/db"
	"backendProject/internal/spotify"
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func TestGetTodaysQuiz(t *testing.T) {
	// Setup the quiz service
	ctx := context.Background()
	db, err := db.NewSQLiteDB(ctx, ":memory:")
	if err != nil {
		log.Fatalf("error connecting to in memory db: %v", err)
	}
	defer db.Close()
	repo := NewRepository(db)
	spotifyService := spotify.NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	quizService := NewService(repo, spotifyService)

	// Get today's quiz
	quiz, err := quizService.GetTodaysQuiz(ctx)
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

	fmt.Printf("Quiz: %v\n", quiz)
}

func TestGetTodaysQuizTwice(t *testing.T) {
	// Setup the quiz service
	ctx := context.Background()
	db, err := db.NewSQLiteDB(ctx, ":memory:")
	if err != nil {
		log.Fatalf("error connecting to in memory db: %v", err)
	}
	defer db.Close()
	repo := NewRepository(db)
	spotifyService := spotify.NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	quizService := NewService(repo, spotifyService)

	// Get today's quiz
	quiz, err := quizService.GetTodaysQuiz(ctx)
	if err != nil {
		log.Fatalf("error getting today's quiz")
	}

	// Get today's quiz again
	quiz2, err := quizService.GetTodaysQuiz(ctx)
	if err != nil {
		log.Fatalf("error getting today's quiz")
	}

	// Check if both quizzes are the same
	if quiz.Track.ID != quiz2.Track.ID {
		t.Errorf("Expected both quizzes to have the same track ID, got different IDs")
	}

	if quiz.CreatedAt.Compare(quiz2.CreatedAt) != 0 {
		t.Errorf("Expected both quizzes to have the same creation time, got different times")
		log.Printf("quiz 1: %v, quiz 2: %v", quiz.CreatedAt, quiz2.CreatedAt)
	}
}

func TestGetRandomTrack(t *testing.T) {
	// Setup the quiz service
	ctx := context.Background()
	db, err := db.NewSQLiteDB(ctx, ":memory:")
	if err != nil {
		log.Fatalf("error connecting to in memory db: %v", err)
	}
	defer db.Close()
	repo := NewRepository(db)
	spotifyService := spotify.NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	quizService := NewService(repo, spotifyService)

	// Get a random track based on Wish You Were Here by pink floyd
	track, err := quizService.getRandomTrack([]string{"0k17h0D3J5VfsdmQ1iZtE9"}, "6mFkJmJqdDVQ1REhVfGgd1")
	if err != nil {
		switch err.(type) {
		case *spotify.ErrRecommendationsEmpty:
			// Do nothing, this is expected
			return
		default:
			t.Errorf("Error getting random track: %v", err)
		}
	}

	// Check if the track has the correct fields
	if track.Name == "" {
		t.Errorf("Expected track to have a name, got empty string")
	}
	if track.ID == "" {
		t.Errorf("Expected track to have an ID, got empty string")
	}
	if len(track.Album.Artists) == 0 {
		t.Errorf("Expected track to have artists, got 0")
	}
	if track.PreviewURL == "" {
		t.Errorf("Expected track to have a preview URL, got empty string")
	}
}
