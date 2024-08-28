package routes

import (
	"fmt"
	"net/http"

	"backendProject/internal/spotify"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello World!"}`)
	})

	// Spotify
	spotifyService := spotify.NewService()
	spotifyHandler := spotify.NewHandler(spotifyService)

	r.Get("/albums", spotifyHandler.GetAlbumsHandler)
	r.Get("/tracks", spotifyHandler.GetTracksHandler)
	r.Get("/artists", spotifyHandler.GetArtistsHandler)
	r.Get("/search", spotifyHandler.SearchHandler)

	return r
}
