package websocket

import (
	"backendProject/internal/spotify"
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
	game        Game
}

// Connection represents a websocket connection to a room.
type Connection struct {
	ws     *websocket.Conn
	room   *Room
	player Player
}

type CreateRoomRequest struct {
	RoomID   string `json:"room"`
	Password string `json:"password"`
}

type WSMessage struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	PlayerID string `json:"from"` // the player's ID who sent the message
}

const (
	// MessageTypeJoin   = "join"
	// MessageTypeLeave  = "leave"
	// MessageTypeChat   = "chat"
	MessageTypeStart  = "start"  // message that represents the start of a new game
	MessageTypeGuess  = "guess"  // message that represents a guess on the game
	MessageTypeResult = "result" // message that represents the outcome of a guess (correct/incorrect)
	MessageTypeError  = "error"
	MessageTypeSystem = "system" // message that represents a system message (e.g. game started, game ended, etc.)
)

// Player represents a player in the room.
type Player struct {
	ID           string        `json:"id"` // the user's spotify ID
	SpotifyToken spotify.Token `json:"spotifyToken"`
	IsAdmin      bool          `json:"isAdmin"`
}

// Game represents a game in the room.
type Game struct {
	Players []Player `json:"players"`
	Rounds  []Round  `json:"rounds"`
}

// Round represents a round in the game.
type Round struct {
	Track  spotify.Track `json:"track"`  // the track that the players will have as a reference to guess the source player
	Source Player        `json:"source"` // the player from whose the track was chosen
	Winner Player        `json:"winner"` // the player who guessed the user correctly first
}
