package spotify

import (
	"encoding/json"
	"log"
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
	json.NewEncoder(w).Encode(albums)
}

func (h *Handler) GetTracksHandler(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.Service.GetTracks([]string{"2vjoV2tKJMfhLCPjPa9dWt", "1301WleyT98MSxVHPZCA6M", "2VQc9orzwE6a5qFfy54P6e"})
	if err != nil {
		log.Printf("error getting tracks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func (h *Handler) GetArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := h.Service.GetArtists([]string{"0OdUWJ0sBjDrqHygGUXeCF", "3dBVyJ7JuOMt4GE9607Qin", "7dGJo4pcD2V6oG8kP0tJRR"})
	if err != nil {
		log.Printf("error getting artists: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	queryType := r.URL.Query().Get("t")

	if queryType == "" {
		http.Error(w, "missing type parameter", http.StatusBadRequest)
		return
	}

	search, err := h.Service.Search(query, queryType)
	if err != nil {
		log.Printf("error searching: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(search)

}
