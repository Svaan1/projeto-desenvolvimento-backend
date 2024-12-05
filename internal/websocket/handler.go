package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles websocket requests from the peer.
type Handler struct {
	upgrader websocket.Upgrader
	hub      *Hub
}

// NewHandler creates a new Handler.
//
// Returns:
//   - A new websocket Handler.
func NewHandler() *Handler {
	hub := &Hub{
		rooms: make(map[string]*Room),
	}

	return &Handler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		hub: hub,
	}
}

// HandleWS handles the websocket connections.
func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room")
	password := r.Header.Get("Room-Password")
	isAdmin := r.Header.Get("Room-Admin") == "true"
	spotifyToken := r.Header.Get("Spotify-Token")

	// TODO: actually use the token
	if spotifyToken == "" {
		spotifyToken = "invalid"
	}

	// check for missing fields
	if roomID == "" {
		http.Error(w, "room not specified", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "password not specified", http.StatusBadRequest)
		return
	}

	// get room from url param
	room := h.hub.getRoom(roomID)
	if room == nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	// check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(room.password), []byte(password))
	if err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	// upgrade connection and add to room
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade connection failed", http.StatusInternalServerError)
		return
	}

	connection := &Connection{
		ws:   conn,
		room: room,
		player: Player{
			ID: "123", // TODO: get user ID from API using the token
			// SpotifyToken: spotify.Token{AccessToken: spotifyToken},
			IsAdmin: isAdmin,
		},
	}
	room.mu.Lock()
	room.connections[connection] = true
	room.mu.Unlock()

	defer func() {
		room.mu.Lock()

		conn.Close()
		delete(room.connections, connection)
		isRoomEmpty := len(room.connections) == 0

		room.mu.Unlock()

		// delete empty rooms
		if isRoomEmpty {
			delete(h.hub.rooms, roomID)
			log.Printf("Deleted room [%s]", roomID)
		}
	}()

	// everytime a message is received, broadcast it to all connections in the room
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		room.broadcast <- msg
	}
}

// ListRoomCodes returns a list of all rooms.
//
// Returns:
//   - A list of all rooms.
func (h *Handler) ListRoomCodes(w http.ResponseWriter, r *http.Request) {
	roomCodes := h.hub.listRoomCodes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomCodes)
}

// CreateRoom creates a new room.
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	roomID := req.RoomID
	password := req.Password
	err := h.hub.createRoom(roomID, password)
	if err != nil {
		http.Error(w, "room already exists", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
