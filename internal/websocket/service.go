package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func (h *Hub) getRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[roomID]
	if exists {
		return room
	}

	room = &Room{
		connections: make(map[*Connection]bool),
		broadcast:   make(chan []byte),
		mu:          sync.Mutex{},
	}
	h.rooms[roomID] = room
	go room.run()

	return room
}

func (r *Room) run() {
	for {
		msg := <-r.broadcast
		log.Printf("Broadcasting message: %s", string(msg))

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
