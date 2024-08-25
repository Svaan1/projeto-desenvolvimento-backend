package spotify

import (
	"fmt"
	"log"
	"net/http"
)

type Service interface {
	GetAlbums(albumIds []string) (AlbumResponse, error)
	GetTracks(trackIds []string) (TrackResponse, error)
	GetArtists(artistIds []string) (ArtistResponse, error)
}

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

//////////////////////////////////////////////////////////////////////////
// ALL OF THESE WILL PROBABLY NOT NEED A HANDLER IN THE FUTURE 			//
// BUT IT'S GOOD TO HAVE IT FOR NOW TO TEST THE SERVICE AND  			//
// TO HAVE AN EXAMPLE OF HOW TO SET THE SERVICE-TO-HANDLER INTERACTION! //
//  																	//
// EVERYTHING IS MOCKED! 												//
//////////////////////////////////////////////////////////////////////////

func (h *Handler) GetAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	albums, err := h.Service.GetAlbums([]string{"7pp0eBrLEcmprISZOmY4ve", "4aQI0KeN26jZ5bPQXKtfAa", "0qWKuuPbJcUOEcwkUPPnrD"})
	if err != nil {
		log.Printf("error getting albums: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{ "albums": %v }`, albums)
}

func (h *Handler) GetTracksHandler(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.Service.GetTracks([]string{"4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M", "2VQc9orzwE6a5qFfy54P6e"})
	if err != nil {
		log.Printf("error getting tracks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{ "tracks": %v }`, tracks)
}

func (h *Handler) GetArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := h.Service.GetArtists([]string{"0OdUWJ0sBjDrqHygGUXeCF", "3dBVyJ7JuOMt4GE9607Qin", "7dGJo4pcD2V6oG8kP0tJRR"})
	if err != nil {
		log.Printf("error getting artists: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{ "artists": %v }`, artists)
}
