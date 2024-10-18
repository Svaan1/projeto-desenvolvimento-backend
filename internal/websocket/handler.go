package websocket

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
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

// Handle handles the websocket request.
//
// Parameters:
//   - w: the response writer.
//   - r: the request.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room")
	isAdmin := r.URL.Query().Get("admin") == "true"

	room := h.hub.getRoom(roomID)
	log.Printf("New connection to room %s", roomID)
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade connection failed", http.StatusInternalServerError)
		return
	}

	connection := &Connection{
		ws:    conn,
		room:  room,
		admin: isAdmin,
	}
	room.mu.Lock()
	room.connections[connection] = true
	room.mu.Unlock()

	defer func() {
		room.mu.Lock()
		conn.Close()
		delete(room.connections, connection)
		room.mu.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		room.broadcast <- msg
	}
}
