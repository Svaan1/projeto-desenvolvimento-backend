package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand/v2"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
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

		r.interpretMessage(msg)

		r.mu.Unlock()
	}
}

func (r *Room) newGame() {
	log.Println("Starting new game")

	var players []Player
	for conn := range r.connections {
		players = append(players, conn.player)
	}

	r.game = Game{
		Rounds:  []Round{},
		Players: players,
	}
	r.game.newRound()
}

func (g *Game) newRound() {
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	source := g.Players[r.IntN(len(g.Players))]
	round := Round{
		Source: source,
	}
	g.Rounds = append(g.Rounds, round)
}

func (r *Room) interpretMessage(msg []byte) {
	var message WSMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("error:", err)
		return
	}

	switch message.Type {
	// case MessageTypeJoin:
	// 	// add user to room
	// case MessageTypeLeave:
	// 	// remove user from room
	// case MessageTypeChat:
	// 	// show message to all users
	case MessageTypeStart:
		r.newGame()
		res, err := json.Marshal(WSMessage{
			Type:     MessageTypeSystem,
			Data:     "game started",
			PlayerID: message.PlayerID,
		})
		if err != nil {
			log.Println("error:", err)
			return
		}

		for conn := range r.connections {
			err := conn.ws.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				conn.ws.Close()
				delete(r.connections, conn)
				log.Println("error:", err)
			}
		}

	case MessageTypeGuess:
		currentRound := r.game.Rounds[len(r.game.Rounds)-1]
		isCorrectGuess := message.Data == currentRound.Source.ID

		var result string
		if isCorrectGuess {
			result = "correct"
		} else {
			result = "incorrect"
		}

		res, err := json.Marshal(WSMessage{
			Type:     MessageTypeResult,
			Data:     result,
			PlayerID: message.PlayerID,
		})
		if err != nil {
			log.Println("error:", err)
			return
		}

		for conn := range r.connections {
			err = conn.ws.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				conn.ws.Close()
				delete(r.connections, conn)
			}
		}
	}
}
