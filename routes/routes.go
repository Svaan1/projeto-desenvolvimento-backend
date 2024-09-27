package routes

import (
	"fmt"
	"net/http"
	"os"

	"backendProject/internal/db"
	"backendProject/internal/quiz"
	"backendProject/internal/spotify"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	baseURL = "/api/v1"
)

func NewRouter(db db.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello World!"}`)
	})

	// Spotify
	spotifyService := spotify.NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	spotifyHandler := spotify.NewHandler(spotifyService)

	r.Get("/albums", spotifyHandler.GetAlbumsHandler)
	r.Get("/tracks", spotifyHandler.GetTracksHandler)
	r.Get("/artists", spotifyHandler.GetArtistsHandler)
	r.Get("/search", spotifyHandler.SearchHandler)

	// Quiz
	quizRepository := quiz.NewRepository(db)
	quizService := quiz.NewService(quizRepository, spotifyService)
	quizHandler := quiz.NewHandler(quizService)

	r.Get(baseURL+"/quiz", quizHandler.GetTodaysQuizHandler)

	return r
}
