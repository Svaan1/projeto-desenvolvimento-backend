package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active rooms and broadcasts messages to the rooms.
type Hub struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

// Room represents a game room that users can join. It contains a set of connections and a broadcast channel.
type Room struct {
	connections map[*Connection]bool
	broadcast   chan []byte
	mu          sync.Mutex
	password    []byte
}

// Connection represents a websocket connection to a room.
type Connection struct {
	ws    *websocket.Conn
	room  *Room
	admin bool
}

type CreateRoomRequest struct {
	RoomID   string `json:"room"`
	Password string `json:"password"`
}
