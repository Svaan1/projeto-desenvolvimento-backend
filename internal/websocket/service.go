package websocket

import (
	"errors"
	"log"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/websocket"
)

func (h *Hub) listRoomCodes() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	var rooms = make([]string, 0, len(h.rooms))
	for roomID := range h.rooms {
		rooms = append(rooms, roomID)
	}

	return rooms
}

func (h *Hub) getRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[roomID]
	return room
}

func (h *Hub) createRoom(roomID string, password string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, exists := h.rooms[roomID]
	if exists {
		return errors.New("Room already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	h.rooms[roomID] = &Room{
		connections: make(map[*Connection]bool),
		broadcast:   make(chan []byte),
		mu:          sync.Mutex{},
		password:    hashedPassword,
	}
	log.Printf("Created room [%s]", roomID)
	go h.rooms[roomID].run()
	return nil
}

func (r *Room) run() {
	for {
		msg := <-r.broadcast

		r.mu.Lock()
		for conn := range r.connections {
			err := conn.ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				conn.ws.Close()
				delete(r.connections, conn)
			}
		}
		r.mu.Unlock()
	}
}
